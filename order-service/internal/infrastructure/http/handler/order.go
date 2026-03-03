package handler

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"
	"wbtech/internal/usecase/order"
	"wbtech/metrics"
)

type OrderHandler struct {
	uc *order.OrderUseCase
}

func NewOrderHandler(Uc *order.OrderUseCase) *OrderHandler {
	return &OrderHandler{
		uc: Uc,
	}
}

func (h *OrderHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	start := time.Now()
	statusCode := http.StatusOK
	defer func() {
		duration := time.Since(start).Seconds()
		log.Printf("DEBUG: recording metrics for status %d", statusCode)
		metrics.ObserveRequest("GET", "/order", duration)
		metrics.IncRequest("GET", "/order", strconv.Itoa(statusCode))
	}()

	id := strings.TrimPrefix(r.URL.Path, "/order/")
	defer func() {
		log.Printf("⏱️ GET /order/%s completed in %v", id, time.Since(start))
	}()

	path := r.URL.Path
	if id == "" || id == path {
		statusCode = http.StatusBadRequest
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(statusCode)
		json.NewEncoder(w).Encode(map[string]string{"error": "order id is missing"})
		return
	}

	orderData, err := h.uc.GetOrder(r.Context(), id)
	if err != nil {
		if errors.Is(err, order.ErrOrderNotFound) {
			statusCode = http.StatusNotFound
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(statusCode)
			json.NewEncoder(w).Encode(map[string]string{"error": "order not found"})
			return
		}
		statusCode = http.StatusInternalServerError
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(statusCode)
		json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(orderData); err != nil {
		statusCode = http.StatusInternalServerError
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(statusCode)
		json.NewEncoder(w).Encode(map[string]string{"error": "failed to encode to json"})
		return
	}
}
