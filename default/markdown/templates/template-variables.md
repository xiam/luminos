# Template variables

These are the variables available on your template.

* ``{{ .Title }}`` Page title, guessed from the current document. (Looks for the first H1, H2, ..., H6 tag)
* ``{{ .Content }}`` HTML of the current document.
* ``{{ .ContentHeader }}`` A special element, the HTML content of the ``_header.md`` file (if any) that is on the same directory of the
current document
* ``{{ .ContentFooter }}`` A special element, the HTML content of the ``_footer.md`` file (if any) that is on the same directory of the
current document.
* ``{{ .SideMenu }}`` An array of maps that contains names and links of all the items on the current document's directory. File names
that begin with "." or "_" are ignored in this list.
* ``{{ .BreadCrumb }}`` An array of maps that contains names and links of the current document's path.
* ``{{ .CurrentPage }}`` A map that contains name and link of the current page.
* ``{{ .FilePath }}`` Absolute path of the current document.
* ``{{ .FileDir }}`` Absolute parent directory of the current document.
* ``{{ .BasePath }}`` Relative path of the current document.
* ``{{ .BaseDir }}`` Relative parent directory of the current document.
* ``{{ .IsHome }}`` True if the current document is / (home).


If you want to hack on those structures, edit ``pages/main.go`` according to your needs.

    type Page struct {

        // Page title, guessed from the current document. (Looks for the first H1, H2, ..., H6 tag)
        Title string

        // The HTML of the current document.
        Content template.HTML

        // The HTML of the _header.md or _header.html file on the current document's directory.
        ContentHeader template.HTML

        // The HTML of the _footer.md or _footer.html file on the current document's directory.
        ContentFooter template.HTML

        // An array of maps that contains names and links of all the items on the current document's directory.
        // Names begginning with "." or "_" are ignored in this list.
        SideMenu []map[string]interface{}

        // An array of maps that contains names and links of the current document's path.
        BreadCrumb []map[string]interface{}

        // A map that contains name and link of the current page.
        CurrentPage map[string]interface{}

        // Absolute path of the current document.
        FilePath string

        // Absolute parent directory of the current document.
        FileDir string

        // Relative path of the current document.
        BasePath string

        // Relative parent directory of the current document.
        BaseDir string

        // True if the current document is / (home).
        IsHome bool
    }

