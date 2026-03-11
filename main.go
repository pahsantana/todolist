package main

import (
	"context"
	"log"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/pahsantana/todolist/config"
	"github.com/pahsantana/todolist/internal/handlers"
	"github.com/pahsantana/todolist/internal/middleware"
	"github.com/pahsantana/todolist/internal/repositories"
	"github.com/pahsantana/todolist/internal/services"
	"go.uber.org/zap"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Println("warning: .env file not found, using system environment variables")
	}

	cfg := config.Load()

	logger, _ := zap.NewProduction()
	defer logger.Sync()

	client, err := repositories.NewMongoClient(cfg.MongoURI)
	if err != nil {
		log.Fatalf("failed to connect to mongodb: %v", err)
	}
	defer client.Disconnect(context.Background())
	logger.Info("connected to mongodb")

	db := client.Database(cfg.MongoDB)
	repo := repositories.NewTaskRepository(db)
	service := services.NewTaskService(repo)
	handler := handlers.NewTaskHandler(service, logger)

	router := gin.New()
	router.Use(gin.Recovery())
	router.Use(middleware.RequestLogger(logger))

	router.GET("/health", handlers.Health)

	tasks := router.Group("/tasks")
	{
		tasks.POST("", handler.Create)
		tasks.GET("", handler.List)
		tasks.GET("/:id", handler.GetByID)
		tasks.PUT("/:id", handler.Update)
		tasks.DELETE("/:id", handler.Delete)
	}

	logger.Info("server starting", zap.String("port", cfg.ServerPort))
	if err := router.Run(":" + cfg.ServerPort); err != nil {
		log.Fatalf("server error: %v", err)
	}
}
