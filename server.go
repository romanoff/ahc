package main

import (
	"crypto/sha1"
	"encoding/json"
	"fmt"
	"github.com/foize/go.sgr"
	"github.com/romanoff/ahc/component"
	"github.com/romanoff/ahc/parse"
	"github.com/romanoff/ahc/server"
	"github.com/romanoff/ahc/view"
	"github.com/romanoff/htmlcompressor"
	"net/http"
	"regexp"
	"strings"
)

type AhcServer struct {
	TemplatesPool   *view.Pool
	TemplatesStyles map[string][]byte
	Dev             bool
	FsParser        *parse.Fs
	HtmlCompressor  *htmlcompressor.HtmlCompressor
}

func (self *AhcServer) ReadComponents() {
	componentsPool := &component.Pool{}
	self.TemplatesPool = view.InitPool()
	fsParser := self.FsParser
	fsParser.ParseIntoPool(componentsPool, "components")
	fsParser.ParseIntoTemplatePool(self.TemplatesPool, "templates")
	self.TemplatesPool.ComponentsPool = componentsPool
	self.TemplatesStyles = make(map[string][]byte)
}

func (self *AhcServer) ViewHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
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
	err = self.AddStylesheets(params, path)
	if err != nil {
		fmt.Fprint(w, err)
		return
	}
	content, err := self.TemplatesPool.Render(path, params)
	if err != nil {
		fmt.Fprint(w, err)
		return
	}
	w.Write(self.HtmlCompressor.Compress(content))
}

func (self *AhcServer) TemplateHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	if self.Dev {
		self.ReadComponents()
	}
	path := strings.TrimPrefix(r.URL.Path, "/t/")
	fmt.Print(sgr.MustParseln(fmt.Sprintf("Rendering  [fg-green]%v[reset]", path)))
	jsonParam := r.FormValue("params")
	params := make(map[string]interface{})
	err := json.Unmarshal([]byte(jsonParam), &params)

	err = self.AddStylesheets(params, path)
	if err != nil {
		fmt.Fprint(w, err)
		return
	}
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

func (self *AhcServer) AddStylesheets(params map[string]interface{}, path string) error {
	content, err := self.getStyleFor(path)
	if err != nil {
		return err
	}
	h := sha1.New()
	h.Write(content)
	hash := fmt.Sprintf("%x", h.Sum(nil))
	params["stylesheets"] = []string{"/s/" + path + "-" + hash + ".css"}
	return nil
}

var stylePathRe *regexp.Regexp = regexp.MustCompile("/s/([\\w|/]+)")

func (self *AhcServer) StyleHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/css; charset=utf-8")
	match := stylePathRe.FindStringSubmatch(r.URL.Path)
	if len(match) != 2 {
		fmt.Fprint(w, "Wrong style path format")
		return
	}
	path := match[1]
	fmt.Print(sgr.MustParseln(fmt.Sprintf("Rendering css for [fg-green]%v[reset]", path)))
	if self.TemplatesStyles[path] != nil {
		w.Write(self.TemplatesStyles[path])
		return
	}
	content, err := self.getStyleFor(path)
	if err != nil {
		fmt.Fprint(w, err)
		return
	}
	self.TemplatesStyles[path] = content
	w.Write(content)
}

func (self *AhcServer) getStyleFor(path string) ([]byte, error) {
	componentsSearch := server.InitComponentSearch(self.TemplatesPool)
	err := componentsSearch.Search(path)
	if err != nil {
		return nil, err
	}
	return server.GetComponentsCss(self.TemplatesPool.ComponentsPool, componentsSearch.Components)
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
	http.HandleFunc("/s/", server.StyleHandler)
	http.ListenAndServe(":"+port, nil)
}
