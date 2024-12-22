# Step 1: Build the Go application
FROM golang:1.20 AS builder

# Set the working directory inside the container
WORKDIR /app

# Copy go mod and sum files
COPY go.mod go.sum ./

# Download dependencies (to cache dependencies layer)
RUN go mod tidy

# Copy the entire source code to the container
COPY . .

# Build the Go application
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o rate_limiter .

# Step 2: Create a smaller image for running the app
FROM alpine:latest

# Install necessary dependencies
RUN apk --no-cache add ca-certificates

# Set the working directory
WORKDIR /root/

# Copy the compiled Go binary and the certificate files (if required)
COPY --from=builder /app/rate_limiter .
COPY --from=builder /app/secure /root/secure

# Expose the port that the gRPC server will run on
EXPOSE 50051

# Command to run the Go application
CMD ["./rate_limiter"]
