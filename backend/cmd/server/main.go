package main

import (
	"context"
	"log"

	"github.com/Kenya-i/twitter-clone/internal/config"
	"github.com/Kenya-i/twitter-clone/internal/handler"
	"github.com/Kenya-i/twitter-clone/internal/middleware"
	"github.com/Kenya-i/twitter-clone/internal/repository"
	"github.com/Kenya-i/twitter-clone/internal/usecase"
	"github.com/gin-contrib/cors"
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

	tweetRepo := repository.NewTweetRepository(db)
	tweetUsecase := usecase.NewTweetUsecase(tweetRepo)
	tweetHandler := handler.NewTweetHandler(tweetUsecase)

	r := gin.Default()

	r.Use(cors.New(cors.Config{
		AllowOrigins: []string{"http://localhost:3000"},
		AllowMethods: []string{"GET", "POST", "DELETE", "OPTIONS"},
		AllowHeaders: []string{"Origin", "Content-Type", "Authorization"},
	}))

	r.POST("/register", userHandler.Register)
	r.POST("/login", userHandler.Login)

	auth := r.Group("/")
	auth.Use(middleware.AuthMiddleware(cfg.JWTSecret))
	{
		auth.GET("/users/:id", userHandler.GetProfile)
		auth.POST("/tweets", tweetHandler.Post)
		auth.GET("/tweets/:id", tweetHandler.GetTweet)
		auth.DELETE("/tweets/:id", tweetHandler.Delete)
	}

	r.Run(":" + cfg.Port)
}
