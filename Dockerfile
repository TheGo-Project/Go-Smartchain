# Support setting various labels on the final image
ARG COMMIT=""
ARG VERSION=""
ARG BUILDNUM=""

# Build Geth in a stock Go builder container
FROM golang:1.22-alpine as builder

ENV GO111MODULE=on \
    CGO_ENABLED=0 \
    GOOS=linux \
    GOARCH=amd64

RUN apk add --no-cache gcc musl-dev linux-headers git

# Get dependencies - will also be cached if we won't change go.mod/go.sum
COPY go.mod /go-ethereum/
COPY go.sum /go-ethereum/
RUN cd /go-ethereum && go mod download

ADD . /go-ethereum
RUN cd /go-ethereum && go run build/ci.go install -static ./cmd/geth

## Pull Geth into a second stage deploy alpine container
#FROM alpine:latest
#
#RUN apk add --no-cache ca-certificates
#COPY --from=builder /go-ethereum/build/bin/geth /usr/local/bin/

## backend
WORKDIR /app
COPY userInterface/back/go.mod .
COPY userInterface/back/go.sum .
RUN go mod download
COPY userInterface/back/. .
RUN go build -o main ./cmd/server.go
#FROM debian:bookworm-slim

FROM alpine:latest

RUN apk add --no-cache ca-certificates
COPY --from=builder /go-ethereum/build/bin/geth /usr/local/bin/

COPY --from=builder /app/main /main

EXPOSE 8545 8546 30303 30303/udp 8080

COPY wrapper.sh /wrapper.sh
RUN chmod +x /wrapper.sh
CMD ["sh", "/wrapper.sh"]
#CMD sleep infinity

# Add some metadata labels to help programmatic image consumption
ARG COMMIT=""
ARG VERSION=""
ARG BUILDNUM=""

LABEL commit="$COMMIT" version="$VERSION" buildnum="$BUILDNUM"
