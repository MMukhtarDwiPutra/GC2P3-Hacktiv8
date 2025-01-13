package config

import (
	"log"
	"os"

	pb "api-gateway/pb/userpb"

	"google.golang.org/grpc"
)

func InitUserServiceClient() (*grpc.ClientConn, pb.UserServiceClient) {
	// Without TLS (use this if the server does not support TLS)
	conn, err := grpc.Dial(os.Getenv("USER_SERVICE_URI"), grpc.WithInsecure())
	if err != nil {
		log.Fatal(err)
	}
	return conn, pb.NewUserServiceClient(conn)
}
