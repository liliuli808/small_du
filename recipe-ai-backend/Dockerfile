# 构建阶段
FROM golang:1.22-alpine AS builder

WORKDIR /app

RUN apk add --no-cache git

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-w -s" -o /app/api ./cmd/api/main.go
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-w -s" -o /app/worker ./cmd/worker/main.go

# API服务运行阶段
FROM alpine:latest AS api

RUN apk --no-cache add ca-certificates tzdata

WORKDIR /app

COPY --from=builder /app/api /app/api
COPY --from=builder /app/config /app/config
COPY --from=builder /app/migrations /app/migrations

EXPOSE 8080

CMD ["/app/api"]

# Worker服务运行阶段
FROM alpine:latest AS worker

RUN apk --no-cache add ca-certificates tzdata

WORKDIR /app

COPY --from=builder /app/worker /app/worker
COPY --from=builder /app/config /app/config
COPY --from=builder /app/migrations /app/migrations

CMD ["/app/worker"]
