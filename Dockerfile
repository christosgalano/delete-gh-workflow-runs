# Stage 1: Build the Go binary
FROM golang:1.23 AS build

# Set the working directory inside the container
WORKDIR /app

# Copy go.mod and go.sum first to cache dependencies
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy the rest of the application source code
COPY . .

# Build the Go binary
RUN CGO_ENABLED=0 GOOS=linux go build -o delete-gh-workflow-runs ./cmd/delete-gh-workflow-runs/main.go

# Stage 2: Create a minimal image for the application
FROM alpine:3.18

# Install bash (optional if your app depends on it)
RUN apk add --no-cache bash

# Copy the Go binary from the builder stage
COPY --from=build /app/delete-gh-workflow-runs /usr/local/bin/delete-gh-workflow-runs

# Make the binary executable
RUN chmod +x /usr/local/bin/delete-gh-workflow-runs

# Set the entrypoint to the Go binary
ENTRYPOINT ["delete-gh-workflow-runs"]
