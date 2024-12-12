FROM golang:latest

WORKDIR /app
COPY ./ ./

RUN go mod download
RUN go build -o /app/bot ./cmd/app

CMD ["/app/bot", "-tg.token=SECRET"]
