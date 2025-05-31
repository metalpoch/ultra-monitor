FROM golang:1.24.2-alpine AS builder

WORKDIR /app

COPY . .

RUN CGO_ENABLED=0 go build -ldflags="-s -w" -o ./server ./cmd/server/main.go

FROM alpine

LABEL name="Ultra Monitor"
LABEL maintainer="keiberup.dev@gmail.com"
LABEL description="Ultra monitor server"
LABEL version="1.0.0"

WORKDIR /app

COPY --from=builder /app/server .

CMD ["./server"]