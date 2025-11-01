FROM golang:1.25-alpine AS builder

RUN apk update && \
    apk add --no-cache -U tzdata make git build-base ca-certificates && \
    rm -rf /var/cache/apk/*

# ENV vars
ENV GO111MODULE=on \
  CGO_ENABLED=0 \
  GOOS=linux \
  GOARCH=amd64

# Move to working directory /build
WORKDIR /build

# Copy the code into the container
COPY . .

# Copy and download dependency using go mod
RUN make mod

# Build the application
RUN make build

# Install DB migration tool, then copy the binary to user bin path
# NOTE: the output binary is placed in a directory like `linux_amd64`, but the name depends on
# the actual os arch. So I use `find` command here.
RUN go install github.com/rubenv/sql-migrate/...@latest && \
    find /go/bin/ -name 'sql-migrate' -exec cp "{}" ./sql-migrate \;

############################################################
# Build a small image
FROM alpine:3.22
WORKDIR /app
RUN apk update && \
    apk add --no-cache ca-certificates make openssl aws-cli postgresql-client && \
    rm -rf /var/cache/apk/*

# For running app
COPY --from=builder /build/app .
COPY --from=builder /build/localpaas_app/db ./localpaas_app/db
COPY --from=builder /build/config ./config
COPY --from=builder /build/scripts ./scripts
COPY --from=builder /build/Makefile .
COPY --from=builder /build/sql-migrate .

# Command to run
CMD ["sh","-c","./app"]
