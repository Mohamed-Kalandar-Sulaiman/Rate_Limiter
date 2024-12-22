# Run protoc to output files in the root folder temporarily
protoc --go_out=. --go-grpc_out=. Proto/source/rate_limiter_service.proto

