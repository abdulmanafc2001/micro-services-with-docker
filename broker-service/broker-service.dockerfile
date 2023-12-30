FROM golang:1.21.4-alpine AS builder
WORKDIR /app
COPY go.mod go.sum ./
COPY . .
RUN CGO_ENABLED=0 go build -o brockerApp ./cmd/api

FROM alpine:latest
WORKDIR /app
COPY --from=builder /app/brockerApp /app
CMD [ "/app/brockerApp" ]