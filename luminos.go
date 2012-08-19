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

package main

import (
	"fmt"
	"flag"
	"github.com/xiam/gosexy/yaml"
	"github.com/xiam/luminos/host"
	"log"
	"net"
	"os"
	"net/http"
	"net/http/fcgi"
)

var version = "0.0999"

type Server struct {
}

var flagHelp			= flag.Bool("help", false, "Shows command line hints.")
var flagSettings	= flag.String("conf", "./settings.yaml", "Path to the settings.yaml file.")
var flagVersion		= flag.Bool("version", false, "Shows software version.")

var settings *yaml.Yaml

var hosts map[string]*host.Host

func route(req *http.Request) *host.Host {
	if _, ok := hosts[req.Host]; ok == false {

		var err error
		var docroot string

		docroot = settings.GetString(fmt.Sprintf("hosts/%s", req.Host))

		if docroot == "" {
			// Trying to serve default host.
			docroot = settings.GetString("hosts/default")
			if docroot == "" {
				// Default host is not defined.
				return nil
			} else {
				// Serving default host.
				hosts[req.Host], err = host.New(req, docroot)
				if err != nil {
					delete(hosts, req.Host)
					log.Printf("Error loading default host.")
					return nil
				}
			}
		} else {
			// Host is defined in settings.yaml
			hosts[req.Host], err = host.New(req, docroot)
			if err != nil {
				delete(hosts, req.Host)
				log.Printf("Requested host %s does not exists.", req.Host)
				return nil
			}
		}

	}
	return hosts[req.Host]
}

func (server Server) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	h := route(req)
	if h != nil {
		h.ServeHTTP(w, req)
	} else {
		log.Printf("Failed to serve host %s.\n", req.Host)
	}
}

func main() {
	flag.Parse()

	if *flagHelp == true {
		fmt.Printf("Showing %v usage.\n", os.Args[0])
		flag.PrintDefaults()
		return
	}

	if *flagVersion == true {
		fmt.Printf("%v version: %s\n", os.Args[0], version)
		return
	}

	hosts = make(map[string]*host.Host)

	settings = yaml.Open("settings.yaml")

	serverType := settings.GetString("server/type")

	domain	:= "unix"
	address := settings.GetString("server/socket")

	if address == "" {
		domain = "tcp"
		address = fmt.Sprintf("%s:%d", settings.GetString("server/bind"), settings.GetInt("server/port"))
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
