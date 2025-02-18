# Step 1: Build the Go app
FROM golang:1.23.1-alpine AS builder

# Set the Current Working Directory inside the container
WORKDIR /app

# Copy the Go Modules manifests
COPY go.mod go.sum ./

# Download all dependencies. Dependencies will be cached if the go.mod and go.sum are not changed
RUN go mod tidy

# Copy the source code into the container
COPY . .

# Build the Go app
RUN go build -o palmsearch2 .

# Step 2: Run the Go app
FROM alpine:latest

# Install necessary dependencies (e.g., SSL certificates, etc.)
RUN apk --no-cache add ca-certificates

# Set the Current Working Directory inside the container
WORKDIR /root/

# Copy the Go app from the builder stage
COPY --from=builder /app/palmsearch2 .

# Expose the port the app will run on
EXPOSE 5555

# Command to run the executable
CMD ["./palmsearch2"]
