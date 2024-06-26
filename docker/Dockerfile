FROM golang:1.22-alpine

WORKDIR /go/src/app

COPY ./app/go.* ./
RUN go mod download && go mod verify

COPY ./app .

RUN go install github.com/air-verse/air@v1.52.3

CMD ["air", "-c", ".air.toml"]