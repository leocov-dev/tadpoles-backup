FROM golang:1.19-alpine3.17 AS builder
ARG VERSION_TAG

WORKDIR /code

# this is broken out from call to `make` to improve docker caching
COPY go.mod go.sum ./
RUN go mod download && go mod verify

# run build for container
COPY . .
RUN  CGO_ENABLED=0 go build -o bin/tadpoles-backup --ldflags="-X 'tadpoles-backup/config.VersionTag=$VERSION_TAG'"

FROM alpine:3.17 AS prod

WORKDIR /app
COPY --from=builder /code/bin/tadpoles-backup .

ENTRYPOINT ["./tadpoles-backup", "--non-interactive"]

CMD ["debug"]
