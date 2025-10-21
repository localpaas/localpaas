# NOTE: use linux/amd64 as target platform solves weird issue with my Apple M1 (not sure why)
FROM --platform=linux/amd64 golang:1.25.3-alpine

RUN apk update && apk upgrade && apk add --update alpine-sdk python3 postgresql-client

RUN wget -O- -nv https://raw.githubusercontent.com/golangci/golangci-lint/HEAD/install.sh | sh -s v2.1.6 &&\
    go install github.com/rubenv/sql-migrate/...@latest &&\
    go install go.uber.org/mock/mockgen@latest &&\
    go install github.com/swaggo/swag/cmd/swag@v1.16.6

WORKDIR /app
