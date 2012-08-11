package main

import (
	"net"
	"net/http"
	"net/http/fcgi"
	"os"
	"path"
	"log"
	"fmt"
	"strings"
	"time"
	"html/template"
	"github.com/xiam/gosexy/yaml"
	md "github.com/russross/blackfriday"
)

type Template struct {
	path string
	t *template.Template
	mtime time.Time
}

type Server struct {

}

type Context struct {
	*http.Request
}

type Page struct {
	Title string
	Header template.HTML
	Footer template.HTML
	Sidebar template.HTML
	Content template.HTML
	MenuContents []map[string]interface{}
	SideMenu []map[string]interface{}
	BreadCrumb []map[string]interface{}
	FilePath string
	FileDir string
	BasePath string
}

const (
	PS = string(os.PathSeparator)
)

var ctx = &Context{}

var settings = yaml.Open("settings.yaml")

var extensions = []string{ ".html", ".md", "" }

var templates = map[string] *Template {}

var templateFuncMap = template.FuncMap {
	"url": templateURL,
}

func templateURL(url string) string {
	url = strings.TrimLeft(url, "/")
	return "/" + url
}

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

		fileExt := path.Ext(item["Text"].(string))

		for _, ext := range extensions {
			if fileExt == ext {
				item["Text"] = item["Text"].(string)[:len(item["Text"].(string)) - len(ext)]
				break
			}
		}

	}
	return item
}

func (p *Page) CreateMenu() {
	var item map[string]interface{}
	p.MenuContents = []map[string]interface{}{}

	files := filterList(p.FileDir, directoryFilter)

	for _, file := range files {
		item = p.CreateLink(file, p.BasePath)
		children := filterList(p.FileDir + PS + file.Name(), directoryFilter)
		if len(children) > 0 {
			item["Children"] = []map[string]interface{}{}
			for _, child := range children {
				childItem := p.CreateLink(child, p.BasePath + file.Name() + "/")
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

func loadTemplates() {
	tplroot := settings.GetString("document.templates")

	fp, err := os.Open(tplroot)
	defer fp.Close()

	if err != nil {
		panic(err)
	}

	files, err := fp.Readdir(-1)

	if err != nil {
		panic(err)
	}

	for _, file := range files {
		if strings.HasSuffix(file.Name(), ".tpl") == true {

			tplpath := tplroot + PS + file.Name()

			if value, loaded := templates[file.Name()]; loaded == true {
				if value.mtime == file.ModTime() {
					continue
				} else {
					log.Printf("(Re)loading %s", tplpath)
				}
			} else {
				log.Printf("Loading %s", tplpath)
			}

			parsed := template.New(file.Name())

			parsed, err = parsed.Funcs(templateFuncMap).ParseFiles(tplpath)

			if err == nil {
				templates[file.Name()] = &Template{
					path: tplpath,
					t: parsed,
					mtime: file.ModTime(),
				}
			} else {
				fmt.Printf("Template error on file %s: %s", tplpath, err.Error())
			}

		}
	}

	if _, ok := templates["index.tpl"]; ok == false {
		fmt.Printf("Template %s could not be found.", tplroot + PS + "index.tpl")
		panic("Could not start without index.tpl.")
	}

}

func (f Server) ServeHTTP(w http.ResponseWriter, req *http.Request) {

	// copying request
	ctx.Request = req

	loadTemplates()

	// Requested path
	reqpath := strings.Trim(req.URL.Path, "/")

	log.Printf("%s %s %s", req.Host, req.Method, reqpath)

	// Trying to match a file on webroot/
	localFile := settings.GetString("document.webroot") + PS + reqpath
	stat, err := os.Stat(localFile)

	if err == nil {
		// File exists
		if stat.IsDir() == false {
			// Exists and it's not a directory, let's just serve it.
			log.Printf("-> Serving file %s.", localFile)
			http.ServeFile(w, req, localFile)
			return
		}
	}

	// File does not exists or it's a directory, trying to find a page that fits.
	if reqpath == "" {
		reqpath = "index"
	}

	localPage := settings.GetString("document.pages") + PS + reqpath

	testFile := ""

	for _, extension := range extensions {

		testFile = localPage + extension

		stat, err := os.Stat(testFile)

		if err == nil {

			contentType := settings.GetString("http.default.content_type")

			page := &Page{}

			fp, _ := os.Open(testFile)
			defer fp.Close()

			page.FilePath = localPage
			page.BasePath = req.URL.Path

			if stat.IsDir() == false {

				page.FileDir	= path.Dir(page.FilePath)

				buf := make([]byte, stat.Size())

				_, err := fp.Read(buf)

				if err != nil {
					panic(err)
				}

				switch extension {
				case ".html":
					contentType = "text/html"
					page.Content = template.HTML(buf)
				case ".md":
					contentType = "text/html"
					page.Content = template.HTML(md.MarkdownCommon(buf))
				default:
					fmt.Errorf("Unhandled extension: %v\n", extension)
				}

			} else {

				contentType = "text/html"
				page.FileDir	= page.FilePath

			}

			page.CreateMenu()
			page.CreateSideMenu()

			log.Printf("-> Serving file %s.", fp.Name())

			w.Header().Set("Content-Type", contentType)

			if err := templates["index.tpl"].t.Execute(w, page); err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
			}

			return
		}
	}

	log.Printf("-> Not found.")

	http.Error(w, "Not found", 404)
}

func main() {

	loadTemplates()

	dtype := settings.GetString("server.type")

	switch dtype {
	case "fcgi":
		address := fmt.Sprintf("%s:%d", settings.GetString("server.bind"), settings.GetInt("server.port"))
		listener, err := net.Listen("tcp", address)
		defer listener.Close()

		if err == nil {
			log.Printf("FCGI server listening at %s.", address)
			fcgi.Serve(listener, &Server{})
		} else {
			log.Printf("Failed to start FCGI server.")
			panic(err)
		}
	default:
		log.Printf("Unknown server type %s.", dtype)
	}
}
