package component

import (
	"testing"
)

func TestComponent(t *testing.T) {
	c := &Component{Namespace: "goog.ui.Button"}
	if c.Namespace != "goog.ui.Button" {
		t.Error("Component namespace has not been saved")
	}

}
