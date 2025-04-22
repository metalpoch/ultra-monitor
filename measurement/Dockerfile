FROM golang:1.22-alpine AS build

RUN apk add --no-cache gcc musl-dev

WORKDIR /app

COPY ./go.work .
COPY ./auth/go.mod ./auth/go.mod
COPY ./report/go.mod ./report/go.mod
COPY ./core/go.mod ./core/go.mod

COPY ./common/ ./common/
COPY ./measurement ./measurement/

RUN CGO_ENABLED=1 go build -a -ldflags "-linkmode external -extldflags '-static' -s -w" -o cli ./measurement/cmd/main.go

FROM scratch
COPY --from=build /app/cli /app/cli
ENTRYPOINT ["/app/cli"]
