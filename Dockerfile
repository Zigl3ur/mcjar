FROM alpine:latest

RUN apk add --no-cache curl openjdk21

# Download mcjar binary and make it executable
RUN curl -fsSL -o /usr/local/bin/mcjar \
    https://github.com/Zigl3ur/mcjar/releases/download/0.0.1/mcjar_0.0.1_linux_amd64_musl \
    && chmod +x /usr/local/bin/mcjar