package parse

import (
	"github.com/romanoff/ahc/component"
	"testing"
)

func TestParseIntoPool(t *testing.T) {
	fs := &Fs{}
	pool := &component.Pool{}
	fs.ParseIntoPool(pool, "test_files")
	if len(pool.Components) != 1 {
		t.Error("Failed to parse button component from test_files folder into pool")
	}
}
