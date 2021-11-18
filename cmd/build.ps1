$git_sha = (git rev-parse --short HEAD)
$git_branch = (git branch --show-current)

go install
go build -o build\lmaa_$git_sha.exe main.go