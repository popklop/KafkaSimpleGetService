package http

import (
	"net/http"

	"wbtech/internal/infrastructure/http/handler"

	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func NewRouter(orderHandler *handler.OrderHandler) http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("/order/", orderHandler.GetByID)
	mux.Handle("/metrics", promhttp.Handler())
	fs := http.FileServer(http.Dir("./internal/infrastructure/web"))
	mux.Handle("/", fs)
	return mux
}
