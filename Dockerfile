FROM golang:1.26-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .


RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 \
    go build -trimpath -ldflags="-s -w" \
    -o /app/app ./cmd/api/main.go


FROM alpine:3.20

RUN apk --no-cache add ca-certificates tzdata && \
    adduser -D -H -u 1000 appuser

WORKDIR /app

COPY --from=builder /app/app .

RUN chown -R appuser:appuser /app

USER appuser

EXPOSE 8080

CMD ["./app"]