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

func TestReadCss(t *testing.T) {
	fs := &Fs{}
	component, err := fs.ParseComponent("test_files/button/button.css")
	if err != nil {
		t.Errorf("Expected not to get error while parsing button component, but got %v", err)
	}
	if component.Css == "" {
		t.Errorf("Expected component css to not be empty after parsing button css")
	}
}
