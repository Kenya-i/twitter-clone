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
	likeRepo := repository.NewLikeRepository(db)
	tweetUsecase := usecase.NewTweetUsecase(tweetRepo, likeRepo)
	tweetHandler := handler.NewTweetHandler(tweetUsecase)

	followRepo := repository.NewFollowRepository(db)
	followUsecase := usecase.NewFollowUsecase(followRepo)
	followHandler := handler.NewFollowHandler(followUsecase)

	r := gin.Default()

	r.Use(cors.New(cors.Config{
		AllowOrigins: []string{"http://localhost:3000"},
		AllowMethods: []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders: []string{"Origin", "Content-Type", "Authorization"},
	}))

	r.POST("/register", userHandler.Register)
	r.POST("/login", userHandler.Login)

	auth := r.Group("/")
	auth.Use(middleware.AuthMiddleware(cfg.JWTSecret))
	{
		auth.GET("/users", userHandler.GetUsers)
		auth.GET("/users/:id", userHandler.GetProfile)
		auth.GET("/users/:id/follow", followHandler.GetFollowInfo)
		auth.POST("/users/:id/follow", followHandler.Follow)
		auth.DELETE("/users/:id/follow", followHandler.Unfollow)
		auth.GET("/tweets", tweetHandler.GetTimeline)
		auth.POST("/tweets", tweetHandler.Post)
		auth.GET("/tweets/:id", tweetHandler.GetTweet)
		auth.PUT("/tweets/:id", tweetHandler.Update)
		auth.DELETE("/tweets/:id", tweetHandler.Delete)
		auth.POST("/tweets/:id/like", tweetHandler.Like)
		auth.DELETE("/tweets/:id/like", tweetHandler.Unlike)
	}

	r.Run(":" + cfg.Port)
}
