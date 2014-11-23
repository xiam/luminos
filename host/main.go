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

package host

import (
	"fmt"
	//"github.com/howeyc/fsnotify"
	md "github.com/russross/blackfriday"
	"html/template"
	"log"
	"menteslibres.net/gosexy/to"
	"menteslibres.net/gosexy/yaml"
	"menteslibres.net/luminos/page"
	"menteslibres.net/luminos/watcher"
	"net/http"
	"os"
	"path"
	"regexp"
	"strings"
	"time"
)

const (
	pathSeparator = string(os.PathSeparator)
	settingsFile  = "site.yaml"
)

var (
	isExternalLinkPattern = regexp.MustCompile(`^[a-zA-Z0-9]+:\/\/`)
)

// Host is the struct that represents virtual hosts.
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
	Templates map[string]*template.Template
	// Function map for templates.
	template.FuncMap
	// Standard request.
	*http.Request
	// Standard response writer.
	http.ResponseWriter
	// Function map
	funcMap template.FuncMap
	// File watcher
	//Watcher *fsnotify.Watcher
	Watcher *watcher.Watcher
	// Template root
	TemplateRoot string
}

// Expected extensions. Elements on the left have precedence.
var extensions = []string{
	".md",
	".html",
	".txt",
}

// Close removes the watcher that is currently associated with the host.
func (host *Host) Close() {
	host.Watcher.Close()
}

// asset returns a relative URL.
func (host *Host) asset(url string) string {
	if host.isExternalLink(url) == false {
		if host.Path == "" {
			return "/" + strings.TrimLeft(url, "/")
		}
		return "/" + host.Path + "/" + strings.TrimLeft(url, "/")
	}
	return url
}

// url returns an absolute URL.
func (host *Host) url(url string) string {
	if host.isExternalLink(url) == false {
		return "//" + host.Request.Host + "/" + strings.TrimLeft(url, "/")
	}
	return url
}

// isExternalLink returns true if the given URL is outside this host.
func (host *Host) isExternalLink(url string) bool {
	return isExternalLinkPattern.MatchString(url)
}

// setting function returns a setting value.
func (host *Host) setting(path string) interface{} {
	route := strings.Split(path, "/")
	args := make([]interface{}, len(route))
	for i := range route {
		args[i] = route[i]
	}
	return host.Settings.Get(args...)
}

// settings is a function that returns an array of settings.
func (host *Host) settings(path string) []interface{} {
	route := strings.Split(path, "/")
	args := make([]interface{}, len(route))
	for i := range route {
		args[i] = route[i]
	}
	val := host.Settings.Get(args...)
	if val == nil {
		return nil
	}
	return val.([]interface{})
}

// jstext is a function for funcMap that writes text as Javascript.
func jstext(text string) template.JS {
	return template.JS(text)
}

// htmltext is a function for funcMap that writes text as plain HTML.
func htmltext(text string) template.HTML {
	return template.HTML(text)
}

// links is a function for funcMap that writes links.
func (host *Host) link(url, text string) template.HTML {
	if host.isExternalLink(url) {
		return template.HTML(fmt.Sprintf(`<a target="_blank" href="%s">%s</a>`, host.asset(url), text))
	}
	return template.HTML(fmt.Sprintf(`<a href="%s">%s</a>`, host.asset(url), text))
}

// guessFile checks for files names and returns a guessed name.
func guessFile(file string, descend bool) (string, os.FileInfo) {
	stat, err := os.Stat(file)

	file = strings.TrimRight(file, pathSeparator)

	if descend {
		if err == nil {
			if stat.IsDir() {
				f, s := guessFile(file+pathSeparator+"index", true)
				if s != nil {
					return f, s
				}
			}
			return file, stat
		}
		for _, extension := range extensions {
			f, s := guessFile(file+extension, false)
			if s != nil {
				return f, s
			}
		}
	}

	if err == nil {
		return file, stat
	}

	return "", nil
}

