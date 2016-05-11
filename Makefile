.DEFAULT_GOAL := test

.PHONY : build

dependencies :
		go get github.com/onsi/ginkgo/ginkgo
		go get golang.org/x/tools/cmd/goimports
		go get -v -t ./...

format :
		goimports -w .
		go fmt .

ginkgo : dependencies
		ginkgo -r -randomizeSuites -randomizeAllSpecs -race

test : format ginkgo

ci : ginkgo

build :
		mkdir -p build
ifeq ($(GOOS),windows)
		go build -o build/cm.exe
else
		go build -o build/cm
endif
