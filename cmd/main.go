package main

import (
	"context"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	_ "task4.2.3/docs"
	"task4.2.3/internal/controller"
	"task4.2.3/internal/handlers"
	"task4.2.3/internal/middleware"
	"task4.2.3/internal/monitoring"
	"task4.2.3/internal/repository"
	"task4.2.3/internal/service"

	"github.com/go-chi/chi/v5"
	"github.com/redis/go-redis/v9"
	httpSwagger "github.com/swaggo/http-swagger"
)

// @title Geo Service API
// @version 1.0
// @description	API для работы с адресами и геокодингом
// @host localhost:8080
// @BasePath /
// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization

func main() {
	monitoring.RegisterMetrics()

	var userStore = handlers.NewUserStore()

	client := &http.Client{}

	redisClient := redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
	})

	baseRepo := repository.NewDaDataRepository("9c667615626123e3c70123efa6ca12e53ae94e06", client)
	repo := repository.NewAddressRepositoryProxy(baseRepo, redisClient, 5*time.Minute)

	service := service.NewAddressService(repo)
	controller := controller.NewAddressController(service)
	addressHandler := handlers.NewAddressHandler(controller)

	r := chi.NewRouter()

	r.Route("/api/address", func(r chi.Router) {
		r.Use(middleware.JWTMiddleware)
		r.Post("/search", addressHandler.SearchHandler)
		r.Post("/geocode", addressHandler.GeocodeHandler)
	})

	r.Post("/api/register", func(w http.ResponseWriter, r *http.Request) {
		handlers.RegisterHandler(w, r, userStore)
	})
	r.Post("/api/login", func(w http.ResponseWriter, r *http.Request) {
		handlers.LoginHandler(w, r, userStore)
	})

	r.Get("/swagger/*", httpSwagger.Handler(
		httpSwagger.URL("http://localhost:8080/swagger/doc.json"),
	))

	monitoring.RegisterPprofRoutes(r)
	monitoring.RegisterGoroutineDumpRoute(r)

	monitoring.RegisterMetricsRoute(r)

	server := &http.Server{
		Addr:         ":8080",
		Handler:      r,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	listner, err := net.Listen("tcp", server.Addr)
	if err != nil {
		log.Fatalf("Error create listner: %v", err)
	}

	stopChan := make(chan os.Signal, 1)
	signal.Notify(stopChan, os.Interrupt, syscall.SIGTERM)

	go func() {
		log.Println("Server started on :8080")
		if err := server.Serve(listner); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Server error: %v", err)
		}
	}()

	<-stopChan
	log.Println("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		log.Fatalf("Server forced to shutdown: %v", err)
	}

	log.Println("Server stopped gracefully")
}
