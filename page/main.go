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
	"fmt"
	"html/template"
	"log"
	"os"
	"path"
	"regexp"
	"sort"
	"strconv"
	"strings"

	"github.com/extemporalgenome/slug"
)

var titlePattern = regexp.MustCompile(`<h([\d])>(.+)</h[\d]>`)

type anchor struct {
	Text     string
	URL      string
	children []anchor
}

var homeAnchor = anchor{Text: "Home", URL: "/"}

var (
	titleReplacePattern = regexp.MustCompile(`[-_]`)
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

	// Titles holds links to all page subtitles.
	Titles map[int][]anchor

	// An array that contains names and links of all the items on the document's
	// root. Names that begin with a dot or an underscore are ignored from the
	// listing.
	Menu []anchor

	// An array that contains names and links of all the items on the current
	// document's directory. Names that begin with a dot or an underscore are
	// ignored from the listing.
	SideMenu []anchor

	// An array of anchors that contain names and URLs of the current document's
	// path.
	BreadCrumb []anchor

	// Contains the name and URL of the current page.
	CurrentPage anchor

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
	s = titleReplacePattern.ReplaceAllString(s, " ")
	return strings.Title(s[:1]) + s[1:]
}

// CreateLink returns a link to another page.
func (p *Page) CreateLink(file os.FileInfo, prefix string) anchor {
	item := anchor{}

	if file.IsDir() == true {
		item.URL = prefix + file.Name()
	} else {
		item.URL = prefix + removeKnownExtension(file.Name())
	}

	item.URL = path.Clean(item.URL)

	item.Text = createTitle(file.Name())

	return item
}

// CreateMenu scans files and directories and builds a list of children links.
func (p *Page) CreateMenu() {
	var item anchor
	p.Menu = []anchor{}

	files := filterList(p.FileDir, directoryFilter)

	for _, file := range files {
		item = p.CreateLink(file, p.BasePath)
		children := filterList(p.FileDir+pathSeparator+file.Name(), directoryFilter)
		if len(children) > 0 {
			item.children = make([]anchor, 0, len(children))
			for _, child := range children {
				childItem := p.CreateLink(child, p.BasePath+file.Name())
				item.children = append(item.children, childItem)
			}
		}
		p.Menu = append(p.Menu, item)
	}
}

// CreateBreadCrumb populates Page.BreadCrumb with links.
func (p *Page) CreateBreadCrumb() {

	chunks := strings.Split(strings.Trim(p.BasePath, "/"), "/")

	p.BreadCrumb = make([]anchor, 0, len(chunks)+1)

	p.BreadCrumb = append(p.BreadCrumb, homeAnchor)

	prefix := ""

	for _, chunk := range chunks {
		if chunk != "" {

			item := anchor{
				URL:  prefix + "/" + chunk,
				Text: createTitle(chunk),
			}

			prefix = prefix + pathSeparator + chunk
			p.BreadCrumb = append(p.BreadCrumb, item)
		}
	}

	p.CurrentPage = p.BreadCrumb[len(p.BreadCrumb)-1]
}

// CreateSideMenu populates Page.SideMenu with files on the current document's
// directory.
func (p *Page) CreateSideMenu() {
	var item anchor

	files := filterList(p.FileDir, dummyFilter)

	p.SideMenu = make([]anchor, 0, len(files))

	for _, file := range files {
		item = p.CreateLink(file, p.BasePath)
		if strings.ToLower(item.Text) != "index" {
			p.SideMenu = append(p.SideMenu, item)
		}
	}

	if strings.Trim(p.BasePath, "/") == "" {
		return
	}

	if len(p.SideMenu) == 0 {

		// Attempt to index parent directory.
		files = filterList(p.FileDir+pathSeparator+"..", dummyFilter)

		for _, file := range files {
			item = p.CreateLink(file, p.BasePath+".."+pathSeparator)
			if strings.ToLower(item.Text) != "index" {
				p.SideMenu = append(p.SideMenu, item)
			}
		}

	}
}

func (p *Page) ProcessContent() {
	content := string(p.Content)
	titles := titlePattern.FindAllStringSubmatch(content, -1)
	for _, title := range titles {
		if p.Titles == nil {
			p.Titles = make(map[int][]anchor)
		}
		if len(title) == 3 {
			if level, _ := strconv.Atoi(title[1]); level > 0 {
				ll := level - 1
				text := title[2]

				id := slug.Slug(text)

				if id == "" {
					id = fmt.Sprintf("%05d", level)
				}

				if p.Titles[ll] == nil {
					p.Titles[ll] = []anchor{}
				}

				r := fmt.Sprintf(`<h%d><a href="#%s" name="%s">%s</a></h%d>`, level, id, id, text, level)
				p.Titles[ll] = append(p.Titles[ll], anchor{Text: text, URL: "#" + id})

				content = strings.Replace(content, title[0], r, 1)
			}
		}
	}
	p.Content = template.HTML(content)
}

func (p *Page) GetTitlesFromLevel(ll int) []anchor {
	if p.Titles == nil || p.Titles[ll] == nil {
		return []anchor{}
	}
	return p.Titles[ll]
}

func (p *Page) URLMatch(s string) bool {
	re, err := regexp.Compile(s)
	if err != nil {
		log.Printf("URLMatch: %q", err)
		return false
	}
	return re.MatchString(p.CurrentPage.URL)
}
