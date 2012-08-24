default_target: all

VERSION=$(shell grep "var version" luminos.go | sed s/'"'//g | sed s/.*=//g | tr -d ' ')
PKG_NAME=luminosd-$(VERSION)-$(shell uname -s | tr '[A-Z]' '[a-z]')-$(shell uname -m | tr '[A-Z]' '[a-z]')

all:
	rm -rf dist/
	mkdir -p dist/
	cp -rf settings.yaml default dist/
	go build luminos.go
	mv luminos dist/luminosd

install:
	mkdir -p /usr/local/luminosd

package:
	rm -rf $(PKG_NAME)
	mkdir -p $(PKG_NAME)
	cp -rf settings.yaml default $(PKG_NAME)/
	go build luminos.go
	mv luminos $(PKG_NAME)/luminosd
	tar cvzf $(PKG_NAME).tar.gz $(PKG_NAME)
	rm -rf $(PKG_NAME)
