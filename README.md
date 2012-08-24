# What is Luminos?

Luminos is a server than transforms [markdown][3] code into HTML code right away, it was highly inspired by [werc][1] but build with less options in mind.

Here are some of Luminos' features:

* It's the perfect tool for publishing documentation.
* Does not require a database.
* Can serve both HTML or markdown files.
* Capable of serving multiple virtual hosts within the same process.
* Works out of the box on Linux, FreeBSD, OSX <s>and Windows</s> with or without an external web server.
* Written in [Go][2] and released as an Open Source [project][4].

# Building Luminos from source

Before building from source, please take a look at the Luminos [project page][5].

    go get -u github.com/xiam/luminos
    chdir $GOPATH/src/github.com/xiam/luminos
    gmake
    cd dist/
    ./luminosd

# License

Luminos is released under the MIT License:

> Copyright (c) 2012 José Carlos Nieto, http://xiam.menteslibres.org/
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
[4]: http://github.com/xiam/luminos
[5]: http://luminos.menteslibres.org/
