# Build stage
# https://docs.docker.com/build/building/multi-stage/
# Use the official Go image as the base image
FROM golang:1.21-alpine AS builder

RUN mkdir -p /app

# Set the working directory inside the container
WORKDIR /app

COPY . .

# Download and install Go module dependencies
RUN go mod download

# Build the Go application for the CLI tool and the Kademlia tool
RUN GOOS=linux GOARCH=amd64 go build -o kadlab ./app/cmd/kademlia/main.go

# Final stage Small OS
FROM alpine:latest 

RUN apk add --no-cache bash

COPY --from=builder /app/kadlab /

# Make the binary files executable
RUN chmod +x /kadlab
# Define the command to run the Kademlia tool by default
CMD ["/kadlab"]