// readFile opens a file and reads its contents, if the file has the .md
// extension the contents are parsed and HTML is returned.
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
		}
		return buf, nil

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

	// TODO: Fix this non-critical race condition.  We need to save some
	// variables in a per request basis, in particular the hostname. It may not
	// always match the host name we gave to it (i.e: the "default" hostname). A
	// per-request context would be useful.
	host.Request = req

	// Default status.
	status := http.StatusNotFound
	size := -1

	// Requested path
	reqpath := strings.Trim(req.URL.Path, "/")

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

	webroot := host.DocumentRoot + pathSeparator + webrootdir

	localFile = webroot + pathSeparator + reqpath

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

		docroot := host.DocumentRoot + pathSeparator + docrootdir

		testFile := docroot + pathSeparator + reqpath

		stat, err = os.Stat(testFile)

		localFile, stat = guessFile(testFile, true)

		if stat != nil {

			if reqpath != "" {
				if stat.IsDir() == false {
					if strings.HasSuffix(req.URL.Path, "/") == true {
						http.Redirect(w, req, "/"+host.Path+"/"+reqpath, 301)
						w.Write([]byte(http.StatusText(301)))
						return
					}
				} else {
					if strings.HasSuffix(req.URL.Path, "/") == false {
						http.Redirect(w, req, req.URL.Path+"/", 301)
						w.Write([]byte(http.StatusText(301)))
						return
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

			p.FileDir = strings.TrimRight(p.FileDir, pathSeparator) + pathSeparator
			p.BasePath = strings.TrimRight(p.BasePath, pathSeparator) + pathSeparator

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

			err = host.Templates["index.tpl"].Execute(w, p)

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

func (host *Host) loadTemplate(file string) error {

	name := path.Base(file)

	parsed := template.New(name)

	parsed, err := parsed.Funcs(host.funcMap).ParseFiles(file)

	if err != nil {
		return err
	}

	host.Templates[name] = parsed

	if host.Watcher != nil {
		host.Watcher.RemoveWatch(file)
		host.Watcher.Watch(file)
	}

	return nil
}

// loadTemplates loads templates with .tpl extension from the templates
// directory. At this moment only index.tpl is expected.
func (host *Host) loadTemplates() error {

	tpldir := to.String(host.Settings.Get("document", "templates"))

	if tpldir == "" {
		tpldir = "templates"
	}

	tplroot := host.DocumentRoot + pathSeparator + tpldir

	fp, err := os.Open(tplroot)

	if err != nil {
		return fmt.Errorf("Error trying to open %s: %s", tplroot, err.Error())
	}

	host.TemplateRoot = tplroot

	defer fp.Close()

	files, err := fp.Readdir(-1)

	if err != nil {
		return fmt.Errorf("Error reading directory %s: %s", tplroot, err.Error())
	}

	for _, fp := range files {

		if strings.HasSuffix(fp.Name(), ".tpl") == true {

			file := host.TemplateRoot + pathSeparator + fp.Name()

			err := host.loadTemplate(file)

			if err != nil {
				log.Printf("%s: Template error in file %s: %s\n", host.Name, file, err.Error())
			}

		}
	}

	if _, ok := host.Templates["index.tpl"]; ok == false {
		return fmt.Errorf("Template %s could not be found.", "index.tpl")
	}

	return nil

}

func (host *Host) fileWatcher() error {

	var err error

	/*
		// File watcher.
		host.Watcher, err = fsnotify.NewWatcher()

		if err == nil {

			go func() {

				for {

					select {

					case ev := <-host.Watcher.Event:

						fmt.Printf("%s: got ev: %v\n", host.Name, ev)

						if ev == nil {
							return
						}

						if ev.IsModify() {
							// Is settings file?
							if ev.Name == host.DocumentRoot+pathSeparator+settingsFile {
								log.Printf("%s: Reloading host settings %s...\n", host.Name, ev.Name)
								err := host.loadSettings()

								if err != nil {
									log.Printf("%s: Could not reload host settings: %s\n", host.Name, host.DocumentRoot+pathSeparator+settingsFile)
								}
							}

							// Is a template?
							if strings.HasPrefix(ev.Name, host.TemplateRoot) == true {

								if strings.HasSuffix(ev.Name, ".tpl") == true {
									log.Printf("%s: Reloading template %s", host.Name, ev.Name)
									host.loadTemplate(ev.Name)

									if err != nil {
										log.Printf("%s: Could not reload template %s: %s", host.Name, ev.Name, err.Error())
									}

								}
							}

						} else if ev.IsDelete() {
							// Attemping to re-add watcher.
							host.Watcher.RemoveWatch(ev.Name)
							host.Watcher.Watch(ev.Name)
						}

					}
				}

			}()

		}
	*/

	// (Stupid) file modification watcher.
	host.Watcher, err = watcher.New()

	if err == nil {

		go func() {

			for {
				select {
				case ev := <-host.Watcher.Event:

					if ev.IsModify() {
						// Is settings file?
						if ev.Name == host.DocumentRoot+pathSeparator+settingsFile {
							log.Printf("%s: Reloading host settings %s...\n", host.Name, ev.Name)
							err := host.loadSettings()

							if err != nil {
								log.Printf("%s: Could not reload host settings: %s\n", host.Name, host.DocumentRoot+pathSeparator+settingsFile)
							}
						}

						// Is a template?
						if strings.HasPrefix(ev.Name, host.TemplateRoot) == true {
							if strings.HasSuffix(ev.Name, ".tpl") == true {
								log.Printf("%s: Reloading template %s", host.Name, ev.Name)
								host.loadTemplate(ev.Name)
								if err != nil {
									log.Printf("%s: Could not reload template %s: %s", host.Name, ev.Name, err.Error())
								}
							}
						}
					}
				}
			}
		}()
	}

	return err

}

// loadSettings loads settings for the host.
func (host *Host) loadSettings() error {

	var settings *yaml.Yaml

	file := host.DocumentRoot + pathSeparator + settingsFile

	_, err := os.Stat(file)

	if err == nil {
		settings, err = yaml.Open(file)
		if err != nil {
			return fmt.Errorf(`Could not parse settings file (%s): %s`, file, err.Error())
		}
	} else {
		return fmt.Errorf(`Error trying to open settings file (%s): %s.`, file, err.Error())
	}

	if host.Watcher != nil {
		host.Watcher.RemoveWatch(file)
		host.Watcher.Watch(file)
	}

	host.Settings = settings

	return nil
}

// New creates and returns a host.
func New(name string, root string) (*Host, error) {

	_, err := os.Stat(root)

	if err == nil {

		name = strings.Trim(name, "/")

		route := "/"

		index := strings.Index(name, "/")

		if index > 0 {
			route = name[index:]
		}

		host := &Host{
			Name:         strings.Trim(name, "/"),
			Path:         strings.Trim(route, "/"),
			DocumentRoot: root,
			Templates:    make(map[string]*template.Template),
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

		// Watcher
		host.fileWatcher()

		// Loading host settings
		err := host.loadSettings()

		if err != nil {
			log.Printf("Could not start host: %s\n", name)
			return nil, err
		}

		// Loading templates.
		err = host.loadTemplates()

		if err != nil {
			log.Printf("Could not start host: %s\n", name)
			return nil, err
		}

		log.Printf("Routing: %s -> %s\n", name, root)

		return host, nil

	}

	log.Printf("Error reading directory %s: %s\n", root, err.Error())
	log.Printf("Checkout an example directory at https://github.com/xiam/luminos/tree/master/default\n")

	return nil, err
}
