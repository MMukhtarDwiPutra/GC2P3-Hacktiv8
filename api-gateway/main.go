package main

import (
	"api-gateway/config"
	"api-gateway/controller"
	"api-gateway/routes"
	"log"
	"os"

	"github.com/joho/godotenv"
)

func main() {
	// Load environment variables
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}
	log.Println(".env file loaded successfully")

	// Initialize gRPC client connection
	userClientConn, userClient := config.InitUserServiceClient()
	defer userClientConn.Close() // Close the connection when the application exits

	// Initialize gRPC client connection
	bookClientConn, bookClient := config.InitBookServiceClient()
	defer bookClientConn.Close() // Close the connection when the application exits

	// Initialize the UserController
	userController := controller.NewUserController(userClient)
	bookController := controller.NewBookController(bookClient)

	// Create the router
	router := routes.NewRouter(userController, bookController)

	// Start the server
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	log.Printf("Starting server on port %s", port)
	router.Logger.Fatal(router.Start(":" + port))
}
