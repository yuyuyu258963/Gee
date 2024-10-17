package gee

import (
	"testing"
)

func TestSplitStr(t *testing.T) {
	if len(splitStr("/a", urlSep)) != 0 {
		t.Fatalf("Waring %v", splitStr("/a", urlSep))
	}
}
