package preprocessor

import (
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"strings"
	"testing"
)

func TestCss(t *testing.T) {
	css := Init()
	filepath.Walk("tests", func(filepath string, f os.FileInfo, err error) error {
		if !f.IsDir() && path.Ext(filepath) == ".acss" {
			css.Content, _ = ioutil.ReadFile(filepath)
			expected, _ := ioutil.ReadFile(strings.Replace(filepath, ".acss", ".css", 1))
			result, err := css.Get()
			if err != nil {
				t.Errorf("%v: Expected to not get error while preprocessing css, but got %v", filepath, err)
			}
			if string(expected) != string(result) {
				t.Errorf("%v: Expected to get:\n%v, but got:\n%v", filepath, string(expected), string(result))
			}
		}
		return nil
	})
}
