package router

import (
	"context"
	"encoding/json"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/fenek-dev-code/backend-tasks/internal/controller/handler"
	"github.com/fenek-dev-code/backend-tasks/internal/storage"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"go.uber.org/zap"
)

type Router struct {
	chi.Router
	hand *handler.Handler
	addr string
	log  *zap.Logger
}

func NewRouter(storage *storage.Storage, addr string, log *zap.Logger) *Router {
	r := chi.NewRouter()

	// Базовая middleware
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.RequestID)

	// CORS middleware
	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"https://*", "http://*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type"},
		AllowCredentials: true,
		MaxAge:           300,
	}))

	handlers := handler.NewHandler(storage)
	return &Router{r, handlers, addr, log}
}

func (r *Router) healthCheck() http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		// Add database health check if needed
		r.log.Info("health check", zap.String("status", "OK"))
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]string{
			"status":  "OK",
			"version": "1.0.0",
			"time":    time.Now().Format(time.RFC3339),
		})
	}
}

func (r *Router) Init() {
	r.Get("/health", r.healthCheck())
}

func (r *Router) Run() error {
	server := &http.Server{
		Addr:         r.addr,
		Handler:      r,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	r.log.Info("server: server started", zap.String("addr", r.addr))

	// Graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)

	go func() {
		<-quit
		r.log.Info("server: shutting down...")

		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()

		if err := server.Shutdown(ctx); err != nil {
			r.log.Fatal("Server forced to shutdown: %v", zap.Error(err))
		}
	}()

	return server.ListenAndServe()
}
