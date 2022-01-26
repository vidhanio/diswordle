FROM golang:1.18-rc-alpine

WORKDIR /app

COPY . .

RUN go mod download

RUN go build ./cmd/bot -o /app/bin/bot


CMD ["/app/bin/bot", "-words", "/app/words/words.txt", "-commonwords", "/app/words/commonwords.txt"]