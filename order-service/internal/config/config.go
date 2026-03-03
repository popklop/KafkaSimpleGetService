package config

import (
	"os"
)

type Config struct {
	DBUser      string
	DBPassword  string
	DBHost      string
	DBPort      string
	DBName      string
	KafkaBroker string
	KafkaTopic  string
	KafkaGroup  string
	HTTPPort    string
}

func Load() Config {
	return Config{
		DBUser:      getEnv("DB_USER", "postgres"),
		DBPassword:  getEnv("DB_PASSWORD", "order_pass"),
		DBHost:      getEnv("DB_HOST", "localhost"),
		DBPort:      getEnv("DB_PORT", "5432"),
		DBName:      getEnv("DB_NAME", "wborders"),
		KafkaBroker: getEnv("KAFKA_BROKER", "localhost:9092"),
		KafkaTopic:  getEnv("KAFKA_TOPIC", "orders"),
		KafkaGroup:  getEnv("KAFKA_GROUP", "order-service"),
		HTTPPort:    getEnv("HTTP_PORT", "8081"),
	}
}

func getEnv(key, fallback string) string {
	if val := os.Getenv(key); val != "" {
		return val
	}
	return fallback
}
