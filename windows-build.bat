@ECHO OFF
ECHO Fetching dependencies...
"go.exe" get -v -u ./...
ECHO Building...
"go.exe" build -o certinject.exe -v -ldflags '-w -s'
PAUSE