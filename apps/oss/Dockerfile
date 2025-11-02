# apps/oss/Dockerfile
FROM golang:1.22-alpine AS build
WORKDIR /app

# Copy only go.mod
COPY go.mod ./

# Generate go.sum and download deps
RUN go mod tidy && go mod download

# Copy source
COPY . .

# Build
RUN go build -o /oss ./cmd/server/main.go

# Runtime
FROM alpine:latest
RUN apk add --no-cache ca-certificates
COPY --from=build /oss /oss
COPY config.yaml /config.yaml
EXPOSE 8080
CMD ["/oss"]
