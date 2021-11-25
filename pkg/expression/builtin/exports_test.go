package builtin

import "testing"

func TestSliceStr(t *testing.T) {
	maps := []map[string]interface{}{
		{"name": "value1", "age": 10},
		{"name": "value2", "age": 10},
	}
	t.Log(SliceStr(maps, "age"))
}
