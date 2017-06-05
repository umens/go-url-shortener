FROM golang:1.8-alpine

ENV APP_DIR /go/src/go-url-shortener

ADD . $APP_DIR
WORKDIR $APP_DIR

RUN apk --update add git && rm -rf /var/cache/apk/*

# Install packages.
# Ideally we'd use `godep` for this, but to keep this short
# we'll just install them manually here.
RUN go get -u github.com/go-redis/redis

RUN go build -o go-url-shortener

EXPOSE 8080
CMD ./go-url-shortener