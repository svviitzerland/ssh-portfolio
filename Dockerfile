# ── Build Stage ──────────────────────────────────────────────────────
FROM golang:1.24-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o ssh-portfolio .

# ── Runtime Stage ───────────────────────────────────────────────────
FROM alpine:3.21

RUN apk --no-cache add ca-certificates

RUN adduser -D -g '' appuser
USER appuser

WORKDIR /app

COPY --from=builder /app/ssh-portfolio .

EXPOSE 2222

ENTRYPOINT ["./ssh-portfolio"]
