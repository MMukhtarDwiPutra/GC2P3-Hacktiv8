package config

import (
	"log"
	"os"

	pb "api-gateway/pb/bookpb"

	"google.golang.org/grpc"
)

func InitBookServiceClient() (*grpc.ClientConn, pb.BookServiceClient) {
	// Without TLS (use this if the server does not support TLS)
	conn, err := grpc.Dial(os.Getenv("BOOK_SERVICE_URI"), grpc.WithInsecure())
	if err != nil {
		log.Fatal(err)
	}
	return conn, pb.NewBookServiceClient(conn)
}
