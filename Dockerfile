# Use the official Go image as the base image
FROM golang:1.23.1-alpine AS builder

# Install build dependencies
RUN apk add --no-cache gcc musl-dev

# Set the working directory inside the container
WORKDIR /app

# Copy go mod and sum files
COPY go.mod go.sum ./

# Download all dependencies
RUN go mod download

# Copy the source code into the container
COPY . .

# Build the application
RUN CGO_ENABLED=1 GOOS=linux go build -a -installsuffix cgo -o receipt-processor ./cmd/server

# Use a smaller base image for the final stage
FROM alpine:latest

# Install SQLite and its dependencies
RUN apk add --no-cache sqlite sqlite-libs

# Set the working directory
WORKDIR /root/

# Copy the binary from the builder stage
COPY --from=builder /app/receipt-processor .

# Expose the port the app runs on
EXPOSE 8080

# Command to run the application
CMD ["./receipt-processor"]