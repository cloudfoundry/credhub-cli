.PHONY: all build ci clean dependencies format ginkgo test

ifeq ($(GOOS),windows)
DEST = build/credhub.exe
else
DEST = build/credhub
endif

ifndef VERSION
VERSION = dev
endif

all: test clean build

clean:
	rm -rf build

format:
	go fmt .

ginkgo:
	go run -mod=mod github.com/onsi/ginkgo/v2/ginkgo -r --randomize-suites --randomize-all --race -p 2>&1

test: format ginkgo

ci: ginkgo

build:
	mkdir -p build
	CGO_ENABLED=0 go build -o $(DEST) -ldflags "-X code.cloudfoundry.org/credhub-cli/version.Version=${VERSION}"
