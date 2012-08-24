# Getting started with Luminos

*It's dangerous to go alone!* Here you have some recommendations that may come in handy after [downloading](/download) Luminos.

## Play it on your favorite web server.

If you want to use a web server like Apache2 or nginx, **Luminos** can be configured to listen on a FastCGI port, edit the
``settings.yaml`` file and set ``server/type`` to *fastcgi*:

    server:
      type: fastcgi
      bind: 127.0.0.1
      port: 9000

Then take a look on how to configure your web server, examples are included on the left menu.

## Running standalone

Luminos can run without an external web server, it has its own standalone server. Just tweak the ``settings.yaml`` file and
change ``server/type`` to *standalone* and a simple HTTP server will start listening on the configured server and port.

    server:
      type: standalone
      bind: 127.0.0.1
      port: 9000

## Starting the Luminos server

Whether you run it as standalone or FastCGI, you start Luminos just the same:

    % cd $LUMINOSPATH
    % ./luminosd
    2012/08/18 12:39:59 FCGI server listening at 127.0.0.1:9000.

