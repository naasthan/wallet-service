FROM golang:1.24.11-alpine AS builder

RUN apk add --no-cache git gcc musl-dev

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o /wallet-app ./cmd/wallet/main.go

FROM alpine:latest

RUN apk --no-cache add ca-certificates tzdata

WORKDIR /app/

COPY --from=builder /wallet-app .
COPY --from=builder /app/config.env .
COPY --from=builder /app/migrations ./migrations

RUN chmod +x ./wallet-app

EXPOSE 8080

ENTRYPOINT ["/app/wallet-app"]