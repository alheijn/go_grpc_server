# Use the official Golang image to create a build artifact.
FROM golang:1.23-alpine AS builder

# Set the working directory inside the container
WORKDIR /app

# Copy go.mod and go.sum files to download dependencies
COPY go.mod go.sum ./
RUN go mod download

# Copy the rest of the application source code
COPY . .

# Compile the .proto file
# Note: This requires protoc to be installed in the build image.
# A simpler approach is to have the generated files checked into source control.
# If you don't have protoc, generate the files locally and comment out this RUN command.
RUN apk add --no-cache protobuf-dev
RUN protoc --go_out=. --go_opt=paths=source_relative \
    --go-grpc_out=. --go-grpc_opt=paths=source_relative \
    ecommerce/ordermanagement/ordermanagement.proto

# Build the Go application, creating a statically linked binary
RUN CGO_ENABLED=0 GOOS=linux go build -o /go/bin/server .

# ---- Final Stage ----
# Use a minimal base image for a small final image size
FROM alpine:latest

# The alpine image already has root certificates for SSL/TLS
RUN apk --no-cache add ca-certificates

# Copy only the compiled application from the builder stage
COPY --from=builder /go/bin/server /server

# Expose the gRPC port
EXPOSE 50051

# Command to run the executable
CMD ["/server"]