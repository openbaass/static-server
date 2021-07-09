package main

import (
	"bytes"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"path"
	"regexp"
	"strings"
)

var (
	dir          = *flag.String("dir", "./public", "Set this to the directory of static files")
	port         = *flag.Int("port", 8000, "HTTP port")
	allowLogging = *flag.Bool("logging", true, "Allow logging")
	cache        map[string][]byte
)

func main() {
	cache = make(map[string][]byte)
	indexHtml, err := ioutil.ReadFile(path.Join(dir, "index.html"))
	if err != nil {
		log.Fatalf("No index.html found!\nSearched at: %s\n%s\n", path.Join(dir, "index.html"), err.Error())
	}
	notFound, err := ioutil.ReadFile(path.Join(dir, "404.html"))
	if err != nil {
		notFound = []byte("<h1> 404 - Page not found </h1>")
	}
	http.HandleFunc("/", func(writer http.ResponseWriter, request *http.Request) {
		writer.Header().Add("Server", "Static-Server")
		if allowLogging {
			log.Printf("%s | %s (%s)", request.Method, request.RequestURI, request.RemoteAddr)
		}
		if ok, _ := regexp.Match(".*\\.\\w{1,5}\\z", []byte(request.RequestURI)); ok {
			//Send back a static file
			//Add headers
			s := strings.Split(request.RequestURI, ".")
			ext := s[len(s)-1]
			switch ext {
			case "css":
				writer.Header().Add("Content-Type", "text/css")
				break
			case "js":
				writer.Header().Add("Content-Type", "application/javascript")
				break
			case "svg":
				writer.Header().Add("Content-Type", "image/svg+xml")
				break

			}

			//Check cache
			cached := cache[request.RequestURI]
			if cached != nil {
				if bytes.Equal(cached, notFound) {
					writer.WriteHeader(404)
					writer.Write(cached)
					return
				}
				writer.WriteHeader(200)
				writer.Write(cached)
				return
			}
			f, err := ioutil.ReadFile(path.Join(dir, request.RequestURI))
			if err != nil {
				writer.WriteHeader(404)
				writer.Write(notFound)
				cache[request.RequestURI] = notFound
				return
			}
			writer.WriteHeader(200)
			writer.Write(f)
			cache[request.RequestURI] = f
			return
		}
		writer.Header().Add("Content-Type", "text/html")
		writer.WriteHeader(200)
		writer.Write(indexHtml)
	})

	log.Printf("Serving on port %d.\n", port)
	log.Fatalln(http.ListenAndServe(fmt.Sprintf(":%d", port), nil))
}
