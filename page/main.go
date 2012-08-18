/*
  Copyright (c) 2012 José Carlos Nieto, http://xiam.menteslibres.org/

  Permission is hereby granted, free of charge, to any person obtaining
  a copy of this software and associated documentation files (the
  "Software"), to deal in the Software without restriction, including
  without limitation the rights to use, copy, modify, merge, publish,
  distribute, sublicense, and/or sell copies of the Software, and to
  permit persons to whom the Software is furnished to do so, subject to
  the following conditions:

  The above copyright notice and this permission notice shall be
  included in all copies or substantial portions of the Software.

  THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND,
  EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF
  MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND
  NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT HOLDERS BE
  LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER IN AN ACTION
  OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN CONNECTION
  WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.
*/

package page

import (
	"html/template"
	"os"
	"path"
	"sort"
	"strings"
)

type Page struct {
	Title         string
	Header        template.HTML
	Footer        template.HTML
	Sidebar       template.HTML
	Content       template.HTML
	ContentHeader template.HTML
	ContentFooter template.HTML
	MenuContents  []map[string]interface{}
	SideMenu      []map[string]interface{}
	BreadCrumb    []map[string]interface{}
	CurrentPage   map[string]interface{}
	FilePath      string
	FileDir       string
	PageTitle     string
	BasePath      string
}

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

type byName struct{ fileList }

var extensions = []string{".html", ".md", ""}

const (
	PS = string(os.PathSeparator)
)

func removeKnownExtension(s string) string {
	fileExt := path.Ext(s)

	for _, ext := range extensions {
		if ext != "" {
			if fileExt == ext {
				return s[:len(s)-len(ext)]
			}
		}
	}

	return s
}

func filterList(directory string, filter func(os.FileInfo) bool) fileList {
	var list fileList

	fp, err := os.Open(directory)
	defer fp.Close()

	if err != nil {
		panic(err)
	}

	ls, err := fp.Readdir(-1)

	if err != nil {
		panic(err)
	}

	for _, file := range ls {
		if filter(file) == true {
			list = append(list, file)
		}
	}

	sort.Sort(byName{list})

	return list
}

func dummyFilter(f os.FileInfo) bool {
	if strings.HasPrefix(f.Name(), ".") == false && strings.HasPrefix(f.Name(), "_") == false {
		return true
	}
	return false
}

func directoryFilter(f os.FileInfo) bool {
	if strings.HasPrefix(f.Name(), ".") == false && strings.HasPrefix(f.Name(), "_") == false {
		return f.IsDir()
	}
	return false
}

func fileFilter(f os.FileInfo) bool {
	if strings.HasPrefix(f.Name(), ".") == false && strings.HasPrefix(f.Name(), "_") == false {
		return (f.IsDir() == false)
	}
	return false
}

func (p *Page) CreateLink(file os.FileInfo, prefix string) map[string]interface{} {
	item := map[string]interface{}{}

	if file.IsDir() == true {
		item["link"] = prefix + file.Name() + "/"
		item["text"] = file.Name()
	} else {
		item["link"] = prefix + file.Name()
		item["text"] = file.Name()
	}

	item["link"] = removeKnownExtension(item["link"].(string))
	item["text"] = removeKnownExtension(item["text"].(string))

	return item
}

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
			item["link"] = prefix + "/" + chunk + "/"
			item["text"] = strings.Title(chunk)
			prefix = prefix + PS + chunk
			p.BreadCrumb = append(p.BreadCrumb, item)
			p.CurrentPage = item
		}
	}

}

func (p *Page) CreateMenu() {
	var item map[string]interface{}
	p.MenuContents = []map[string]interface{}{}

	files := filterList(p.FileDir, directoryFilter)

	for _, file := range files {
		item = p.CreateLink(file, p.BasePath)
		children := filterList(p.FileDir+PS+file.Name(), directoryFilter)
		if len(children) > 0 {
			item["Children"] = []map[string]interface{}{}
			for _, child := range children {
				childItem := p.CreateLink(child, p.BasePath+file.Name()+"/")
				item["Children"] = append(item["Children"].([]map[string]interface{}), childItem)
			}
		}
		p.MenuContents = append(p.MenuContents, item)
	}
}

func (p *Page) CreateSideMenu() {
	var item map[string]interface{}
	p.SideMenu = []map[string]interface{}{}

	files := filterList(p.FileDir, dummyFilter)

	for _, file := range files {
		item = p.CreateLink(file, p.BasePath)
		// Ignoring index.
		if item["text"].(string) != "index" {
			item["text"] = strings.Title(item["text"].(string))
			p.SideMenu = append(p.SideMenu, item)
		}
	}
}
