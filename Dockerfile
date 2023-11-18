ARG BUILD_IMAGE=golang:alpine
FROM ${BUILD_IMAGE} as Build

COPY . . 

RUN GOPATH= go build -o /main


# Just copy the built artefact to a small image.
# This should clock in below 50MB.
FROM alpine:latest

COPY --from=Build /main /main

# We always listen on 0.0.0.0 inside the docker container.
# By setting this explicit we override the
# SPACE_PUBSUB_DAEMON_ADDR which, has lower priority
ENTRYPOINT [ "/main"]
