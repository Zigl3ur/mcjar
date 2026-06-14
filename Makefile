ifndef VERSION
  $(error VERSION arg is required. Usage: make <target> VERSION=1.2.3)
endif

LDFLAGS := -ldflags "-X github.com/Zigl3ur/mcjar/cmd.version=$(VERSION) -s -w"

build-linux-glibc:
	GOOS=linux GOARCH=amd64 CGO_ENABLED=0 \
	  go build $(LDFLAGS) -o ./bin/mcjar-linux-amd64

build-linux-musl:
	GOOS=linux GOARCH=amd64 CGO_ENABLED=0 \
	  go build $(LDFLAGS) -ldflags='-extldflags "-static"' \
	  -o ./bin/mcjar-linux-musl-amd64 .

build-windows:
	GOOS=windows GOARCH=amd64 CGO_ENABLED=0 \
	  go build $(LDFLAGS) -o ./bin/mcjar-windows-amd64.exe .

checksums:
	sha256sum \
	  ./bin/mcjar-linux-amd64 \
	  ./bin/mcjar-linux-musl-amd64 \
	  ./bin/mcjar-windows-amd64.exe \
	  > ./bin/mcjar-checksums.txt

all: build-linux-glibc build-linux-musl build-windows checksums