FROM golang:1.19 AS builder

WORKDIR /code

# this is broken out from call to `make` to improve docker caching
COPY go.mod go.sum ./
RUN go mod download && go mod verify

# run main build
COPY . .
RUN  GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -o bin/tadpoles-backup

FROM alpine:latest AS prod

WORKDIR /app
COPY --from=builder /code/bin/tadpoles-backup .

ENTRYPOINT ["./tadpoles-backup"]

CMD ["--non-interactive debug"]
