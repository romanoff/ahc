package parse

import (
	"github.com/romanoff/ahc/component"
	"testing"
)

func TestParseComponentTest(t *testing.T) {
	fs := &Fs{}
	pool := &component.Pool{}
	fs.ParseIntoPool(pool, "test_files")
	testSuite, err := fs.ParseComponentTest("test_files/button/button.test", pool)
	if err != nil {
		t.Errorf("Didn't expect error while parsing component test, but got %v", err)
	}
	if len(testSuite.Tests) != 3 {
		t.Errorf("Expected to get 3 tests in suite, but got %v", len(testSuite.Tests))
	}
}
