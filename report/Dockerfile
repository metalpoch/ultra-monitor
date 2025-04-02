FROM golang:1.22.3-alpine

WORKDIR /workspace

COPY auth/go.mod ./auth/go.mod
COPY core/go.mod ./core/go.mod
COPY measurement/go.mod ./measurement/go.mod
COPY common/ ./common/
COPY report/ ./report/
COPY go.work.sum ./
COPY go.work ./
COPY config.product.json ./

WORKDIR /workspace/report
ENV CONFIG_JSON=/workspace/config.product.json

RUN  CGO_ENABLED=0 GOOS=linux go build -a -ldflags "-s -w" -o ./main ./cmd/main.go
EXPOSE 3002

CMD ["./main"]
