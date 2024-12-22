package main

import (
	"context"
	"fmt"
	"log"
	"net"

	// 3RD PARTY LIBS
	"github.com/go-redis/redis/v8"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"

	"github.com/Mohamed-Kalandar-Sulaiman/Rate_Limiter/src/interceptors"
	"github.com/Mohamed-Kalandar-Sulaiman/Rate_Limiter/src/repository"
	"github.com/Mohamed-Kalandar-Sulaiman/Rate_Limiter/src/utils"

	// SELF LIBS
	proto "github.com/Mohamed-Kalandar-Sulaiman/Rate_Limiter/Proto/generated"
	server "github.com/Mohamed-Kalandar-Sulaiman/Rate_Limiter/src/services"
)

const (
	port      = ":50051"
	redisAddr = "localhost:6379"
)

func main() {
	// Configure Redis options
	redisOptions := &redis.Options{
		Addr:     redisAddr, 
		Password: "",        
		DB:       0,        
	}

	redisClient := redis.NewClient(redisOptions)

	ctx := context.Background()
	_, err := redisClient.Ping(ctx).Result()
	if err != nil {
		log.Fatalf("Failed to connect to Redis: %v", err)
	}
	fmt.Println("Successfully connected to Redis")



	repo      := repository.NewRateLimiterRepository(redisClient)
	configMap := utils.NewConfigMap()

	// Load configuration
	err = configMap.LoadConfig("./config.yaml")
	if err != nil {
		log.Fatalf("Error loading config: %v", err)
	}

	// Create a gRPC server
	// Load server's certificate and private key
	certFile := "./secure/server.crt"
	keyFile  := "./secure/server.key"

	// Create a credentials object using the certificate and key
	creds, err := credentials.NewServerTLSFromFile(certFile, keyFile)
	if err != nil {
		log.Fatalf("failed to load TLS credentials: %v", err)
	}

	// Create a gRPC server with TLS credentials
	grpcServer        := grpc.NewServer(
											grpc.Creds(creds),
											grpc.ChainUnaryInterceptor(
																		interceptors.AuthInterceptor(),     
																		interceptors.MetadataInterceptor(),   
																	),
										)

	rateLimiterServer := server.NewRateLimiterServer(*repo, configMap)
	proto.RegisterRateLimitServiceServer(grpcServer, rateLimiterServer)

	// Start listening
	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}
	fmt.Printf("gRPC server listening on %s\n", port)

	// Serve the gRPC server
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}
