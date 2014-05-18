package main

import (
	"github.com/romanoff/gofh"
	"os"
)

// Ahc command line tool
func main() {
	flagHandler := gofh.Init()
	serverOptions := []*gofh.Option{
		&gofh.Option{Name: "port", Boolean: false},
		&gofh.Option{Name: "dev", Boolean: true},
	}
	flagHandler.HandleCommandWithOptions("server", serverOptions, StartServer)
	flagHandler.Parse(os.Args[1:])
}
