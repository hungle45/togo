package main

import (
	"log"
	"togo/config"
	"togo/server"

	"github.com/joho/godotenv"
)

func init() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
}

func main() {

	cfg := config.LoadConfig("config.yml")
	server.StartServer(cfg)
}
