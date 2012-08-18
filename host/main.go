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
	"regexp"
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

var extensions = []string{".md", ".html", ".txt"}

func (host *Host) url(url string) string {
	//return "/" + strings.TrimLeft(url, "/")
	return url
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

func guessFile(file string, descend bool) (string, os.FileInfo) {
	stat, err := os.Stat(file)

	file = strings.TrimRight(file, PS)

	//fmt.Printf("%v\n", file)

	if descend == true {
		if err == nil {
			if stat.IsDir() {
				f, s := guessFile(file+PS+"index", true)
				if s != nil {
					return f, s
				}
			}
			return file, stat
		} else {
			for _, extension := range extensions {
				f, s := guessFile(file+extension, false)
				if s != nil {
					return f, s
				}
			}
		}
	} else {
		if err == nil {
			return file, stat
		}
	}

	return "", nil
}

func (host *Host) readFile(file string) ([]byte, error) {
	stat, err := os.Stat(file)

	if err != nil {
		return nil, err
	}

	if stat.IsDir() == false {

		fp, _ := os.Open(file)
		defer fp.Close()

		buf := make([]byte, stat.Size())

		_, err := fp.Read(buf)

		if err != nil {
			return nil, err
		}

		if strings.HasSuffix(file, ".md") {
			return md.MarkdownCommon(buf), nil
		} else {
			return buf, nil
		}

	}

	return nil, nil
}

func (host *Host) ServeHTTP(w http.ResponseWriter, req *http.Request) {

	var localFile string

	// Updating settings and template files that have changed.
	host.Update()

	log.Printf("%s: Routing request %s %s", req.Host, req.Method, req.URL.Path,)

	// Requested path
	reqpath := strings.Trim(req.URL.Path, "/")

	// Trying to match a file on webroot/
	localFile = host.DocumentRoot + PS + host.Settings.Get("document.webroot", "webroot").(string) + PS + reqpath
	stat, err := os.Stat(localFile)

	if err == nil {
		// File exists
		if stat.IsDir() == false {
			// Exists and it's not a directory, let's serve it.
			log.Printf("%s: Serving file %s.", host.Name, localFile)
			http.ServeFile(w, req, localFile)
			return
		}
	}

	tryName := host.DocumentRoot + PS + strings.TrimRight(host.Settings.Get("document.htdocs", "htdocs").(string), PS) + PS + reqpath

	stat, err = os.Stat(tryName)

	localFile, stat = guessFile(tryName, true)

	if stat != nil {

		if reqpath != "" {
			if stat.IsDir() == false {
				if strings.HasSuffix(req.URL.Path, "/") == true {
					http.Redirect(w, req, "/"+reqpath, 301)
				}
			} else {
				if strings.HasSuffix(req.URL.Path, "/") == false {
					http.Redirect(w, req, req.URL.Path+"/", 301)
				}
			}
		}

		p := &page.Page{}

		p.FilePath = localFile
		p.BasePath = req.URL.Path

		if stat.IsDir() == false {
			p.FileDir = path.Dir(localFile)
			p.BasePath = path.Dir(req.URL.Path)
		} else {
			p.FileDir = p.FilePath
		}

		content, err := host.readFile(localFile)

		if err == nil {
			p.Content = template.HTML(content)
		}

		p.BasePath = strings.TrimRight(p.BasePath, "/") + "/"
		p.FileDir = strings.TrimRight(p.FileDir, "/") + "/"

		// werc-like header and footer.
		hfile, hstat := guessFile(p.FileDir+"_header", true)

		if hstat != nil {
			hcontent, herr := host.readFile(hfile)
			if herr == nil {
				p.ContentHeader = template.HTML(hcontent)
			}
		}

		// werc-like header and footer.
		ffile, fstat := guessFile(p.FileDir+"_footer", true)

		if fstat != nil {
			fcontent, ferr := host.readFile(ffile)
			if ferr == nil {
				p.ContentFooter = template.HTML(fcontent)
			}
		}

		if p.Content != "" {
			title, _ := regexp.Compile(`<h[\d]>(.+)</h`)
			found := title.FindStringSubmatch(string(p.Content))
			if len(found) > 0 {
				p.PageTitle = found[1]
			}
		}

		p.CreateBreadCrumb()
		p.CreateMenu()
		p.CreateSideMenu()

		log.Printf("%s: Serving file %s.", host.Name, localFile)

		if err := host.Templates["index.tpl"].t.Execute(w, p); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}

		return
	}

	log.Printf("%s: Resource was not found.", host.Name)

	http.Error(w, "Not found", 404)
}

func (host *Host) loadTemplates() bool {
	dir := host.DocumentRoot + PS + host.Settings.Get("document.templates", "templates").(string)

	fp, err := os.Open(dir)

	if err != nil {
		log.Printf("%s: %s", host.Name, err)
		return false
	}

	defer fp.Close()

	files, err := fp.Readdir(-1)

	if err != nil {
		log.Printf("%s: %s", host.Name, err)
		return false
	}

	for _, file := range files {
		if strings.HasSuffix(file.Name(), ".tpl") == true {

			tpl := dir + PS + file.Name()

			if value, loaded := host.Templates[file.Name()]; loaded == true {
				if value.mtime == file.ModTime() {
					continue
				} else {
					log.Printf("%s: (Re)loading %s", host.Name, tpl)
				}
			} else {
				log.Printf("%s: Loading %s", host.Name, tpl)
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
				log.Printf("%s: Template error on file %s: %s", host.Name, tpl, err.Error())
			}

		}
	}

	if _, ok := host.Templates["index.tpl"]; ok == false {
		log.Printf("%s: template %s could not be found.", host.Name, dir+PS+"index.tpl")
		return false
	}

	return true

}

func (host *Host) loadSettings() bool {
	file := host.DocumentRoot + PS + "luminos.yaml"
	_, err := os.Stat(file)
	if err == nil {
		host.Settings = yaml.Open(file)
		return true
	} else {
		log.Printf("%s: seems like %s does not exists.\n", host.Name, file)
	}
	return false
}

func (host *Host) Update() bool {
	if host.loadSettings() == false {
		return false
	}
	if host.loadTemplates() == false {
		return false
	}
	return true
}

func New(req *http.Request, docroot string) (*Host, error) {
	host := &Host{}
	host.Name = req.Host

	_, err := os.Stat(docroot)

	if err == nil {
		host.DocumentRoot = docroot

		host.Templates = make(map[string]*Template)

		if host.Update() == true {
			return host, nil
		} else {
			return nil, errors.New(fmt.Sprintf("Could not start host: %s", host.Name))
		}
	} else {
		return nil, err
	}

	return nil, nil
}
