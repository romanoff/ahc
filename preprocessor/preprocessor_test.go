package preprocessor

import (
	"testing"
)

func TestCss(t *testing.T) {
	content := []byte(".a{color:red;}")
	css := &Css{Content: content}
	result, err := css.Get()
	if err != nil {
		t.Errorf("Expected to not get error while preprocessing css, but got %v", err)
	}
	if string(content) != string(result) {
		t.Errorf("Expected to get:\n%v\n, but got:\n%v", string(content), string(result))
	}
}
