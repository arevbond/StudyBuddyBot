FROM golang:1.23 AS build-stage

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o /tg-bot

FROM alpine

COPY --from=build-stage /tg-bot /tg-bot

COPY ./.env ./.env

COPY ./config ./config/

ENTRYPOINT ["/tg-bot"]