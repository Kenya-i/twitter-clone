package main

import (
	"context"
	"log"

	"github.com/Kenya-i/twitter-clone/internal/config"
	"github.com/Kenya-i/twitter-clone/internal/handler"
	"github.com/Kenya-i/twitter-clone/internal/middleware"
	"github.com/Kenya-i/twitter-clone/internal/repository"
	"github.com/Kenya-i/twitter-clone/internal/usecase"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
)

func main() {
	cfg := config.Load()

	db, err := pgxpool.New(context.Background(), cfg.DatabaseURL)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	userRepo := repository.NewUserRepository(db)
	userUsecase := usecase.NewUserUsecase(userRepo, cfg.JWTSecret)
	userHandler := handler.NewUserHandler(userUsecase)

	r := gin.Default()

	r.POST("/register", userHandler.Register)
	r.POST("/login", userHandler.Login)

	auth := r.Group("/")
	auth.Use(middleware.AuthMiddleware(cfg.JWTSecret))
	{
		auth.GET("/users/:id", userHandler.GetProfile)
	}

	r.Run(":" + cfg.Port)
}
