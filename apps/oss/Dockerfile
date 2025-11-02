# apps/oss/Dockerfile
FROM golang:1.22-alpine AS build

# Build context is: neftac-storage-monorepo/apps/oss/
WORKDIR /app

# Copy go.mod (at apps/oss/go.mod)
COPY go.mod ./

# Generate go.sum + download deps
RUN go mod tidy && go mod download

# Copy source files from apps/oss/
COPY cmd ./cmd
COPY internal ./internal
COPY scripts ./scripts
COPY config.yaml ./

# Build from correct path
RUN go build -o /oss ./cmd/server/main.go

# Runtime
FROM alpine:latest
RUN apk add --no-cache ca-certificates
WORKDIR /app
COPY --from=build /oss /app/oss
COPY config.yaml /app/config.yaml
EXPOSE 8080
CMD ["/app/oss"]
