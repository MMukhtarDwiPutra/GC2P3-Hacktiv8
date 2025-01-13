package config

import (
	"book-service/controller"
	"book-service/middlewares"
	pb "book-service/pb"
	"book-service/repository"
	"book-service/service"
	"context"
	"fmt"
	"log"
	"net"
	"os"

	grpc_auth "github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/auth"
	"github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/logging"
	"google.golang.org/grpc"
)

func ListenAndServeGrpc() {
	port := os.Getenv("PORT")

	lis, err := net.Listen("tcp", ":"+port)
	if err != nil {
		log.Fatal(err)
	}

	collection, err := ConnectionDatabase(context.Background())
	if err != nil {
		fmt.Sprintf("Error database: %v", err)
	}

	bookRepository := repository.NewBookRepository(collection)
	bookService := service.NewBookService(bookRepository)
	bookController := controller.NewBookController(bookService)

	// Define a custom interceptor for JWT that conditionally skips authentication for register endpoint
	grpcServer := grpc.NewServer(
		grpc.ChainUnaryInterceptor(
			logging.UnaryServerInterceptor(middlewares.NewInterceptorLogger()),
			grpc_auth.UnaryServerInterceptor(func(ctx context.Context) (context.Context, error) {
				return middlewares.JWTAuth(ctx)
			}),
		),
	)

	pb.RegisterBookServiceServer(grpcServer, bookController)

	log.Println("\033[36mGRPC server is running on port:", port, "\033[0m")
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatal("Failed to server gRPC:", err)
	}
}
