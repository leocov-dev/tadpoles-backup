FROM golang:1.19 AS builder
ARG ARCH=amd64
ARG OS=linux
ARG VERSION_TAG

WORKDIR /code

# this is broken out from call to `make` to improve docker caching
COPY go.mod go.sum ./
RUN go mod download && go mod verify

# run build for container
COPY . .
RUN  make container GOOS=$OS GOARCH=$ARCH VERSION_TAG=$VERSION_TAG

FROM alpine:latest AS prod

WORKDIR /app
COPY --from=builder /code/bin/tadpoles-backup .

ENTRYPOINT ["./tadpoles-backup", "--non-interactive"]

CMD ["debug"]
