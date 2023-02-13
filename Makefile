build:
	CGO_ENABLE=0 go build -ldflags "-w -s" -o bin/mde
copy: build
	cp bin/mde /usr/local/bin

