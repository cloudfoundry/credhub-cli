powershell -command set-executionpolicy remotesigned

set work_dir = "C:\concourse"
set download_dir = "C:\downloads"

md %download_dir% -Force

set git_url = "https://github.com/git-for-windows/git/releases/download/v2.9.0.windows.1/Git-2.9.0-64-bit.exe"
set git_path = "%download_dir%\Git-2.9.0-64-bit.exe"

set go_url = "https://storage.googleapis.com/golang/go1.6.2.windows-amd64.msi"
set go_path = "%download_dir%\go1.6.2.windows-amd64.msi"

set wc = New-Object System.Net.WebClient

Write-Output "Downloading Git..."
%wc%.DownloadFile($git_url, $git_path)

Write-Output "Downloading Go..."
%wc%.DownloadFile(%go_url%, %go_path%)

%git_path% /SP- /SILENT

msiexec /passive /i $go_path

md %work_dir% -Force

SET GOPATH=%CD%\go
SET PATH=C:\Go\bin;%GOPATH%\bin;C:\Program Files\Git\cmd;%PATH%

cd %GOPATH%\src\github.com\pivotal-cf\credhub-cli

go get github.com/onsi/ginkgo/ginkgo
go get -v -t ./...

ginkgo -r -randomizeSuites -randomizeAllSpecs -race