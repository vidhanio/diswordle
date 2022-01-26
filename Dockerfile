FROM golang:1.18-rc-alpine

WORKDIR /app

COPY . .

RUN go mod download

RUN go build -o /app/bin/bot ./cmd/bot


CMD ["/app/bin/bot", "-dictionary", "/app/words/dictionary.txt", "-common", "/app/words/common.txt"]