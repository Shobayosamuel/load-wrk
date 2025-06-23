FROM golang:1.21-alpine

WORKDIR /app
COPY . .

RUN go build -o load-wrk ./cmd/load-wrk

ENTRYPOINT ["./load-wrk"]