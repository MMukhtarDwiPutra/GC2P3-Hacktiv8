package config

import (
	"log"
	"os"

	pb "api-gateway/pb/borrowpb"

	"google.golang.org/grpc"
)

func InitBorrowServiceClient() (*grpc.ClientConn, pb.BorrowServiceClient) {
	// Without TLS (use this if the server does not support TLS)
	conn, err := grpc.Dial(os.Getenv("BORROW_SERVICE_URI"), grpc.WithInsecure())
	if err != nil {
		log.Fatal(err)
	}
	return conn, pb.NewBorrowServiceClient(conn)
}
