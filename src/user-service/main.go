package main

import (
	"user-service/config"

	_ "github.com/joho/godotenv/autoload"
)

func main() {
	config.ListenAndServeGrpc()
}
