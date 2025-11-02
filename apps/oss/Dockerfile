# apps/oss/Dockerfile
FROM golang:1.22-alpine AS build
WORKDIR /app

# Copy go.mod
COPY go.mod ./

# Generate go.sum + download deps (inside Docker)
RUN go mod tidy && go mod download

# Copy all source files
COPY . .

# Build
RUN go build -o /oss ./cmd/server/main.go

# Runtime
FROM alpine:latest
RUN apk add --no-cache ca-certificates
WORKDIR /app
COPY --from=build /oss /app/oss
COPY config.yaml /app/config.yaml
EXPOSE 8080
CMD ["/app/oss"]
