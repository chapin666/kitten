package util

import (
	"testing"
)

func TestStringToInt(t *testing.T) {
	r, err := StringToInt("1")
	if err != nil {
		t.Error(err.Error())
	}
	t.Log(r)
}
