# Download Luminos

Here are some precompiled packages that you can use right out the box.

* Luminos for **MacOSX** <s>32</s> / [64](https://github.com/downloads/xiam/luminos/luminosd-0.1-darwin-x86_64.tar.gz)
* Luminos for **Linux** <s>32</s> / [64](https://github.com/downloads/xiam/luminos/luminosd-0.1-linux-x86_64.tar.gz)
* Luminos for **FreeBSD** <s>32</s> / [64](https://github.com/downloads/xiam/luminos/luminosd-0.1-freebsd-amd64.tar.gz)
* Luminos for **Windows** <s>32</s> / <s>64</s>

You can see all packages in the Luminos' [downloads](https://github.com/xiam/luminos/downloads) page at github.

What's next? Please read the [getting started](/getting-started) page.

# Didn't work? compile from source

Make sure you have a working [Go](http://golang.org) installation before you cast these spells.

    go get -u github.com/xiam/luminos
    chdir $GOPATH/src/github.com/xiam/luminos
    gmake
    cd dist/
    ./luminosd

