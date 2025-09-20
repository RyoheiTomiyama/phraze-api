FROM golang:1.25-alpine

WORKDIR /go/src/app

COPY ./app/go.* ./
RUN go mod download && go mod verify

COPY ./app .

RUN curl -sSf https://atlasgo.sh | sh

RUN go install github.com/air-verse/air@v1.52.3

CMD ["air", "-c", ".air.toml"]