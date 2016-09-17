powershell -command set-executionpolicy remotesigned

powershell task-repo\ci\tasks\unit-tests-windows\install_go_and_git.ps1

SET GOPATH=%CD%\go
SET PATH=C:\Go\bin;%GOPATH%\bin;C:\Program Files\Git\cmd;%PATH%

cd %GOPATH%\src\github.com\pivotal-cf\credhub-cli

go get github.com/onsi/ginkgo/ginkgo
go get -v -t ./...

ginkgo -r -randomizeSuites -randomizeAllSpecs -race