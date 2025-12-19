package handler

import (
	"log/slog"
	"net/http"
	"time"
)

func InitRouter(handler *WalletHandler) http.Handler {
	mux := http.NewServeMux()

	mux.HandleFunc("POST /api/v1/wallet", handler.HandleOperation)
	mux.HandleFunc("GET /api/v1/wallets/{id}", handler.HandleGetBalance)

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		mux.ServeHTTP(w, r)
		
		slog.Info("request",
			"method", r.Method,
			"path", r.URL.Path,
			"duration", time.Since(start),
		)
	})
}
