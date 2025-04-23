FROM golang:1.22-alpine AS build

WORKDIR /app

COPY ./auth/go.mod ./auth/go.mod
COPY ./report/go.mod ./report/go.mod
COPY ./core/go.mod ./core/go.mod
COPY ./go.work .

COPY ./common/ ./common/
COPY ./measurement ./measurement/

RUN go mod download
RUN CGO_ENABLED=0 go build -ldflags="-s -w" -o cli ./measurement/cmd/main.go

FROM scratch

LABEL maintainer="keiberup.dev@gmail.com"
LABEL version="1.0"
LABEL description="Simple Olt Scanner"

COPY --from=build /app/cli /cli

ENTRYPOINT ["/cli"]

