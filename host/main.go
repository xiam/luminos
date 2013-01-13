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
	"github.com/gosexy/to"
	"github.com/gosexy/yaml"
	md "github.com/russross/blackfriday"
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

// Virtual host that serves document of a given directory.
type Host struct {
	// Host name
	Name string
	// Main directory
	DocumentRoot string
	// Settings
	Settings *yaml.Yaml
	// Templates (not fully functional yet)
	Templates map[string]*Template
	// Function map for templates.
	template.FuncMap
	// Standard request.
	*http.Request
	// Standard response writer.
	http.ResponseWriter
}

// A template
type Template struct {
	path  string
	t     *template.Template
	mtime time.Time
}

var extensions = []string{".md", ".html", ".txt"}

// Function for funcMap that returns an absolute or relative URL.
func (host *Host) url(url string) string {
	if host.isExternalLink(url) == false {
		return "/" + strings.TrimLeft(url, "/")
	}
	return url
}

func (host *Host) isExternalLink(url string) bool {
	test, _ := regexp.Compile(`^[a-z0-9]+:\/\/`)
	return test.MatchString(url)
}

// Function for funcMap that returns a setting value.
func (host *Host) setting(path string) interface{} {
	return host.Settings.Get(path)
}

// Function for funcMap that returns an array of settings.
func (host *Host) settings(path string) []interface{} {
	val := host.Settings.Get(path).([]interface{})
	return val
}

// Function for funcMap that writes text as Javascript.
func jstext(text string) template.JS {
	return template.JS(text)
}

// Function for funcMap that writes text as plain HTML.
func htmltext(text string) template.HTML {
	return template.HTML(text)
}

// Function for funcMap that writes links.
func (host *Host) link(url, text string) template.HTML {
	if host.isExternalLink(url) {
		return template.HTML(fmt.Sprintf(`<a target="_blank" href="%s">%s</a>`, host.url(url), text))
	}
	return template.HTML(fmt.Sprintf(`<a href="%s">%s</a>`, host.url(url), text))
}

// Checks for files names and returns a guessed name.
func guessFile(file string, descend bool) (string, os.FileInfo) {
	stat, err := os.Stat(file)

	file = strings.TrimRight(file, PS)

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

// Reads a file, if the file has the .md extension the contents are parsed and HTML is returned.
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

// A simple ServeHTTP.
func (host *Host) ServeHTTP(w http.ResponseWriter, req *http.Request) {

	var localFile string

	// Updating settings and template files that have changed.
	host.Update()

	log.Printf("%s: Routing request %s %s", req.Host, req.Method, req.URL.Path)

	// Requested path
	reqpath := strings.Trim(req.URL.Path, "/")

	// Trying to match a file on webroot/
	var webroot string

	if host.Settings.Get("document.webroot") == nil {
		webroot = host.DocumentRoot + PS + "webroot"
	} else {
		webroot = host.DocumentRoot + PS + to.String(host.Settings.Get("document.webroot"))
	}

	localFile = webroot + PS + reqpath

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

	var docroot string

	if host.Settings.Get("document.markdown") == nil {
		docroot = host.DocumentRoot + PS + "markdown"
	} else {
		docroot = host.DocumentRoot + PS + to.String(host.Settings.Get("document.markdown"))
	}

	testFile := docroot + PS + reqpath

	stat, err = os.Stat(testFile)

	localFile, stat = guessFile(testFile, true)

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

		relPath := localFile[len(docroot):]

		if stat.IsDir() == false {
			p.FileDir = path.Dir(localFile)
			p.BasePath = path.Dir(relPath)
		} else {
			p.FileDir = localFile
			p.BasePath = relPath
		}

		content, err := host.readFile(localFile)

		if err == nil {
			p.Content = template.HTML(content)
		}

		p.FileDir = strings.TrimRight(p.FileDir, PS) + PS
		p.BasePath = strings.TrimRight(p.BasePath, PS) + PS

		// werc-like header and footer.
		hfile, hstat := guessFile(p.FileDir+"_header", true)

		if hstat != nil {
			hcontent, herr := host.readFile(hfile)
			if herr == nil {
				p.ContentHeader = template.HTML(hcontent)
			}
		}

		if p.BasePath == "/" {
			p.IsHome = true
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
				p.Title = found[1]
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

// Loads templates with .tpl extension from the templates directory. At this moment only index.tpl is expected.
func (host *Host) loadTemplates() bool {
	var dir string

	if host.Settings.Get("document.templates") == nil {
		dir = host.DocumentRoot + PS + "templates"
	} else {
		dir = host.DocumentRoot + PS + to.String(host.Settings.Get("document.templates"))
	}

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

// Loads settings into Host.Settings
func (host *Host) loadSettings() bool {
	file := host.DocumentRoot + PS + "site.yaml"
	_, err := os.Stat(file)
	if err == nil {
		host.Settings, err = yaml.Open(file)
		if err != nil {
			log.Fatalf("Could not open settings file: %s", err.Error())
		}
		return true
	} else {
		log.Printf("%s: %s\n", host.Name, err.Error())
		log.Printf("See http://luminos.menteslibres.org/getting-started/directory-structure")
	}
	return false
}

// Reloads templates and settings.
func (host *Host) Update() bool {
	if host.loadSettings() == false {
		return false
	}
	if host.loadTemplates() == false {
		return false
	}
	return true
}

// Creates and returns a host.
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
		log.Printf("Error reading directory %s: %s\n", docroot, err.Error())
		log.Printf("Checkout an example directory at https://github.com/xiam/luminos/tree/master/default\n")
		return nil, err
	}

	return nil, nil
}
