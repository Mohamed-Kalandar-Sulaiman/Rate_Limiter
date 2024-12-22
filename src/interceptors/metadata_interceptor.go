package interceptors

import (
	"context"
	"fmt"
	"log"

	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

// MetadataInterceptor logs all metadata passed with the request
func MetadataInterceptor() grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context,
		req interface{},
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (interface{}, error) {
		// Extract metadata from the context
		md, ok := metadata.FromIncomingContext(ctx)
		if !ok {
			return nil, fmt.Errorf("missing metadata")
		}

		// Loop through metadata and print each key-value pair
		for key, values := range md {
			for _, value := range values {
				// Print the metadata key and value
				log.Printf("Metadata - Key: %s, Value: %s", key, value)
			}
		}

		// Continue with the request processing
		return handler(ctx, req)
	}
}
