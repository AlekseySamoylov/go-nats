############################
# STEP 1 build executable binary
############################
FROM golang:alpine AS builder
RUN apk update && apk add --no-cache git
WORKDIR $GOPATH/src/github.com/alekseysamoylov/dealer
COPY . .

WORKDIR $GOPATH/src/github.com/alekseysamoylov/dealer
RUN go build -o /go/bin/dealer

############################
# STEP 2 build a small image
############################
FROM alpine:3.10
# Copy our static executable.
RUN mkdir /usr/share/dealer
RUN mkdir /var/log/dealer
COPY --from=builder /go/bin/dealer /usr/share/dealer/dealer
WORKDIR /usr/share/dealer
# Run the metricbeat binary.
ENTRYPOINT ["./dealer"]
