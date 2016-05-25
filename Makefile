.DEFAULT_GOAL := test

.PHONY : build

ifeq ($(GOOS),windows)
DEST = build/cm.exe
else
DEST = build/cm
endif

BUILD_NUMBER := $(shell date +%s)

GOFLAGS := -o $(DEST)
GOFLAGS := $(GOFLAGS) -ldflags "-X github.com/pivotal-cf/cm-cli/version.BuildNumber=$(BUILD_NUMBER)"

dependencies :
		go get github.com/onsi/ginkgo/ginkgo
		go get golang.org/x/tools/cmd/goimports
		go get github.com/maxbrunsfeld/counterfeiter
		go get -v -t ./...

format : dependencies
		goimports -w .
		go fmt .

ginkgo : dependencies
		ginkgo -r -randomizeSuites -randomizeAllSpecs -race

test : format ginkgo

ci : ginkgo

build :
		mkdir -p build
		go build $(GOFLAGS)
