FROM golang:latest

WORKDIR /app/aviasalesbot
COPY ./ ./

RUN go mod download
RUN go build -o /app/bot ./cmd/app

CMD ["/app/bot"]