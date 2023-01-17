# Build the service to a binary
FROM golang:1.19.5-alpine AS builder

# Install packages
RUN apk add --no-cache bash gcc musl-dev linux-headers git

# Compile application
WORKDIR /go/src/github.com/Guillembonet/nginx-wg-proxy
ADD . .
RUN go build -o build/main main.go

# Copy and run the made binary
FROM alpine:3.17

RUN apk add --no-cache --update ca-certificates nginx wireguard-tools

COPY --from=builder /go/src/github.com/Guillembonet/nginx-wg-proxy/build/main /usr/bin/api

ENTRYPOINT ["/usr/bin/api"]
