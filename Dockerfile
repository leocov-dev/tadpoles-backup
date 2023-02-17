FROM golang:1.19 AS builder

WORKDIR /code

# this is broken out from call to `make` to improve docker caching
COPY go.mod go.sum ./
RUN go mod download && go mod verify

# run main build
COPY . .
RUN  make dev

FROM builder

WORKDIR /app
COPY --from=builder /code/bin/tadpoles-backup .

ENTRYPOINT ["./tadpoles-backup"]

CMD ["version"]
