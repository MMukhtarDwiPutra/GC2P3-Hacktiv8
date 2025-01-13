package main

import (
	"borrow-service/config"

	_ "github.com/joho/godotenv/autoload"
)

func main() {
	config.ListenAndServeGrpc()
}
