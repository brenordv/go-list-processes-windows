echo off

echo Checking required build folders...
if not exist ".tmp" mkdir .tmp
if not exist ".tmp\win64" mkdir .tmp\win64

echo Setting OS and Architecture for windows / amd64...
set GOOS=windows
set GOARCH=amd64


echo Building...
go build -o .tmp\win64\go-list-prc.exe .\cmd\go-list-process

echo Creating archive for version %*...
7z a -tzip .\.dist\go-list-processes--windows-amd64--%*.zip .\.tmp\win64\*.exe readme.md

echo Copying application to current directory...
copy /b/v/y .\.tmp\win64\*.exe .\

echo All done!
