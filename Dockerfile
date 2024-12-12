FROM golang:latest

WORKDIR /app
COPY ./ ./

RUN go mod download
RUN go build -o /app/bot ./cmd/app

CMD ["/app/bot", "-tg.token=7554451672:AAG-T9evjULW8DCHSZDkVltg5HN-uimBeEg"]
