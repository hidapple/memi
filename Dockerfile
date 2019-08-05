FROM golang:1.12-alpine

RUN apk add --update --no-cache git

ENV GO111MODULE=on

RUN mkdir /app
WORKDIR /app
COPY . .

RUN go mod download
RUN CGO_ENABLED=0 go build -o /app/bin/memi

ENTRYPOINT ["/app/bin/memi"]
