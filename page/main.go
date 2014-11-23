// Copyright (c) 2012-2014 José Carlos Nieto, https://menteslibres.net/xiam
//
// Permission is hereby granted, free of charge, to any person obtaining
// a copy of this software and associated documentation files (the
// "Software"), to deal in the Software without restriction, including
// without limitation the rights to use, copy, modify, merge, publish,
// distribute, sublicense, and/or sell copies of the Software, and to
// permit persons to whom the Software is furnished to do so, subject to
// the following conditions:
//
// The above copyright notice and this permission notice shall be
// included in all copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND,
// EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF
// MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND
// NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT HOLDERS BE
// LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER IN AN ACTION
// OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN CONNECTION
// WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.

package page

import (
	"html/template"
	"os"
	"path"
	"regexp"
	"sort"
	"strings"
)

// Page struct holds information on the current document being served.
type Page struct {

	// Title of the page. This is guessed from the current document. (It looks
	// for the first H1, H2, ..., H6 tag).
	Title string

	// The HTML source of the current document.
	Content template.HTML

	// The HTML source of the _header.md or _header.html file on the current
	// document's directory.
	ContentHeader template.HTML

	// The HTML source of the _footer.md or _footer.html file on the current
	// document's directory.
	ContentFooter template.HTML

	// An array of maps that contains names and links of all the items on the
	// document root.  Names that begin with "." or "_" are ignored in this list.
	Menu []map[string]interface{}

	// An array of maps that contains names and links of all the items on the
	// current document's directory.  Names that begin with "." or "_" are
	// ignored in this list.
	SideMenu []map[string]interface{}

	// An array of maps that contains names and links of the current document's
	// path.
	BreadCrumb []map[string]interface{}

	// A map that contains the name and link of the current page.
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

const (
	pathSeparator = string(os.PathSeparator)
)

// List of known extensions.
var knownExtensions = []string{".html", ".md", ""}

// fileList struct is a sorted list of files.
type fileList []os.FileInfo

func (f fileList) Len() int {
	return len(f)
}

func (f fileList) Less(i, j int) bool {
	return f[i].Name() < f[j].Name()
}

func (f fileList) Swap(i, j int) {
	f[i], f[j] = f[j], f[i]
}

// removeKnownExtension strips out known extensions from a given file name.
func removeKnownExtension(s string) string {
	fileExt := path.Ext(s)

	for _, ext := range knownExtensions {
		if ext != "" {
			if fileExt == ext {
				return s[:len(s)-len(ext)]
			}
		}
	}

	return s
}

// filterList returns files in a directory passed through a filter.
func filterList(directory string, filter func(os.FileInfo) bool) fileList {
	var list fileList
	var err error

	// Attempt to open directory.
	var fp *os.File
	if fp, err = os.Open(directory); err != nil {
		panic(err)
	}

	defer fp.Close()

	// Listing directory contents.
	var dirContents []os.FileInfo
	if dirContents, err = fp.Readdir(-1); err != nil {
		panic(err)
	}

	// Looping over directory contents.
	for _, file := range dirContents {
		if filter(file) == true {
			list = append(list, file)
		}
	}

	// Sorting file list.
	sort.Sort(struct{ fileList }{list})

	return list
}

// dummyFilter is a filter for filterList. Returns all files except for those
// that begin with "." or "_".
func dummyFilter(f os.FileInfo) bool {
	if !strings.HasPrefix(f.Name(), ".") && !strings.HasPrefix(f.Name(), "_") {
		return true
	}
	return false
}

// directoryFilter is a filter for filterList. Returns all directories except
// those that begin with "." or "_".
func directoryFilter(f os.FileInfo) bool {
	if !strings.HasPrefix(f.Name(), ".") && !strings.HasPrefix(f.Name(), "_") {
		return f.IsDir()
	}
	return false
}

// fileFilter is a filter for filterList. Returns all files except for those
// that begin with "." or "_".
func fileFilter(f os.FileInfo) bool {
	if !strings.HasPrefix(f.Name(), ".") && !strings.HasPrefix(f.Name(), "_") {
		return (f.IsDir() == false)
	}
	return false
}

// createTitle expects a filename and returns a stylized human title.
func createTitle(s string) string {
	s = removeKnownExtension(s)
	s = regexp.MustCompile("[-_]").ReplaceAllString(s, " ")
	return strings.Title(s[:1]) + s[1:]
}

// CreateLink returns a link to another page.
func (p *Page) CreateLink(file os.FileInfo, prefix string) map[string]interface{} {
	item := map[string]interface{}{}

	if file.IsDir() == true {
		item["link"] = prefix + file.Name()
	} else {
		item["link"] = prefix + removeKnownExtension(file.Name())
	}

	item["text"] = createTitle(file.Name())

	return item
}

// CreateMenu scans files and directories and builds a list of children links.
func (p *Page) CreateMenu() {
	var item map[string]interface{}
	p.Menu = []map[string]interface{}{}

	files := filterList(p.FileDir, directoryFilter)

	for _, file := range files {
		item = p.CreateLink(file, p.BasePath)
		children := filterList(p.FileDir+pathSeparator+file.Name(), directoryFilter)
		if len(children) > 0 {
			item["children"] = []map[string]interface{}{}
			for _, child := range children {
				childItem := p.CreateLink(child, p.BasePath+file.Name())
				item["children"] = append(item["children"].([]map[string]interface{}), childItem)
			}
		}
		p.Menu = append(p.Menu, item)
	}
}

// CreateBreadCrumb populates Page.BreadCrumb with links.
func (p *Page) CreateBreadCrumb() {

	p.BreadCrumb = []map[string]interface{}{
		map[string]interface{}{
			"link": "/",
			"text": "Home",
		},
	}

	chunks := strings.Split(strings.Trim(p.BasePath, "/"), "/")

	prefix := ""

	for _, chunk := range chunks {
		if chunk != "" {
			item := map[string]interface{}{}
			item["link"] = prefix + "/" + chunk
			item["text"] = createTitle(chunk)
			prefix = prefix + pathSeparator + chunk
			p.BreadCrumb = append(p.BreadCrumb, item)
			p.CurrentPage = item
		}
	}

}

// CreateSideMenu populates Page.SideMenu with files on the current document's
// directory.
func (p *Page) CreateSideMenu() {
	var item map[string]interface{}
	p.SideMenu = []map[string]interface{}{}

	files := filterList(p.FileDir, dummyFilter)

	for _, file := range files {
		item = p.CreateLink(file, p.BasePath)
		if strings.ToLower(item["text"].(string)) != "index" {
			p.SideMenu = append(p.SideMenu, item)
		}
	}
}
