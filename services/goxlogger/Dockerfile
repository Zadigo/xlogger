# Use the official Golang image as the base
FROM golang:1.26.3-alpine3.23 AS builder

# Set environment variables
ENV CGO_ENABLED=0 \
    GOOS=linux \
    GOARCH=amd64

# Set working directory inside the container
WORKDIR /app

# Copy go.mod and go.sum files for dependency installation
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy the entire application source
COPY . .

# Build the Go binary
RUN go build -o /app/server .

# Final lightweight stage
FROM alpine:3.21 AS final

# Copy the compiled binary from the builder stage
COPY --from=builder /app/server /bin/app

# ENV GOROOT=/go

# ENV PATH=$PATH:$GOROOT/bin

# Expose the application's port
EXPOSE 8000

# Run the application
CMD ["/bin/app"]
