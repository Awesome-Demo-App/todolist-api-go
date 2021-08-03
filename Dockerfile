# Build
FROM golang:1.16.6-alpine3.14 as builder

WORKDIR /src

RUN apk add --update git musl-dev gcc && rm -rf /var/cache/apk/*

COPY go.mod .
COPY go.sum .

RUN go mod download

COPY main.go .
COPY core ./core
COPY handlers ./handlers
COPY repositories ./repositories

RUN go build -a -ldflags "-linkmode external -extldflags '-static' -s -w" -o /dist/todolist-api .
# Runtime
FROM scratch

COPY --from=builder /dist/todolist-api /todolist-api

ENTRYPOINT [ "/todolist-api"]
