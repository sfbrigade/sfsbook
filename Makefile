
help:
	@echo Available targets are listed in the Makefile:
	@cat Makefile

build:
	go get
	go build

start:
	./sfsbook -init_passwords

open:
	open "https://localhost:10443/index.html" #works on Mac

clean:
	rm -f state/{*.dat,*.pem}
	rm -rf state/*.bleve