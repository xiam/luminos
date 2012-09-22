# Getting started with Luminos

*It's dangerous to go alone!* Here you have some recommendations that may come in handy after [downloading](/download) Luminos.

## Install Luminos from a precompiled package

The quickest way to get started and it does not require Go to be installed, just download the appropriate executable package from the [downloads](/download) page and run the example from a Terminal:

```bash
$ cd ~/projects
$ git clone https://github.com/xiam/luminos-doc.git
$ cd ~/luminos-doc
$ luminos
```


## Install Luminos from source (using Go)

Make sure you have a working [Go](http://golang.org) development environment before installing from source.

```bash
$ go get -u github.com/xiam/luminos
$ go install github.com/xiam/luminos
$ cd ~/projects
$ git clone https://github.com/xiam/luminos-doc.git
$ cd ~/luminos-doc
$ luminos
```

## Use Luminos with your favorite web server

If you want to use a web server like **Apache2** or **nginx**, Luminos can be configured to listen on a FastCGI port, edit the
``settings.yaml`` file and set ``server/type`` to *fastcgi*:

```yaml
server:
  type: "fastcgi"
  bind: 127.0.0.1
  port: 9000
```

Then take a look on how to configure your web server, examples are included on the left menu.

## Running Luminos standalone

Luminos can run without an external web server, it has its own standalone server. Just tweak the ``settings.yaml`` file and
change ``server/type`` to *standalone* and a simple HTTP server will start listening on the configured server and port.

```yaml
server:
  type: "standalone"
  bind: 127.0.0.1
  port: 9000
 ```

## Starting the Luminos server

Whether you run it as standalone or FastCGI, you start Luminos just the same:

```bash
$ cd ~/projects
$ git clone https://github.com/xiam/luminos-doc.git
$ cd ~/luminos-doc
$ luminos
2012/08/18 12:39:59 FCGI server listening at 127.0.0.1:9000.
```

