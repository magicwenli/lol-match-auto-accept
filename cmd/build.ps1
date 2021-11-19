$git_sha = (git rev-parse --short HEAD)
$git_branch = (git branch --show-current)

go install
go get -u github.com/hallazzang/syso/...
syso
go build -ldflags -H=windowsgui -o build\lmaa_$git_sha.exe

Write-Host "::set-output name=FileName::.\build\lmaa_$git_sha.exe"
Write-Host "::set-output name=FileNamePartial::lmaa_$git_sha"