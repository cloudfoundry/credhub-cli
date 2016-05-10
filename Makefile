dependencies :
		go get github.com/onsi/ginkgo/ginkgo
		go get golang.org/x/tools/cmd/goimports

build :
		go build -v .

format :
		goimports -w .
		go fmt .

test : dependencies format build
		go get -v -t ./...
		ginkgo -r -randomizeSuites -randomizeAllSpecs -race

