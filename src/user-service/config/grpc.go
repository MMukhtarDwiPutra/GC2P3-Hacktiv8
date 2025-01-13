package config

import (
	"context"
	"fmt"
	"log"
	"net"
	"os"
	"user-service/controller"
	"user-service/middlewares"
	pb "user-service/pb/userpb"
	"user-service/repository"
	"user-service/service"

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

	userRepository := repository.NewUserRepository(collection)
	userService := service.NewUserService(userRepository)
	userController := controller.NewUserController(userService)

	// Define a custom interceptor for JWT that conditionally skips authentication for register endpoint
	grpcServer := grpc.NewServer(
		grpc.ChainUnaryInterceptor(
			logging.UnaryServerInterceptor(middlewares.NewInterceptorLogger()),
			grpc_auth.UnaryServerInterceptor(func(ctx context.Context) (context.Context, error) {
				return middlewares.JWTAuth(ctx)
			}),
		),
	)

	pb.RegisterUserServiceServer(grpcServer, userController)

	log.Println("\033[36mGRPC server is running on port:", port, "\033[0m")
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatal("Failed to server gRPC:", err)
	}
}
