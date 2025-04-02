FROM golang:1.22.3-alpine

WORKDIR /app

COPY auth/go.mod ./auth/go.mod
COPY core/go.mod ./core/go.mod
COPY measurement/go.mod ./measurement/go.mod
COPY common/ ./common/
COPY report/ ./report/
COPY go.work ./

WORKDIR /app/report

RUN CGO_ENABLED=0 go build -o ./main ./cmd/main.go
EXPOSE 3002

CMD ["./main"]
