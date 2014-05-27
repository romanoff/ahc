package parse

import (
	"github.com/romanoff/ahc/component"
	"github.com/romanoff/ahc/view"
	"github.com/romanoff/htmlcompressor"
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

func (self *Fs) ParseIntoTestPool(testPool *component.TestPool, dirpath string) error {
	pool := &component.Pool{}
	err := self.ParseIntoPool(pool, dirpath)
	if err != nil {
		return err
	}
	htmlcompressor := htmlcompressor.InitAll()
	filepath.Walk(dirpath, func(path string, f os.FileInfo, err error) error {
		if !f.IsDir() && filepath.Ext(f.Name()) == ".test" {
			testSuite, err := self.ParseComponentTest(path, pool)
			testSuite.Compressor = htmlcompressor
			if err != nil {
				return err
			}
			testPool.TestSuites = append(testPool.TestSuites, testSuite)
		}
		return nil
	})
	return nil
}
