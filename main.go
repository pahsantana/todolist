package main

import (
	"context"
	"log"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/pahsantana/todolist/config"
	_ "github.com/pahsantana/todolist/docs"
	"github.com/pahsantana/todolist/internal/handlers"
	"github.com/pahsantana/todolist/internal/middleware"
	"github.com/pahsantana/todolist/internal/repositories"
	"github.com/pahsantana/todolist/internal/services"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
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

	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	tasks := router.Group("/tasks")
	{
		tasks.POST("", handler.Create)
		tasks.GET("", handler.List)
		tasks.GET("/:id", handler.GetByID)
		tasks.PUT("/:id", handler.Update)
		tasks.DELETE("/:id", handler.Delete)
		tasks.GET("/summary", handler.Summary)
	}

	logger.Info("server starting", zap.String("port", cfg.ServerPort))
	if err := router.Run(":" + cfg.ServerPort); err != nil {
		log.Fatalf("server error: %v", err)
	}
}
