package parse

import (
	"fmt"
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
	if len(testSuite.Tests) != 1 {
		t.Errorf("Expected to get 1 tests in suite, but got %v", len(testSuite.Tests))
	}
	fmt.Println(string(testSuite.Tests[0].Expected))
}
