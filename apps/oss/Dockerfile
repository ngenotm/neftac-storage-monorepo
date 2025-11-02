FROM golang:1.22-alpine AS build
WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN go build -o /oss ./cmd/server/main.go

FROM alpine:latest
RUN apk add --no-cache ca-certificates
COPY --from=build /oss /oss
COPY config.yaml /config.yaml
EXPOSE 8080
CMD ["/oss"]
