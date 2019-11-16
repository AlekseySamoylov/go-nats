############################
# STEP 1 build executable binary
############################
FROM golang:alpine AS builder
RUN apk update && apk add --no-cache git
WORKDIR $GOPATH/src/github.com/alekseysamoylov/dialer
COPY . .

WORKDIR $GOPATH/src/github.com/alekseysamoylov/dialer
RUN go build -o /go/bin/dialer

############################
# STEP 2 build a small image
############################
FROM alpine:3.10
# Copy our static executable.
RUN mkdir /usr/share/dialer
RUN mkdir /var/log/dialer
COPY --from=builder /go/bin/dialer /usr/share/dialer/dialer
WORKDIR /usr/share/dialer
# Run the metricbeat binary.
ENTRYPOINT ["./dialer"]
