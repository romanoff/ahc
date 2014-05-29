package main

import (
	"fmt"
	"github.com/romanoff/ahc/component"
	"github.com/romanoff/ahc/parse"
	"os"
)

func TestComponents(options map[string]string) {
	folder := "."
	if fi, err := os.Stat("components"); err == nil && fi.IsDir() {
		folder = "components"
	}
	testPool := &component.TestPool{}
	fs := &parse.Fs{}
	err := fs.ParseIntoTestPool(testPool, folder)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	testPool.Run(true)
}
