FROM alpine:latest

RUN apk add --no-cache curl openjdk21

RUN curl -fsSL https://github.com/Zigl3ur/mcjar/releases/latest/download/mcjar_linux_amd64_musl \
    -o /usr/local/bin/mcjar \
    && chmod +x /usr/local/bin/mcjar