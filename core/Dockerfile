FROM golang:1.22.3-alpine

WORKDIR /workspace

COPY go.work.sum ./
COPY go.work ./
COPY config.product.json ./
COPY auth/go.mod ./auth/go.mod
COPY report/go.mod ./report/go.mod
COPY measurement/go.mod ./measurement/go.mod

COPY common/ ./common/
COPY core/ ./core/

WORKDIR /workspace/core
ENV CONFIG_JSON=/workspace/config.product.json

RUN CGO_ENABLED=0 go build -o ./main ./cmd/main.go

EXPOSE 3000

CMD ["./main"]
