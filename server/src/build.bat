SET CGO_ENABLED=0
SET GOOS=linux
SET GOARCH=amd64
SET GOPATH=C:\work\go\src\server
go build -o .\bin\admin C:\work\go\src\server\src\admin\main.go
pause