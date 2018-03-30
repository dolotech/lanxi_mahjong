set GOPATH=E:\lanxi-mahjong-server\trunk
set GOARCH=amd64
set GOOS=linux
cd bin

go build -o server -ldflags "-X main.VERSION=1.0.4 -X 'main.BUILD_TIME=`date`' -s -w" ../src/server.go
go build -o robot -ldflags "-w -s" ../src/robot.go


pause


