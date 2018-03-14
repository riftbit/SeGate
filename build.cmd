set CGO_ENABLED=0
set GOOS=linux
set GOARCH=amd64
set SGVERS=v1.0.1

FOR /F "tokens=*" %%a in ('git rev-parse --short HEAD') do SET SGBUILD=%%a

set dd=%date:~0,2%
set mm=%date:~3,2%
set yyyy=%date:~6,4%

go build -a -installsuffix cgo -ldflags "-s -w -X 'main.version=%SGVERS%' -X 'main.build=%SGBUILD%' -X 'main.buildDate=%dd%-%mm%-%yyyy%'" ./