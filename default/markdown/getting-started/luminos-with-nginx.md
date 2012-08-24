# Setting up nginx for Luminos

This is an example configuration setting that allows [nginx](http://nginx.org) to play along with Luminos.

It requires Luminos to be configured as a FastCGI server.

    server {
      server_name example.org;
      location / {
        fastcgi_pass   127.0.0.1:9000;
        include        fastcgi_params;
      }
    }


