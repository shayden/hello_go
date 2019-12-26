# start from golang base
FROM golang:latest

LABEL maintainer="Sean <sean@kitbanger.com>"

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN go build -o main .

EXPOSE 8080

CMD ["./main"]