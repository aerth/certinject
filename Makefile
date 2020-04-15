buildflags = -v -ldflags '-w -s'
certinject.exe: *.go
	env GOOS=windows GOARCH=amd64 go build $(buildflags) -o $@
	strip $@
certinject-linux-amd64: *.go
	env GOOS=linux GOARCH=amd64 go build $(buildflags) -o $@
	strip $@
certinject-osx-amd64: *.go
	env GOOS=darwin GOARCH=amd64 go build $(buildflags) -o $@
	#strip $@
all: certinject.exe certinject-linux-amd64 certinject-osx-amd64
clean:
	rm -vf certinject.exe certinject-linux-amd64 certinject-osx-amd64
