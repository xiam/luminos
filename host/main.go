package host

import (
	"errors"
	"fmt"
	md "github.com/russross/blackfriday"
	"github.com/xiam/gosexy/yaml"
	"github.com/xiam/luminos/page"
	"html/template"
	"log"
	"net/http"
	"os"
	"path"
	"strings"
	"time"
)

const (
	PS = string(os.PathSeparator)
)

type Host struct {
	Name         string
	DocumentRoot string
	Settings     *yaml.Yaml
	Templates    map[string]*Template
	template.FuncMap
	*http.Request
	http.ResponseWriter
}

type Template struct {
	path  string
	t     *template.Template
	mtime time.Time
}

var extensions = []string{".html", ".md", ""}

func (host *Host) url(url string) string {
	return "/" + strings.TrimLeft(url, "/")
}

func (host *Host) setting(path string) interface{} {
	return host.Settings.Get(path, nil)
}

func (host *Host) settings(path string) []interface{} {
	val := host.Settings.Get(path, nil).([]interface{})
	return val
}

func jstext(text string) template.JS {
	return template.JS(text)
}

func htmltext(text string) template.HTML {
	return template.HTML(text)
}

func (host *Host) link(url, text string) template.HTML {
	return template.HTML(fmt.Sprintf("<a href=\"%s\">%s</a>", host.url(url), text))
}

func (host *Host) ServeHTTP(w http.ResponseWriter, req *http.Request) {

	host.Update()

	// Requested path
	reqpath := strings.Trim(req.URL.Path, "/")

	log.Printf("%s %s %s", req.Host, req.Method, reqpath)

	// Trying to match a file on webroot/
	localFile := host.DocumentRoot + PS + host.Settings.Get("document.webroot", "webroot").(string) + PS + reqpath
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

	localPage := host.DocumentRoot + PS + host.Settings.Get("document.pages", "pages").(string) + PS + reqpath

	testFile := ""

	for _, extension := range extensions {

		testFile = localPage + extension

		stat, err := os.Stat(testFile)

		if err == nil {

			contentType := host.Settings.GetString("http.default.content_type")

			p := &page.Page{}

			fp, _ := os.Open(testFile)
			defer fp.Close()

			p.FilePath = localPage
			p.BasePath = req.URL.Path

			if stat.IsDir() == false {

				p.FileDir = path.Dir(p.FilePath)

				buf := make([]byte, stat.Size())

				_, err := fp.Read(buf)

				if err != nil {
					panic(err)
				}

				switch extension {
				case ".html":
					contentType = "text/html"
					p.Content = template.HTML(buf)
				case ".md":
					contentType = "text/html"
					p.Content = template.HTML(md.MarkdownCommon(buf))
				default:
					fmt.Errorf("Unhandled extension: %v\n", extension)
				}

			} else {

				contentType = "text/html"
				p.FileDir = p.FilePath

			}

			p.CreateMenu()
			p.CreateSideMenu()

			log.Printf("-> Serving file %s.", fp.Name())

			w.Header().Set("Content-Type", contentType)

			if err := host.Templates["index.tpl"].t.Execute(w, p); err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
			}

			return
		}
	}

	log.Printf("-> Not found.")

	http.Error(w, "Not found", 404)
}

func (host *Host) loadTemplates() bool {
	dir := host.DocumentRoot + PS + host.Settings.Get("document.templates", "templates").(string)

	fp, err := os.Open(dir)
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

			tpl := dir + PS + file.Name()

			if value, loaded := host.Templates[file.Name()]; loaded == true {
				if value.mtime == file.ModTime() {
					continue
				} else {
					log.Printf("(Re)loading %s", tpl)
				}
			} else {
				log.Printf("Loading %s", tpl)
			}

			parsed := template.New(file.Name())

			parsed, err = parsed.Funcs(template.FuncMap{
				"url":      func(s string) string { return host.url(s) },
				"setting":  func(s string) interface{} { return host.setting(s) },
				"settings": func(s string) []interface{} { return host.settings(s) },
				"jstext":   jstext,
				"htmltext": htmltext,
				"link":     func(a, b string) template.HTML { return host.link(a, b) },
			}).ParseFiles(tpl)

			if err == nil {
				host.Templates[file.Name()] = &Template{
					path:  tpl,
					t:     parsed,
					mtime: file.ModTime(),
				}
			} else {
				log.Printf("Template error on file %s: %s", tpl, err.Error())
			}

		}
	}

	if _, ok := host.Templates["index.tpl"]; ok == false {
		log.Printf("Template %s could not be found.", dir+PS+"index.tpl")
		panic("Could not start without index.tpl.")
	}

	return true

}

func (host *Host) loadSettings() bool {
	file := host.DocumentRoot + PS + "settings.yaml"
	_, err := os.Stat(file)
	if err == nil {
		host.Settings = yaml.Open(file)
		return true
	}
	return false
}

func (host *Host) Update() bool {
	settings := host.loadSettings()
	templates := host.loadTemplates()
	return settings && templates
}

func New(req *http.Request) (*Host, error) {
	host := &Host{}
	host.Name = req.Header.Get("Host")

	// Great security risk.
	host.DocumentRoot = req.Header.Get("Document-Root")

	if host.DocumentRoot == "" {
		return nil, errors.New("Document-Root is null.")
	}

	host.Templates = make(map[string]*Template)

	if host.Update() == true {
		return host, nil
	} else {
		return nil, errors.New(fmt.Sprintf("Could not start host: %s", host.Name))
	}

	return nil, nil
}
