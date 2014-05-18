package main

import (
	"encoding/json"
	"fmt"
	"github.com/foize/go.sgr"
	"github.com/romanoff/ahc/component"
	"github.com/romanoff/ahc/parse"
	"github.com/romanoff/ahc/view"
	"github.com/romanoff/htmlcompressor"
	"net/http"
	"strings"
)

type AhcServer struct {
	TemplatesPool  *view.Pool
	Dev            bool
	FsParser       *parse.Fs
	HtmlCompressor *htmlcompressor.HtmlCompressor
}

func (self *AhcServer) ReadComponents() {
	componentsPool := &component.Pool{}
	self.TemplatesPool = view.InitPool()
	fsParser := self.FsParser
	fsParser.ParseIntoPool(componentsPool, "components")
	fsParser.ParseIntoTemplatePool(self.TemplatesPool, "templates")
	self.TemplatesPool.ComponentsPool = componentsPool
}

func (self *AhcServer) ViewHandler(w http.ResponseWriter, r *http.Request) {
	if self.Dev {
		self.ReadComponents()
	}
	path := strings.TrimPrefix(r.URL.Path, "/v/")
	templateJson, err := self.FsParser.ParseTemplateJson("templates/" + path + ".json")
	if err != nil {
		fmt.Fprint(w, err)
		return
	}
	params := templateJson.JsonGroups[0].Params
	content, err := self.TemplatesPool.Render(path, params)
	if err != nil {
		fmt.Fprint(w, err)
		return
	}
	w.Write(self.HtmlCompressor.Compress(content))
}

func (self *AhcServer) TemplateHandler(w http.ResponseWriter, r *http.Request) {
	if self.Dev {
		self.ReadComponents()
	}
	path := strings.TrimPrefix(r.URL.Path, "/t/")
	fmt.Print(sgr.MustParseln(fmt.Sprintf("Rendering  [fg-green]%v[reset]", path)))
	jsonParam := r.FormValue("params")
	params := make(map[string]interface{})
	err := json.Unmarshal([]byte(jsonParam), &params)
	if err != nil {
		fmt.Fprintf(w, "Json unmarshaling error: %v", err)
		return
	}
	content, err := self.TemplatesPool.Render(path, params)
	if err != nil {
		fmt.Fprint(w, err)
		return
	}
	w.Write(self.HtmlCompressor.Compress(content))
}

func (self *AhcServer) IndexHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "List of all components and views should be shown here")
}

func StartServer(options map[string]string) {
	server := &AhcServer{FsParser: &parse.Fs{}, HtmlCompressor: htmlcompressor.InitAll()}
	server.Dev = (options["dev"] == "true")
	server.ReadComponents()
	port := options["port"]
	if port == "" {
		port = "8080"
	}
	http.HandleFunc("/", server.IndexHandler)
	http.HandleFunc("/v/", server.ViewHandler)
	http.HandleFunc("/t/", server.TemplateHandler)
	http.ListenAndServe(":"+port, nil)
}
