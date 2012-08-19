all:
	rm -rf dist/
	mkdir -p dist/
	cp -rf settings.yaml default dist/
	go build luminos.go
	mv luminos dist/luminosd
install:
	mkdir -p /usr/local/luminosd
