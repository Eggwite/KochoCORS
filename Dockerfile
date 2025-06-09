# Stage 1: Build the Go application
FROM golang:1.24-alpine AS builder

WORKDIR /app

# Copy go.mod and go.sum first to leverage Docker cache for dependencies
COPY go.mod go.sum ./
RUN go mod download
RUN go mod verify

# Copy the rest of the application source code
COPY . .

# Build the Go app, creating a staticly linked binary for Linux
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o /app/proxy .

# Stage 2: Create the final lightweight image
FROM alpine:latest

WORKDIR /root/

# Copy only the compiled binary from the builder stage
COPY --from=builder /app/proxy .

# Expose the port the app runs on.
# Your app uses the PORT env var, defaulting to 3000.
# This informs Docker that the container listens on this port.
EXPOSE 3000

# Command to run the executable.
# The application will pick up the PORT from the environment or default to 3000.
CMD ["./proxy"]