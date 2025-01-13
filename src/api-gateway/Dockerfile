# Use an official Go runtime as a parent image
FROM golang:1.22.6-alpine

# Install necessary dependencies: make, protoc, etc.
RUN apk add --no-cache \
    make \
    protobuf-dev \
    bash \
    && go install google.golang.org/protobuf/cmd/protoc-gen-go@latest \
    && go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest

# Set the Current Working Directory inside the container
WORKDIR /app

# Copy go.mod and go.sum files to download dependencies
COPY go.mod go.sum ./
RUN go mod download

# Copy the rest of the application code (including gRPC generated files)
COPY . .

# Generate Go files from protobuf files (if necessary)
RUN protoc --proto_path=proto proto/*.proto --go_out=. --go-grpc_out=.

# Expose the API Gateway on port 5000 or your preferred port
EXPOSE 8080

# Install the necessary Go dependencies for your REST API (e.g., echo, etc.)
RUN go get github.com/labstack/echo/v4

# Build the Go app (API Gateway)
RUN go build -o main .

# Command to run the application (API Gateway)
CMD ["./main"]
