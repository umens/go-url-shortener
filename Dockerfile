# Build stage
FROM golang:1.8-alpine as builder
ENV APP_DIR /go/src/github.com/umens/go-url-shortener
ADD . $APP_DIR
WORKDIR $APP_DIR
RUN apk --update add git && rm -rf /var/cache/apk/*
RUN go get -u github.com/go-redis/redis
RUN go build -o go-url-shortener

# run binary
FROM alpine:latest
WORKDIR /app
COPY --from=builder /go/src/github.com/umens/go-url-shortener/go-url-shortener .
EXPOSE 8080
CMD ["./go-url-shortener"]