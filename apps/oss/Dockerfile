# apps/oss/Dockerfile
FROM golang:1.22-alpine AS build
WORKDIR /app

# Copy go.mod from the module root
COPY go.mod ./

# Generate go.sum and download deps
RUN go mod tidy && go mod download

# Copy ALL source files
COPY . .

# Build from the correct module root
RUN go build -o /oss ./cmd/server/main.go

# Runtime
FROM alpine:latest
RUN apk add --no-cache ca-certificates
WORKDIR /app
COPY --from=build /oss /app/oss
COPY config.yaml /app/config.yaml
EXPOSE 8080
CMD ["/app/oss"]
