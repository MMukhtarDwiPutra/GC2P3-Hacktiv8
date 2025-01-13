package helpers

import (
	"context"
	"fmt"
	"time"

	grpcMetadata "google.golang.org/grpc/metadata"
)

// Creates a new context embedded with auth for gRPC services
func NewServiceContext() (context.Context, context.CancelFunc, error) {
	token, err := SignJwtForGrpc()
	if err != nil {
		return nil, nil, fmt.Errorf("Error sign jwt for grpc")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	ctxWithAuth := grpcMetadata.AppendToOutgoingContext(ctx, "authorization", "Bearer "+token)
	return ctxWithAuth, cancel, nil
}
