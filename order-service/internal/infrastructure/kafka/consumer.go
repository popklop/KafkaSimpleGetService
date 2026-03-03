package kafka

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/segmentio/kafka-go"
	"log"
	"time"
	orderr "wbtech/internal/domain/order"
	"wbtech/internal/usecase/order"
)

const maxRetries = 3

type Consumer struct {
	reader    *kafka.Reader
	dlqWriter *kafka.Writer
	useCase   *order.OrderUseCase
	done      chan struct{}
}

func NewConsumer(broker, topic, group string, uc *order.OrderUseCase) *Consumer {
	reader := kafka.NewReader(kafka.ReaderConfig{
		Brokers:  []string{broker},
		Topic:    topic,
		GroupID:  group,
		MinBytes: 1e3,
		MaxBytes: 10e6,
	})

	dlqWriter := &kafka.Writer{
		Addr:     kafka.TCP(broker),
		Topic:    topic + "_dlq",
		Balancer: &kafka.LeastBytes{},
	}

	return &Consumer{
		reader:    reader,
		dlqWriter: dlqWriter,
		useCase:   uc,
		done:      make(chan struct{}),
	}
}

func (c *Consumer) Start(ctx context.Context) {
	go func() {
		log.Println("Kafka consumer started")
		for {
			select {
			case <-ctx.Done():
				log.Println("Kafka consumer stopped by context")
				return
			case <-c.done:
				return
			default:
				c.readMessage(ctx)
			}
		}
	}()
}

func (c *Consumer) readMessage(ctx context.Context) {
	msg, err := c.reader.ReadMessage(ctx)
	if err != nil {
		if errors.Is(err, context.Canceled) || errors.Is(err, context.DeadlineExceeded) {
			return
		}
		log.Printf("kafka read error: %v", err)
		time.Sleep(time.Second)
		return
	}

	log.Printf("message received: key=%s, offset=%d", string(msg.Key), msg.Offset)

	var dto order.OrderDTO
	if err := json.Unmarshal(msg.Value, &dto); err != nil {
		log.Printf("invalid JSON, sending to DLQ: %v", err)
		c.sendToDLQ(ctx, msg)
		c.commitMessage(ctx, msg)
		return
	}

	var lastErr error
	for attempt := 1; attempt <= maxRetries; attempt++ {
		select {
		case <-ctx.Done():
			log.Println("context cancelled, stopping retry loop")
			return
		default:
		}

		err = c.useCase.SaveOrder(ctx, &dto)
		if err == nil {
			log.Printf("order saved: %s, committing", dto.OrderUID)
			c.commitMessage(ctx, msg)
			return
		}
		lastErr = err
		log.Printf("save order error (attempt %d/%d): %v", attempt, maxRetries, err)
		if errors.Is(err, orderr.ErrInvalidMail) ||
			errors.Is(err, orderr.ErrNoItems) ||
			errors.Is(err, orderr.ErrEmptyID) ||
			errors.Is(err, orderr.ErrInvalidPayment) {
			log.Printf("permanent error, sending to DLQ")
			c.sendToDLQ(ctx, msg)
			c.commitMessage(ctx, msg)
			return
		}

		if attempt < maxRetries {
			wait := time.Duration(attempt*attempt) * time.Second
			log.Printf("waiting %v before retry", wait)
			time.Sleep(wait)
		}
	}

	log.Printf("save order failed after %d attempts, sending to DLQ: %v", maxRetries, lastErr)
	c.sendToDLQ(ctx, msg)
	c.commitMessage(ctx, msg)
}

func (c *Consumer) sendToDLQ(ctx context.Context, msg kafka.Message) {
	dlqMsg := kafka.Message{
		Key:   msg.Key,
		Value: msg.Value,
		Headers: append(msg.Headers, kafka.Header{
			Key:   "original_topic",
			Value: []byte(c.reader.Config().Topic),
		}),
	}
	if err := c.dlqWriter.WriteMessages(ctx, dlqMsg); err != nil {
		log.Printf("failed to send message to DLQ: %v", err)
	}
}

func (c *Consumer) commitMessage(ctx context.Context, msg kafka.Message) {
	if err := c.reader.CommitMessages(ctx, msg); err != nil {
		log.Printf("failed to commit message: %v", err)
	}
}

func (c *Consumer) Close() error {
	close(c.done)
	if err := c.reader.Close(); err != nil {
		return err
	}
	return c.dlqWriter.Close()
}
