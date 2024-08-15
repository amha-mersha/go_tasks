package main

import (
	"log"
	"os"
	"time"

	route "github.com/amha-mersha/go_tasks/test-go-backend-task-manager/delivery/routers"
	"github.com/amha-mersha/go_tasks/test-go-backend-task-manager/infrastructure"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	godotenv.Load()
	port := 8080
	router := gin.Default()
	database, err := infrastructure.ConnectDatabase()
	if err != nil {
		log.Println("Error", err)
	}
	route.Run(port, *database, time.Second, router, os.Getenv("DB_USER_COLLECTION_NAME"), os.Getenv("DB_TASK_COLLECTION_NAME"))
}
