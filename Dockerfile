# Pull the Go image (alpine version is smaller version) as the base image
FROM golang:1.21-alpine

# Set the working derctory inside the container
WORKDIR /app

# Copy the go.mod and go.sum from the working directory
COPY go.mod go.sum ./

# Donwload and install Go module dependencies
RUN go mod download

# Copy the rest of your application ocde to the container
COPY . .

# Build the Go application
RUN go build -o kadlab

# Define the command to ru nthe binary
CMD ["./kadlab"]

