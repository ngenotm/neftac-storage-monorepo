# apps/oss/Dockerfile
FROM golang:1.22-alpine AS build
WORKDIR /app

# Copy go.mod only
COPY go.mod ./

# Auto-generate go.sum + download deps
RUN go mod tidy && go mod download

# Copy all source
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
