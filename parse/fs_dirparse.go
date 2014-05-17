package parse

import (
	"github.com/romanoff/ahc/component"
	"github.com/romanoff/ahc/view"
	"os"
	"path/filepath"
)

func (self *Fs) ParseIntoPool(pool *component.Pool, dirpath string) error {
	filepath.Walk(dirpath, func(path string, f os.FileInfo, err error) error {
		if !f.IsDir() && filepath.Ext(f.Name()) == ".css" {
			component, err := self.ParseComponent(path)
			if err != nil {
				return err
			}
			err = pool.AppendComponent(component)
			if err != nil {
				return err
			}
		}
		return nil
	})
	return nil
}

func (self *Fs) ParseIntoTemplatePool(pool *view.Pool, dirpath string) error {
	filepath.Walk(dirpath, func(path string, f os.FileInfo, err error) error {
		if !f.IsDir() && filepath.Ext(f.Name()) == ".html" {
			template, err := self.ParseTemplate(path, dirpath)
			if err != nil {
				return err
			}
			err = pool.AddTemplate(template)
			if err != nil {
				return err
			}
		}
		return nil
	})
	return nil
}
