package main

import (
	"net/http"
	"os"
	"text/template"
	"bufio"
	"strings"
	"controllers"
	"google.golang.org/appengine"
)

func main() {
	templates := populateTemplates()

	controllers.Register(templates)

	http.ListenAndServe(":8000", nil)

	appengine.Main()
}

func serveResource(w http.ResponseWriter, req *http.Request) {
	path := "public" + req.URL.Path
	var contentType string
	if strings.HasSuffix(path, ".css") {
		contentType = "text/css"
	} else if strings.HasSuffix(path, ".png") {
		contentType = "image/png"
	} else {
		contentType = "text/plain"
	}

	f, err := os.Open(path)

	if err == nil {
		defer f.Close()
		w.Header().Add("Content Type", contentType)

		br := bufio.NewReader(f)
		br.WriteTo(w)
	} else {
		w.WriteHeader(404)
	}
}

func populateTemplates() *template.Template {
	result := template.New("templates")

	basePath := "templates"
	templateFolder, _ := os.Open(basePath)
	defer templateFolder.Close()

	templatePathsRaw, _ := templateFolder.Readdir(-1)

	templatePaths := new([]string)
	for _, pathInfo := range templatePathsRaw {
		if !pathInfo.IsDir() {
			*templatePaths = append(*templatePaths,
						basePath + "/" + pathInfo.Name())
		}
	}

	result.ParseFiles(*templatePaths...)

	return result
}