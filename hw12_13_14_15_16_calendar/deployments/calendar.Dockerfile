FROM golang:1.23-alpine AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o calendar ./cmd/calendar

FROM alpine:latest
RUN apk --no-cache add ca-certificates

RUN addgroup -g 10001 -S appgroup && \
    adduser -u 10001 -S appuser -G appgroup

WORKDIR /app
RUN chown -R appuser:appgroup /app

COPY --from=builder --chown=appuser:appgroup /app/calendar .
COPY --from=builder --chown=appuser:appgroup /app/migrations ./migrations
COPY --from=builder --chown=appuser:appgroup /app/configs/config.yaml ./configs/config.yaml
USER appuser

EXPOSE 8080 50051
CMD ["./calendar", "-config", "configs/config.yaml"]
