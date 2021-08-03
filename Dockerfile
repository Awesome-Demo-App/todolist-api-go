# Build
FROM golang:1.16.6-alpine3.14 as builder

WORKDIR /src

RUN apk add --update git && rm -rf /var/cache/apk/*

COPY go.mod .
COPY go.sum .

RUN go mod download

COPY main.go .
COPY core ./core
COPY handlers ./handlers
COPY repositories ./repositories

RUN CGO_ENABLED=0 \
    GOOS=linux \
    go build -a -installsuffix cgo -o /dist/todolist-api .

# Runtime
FROM scratch

COPY --from=builder /dist/todolist-api /todolist-api

ENTRYPOINT ["/todolist-api"]

