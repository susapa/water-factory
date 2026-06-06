FROM golang:1.24-alpine AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o water-api ./cmd/server

FROM alpine:3.20
RUN apk add --no-cache ca-certificates tzdata
ENV TZ=Asia/Bangkok
WORKDIR /app
COPY --from=builder /app/water-api .
COPY --from=builder /app/internal/database/migrations ./internal/database/migrations
EXPOSE 8080
CMD ["./water-api"]
