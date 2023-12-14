FROM golang:latest

WORKDIR /tg_ics_useful_bot

COPY go.mod .

RUN go mod download

COPY . .

RUN curl -fsSL https://raw.githubusercontent.com/pressly/goose/master/install.sh | sh

WORKDIR /tg_ics_useful_bot

RUN go build -o /tg_ics_useful_bot

EXPOSE 8080

CMD [ "./tg_ics_useful_bot" ]