FROM golang:1.25-alpine

WORKDIR /app

# Install git for go mod download
RUN apk add --no-cache git

# Copy go mod files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy source code
COPY . .

# Expose port
EXPOSE 8080

# Default command (overridden in docker-compose for dev)
CMD ["go", "run", "cmd/api/main.go"]
