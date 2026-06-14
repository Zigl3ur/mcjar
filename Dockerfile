FROM alpine:latest AS setup

RUN apk add --no-cache curl

RUN curl -fsSL https://github.com/Zigl3ur/mcjar/releases/latest/download/mcjar-linux-musl-amd64 \
    -o /mcjar && chmod +x /mcjar

FROM alpine:3.20

RUN apk add --no-cache openjdk21-jre

COPY --from=setup /mcjar /usr/local/bin/mcjar

ENTRYPOINT ["/usr/local/bin/mcjar"]
