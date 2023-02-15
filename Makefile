build:
	CGO_ENABLE=0 go build -ldflags "-w -s" -o bin/mde
copy: build
	cp bin/mde /usr/local/bin
test:
	go test ./...
cp-pre-commit:
	cp .github/pre-commit .git/hooks/
pre-commit: test
