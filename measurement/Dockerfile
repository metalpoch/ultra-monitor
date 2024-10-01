FROM golang:alpine AS build

RUN apk add --no-cache gcc musl-dev

WORKDIR /app

COPY . .
RUN CGO_ENABLED=1 go build -a -ldflags "-linkmode external -extldflags '-static' -s -w" -o cli ./measurement/cmd/cli

FROM scratch
COPY --from=build /app/cli /app/cli
ENTRYPOINT ["/app/cli"]