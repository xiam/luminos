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
	name := req.Header.Get("Host")
	if _, ok := hosts[name]; ok == false {
		var err error
		hosts[name], err = host.New(req)
		if err != nil {
			log.Printf("Requested host %s does not exists.", name)
			return nil
		}
	}
	return hosts[name]
}

func (server Server) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	h := route(req)
	if h != nil {
		h.ServeHTTP(w, req)
	}
}

func main() {

	dtype := settings.GetString("server.type")

	switch dtype {
	case "fcgi":
		address := fmt.Sprintf("%s:%d", settings.GetString("server.bind"), settings.GetInt("server.port"))
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
