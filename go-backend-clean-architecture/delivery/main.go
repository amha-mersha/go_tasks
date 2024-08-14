package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	route "github.com/amha-mersha/go_tasks/go-backend-clean-architecture/delivery/routers"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func main() {
	godotenv.Load()
	port := 8080
	connectionString := os.Getenv("DB_CONNECTION_STRING")
	clientOptions := options.Client().ApplyURI(connectionString)
	client, err := mongo.Connect(context.TODO(), clientOptions)
	if err != nil {
		fmt.Println("Connection Error: ", err)
	}
	err = client.Ping(context.TODO(), nil)
	if err != nil {
		fmt.Println("Ping Error: ", err)
	}

	database := client.Database(os.Getenv("DB_NAME"))
	log.Println("Database successfuly connected.")
	router := gin.Default()
	route.Run(port, *database, time.Second, router, os.Getenv("DB_USER_COLLECTION_NAME"), os.Getenv("DB_TASK_COLLECTION_NAME"))
}
