FROM golang:1.16.3-alpine3.13 AS builder

WORKDIR /go/src/github.com/blinkops/blink-crowdstrike
COPY .. .

RUN apk --no-cache add ca-certificates

# compress the yaml files.
RUN find . -maxdepth 1 -type f -name '*.yaml' -exec sh -c 'echo "${1%}";gzip "${1%}"' sh {} \;
RUN go mod tidy
# build the plugin.
RUN go build -tags netgo -ldflags '-w -s -extldflags "-static"' -o /go/bin/blink-crowdstrike ./cmd/main

FROM scratch AS plugin
WORKDIR /blink-crowdstrike

# certificates need to run in scratch
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/

COPY --from=builder /go/bin/blink-crowdstrike .
COPY --from=builder /go/src/github.com/blinkops/blink-crowdstrike/custom_actions ./custom_actions
COPY --from=builder /go/src/github.com/blinkops/blink-crowdstrike/crowdstrike-openapi.yaml.gz .
COPY --from=builder /go/src/github.com/blinkops/blink-crowdstrike/mask.yaml.gz .

# Expose the gRPC port
EXPOSE 1337

ENV PROD true

ENTRYPOINT ["./blink-crowdstrike"]