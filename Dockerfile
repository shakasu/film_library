FROM golang:latest

RUN go version

COPY . /film_library/
WORKDIR /film_library/

RUN go mod download
RUN GOOS=linux go build -o app ./cmd/main.go

CMD ["./app"]