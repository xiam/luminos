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

package main

import (
	"errors"
	"fmt"
	//"github.com/howeyc/fsnotify"
	"log"
	"net/http"
	"os"
	"strings"

	"menteslibres.net/gosexy/to"
	"menteslibres.net/gosexy/yaml"
	"menteslibres.net/luminos/host"
	"menteslibres.net/luminos/watcher"
)

// Map of hosts.
var hosts map[string]*host.Host

// File watcher.
var watch *watcher.Watcher

type server struct {
}

func init() {
	// Allocating map.
	hosts = make(map[string]*host.Host)
}

// Finds the appropriate hosts for a request.
func route(req *http.Request) *host.Host {

	// Request's hostname.
	name := req.Host

	// Removing the port part of the host.
	if strings.Contains(name, ":") {
		name = name[0:strings.Index(name, ":")]
	}

	// Host and path.
	path := name + req.URL.Path

	// Searching for the host that best matches this request.
	match := ""

	for key := range hosts {
		lkey := len(key)
		if lkey >= len(match) && lkey <= len(path) {
			if path[0:lkey] == key {
				match = key
			}
		}
	}

	// No host matched, let's use the default host.
	if match == "" {
		log.Printf("Could not match any host: %s, falling back to the default.\n", req.Host)
		match = "default"
	}

	// Let's verify and return the host.

	if _, ok := hosts[match]; !ok {
		// Host was not found.
		log.Printf("Request for unknown host: %s\n", req.Host)
		return nil
	}

	return hosts[match]
}

// Routes a request and lets the host handle it.
func (s server) ServeHTTP(wri http.ResponseWriter, req *http.Request) {
	r := route(req)
	if r != nil {
		r.ServeHTTP(wri, req)
	} else {
		log.Printf("Failed to serve host %s.\n", req.Host)
		http.Error(wri, "Not found", http.StatusNotFound)
	}
}

// Loads settings
func loadSettings(file string) (*yaml.Yaml, error) {

	var entries map[interface{}]interface{}
	var ok bool

	// Trying to read settings from file.
	y, err := yaml.Open(file)

	if err != nil {
		return nil, err
	}

	// Loading and verifying host entries
	if entries, ok = y.Get("hosts").(map[interface{}]interface{}); ok == false {
		return nil, errors.New("Missing \"hosts\" entry.")
	}

	h := map[string]*host.Host{}

	// Populating host entries.
	for key := range entries {
		name := to.String(key)
		path := to.String(entries[name])

		info, err := os.Stat(path)
		if err != nil {
			return nil, fmt.Errorf("Failed to validate host %s: %q.", name, err)
		}
		if info.IsDir() == false {
			return nil, fmt.Errorf("Host %s does not point to a directory.", name)
		}

		h[name], err = host.New(name, path)

		if err != nil {
			return nil, fmt.Errorf("Failed to initialize host %s: %q.", name, err)
		}
	}

	for name := range hosts {
		hosts[name].Close()
	}

	hosts = h

	if _, ok := hosts["default"]; ok == false {
		log.Printf("Warning: default host was not provided.")
	}

	return y, nil
}

func settingsWatcher() error {

	var err error

	/*
		// Watching settings file for changes.
		// Was not properly returning events on OSX.
		// https://github.com/howeyc/fsnotify/issues/34

		watcher, err := fsnotify.NewWatcher()

		if err == nil {
			defer watcher.Close()

			go func() {
				for {
					select {
					case ev := <-watcher.Event:
						if ev == nil {
							return
						}
						if ev.IsModify() {
							log.Printf("Trying to reload settings file %s...\n", ev.Name)
							y, err := loadSettings(ev.Name)
							if err != nil {
								log.Printf("Error loading settings file %s: %q\n", ev.Name, err)
							} else {
								settings = y
							}
						} else if ev.IsDelete() {
							watcher.RemoveWatch(ev.Name)
							watcher.Watch(ev.Name)
						}
					case err := <-watcher.Error:
						log.Printf("Watcher error: %q\n", err)
					}
				}
			}()

			watcher.Watch(settingsFile)
		}
	*/

	// (Stupid) time based file modification watcher.
	watch, err = watcher.New()

	if err == nil {
		go func() {
			defer watch.Close()
			for {
				select {
				case ev := <-watch.Event:
					if ev.IsModify() {
						y, err := loadSettings(ev.Name)
						if err != nil {
							log.Printf("Error loading settings file %s: %q\n", ev.Name, err)
						} else {
							log.Printf("Reloading settings file %s.\n", ev.Name)
							settings = y
						}
					}
				}
			}
		}()

	}

	return err
}
