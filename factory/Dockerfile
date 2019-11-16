############################
# STEP 1 build executable binary
############################
FROM golang:alpine AS builder
RUN apk update && apk add --no-cache git
WORKDIR $GOPATH/src/github.com/alekseysamoylov/factory
COPY . .

WORKDIR $GOPATH/src/github.com/alekseysamoylov/factory
RUN go build -o /go/bin/factory

############################
# STEP 2 build a small image
############################
FROM alpine:3.10
# Copy our static executable.
RUN mkdir /usr/share/factory
RUN mkdir /var/log/factory
COPY --from=builder /go/bin/factory /usr/share/factory/factory
WORKDIR /usr/share/factory
# Run the metricbeat binary.
ENTRYPOINT ["./factory"]
