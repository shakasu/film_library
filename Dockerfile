FROM golang:latest

RUN go version

COPY . /avitoTech/
WORKDIR /avitoTech/

# build go app
RUN go mod download
RUN GOOS=linux go build -o app ./cmd/main.go

CMD ["./app"]