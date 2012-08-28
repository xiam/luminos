# Luminos directory structure

Luminos requires an special directory structure to serve documents and a configuration
file that rules the current host.

## luminos.yaml

A [YAML](http://www.yaml.org/) formatted file. This is an example.

    page:
      brand: "Luminos"
      head:
        title: "Luminos, markdown server"
      body:
        title: "Luminos project"
        menu:
          - { text: "Getting started", url: "/getting-started" }
          - { text: "Templates", url: "/templates" }
          - { text: "Source code", url: "https://github.com/xiam/luminos" }
        menu_pull:
          - { text: "Home", url: "/" }

Be aware to use two-space tabs instead of the tab character.

## markdown/

A directory with markdown or HTML files on it.

Add the ``.md`` extension to your text files to be parsed as markdown.

Every time Luminos receives an URL, it will look for the most appropriate file in the markdown
directory, for example, if the user request and URL like ``/path/to/resource``, Luminos
will try to find and serve these files:

    markdown/path/to/resource.md
    markdown/path/to/resource.html
    markdown/path/to/resource

If a directory is found instead, Luminos will look for an ``index.md`` or ``index.html``
file on that directory.

    markdown/path/to/resource/index.md
    markdown/path/to/resource/index.html

## templates/

A directory that contains ``template/html`` templates. Tipically just an ``index.tpl`` file.

The ``index.tpl`` file determines how the whole page will be seen according to some
[template variables](/templates/template-variables).

## webroot/

This directory contains all the files that must be served without being parsed. If you're using CSS or Javascript,
you should drop it here.

Luminos checks here first, even before looking into the ``markdown`` directory.

If you want to add a ``favicon.ico`` put it there and access it with ``/favicon.ico``.

# Examples

Luminos has a ``default/`` directory that contains this whole site, you can give it a look to have a starting point.

Here are other examples:

* The gosexy.org site [source](https://github.com/gosexy/gosexy.org).
