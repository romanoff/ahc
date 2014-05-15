package parse

import (
	"github.com/romanoff/ahc/component"
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
