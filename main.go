package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"os"

	// 3RD PARTY LIBS
	"github.com/go-redis/redis/v8"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"

	"github.com/joho/godotenv"

	"github.com/Mohamed-Kalandar-Sulaiman/Rate_Limiter/src/interceptors"
	"github.com/Mohamed-Kalandar-Sulaiman/Rate_Limiter/src/repository"
	"github.com/Mohamed-Kalandar-Sulaiman/Rate_Limiter/src/utils"

	// SELF LIBS
	proto "github.com/Mohamed-Kalandar-Sulaiman/Rate_Limiter/Proto/generated"
	server "github.com/Mohamed-Kalandar-Sulaiman/Rate_Limiter/src/services"
)

func main() {
	// Load environment variables from .env file
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatalf("Error loading .env file")
	}
	// Configure Redis options
	redisOptions := &redis.Options{
								Addr:     os.Getenv("REDIS_HOST_URL") + ":" + os.Getenv("REDIS_PORT"), 
								Password: os.Getenv("REDIS_PASSWORD"),
								Username: os.Getenv("REDIS_USERNAME"),        
								DB:       0,        
							}

	redisClient := redis.NewClient(redisOptions)

	ctx := context.Background()
	_, err = redisClient.Ping(ctx).Result()
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
																		interceptors.LogInterceptor(),
																		interceptors.AuthInterceptor(),     
																		
																	),
										)

	rateLimiterServer := server.NewRateLimiterServer(*repo, configMap)
	proto.RegisterRateLimitServiceServer(grpcServer, rateLimiterServer)

	// Start listening
	lis, err := net.Listen("tcp",os.Getenv("HOSTED_URL") + ":" + os.Getenv("HOSTED_PORT"), )
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}
	fmt.Printf("gRPC server listening on %s\n", os.Getenv("HOSTED_PORT"))

	// Serve the gRPC server
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}
