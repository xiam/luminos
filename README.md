# Luminos, markdown server

[Luminos][5] is a tiny HTTP/FastCGI server that transforms [markdown][3]
files into HTML on the fly, it was highly inspired by [werc][1].

[Luminos][5] is written in the Go programming language.

```
Luminos Markdown Server (0.7) - https://menteslibres.net/luminos
by Carlos Reventlov <carlos@reventlov.com>

Usage: luminos <arguments> <command>

Available commands for luminos:

        help            Shows information about the given command.
        init            Initializes a working directory with a Luminos base project.
        run             Runs a luminos server.
        version         Prints software version.

Use "luminos help <command>" to view more information about a command.
```

## Features

* Luminos does not require a database.
* It is able to server serve multiple virtual hosts.
* Works out of the box on Linux, FreeBSD, OpenBSD, OSX <s>and Windows</s> with
  or without using an external web server.
* Simple configuration changes do not require a restart.

## Getting Luminos

If you want to build from source, a [Go][2] development environment is
required.

Use `go get ` to download, compile and install Luminos to your
`$GOPATH/bin`

```sh
$ go get menteslibres.net/luminos
```

Then use `luminos init` to create an example project and `luminos run` to serve
it.


```sh
$ mkdir -p ~/projects/
$ cd ~/projects/
$ luminos init luminos-example-project
$ cd luminos-example-project
$ luminos run
```

## Documentation

See the [project's page][5] for documentation, tips and tricks.

## License

Luminos is released under the MIT License:

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

[1]: http://werc.cat-v.org
[2]: http://golang.org
[3]: http://daringfireball.net/projects/markdown/
[4]: https://github.com/xiam/luminos
[5]: https://menteslibres.net/luminos
