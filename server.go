package main

import (
	"fmt"
	"github.com/romanoff/ahc/component"
	"github.com/romanoff/ahc/parse"
	"github.com/romanoff/ahc/view"
	"net/http"
	"strings"
)

type AhcServer struct {
	TemplatesPool *view.Pool
	Dev           bool
}

func (self *AhcServer) ReadComponents() {
	componentsPool := &component.Pool{}
	self.TemplatesPool = view.InitPool()
	fsParser := &parse.Fs{}
	fsParser.ParseIntoPool(componentsPool, "components")
	fsParser.ParseIntoTemplatePool(self.TemplatesPool, "templates")
	self.TemplatesPool.ComponentsPool = componentsPool
}

func (self *AhcServer) ViewHandler(w http.ResponseWriter, r *http.Request) {
	if self.Dev {
		self.ReadComponents()
	}
	path := strings.TrimPrefix(r.URL.Path, "/v/")
	params := make(map[string]interface{})
	content, err := self.TemplatesPool.Render(path, params)
	if err != nil {
		fmt.Fprint(w, err)
		return
	}
	w.Write(content)
}

func (self *AhcServer) IndexHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "List of all components and views should be shown here")
}

func StartServer(options map[string]string) {
	server := &AhcServer{}
	server.Dev = (options["dev"] == "true")
	server.ReadComponents()
	port := options["port"]
	if port == "" {
		port = "8080"
	}
	http.HandleFunc("/", server.IndexHandler)
	http.HandleFunc("/v/", server.ViewHandler)
	http.ListenAndServe(":"+port, nil)
}
