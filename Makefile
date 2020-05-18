PREFIX ?= /usr/local

build:
	echo "Building pacheck binary"
	go build -o bin/pacheck main.go

run:
	go run main.go

install: build
	echo "$(PREFIX)"
	cp ./bin/pacheck $(PREFIX)/bin/pacheck

uninstall:
	echo "Removing pacheck"
	rm -vf $(PREFIX)/bin/pacheck