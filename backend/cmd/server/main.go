package main

import (
	"context"
	"log"

	"github.com/Kenya-i/twitter-clone/internal/handler"
	"github.com/Kenya-i/twitter-clone/internal/repository"
	"github.com/Kenya-i/twitter-clone/internal/usecase"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

func main() {
	client, err := mongo.Connect(options.Client().ApplyURI("mongodb://localhost:27017"))
	if err != nil {
		log.Fatal(err)
	}
	defer client.Disconnect(context.Background())

	db := client.Database("twitter_clone")

	userRepo := repository.NewUserRepository(db)
	userUsecase := usecase.NewUserUsecase(userRepo)
	userHandler := handler.NewUserHandler(userUsecase)

	r := gin.Default()

	r.POST("/register", userHandler.Register)
	r.POST("/login", userHandler.Login)
	r.GET("/users/:id", userHandler.GetProfile)

	r.Run(":8080")
}
