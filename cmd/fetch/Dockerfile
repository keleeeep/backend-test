FROM golang:alpine

RUN apk update && apk add --no-cache git gcc musl-dev

WORKDIR /app

COPY . .

RUN go mod tidy

RUN go build -o binary github.com/keleeeep/test/cmd/fetch

EXPOSE 8080

ENTRYPOINT ["/app/binary"]