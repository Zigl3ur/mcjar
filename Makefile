VERSION ?= 0.0.1

LDFLAGS := -ldflags "-X github.com/Zigl3ur/mcjar/cmd.version=2 -s -w"

build-linux-glibc:
	GOOS=linux GOARCH=amd64 CGO_ENABLED=0 \
	  go build $(LDFLAGS) -o ./bin/mcjar-$(VERSION)-linux-amd64

build-linux-musl:
	GOOS=linux GOARCH=amd64 CGO_ENABLED=0 \
	  go build $(LDFLAGS) -ldflags='-extldflags "-static"' \
	  -o ./bin/mcjar-$(VERSION)-linux-musl-amd64 .

build-windows:
	GOOS=windows GOARCH=amd64 CGO_ENABLED=0 \
	  go build $(LDFLAGS) -o ./bin/mcjar-$(VERSION)-windows-amd64.exe .

checksums:
	sha256sum \
	  ./bin/mcjar-$(VERSION)-linux-amd64 \
	  ./bin/mcjar-$(VERSION)-linux-musl-amd64 \
	  ./bin/mcjar-$(VERSION)-windows-amd64.exe \
	  > ./bin/mcjar-$(VERSION)-checksums.txt

all: build-linux-glibc build-linux-musl build-windows checksums