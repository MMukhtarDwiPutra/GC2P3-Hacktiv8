package main

import (
	"book-service/config"

	_ "github.com/joho/godotenv/autoload"
)

func main() {
	config.ListenAndServeGrpc()
}
