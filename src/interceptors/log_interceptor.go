package interceptors

import (
	"context"
	"log"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

func LogInterceptor() grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context,
		req interface{},
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (interface{}, error) {
		
		startTime := time.Now()
		log.Printf("Handling gRPC request: Method=%s", info.FullMethod)
		log.Printf("Request: %+v", req)
		md, _ := metadata.FromIncomingContext(ctx)

		for key, values := range md {
			for _, value := range values {
				log.Printf("Metadata - Key: %s, Value: %s", key, value)
				
			}
		}

		resp, err := handler(ctx, req)
		log.Printf("Handled gRPC request: Method=%s, Duration=%v, Error=%v", info.FullMethod, time.Since(startTime), err)

		return resp, err
	}
}