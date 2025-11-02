FROM golang:1.22-alpine AS build
WORKDIR /app

COPY go.mod ./
RUN go mod tidy && go mod download

COPY . .
RUN go build -o /oss ./cmd/server/main.go

FROM alpine:latest
RUN apk add --no-cache ca-certificates
WORKDIR /app
COPY --from=build /oss /app/oss
COPY config.yaml /app/config.yaml
EXPOSE 8080
CMD ["/app/oss"]
