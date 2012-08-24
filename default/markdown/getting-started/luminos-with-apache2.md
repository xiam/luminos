# Setting up Apache 2 for Luminos

This configuration relies on the ``mod_fastcgi`` module and requires Luminos to be configured as a FastCGI server.

    LoadModule fastcgi_module modules/mod_fastcgi.so

You can configure your virtual host like this:

    <VirtualHost *:80>
      ServerName example.org

      DocumentRoot    /opt/luminos
      FastCgiExternalServer /opt/luminos/luminos -host 127.0.0.1:9000

      RewriteEngine On
      RewriteCond %{DOCUMENT_ROOT}%{REQUEST_FILENAME} !-f
      RewriteRule ^(.*)$ /luminos [QSA,L]

      <Directory /opt/luminos/luminos/>
        SetHandler fastcgi-script
        Order Deny,Allow
        Allow from all
      </Directory>

    </VirtualHost>

