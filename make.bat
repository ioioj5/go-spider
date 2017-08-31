@echo off

setlocal

if exist make.bat goto run
echo make.bat must be run from its folder
goto END

: run

chcp 437

SET GOROOT=D:\server\go
SET GOPATH=%GOPATH%;d:\git\go-spider

go get github.com/go-sql-driver/mysql

go build -v -o bin/parse parse


:end
echo finished
pause
