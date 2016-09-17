$ErrorActionPreference = "Stop"

Move-Item .\go\src\github.com\pivotal-cf\credhub-cli .\credhub-cli
Remove-Item -Recurse -Force .\go
echo "After removing go directory, expecting just credhub-cli: "
dir
mkdir go\src\github.com\pivotal-cf\
Move-Item credhub-cli go\src\github.com\pivotal-cf\credhub-cli

$download_dir = "C:\downloads"

md $download_dir -Force

$git_url = "https://github.com/git-for-windows/git/releases/download/v2.9.0.windows.1/Git-2.9.0-64-bit.exe"
$git_path = "$download_dir\Git-2.9.0-64-bit.exe"

$go_url = "https://storage.googleapis.com/golang/go1.6.2.windows-amd64.msi"
$go_path = "$download_dir\go1.6.2.windows-amd64.msi"

$wc = New-Object System.Net.WebClient

Write-Output "Downloading Git..."
$wc.DownloadFile($git_url, $git_path)

Write-Output "Downloading Go..."
$wc.DownloadFile($go_url, $go_path)

& $git_path /SP- /SILENT

Start-Process -FilePath msiexec -ArgumentList /i, $go_path, /quiet -Wait