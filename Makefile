default_target: all

VERSION=$(shell grep "var version" luminos.go | sed s/'"'//g | sed s/.*=//g | tr -d ' ')
PKG_NAME=luminos-$(VERSION)-$(shell uname -s | tr '[A-Z]' '[a-z]')-$(shell uname -m | tr '[A-Z]' '[a-z]')

install:
	sudo install

package:
	rm -rf $(PKG_NAME)
	mkdir -p $(PKG_NAME)
	go build luminos.go
	mv luminos $(PKG_NAME)/luminos
	cp LICENSE $(PKG_NAME)/
	tar cvzf $(PKG_NAME).tar.gz $(PKG_NAME)
	rm -rf $(PKG_NAME)
