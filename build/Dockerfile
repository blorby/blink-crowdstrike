FROM golang:1.16.3 AS builder

WORKDIR /go/src/github.com/blinkops/blink-EXAMPLE
COPY .. .

# install certificates
RUN apk --no-cache add ca-certificates

# gzip all the yaml files.
RUN find . -type f -name '*.yaml' -exec sh -c 'echo "${1%}";gzip "${1%}"' sh {} \;

# build the plugin
RUN go build -tags netgo -ldflags '-w -extldflags "-static"' -o /go/bin/blink-EXAMPLE ./cmd/main

##################
FROM scratch AS plugin

WORKDIR /blink-EXAMPLE

# copy the certificates
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/

# copy the files that are needed in order run.
COPY --from=builder /go/bin/blink-EXAMPLE .
COPY --from=builder /go/src/github.com/blinkops/blink-EXAMPLE/EXAMPLE-openapi.yaml.gz .
COPY --from=builder /go/src/github.com/blinkops/blink-EXAMPLE/mask.yaml.gz .

# Expose the gRPC port
EXPOSE 1337

ENV PROD true

ENTRYPOINT ./blink-EXAMPLE
