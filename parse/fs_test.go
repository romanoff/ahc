package parse

import (
	"testing"
)

func TestNonExistingFileParse(t *testing.T) {
	fs := &Fs{}
	_, err := fs.ParseComponent("test_files/non_existing_component/non_existing_component.css")
	if err == nil {
		t.Error("Expected to get file doesn't exist error, but got nil")
	}
}
