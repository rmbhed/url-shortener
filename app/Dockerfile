# Stage 1: Build the binary
FROM golang:1.25-alpine AS builder

WORKDIR /app

# Download dependencies first (helps with caching)
COPY go.mod go.sum ./
RUN go mod download

# Copy source and build
COPY . .
# CGO_ENABLED=0 makes the binary statically linked
RUN CGO_ENABLED=0 GOOS=linux go build -o server .

# Stage 2: Create the final lean image
FROM alpine:latest
RUN apk --no-cache add ca-certificates

WORKDIR /root/

# Copy the binary from the builder
COPY --from=builder /app/server .
# Copy your static folder
COPY --from=builder /app/static ./static

# Run the web service on container startup.
CMD ["./server"]