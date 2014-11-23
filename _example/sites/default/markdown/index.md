# It works!

This is [Luminos][2] a markdown server written in [Go][1]. You can find some
tips and tricks at the [getting started][5] page.

[Luminos][2] is an Open Source project, feel free to [browse and hack][2] the
source and, if you find this project useful, please consider [making a
donation][4] to the [author][6].

Thanks for using [Luminos][3]!

## A few markdown examples

[Markdown](http://daringfireball.net/projects/markdown/) is a very comfortable
format for writing documents in plain text format.

Here are some examples on how your markdown code would be translated into HTML
by [Luminos][3].

<table class="table">
  <thead>
    <tr>
      <th>Markdown code</th>
      <th>Result</th>
    </tr>
  </thead>
  <tbody>
    <tr>
      <td>
        <code>**Bold text**</code>
      </td>
      <td>
        <strong>Bold text</strong>
      </td>
    </tr>
    <tr>
      <td>
        <code>*Italics*</code>
      </td>
      <td>
        <em>Italics</em>
      </td>
    </tr>
    <tr>
      <td>
        <code>~~Striked-through~~</code>
      </td>
      <td>
        <del>Striked-through</del>
      </td>
    </tr>
    <tr>
      <td>
        <code># First level header</code>
      </td>
      <td>
        <h1>First level header</h1>
      </td>
    </tr>
    <tr>
      <td>
        <code>## Second level header</code>
      </td>
      <td>
        <h2>Second level header</h2>
      </td>
    </tr>
    <tr>
      <td>
        <code>### Third level header</code>
      </td>
      <td>
        <h3>Third level header</h3>
      </td>
    </tr>
    <tr>
      <td>
        <code>#### Fourth level header</code>
      </td>
      <td>
        <h4>Fourth level header</h4>
      </td>
    </tr>
    <tr>
      <td>
        <code>##### Fifth level header</code>
      </td>
      <td>
        <h5>Fifth level header</h5>
      </td>
    </tr>
    <tr>
      <td>
        <code>[The Go Programming Language](http://golang.org)</code>
      </td>
      <td>
        <a href="http://golang.org">The Go Programming Language</a>
      </td>
    </tr>
    <tr>
      <td>
        <code>![A gopher](http://bit.ly/SLqdv6)</code>
      </td>
      <td>
        <img src="http://bit.ly/SLqdv6" alt="A gopher!" />
      </td>
    </tr>
    <tr>
      <td>
<pre><code>* List item 1
* List item 2
* List item 3</code></pre>
      </td>
      <td>
        <ul>
          <li>List item 1</li>
          <li>List item 2</li>
          <li>List item 3</li>
        </ul>
      </td>
    </tr>
    <tr>
      <td>
<pre><code>1. List item 1
2. List item 2
3. List item 3</code></pre>
      </td>
      <td>
        <ol>
          <li>List item 1</li>
          <li>List item 2</li>
          <li>List item 3</li>
        </ol>
      </td>
    </tr>
    <tr>
      <td>
<pre><code>Name    | Age
--------|------
Bob     | 27
Alice   | 23</code></pre>
      </td>
      <td>
        <table>
          <thead>
            <tr>
              <td>Name</td>
              <td>Age</td>
            </tr>
          </thead>
          <tbody>
            <tr>
              <td>Bob</td>
              <td>27</td>
            </tr>
            <tr>
              <td>Alice</td>
              <td>23</td>
            </tr>
          </tbody>
        </table>
      </td>
    </tr>
    <tr>
      <td>
<pre><code>```go
import "foo"

func main() {
  foo.Bar()
}
```</code></pre>
      </td>
      <td>
<pre><code class="go">import &quot;foo&quot;

func main() {
  foo.Bar()
}
</code></pre>
      </td>
    </tr>
    <tr>
      <td>
<pre><code>```latex
\LaTeX
```</code></pre>
      </td>
      <td>
<pre><code class="latex">\LaTeX
</code></pre>
      </td>
    </tr>
  </tbody>
</table>

[1]: http://golang.org
[2]: https://github.com/xiam/luminos
[3]: https://menteslibres.net/luminos
[4]: https://menteslibres.net/xiam/donate
[5]: https://menteslibres.net/luminos/getting-started
[6]: https://menteslibres.net/xiam
