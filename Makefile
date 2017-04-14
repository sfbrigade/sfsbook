
help:
	@echo Available targets are listed in the Makefile:
	@cat Makefile

get-deps:
	go get -t ./...

build: get-deps
	go build

test:
	go test -v ./...

start:
	./sfsbook -init_passwords

open:
	open "https://localhost:10443/index.html" #works on Mac

clean:
	rm -f state/{*.dat,*.pem}
	rm -rf state/*.bleve