# Use an official Go runtime as a parent image with the correct Go version
FROM golang:1.22.6-alpine

# Install required dependencies including make and protoc
RUN apk add --no-cache \
    make \
    protobuf-dev \
    && go install google.golang.org/protobuf/cmd/protoc-gen-go@latest \
    && go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest

# Set the Current Working Directory inside the container
WORKDIR /app

# Copy go.mod and go.sum files to download dependencies
COPY go.mod go.sum ./
RUN go mod download

# Copy the rest of the application code
COPY . .

# Generate Go files from protobuf files
RUN make protoc

# Build the Go app
RUN go build -o main .

# Expose the application on port 50051
EXPOSE 50053

# Command to run the application
CMD ["./main"]
