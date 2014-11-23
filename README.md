# Luminos, markdown server

[Luminos][5] is a tiny HTTP/FastCGI server written in [Go][2] that transforms
[markdown][3] files into HTML on the fly. It plays well with NGINX and Apache2.

## Getting Luminos

In order to download and install Luminos, a [working Go
environment](https://golang.org/doc/install) is required.

If you already have Go, you may install Luminos by issuing the following
command:

```sh
go get menteslibres.net/luminos
```

## Usage

```
rev@localhost > luminos

Luminos Markdown Server (0.9) - https://menteslibres.net/luminos
by J. Carlos Nieto <jose.carlos@menteslibres.net>

Usage: luminos <arguments> <command>

Available commands for luminos:

        help            Shows information about the given command.
        init            Creates a new Luminos site scaffold in the given PATH.
        run             Runs a luminos server.
        version         Prints software version.

Use "luminos help <command>" to view more information about a command.
```

Use `luminos init` to create an empty site:

```sh
cd ~/projects
luminos init test-site
```

then you can ask `luminos run` to serve it:

```sh
cd test-site
luminos run
```

If you want to use Luminos with Apache or NGINX see the [Getting
started](https://menteslibres.net/luminos/getting-started) page.

## Documentation

See the [project's page][5] for documentation, tips and tricks.

## Licenses and acknowledgements

Luminos is released under the MIT License.

> Copyright (c) 2012-2013 José Carlos Nieto, https://menteslibres.net/luminos
>
> Permission is hereby granted, free of charge, to any person obtaining
> a copy of this software and associated documentation files (the
> "Software"), to deal in the Software without restriction, including
> without limitation the rights to use, copy, modify, merge, publish,
> distribute, sublicense, and/or sell copies of the Software, and to
> permit persons to whom the Software is furnished to do so, subject to
> the following conditions:
>
> The above copyright notice and this permission notice shall be
> included in all copies or substantial portions of the Software.
>
> THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND,
> EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF
> MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND
> NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT HOLDERS BE
> LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER IN AN ACTION
> OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN CONNECTION
> WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.

The [Hyde theme](https://github.com/poole/hyde) was created by [Mark
Otto](http://jekyllrb.com/) for [Jekyll](http://jekyllrb.com/) and released
under the MIT License.

> Copyright (c) 2013 Mark Otto.
>
> Permission is hereby granted, free of charge, to any person obtaining a copy of
> this software and associated documentation files (the "Software"), to deal in
> the Software without restriction, including without limitation the rights to
> use, copy, modify, merge, publish, distribute, sublicense, and/or sell copies
> of the Software, and to permit persons to whom the Software is furnished to do
> so, subject to the following conditions:
>
> The above copyright notice and this permission notice shall be included in all
> copies or substantial portions of the Software.
>
> THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
> IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
> FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
> AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
> LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
> OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
> SOFTWARE.

[1]: http://werc.cat-v.org
[2]: http://golang.org
[3]: http://daringfireball.net/projects/markdown/
[4]: https://github.com/xiam/luminos
[5]: https://menteslibres.net/luminos
