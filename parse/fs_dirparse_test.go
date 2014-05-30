package parse

import (
	"github.com/romanoff/ahc/component"
	"github.com/romanoff/ahc/view"
	"testing"
)

func TestParseIntoPool(t *testing.T) {
	fs := &Fs{}
	pool := &component.Pool{}
	fs.ParseIntoPool(pool, "test_files")
	if len(pool.Components) == 0 {
		t.Error("Failed to parse button component from test_files folder into pool")
	}
}

func TestParseIntoTemplatePool(t *testing.T) {
	fs := &Fs{}
	pool := view.InitPool()
	fs.ParseIntoTemplatePool(pool, "test_templates")
	if len(pool.Templates) != 1 {
		t.Error("Failed to parse template from test_templates folder into pool")
	}
}
