# Use the official Go image as the base image
FROM golang:1.21-alpine

# Set the working directory inside the container
WORKDIR /app

# Copy go.mod and go.sum from the top-level directory to the working directory
# COPY go.mod go.sum ./
COPY go.mod ./

# Download and install Go module dependencies
RUN go mod download

# Copy the rest of your application code to the container
COPY ./app ./

# Build the Go application for the CLI tool
RUN go build -o /output/cli ./cmd/cli

# Build the Go application for the Kademlia tool
RUN go build -o /output/kadlab ./cmd/kademlia

# Make the binary files executable
RUN chmod +x /output/cli /output/kadlab

# Create a directory for the binary output
RUN mkdir -p /output

# Copy the binaries to the output directory
RUN cp /output/cli /output/kadlab

# Define the command to run the Kademlia tool by default
CMD ["/output/kadlab"]
