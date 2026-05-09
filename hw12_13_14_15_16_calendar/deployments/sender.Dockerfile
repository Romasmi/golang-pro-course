FROM golang:1.23-alpine AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o sender ./cmd/calendar_sender

FROM alpine:latest
RUN apk --no-cache add ca-certificates

RUN addgroup -g 10001 -S appgroup && \
    adduser -u 10001 -S appuser -G appgroup

WORKDIR /app
RUN chown -R appuser:appgroup /app

COPY --from=builder --chown=appuser:appgroup /app/sender .
COPY --from=builder --chown=appuser:appgroup /app/configs/sender_config.yaml ./configs/sender_config.yaml
USER appuser

CMD ["./sender", "-config", "configs/sender_config.yaml"]
