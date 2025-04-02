FROM golang:1.22.3-alpine

WORKDIR /app

COPY go.work.sum ./
COPY go.work ./
COPY core/go.mod ./core/go.mod
COPY report/go.mod ./report/go.mod
COPY measurement/go.mod ./measurement/go.mod

COPY common/ ./common/
COPY auth/ ./auth/

WORKDIR /app/auth

RUN CGO_ENABLED=0 go build -o ./main ./cmd/main.go

EXPOSE 3001

CMD ["./main"]
