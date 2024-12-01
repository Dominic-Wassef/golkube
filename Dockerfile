# Stage 1: Build the application
FROM golang:1.23-alpine AS builder

# Set the Current Working Directory inside the container
WORKDIR /app

# Copy go.mod and go.sum files
COPY go.mod go.sum ./

# Download all dependencies
RUN go mod download

# Copy the source from the current directory
COPY . .

# Build the Go app
RUN go build -o bin/golkube ./cmd/golkube

# Stage 2: Run the application
FROM alpine:latest

WORKDIR /app

# Copy the pre-built binary file
COPY --from=builder /app/bin/golkube .

# Copy configuration files
COPY --from=builder /app/configs ./configs

EXPOSE 8080

ENTRYPOINT ["./golkube"]
CMD ["--config", "configs/default.yaml"]
