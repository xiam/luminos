package page

import (
	"html/template"
	"os"
	"strings"
)

type Page struct {
	Title        string
	Header       template.HTML
	Footer       template.HTML
	Sidebar      template.HTML
	Content      template.HTML
	MenuContents []map[string]interface{}
	SideMenu     []map[string]interface{}
	BreadCrumb   []map[string]interface{}
	FilePath     string
	FileDir      string
	BasePath     string
}

const (
	PS = string(os.PathSeparator)
)

func filterList(directory string, filter func(os.FileInfo) bool) []os.FileInfo {
	var list []os.FileInfo

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

	return list
}

func directoryFilter(f os.FileInfo) bool {
	if strings.HasPrefix(f.Name(), ".") == false {
		return f.IsDir()
	}
	return false
}

func fileFilter(f os.FileInfo) bool {
	if strings.HasPrefix(f.Name(), ".") == false {
		return (f.IsDir() == false)
	}
	return false
}

func (p *Page) CreateLink(file os.FileInfo, prefix string) map[string]interface{} {
	item := map[string]interface{}{}

	if file.IsDir() == true {
		item["Link"] = prefix + file.Name() + "/"
		item["Text"] = file.Name()
	} else {
		item["Link"] = prefix + file.Name()
		item["Text"] = file.Name()

		/*
			fileExt := path.Ext(item["Text"].(string))

			for _, ext := range extensions {
				if fileExt == ext {
					item["Text"] = item["Text"].(string)[:len(item["Text"].(string)) - len(ext)]
					break
				}
			}
		*/

	}
	return item
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

	files := filterList(p.FileDir, fileFilter)

	for _, file := range files {
		item = p.CreateLink(file, p.BasePath)
		// Ignoring index.
		if item["Text"].(string) != "index" {
			p.SideMenu = append(p.SideMenu, item)
		}
	}
}
