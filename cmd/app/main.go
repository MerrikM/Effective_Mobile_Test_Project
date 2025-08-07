package main

import (
	_ "Effective_Mobile_Test_Project/docs"
	"Effective_Mobile_Test_Project/internal/config"
	"Effective_Mobile_Test_Project/internal/handler"
	"Effective_Mobile_Test_Project/internal/repository"
	"Effective_Mobile_Test_Project/internal/service"
	"context"
	"fmt"
	"github.com/go-chi/chi/v5"
	httpSwagger "github.com/swaggo/http-swagger"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

// @title           Subscription API
// @version         1.0
// @description     Сервис управления подписками
// @host      localhost:8080
// @BasePath  /
// @schemes http
func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	cfg, err := config.LoadConfig("config.yaml")
	if err != nil {
		log.Fatalf("ошибка загрузки конфига: %v", err)
	}

	database, err := config.SetupDatabase(cfg.DatabaseConfig.DSN)
	if err != nil {
		log.Fatalf("не удалось подключиться к БД: %v", err)
	}
	defer func() {
		if err := database.Close(); err != nil {
			log.Printf("ошибка при закрытии БД: %v", err)
		}
	}()

	subscriptionRepository := repository.NewSubscriptionRepository(database)
	subscriptionService := service.NewSubscriptionService(subscriptionRepository)
	subscriptionHandler := handler.NewSubscriptionHandler(subscriptionService)

	restServer, router := config.SetupRestServer(cfg.ServerAddr)

	router.Get("/swagger/*", httpSwagger.WrapHandler)

	router.Route("/subscriptions", func(r chi.Router) {
		r.Post("/create", subscriptionHandler.Create)
		r.Get("/user/{uuid}", subscriptionHandler.GetByUserUUID)
		r.Get("/get/{id}", subscriptionHandler.GetByID)
		r.Put("/update/{id}", subscriptionHandler.UpdateByID)
		r.Delete("/delete/{id}", subscriptionHandler.DeleteByID)
		r.Get("/total-cost", subscriptionHandler.GetTotalCost)
	})

	runServer(ctx, restServer)
}

func runServer(ctx context.Context, server *http.Server) {
	serverErrors := make(chan error, 1)
	go func() {
		log.Println("Сервер запущен на " + server.Addr)
		fmt.Println("____________________________________________________")
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			serverErrors <- err
		}
	}()

	signalChannel := make(chan os.Signal, 1)
	signal.Notify(signalChannel, syscall.SIGINT, syscall.SIGTERM)

	select {
	case err := <-serverErrors:
		log.Fatalf("ошибка работы сервера: %v", err)
	case sig := <-signalChannel:
		log.Printf("получен сигнал %v, завершаем работу...", sig)
	}

	shutdownCtx, shutdownCancel := context.WithTimeout(ctx, 5*time.Second)
	defer shutdownCancel()

	if err := server.Shutdown(shutdownCtx); err != nil {
		log.Printf("ошибка при остановке сервера: %v", err)
	} else {
		log.Println("сервер успешно остановлен")
	}
}
