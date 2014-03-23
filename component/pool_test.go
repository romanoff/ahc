package component

import (
	"testing"
)

func TestGetComponent(t *testing.T) {
	button1 := &Component{Namespace: "mp.button"}
	button2 := &Component{Namespace: "goog.button"}
	pool := &Pool{Components: []*Component{button1, button2}}
	if pool.GetComponent("button") != button1 {
		t.Errorf("Expected to get mp.button from pool, but got %v", pool.GetComponent("button").Namespace)
	}
	if pool.GetComponent("goog.button") != button2 {
		t.Errorf("Expected to get goog.button from pool, but got %v", pool.GetComponent("button").Namespace)
	}
}
