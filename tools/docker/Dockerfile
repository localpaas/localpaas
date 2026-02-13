# NOTE: use linux/amd64 as target platform solves weird issue with my Apple M1 (not sure why)
FROM --platform=linux/amd64 golang:1.26-alpine

RUN apk update && apk upgrade && apk add --update alpine-sdk postgresql-client

RUN wget -O- -nv https://raw.githubusercontent.com/golangci/golangci-lint/HEAD/install.sh | sh -s v2.9.0 &&\
    go install github.com/rubenv/sql-migrate/...@latest &&\
    go install github.com/swaggo/swag/cmd/swag@v1.16.6

WORKDIR /app
