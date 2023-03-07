FROM golang:1.19-alpine AS builder

RUN mkdir /app
WORKDIR /app

COPY . .
ENV GOOS=linux
ENV GOARCH=amd64
RUN go build -ldflags="-w -s" -o /go/bin/app

FROM alpine:latest

COPY --from=builder /go/bin/app ./app

CMD [ "./app" ]