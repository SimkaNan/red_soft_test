FROM golang:1.22.6-alpine3.20

WORKDIR /usr/local/src

COPY ./ ./

# install psql
RUN apk update
RUN apk add postgresql-client

RUN chmod +x wait_for_postgres.sh


RUN go mod download
RUN go build -o main ./cmd/main.go
CMD ["./main"]