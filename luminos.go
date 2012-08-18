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
	"github.com/xiam/gosexy/yaml"
	"github.com/xiam/luminos/host"
	"log"
	"net"
	"net/http"
	"net/http/fcgi"
)

type Server struct {
}

var hosts = make(map[string]*host.Host)

var settings = yaml.Open("settings.yaml")

func route(req *http.Request) *host.Host {
	if _, ok := hosts[req.Host]; ok == false {

		var err error
		var docroot string

		docroot = settings.GetString(fmt.Sprintf("hosts/%s", req.Host))

		if docroot != "" {
			// Host is defined in settings.yaml
			hosts[req.Host], err = host.New(req, docroot)
			if err != nil {
				log.Printf("Requested host %s does not exists.", req.Host)
				return nil
			}
		} else {
			// Serving default host.
			docroot = settings.GetString("hosts/default")
			if docroot != "" {
				hosts[req.Host], err = host.New(req, docroot)
				if err != nil {
					log.Printf("Default host was not found.")
					return nil
				}
			} else {
				// Default host is not defined.
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
		log.Printf("Could not route host %s.\n", req.Host)
	}
}

func main() {

	dtype := settings.GetString("server/type")

	switch dtype {
	case "fcgi":
		address := fmt.Sprintf("%s:%d", settings.GetString("server/bind"), settings.GetInt("server/port"))
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
