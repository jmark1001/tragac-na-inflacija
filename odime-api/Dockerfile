# Use golang:1.22 as the base image for building the Go server binary
FROM golang:1.22 as builder

# Set the working directory inside the builder container
WORKDIR /app

# Copy go.mod and go.sum to fetch dependencies
COPY ./odime-api/go.mod ./odime-api/go.sum ./
RUN go mod tidy

# Copy the Go source code into the container
COPY ./odime-api /app

# Build the Go server binary
RUN go build -o server ./main.go


# Final stage:
FROM debian:bookworm-slim

WORKDIR /app

# Copy the compiled Go binary from the builder stage
COPY --from=builder /app/server /app/

# Install necessary dependencies
RUN apt-get update && apt-get install -y ca-certificates libc6

# Copy the wait-for-it script into the container
COPY ./wait-for-it.sh /app/wait-for-it.sh
RUN chmod +x /app/wait-for-it.sh

# Copy the .env file into the container's /app directory
COPY ./odime-api/.env /app/.env

# Expose the necessary port
EXPOSE 8080

# Run the Go server (use wait-for-it.sh to ensure dependencies are ready)
CMD ["sh", "-c", "/app/wait-for-it.sh ${DB_HOST}:${DB_PORT} -- /app/wait-for-it.sh ${MQ_HOST}:${MQ_PORT} -- ./server"]
