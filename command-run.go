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
	"flag"
	"fmt"
	//"github.com/howeyc/fsnotify"
	"log"
	"menteslibres.net/gosexy/cli"
	"menteslibres.net/gosexy/to"
	"menteslibres.net/gosexy/yaml"
	"net"
	"net/http"
	"net/http/fcgi"
	"os"
)

// Default values
const (
	envSettingsFile   = "./settings.yaml"
	envServerDomain   = "unix"
	envServerProtocol = "tcp"
)

// Global software settings.
var settings *yaml.Yaml

// Command line settings.
var flagSettings = flag.String("c", envSettingsFile, "Path to the settings.yaml file.")

// runCommand is the structure that provides instructions for the "luminos
// run" subcommand.
type runCommand struct {
}

// Execute runs a luminos server using a settings file.
func (c *runCommand) Execute() (err error) {
	var stat os.FileInfo

	// If no settings file was specified, use the default.
	if *flagSettings == "" {
		*flagSettings = envSettingsFile
	}

	// Attempt to stat the settings file.
	stat, err = os.Stat(*flagSettings)

	// It must not return an error.
	if err != nil {
		return fmt.Errorf("Error while opening %s: %q", *flagSettings, err)
	}

	// We must have a value in stat.
	if stat == nil {
		return fmt.Errorf("Could not load settings file: %s.", *flagSettings)
	}

	// And the file must not be a directory.
	if stat.IsDir() {
		return fmt.Errorf("Could not open %s: it's a directory!", *flagSettings)
	}

	// Now that we're positively sure that we have a valid file, let's try to
	// read settings from it.
	if settings, err = loadSettings(*flagSettings); err != nil {
		return fmt.Errorf("Error while reading settings file %s: %q", *flagSettings, err)
	}

	// Starting settings watcher.
	if err = settingsWatcher(); err == nil {
		watch.Watch(*flagSettings)
	}

	// Reading setttings.
	serverType := to.String(settings.Get("server", "type"))

	domain := envServerDomain
	address := to.String(settings.Get("server", "socket"))

	if address == "" {
		domain = envServerProtocol
		address = fmt.Sprintf("%s:%d", to.String(settings.Get("server", "bind")), to.Int64(settings.Get("server", "port")))
	}

	// Creating a network listener.
	var listener net.Listener

	if listener, err = net.Listen(domain, address); err != nil {
		return fmt.Errorf("Could not create network listener: %q", err)
	}

	// Listener must be closed when the function exits.
	defer listener.Close()

	// Attempt to start a server.
	switch serverType {
	case "fastcgi":
		if err == nil {
			log.Printf("Starting FastCGI server. Listening at %s.", address)
			fcgi.Serve(listener, &server{})
		} else {
			return fmt.Errorf("Failed to start FastCGI server: %q", err)
		}
	case "standalone":
		if err == nil {
			log.Printf("Starting HTTP server. Listening at %s.", address)
			http.Serve(listener, &server{})
		} else {
			return fmt.Errorf("Failed to start HTTP server: %q", err)
		}
	default:
		return fmt.Errorf("Unknown server type: %s", serverType)
	}

	return nil
}

func init() {
	// Describing the "run" subcommand.

	cli.Register("run", cli.Entry{
		Name:        "run",
		Description: "Runs a luminos server.",
		Arguments:   []string{"c"},
		Command:     &runCommand{},
	})

}
