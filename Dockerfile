FROM golang:1.23 AS build-stage

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o /tg-bot

FROM alpine

RUN apk add --no-cache tzdata
ENV TZ="Europe/Moscow"

COPY --from=build-stage /tg-bot /tg-bot

COPY ./.env ./.env

COPY ./config ./config/
COPY ./lib/anecdots/anecdots.yaml ./lib/anecdots/anecdots.yaml
COPY ./lib/motivation/quotes.txt ./lib/motivation/quotes.txt

ENTRYPOINT ["/tg-bot"]