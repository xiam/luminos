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

package main

import (
	"flag"
	"fmt"
	"github.com/gosexy/to"
	"github.com/gosexy/cli"
	"github.com/gosexy/yaml"
	"github.com/xiam/luminos/host"
	"log"
	"net"
	"net/http"
	"net/http/fcgi"
	"os"
	"strings"
)

var version = "0.3"

type command struct {
	Help error
}

type Server struct {
}

// Command line flags.
var flagHelp = flag.Bool("help", false, "Shows command line hints.")
var flagSettings = flag.String("conf", "./settings.yaml", "Path to the settings.yaml file.")
var flagVersion = flag.Bool("version", false, "Shows software version.")

// Global settings.
var settings *yaml.Yaml

// Host map.
var hosts map[string]*host.Host

func help() {
	fmt.Printf("Usage: %s\n", os.Args[0])
	fmt.Println("")
	flag.PrintDefaults()
	fmt.Println("")
}

// Dispatches a request and returns the appropriate host.
func route(req *http.Request) *host.Host {

	name := req.Host

	if strings.Contains(name, ":") {
		name = name[0:strings.LastIndex(name, ":")]
	}

	if _, ok := hosts[name]; ok == false {

		var err error
		var docroot string

		docroot = to.String(settings.Get(fmt.Sprintf("hosts/%s", name)))

		if docroot == "" {
			// Trying to serve default host.
			docroot = to.String(settings.Get("hosts/default"))
			if docroot == "" {
				// Default host is not defined.
				return nil
			} else {
				// Serving default host.
				hosts[name], err = host.New(req, docroot)
				if err != nil {
					delete(hosts, name)
					log.Printf("Could not find default host.")
					return nil
				}
			}
		} else {
			// Host is defined in settings.yaml
			hosts[name], err = host.New(req, docroot)
			if err != nil {
				delete(hosts, name)
				log.Printf("Requested host %s does not exists.", name)
				return nil
			}
		}

	}
	return hosts[name]
}

// Routes a request and lets the host handle it.
func (server Server) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	h := route(req)
	if h != nil {
		h.ServeHTTP(w, req)
	} else {
		log.Printf("Failed to serve host %s.\n", req.Host)
	}
}

// Starts up Luminos.
func main() {
	var err error

	fmt.Println("Luminos Markdown Server http://luminos.menteslibres.org")
	fmt.Println("Copyright (c) 2012 José Carlos Nieto")
	fmt.Println("")

	flag.Parse()

	err = cli.Dispatch()

	if err != nil {
		fmt.Println(err.Error())
		return
	}

	if *flagHelp == true {
		help()
		return
	}

	if *flagVersion == true {
		fmt.Printf("%v version: %s\n", os.Args[0], version)
		return
	}

	// Settings file exists?

	_, err = os.Stat(*flagSettings)

	if err != nil {
		log.Printf("Missing settings file: %s\n", *flagSettings)
		log.Printf("Download a sample from https://raw.github.com/xiam/luminos/master/settings.yaml\n")
		fmt.Println("")
		help()
		return
	}

	hosts = make(map[string]*host.Host)

	settings, err = yaml.Open(*flagSettings)

	if err != nil {
		log.Fatalf("Failed to open settings file: %s", err.Error())
	}

	serverType := to.String(settings.Get("server/type"))

	domain := "unix"
	address := to.String(settings.Get("server/socket"))

	if address == "" {
		domain = "tcp"
		address = fmt.Sprintf("%s:%d", to.String(settings.Get("server/bind")), to.Int(settings.Get("server/port")))
	}

	listener, err := net.Listen(domain, address)

	if err != nil {
		log.Fatalf("Failed to bind on %s.", address)
	}

	defer listener.Close()

	switch serverType {
	case "fastcgi":
		if err == nil {
			log.Printf("FastCGI server listening at %s.", address)
			fcgi.Serve(listener, &Server{})
		} else {
			log.Fatalf("Failed to start FastCGI server.")
		}
	case "standalone":
		if err == nil {
			log.Printf("HTTP server listening at %s.", address)
			http.Serve(listener, &Server{})
		} else {
			log.Fatalf("Failed to start HTTP server.")
		}
	default:
		log.Fatalf("Unknown server type %s.", serverType)
	}

	log.Printf("Exiting...")
}
