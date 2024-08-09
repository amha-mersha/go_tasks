package main

import (
	"github.com/amha-mersha/go_taskmanager_mongo/route"
	"github.com/joho/godotenv"
)

func main() {
	godotenv.Load()
	port := 8080
	route.Run(port)
}
