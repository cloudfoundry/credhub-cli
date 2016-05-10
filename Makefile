.DEFAULT_GOAL := test

build :
		go build

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