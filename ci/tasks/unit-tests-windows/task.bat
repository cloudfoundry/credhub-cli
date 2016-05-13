SET GOPATH=%CD%\go
SET PATH=C:\Go\bin;%GOPATH%\bin;C:\Program Files\Git\cmd;%PATH%

cd %GOPATH%\src\github.com\pivotal-cf\cm-cli

powershell -command set-executionpolicy remotesigned

go get github.com/onsi/ginkgo/ginkgo
go get -v -t ./...

ginkgo -r -randomizeSuites -randomizeAllSpecs -race