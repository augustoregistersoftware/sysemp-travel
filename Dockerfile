FROM golang:1.25

WORKDIR /go/src/app

COPY go.mod go.sum ./

RUN go mod download

COPY . .

EXPOSE 8080

RUN go build -o /usr/local/bin/api ./cmd/api
# RUN go build -o /usr/local/bin/worker ./cmd/worker

CMD ["/usr/local/bin/api"]