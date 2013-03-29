/*
  Copyright (c) 2012-2013 José Carlos Nieto, http://xiam.menteslibres.org/

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
	// Main path
	Path string
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
	// Function map
	funcMap template.FuncMap
}

// A template
type Template struct {
	path  string
	t     *template.Template
	mtime time.Time
}

var extensions = []string{
	".md",
	".html",
	".txt",
}

// Returns a relative URL.
func (self *Host) asset(url string) string {
	if self.isExternalLink(url) == false {
		return self.Path + "/" + strings.TrimLeft(url, "/")
	}
	return url
}

// Returns an absolute URL.
func (self *Host) url(url string) string {
	if self.isExternalLink(url) == false {
		return self.Name + self.asset(url)
	}
	return url
}

func (self *Host) isExternalLink(url string) bool {
	test, _ := regexp.Compile(`^[a-z0-9]+:\/\/`)
	return test.MatchString(url)
}

// Function for funcMap that returns a setting value.
func (self *Host) setting(path string) interface{} {
	route := strings.Split(path, "/")
	args := make([]interface{}, len(route))
	for i, _ := range route {
		args[i] = route[i]
	}
	return self.Settings.Get(args...)
}

// Function for funcMap that returns an array of settings.
func (self *Host) settings(path string) []interface{} {
	route := strings.Split(path, "/")
	args := make([]interface{}, len(route))
	for i, _ := range route {
		args[i] = route[i]
	}
	val := self.Settings.Get(args...)
	if val == nil {
		return nil
	}
	return val.([]interface{})
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

func chunk(value string) string {
	if value == "" {
		return "-"
	}
	return value
}

// A simple ServeHTTP.
func (host *Host) ServeHTTP(w http.ResponseWriter, req *http.Request) {

	var localFile string

	// Default status.
	status := http.StatusNotFound
	size := -1

	// Updating settings and template files that have changed.
	host.Update()

	// Requested path
	reqpath := req.URL.Path

	// Stripping path
	index := len(host.Path)

	if reqpath[0:index] == host.Path {
		reqpath = reqpath[index:]
	}

	reqpath = strings.Trim(reqpath, "/")

	// Trying to match a file on webroot/
	webrootdir := to.String(host.Settings.Get("document", "webroot"))

	if webrootdir == "" {
		webrootdir = "webroot"
	}

	webroot := host.DocumentRoot + PS + webrootdir

	localFile = webroot + PS + reqpath

	stat, err := os.Stat(localFile)

	if err == nil {
		// File exists
		if stat.IsDir() == false {
			// Exists and it's not a directory, let's serve it.
			status = http.StatusOK
			http.ServeFile(w, req, localFile)
			size = int(stat.Size())
		}
	}

	if status == http.StatusNotFound {

		docrootdir := to.String(host.Settings.Get("document", "markdown"))

		if docrootdir == "" {
			docrootdir = "markdown"
		}

		docroot := host.DocumentRoot + PS + docrootdir

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

			err = host.Templates["index.tpl"].t.Execute(w, p)

			if err == nil {
				status = http.StatusOK
			} else {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				status = http.StatusInternalServerError
			}

		}
	}

	if status == http.StatusNotFound {
		http.Error(w, "Not found", http.StatusNotFound)
	}

	fmt.Println(strings.Join([]string{
		chunk(req.RemoteAddr),
		chunk(""),
		chunk(""),
		chunk("[" + time.Now().Format("02/Jan/2006:15:04:05 -0700") + "]"),
		chunk("\"" + fmt.Sprintf("%s %s %s", req.Method, req.RequestURI, req.Proto) + "\""),
		chunk(fmt.Sprintf("%d", status)),
		chunk(fmt.Sprintf("%d", size)),
	},
		" "),
	)
}

// Loads templates with .tpl extension from the templates directory. At this
// moment only index.tpl is expected.
func (host *Host) loadTemplates() error {
	tpldir := to.String(host.Settings.Get("document", "templates"))

	if tpldir == "" {
		tpldir = "templates"
	}

	tplroot := host.DocumentRoot + PS + tpldir

	fp, err := os.Open(tplroot)

	if err != nil {
		return fmt.Errorf("Error trying to open %s: %s", tplroot, err.Error())
	}

	defer fp.Close()

	files, err := fp.Readdir(-1)

	if err != nil {
		return fmt.Errorf("Error reading directory %s: %s", tplroot, err.Error())
	}

	for _, file := range files {
		if strings.HasSuffix(file.Name(), ".tpl") == true {

			tpl := tplroot + PS + file.Name()

			if value, loaded := host.Templates[file.Name()]; loaded == true {
				if value.mtime == file.ModTime() {
					continue
				} else {
					log.Printf("%s: Reloading template %s...\n", host.Name, tpl)
				}
			} else {
				log.Printf("%s: Loading template %s...\n", host.Name, tpl)
			}

			parsed := template.New(file.Name())

			parsed, err = parsed.Funcs(host.funcMap).ParseFiles(tpl)

			if err == nil {
				host.Templates[file.Name()] = &Template{
					path:  tpl,
					t:     parsed,
					mtime: file.ModTime(),
				}
			} else {
				log.Printf("%s: Template error on file %s: %s\n", host.Name, tpl, err.Error())
			}

		}
	}

	if _, ok := host.Templates["index.tpl"]; ok == false {
		return fmt.Errorf("Template %s could not be found.", tplroot+PS+"index.tpl")
	}

	return nil

}

// Loads settings into Host.Settings
func (host *Host) loadSettings() error {

	file := host.DocumentRoot + PS + "site.yaml"

	_, err := os.Stat(file)

	if err == nil {
		host.Settings, err = yaml.Open(file)
		log.Printf("%s: Loading settings file %s...\n", host.Name, file)
		if err != nil {
			return fmt.Errorf(`Could not parse settings file (%s): %s`, file, err.Error())
		}
	} else {
		return fmt.Errorf(`Error trying to open settings file (%s): %s.`, file, err.Error())
	}

	return nil
}

// Reloads templates and settings.
func (host *Host) Update() error {
	var err error

	err = host.loadSettings()
	if err != nil {
		return err
	}

	err = host.loadTemplates()
	if err != nil {
		return err
	}

	return nil
}

// Creates and returns a host.
func New(name string, root string) (*Host, error) {

	_, err := os.Stat(root)

	if err == nil {

		name = strings.Trim(name, "/")

		path := "/"

		index := strings.Index(name, "/")

		if index > 0 {
			path = name[index:]
			name = name[0:index]
		}

		host := &Host{
			Name:         name,
			Path:         path,
			DocumentRoot: root,
			Templates:    make(map[string]*Template),
		}

		host.funcMap = template.FuncMap{
			"url":      func(s string) string { return host.url(s) },
			"asset":    func(s string) string { return host.asset(s) },
			"setting":  func(s string) interface{} { return host.setting(s) },
			"settings": func(s string) []interface{} { return host.settings(s) },
			"jstext":   jstext,
			"htmltext": htmltext,
			"link":     func(a, b string) template.HTML { return host.link(a, b) },
		}

		err = host.Update()

		if err != nil {
			log.Printf("Could not start host: %s\n", name)
			return nil, err
		}

		return host, nil

	} else {
		log.Printf("Error reading directory %s: %s\n", root, err.Error())
		log.Printf("Checkout an example directory at https://github.com/xiam/luminos/tree/master/default\n")
		return nil, err
	}

	return nil, nil
}
