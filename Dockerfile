FROM golang:1.18-rc-alpine

WORKDIR /bot
COPY . .

RUN cd cmd/bot
RUN go build
CMD [ "cmd/bot/bot" ]