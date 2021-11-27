package xml

import (
	"context"
	"encoding/json"
	"fmt"
	"goworkflow/pkg/util"
	"testing"
)

func TestParseBasicBpmn(t *testing.T) {
	data, err := util.ReadFile("../../../test_data/basic.xml")
	if err != nil {
		t.Errorf("read file failed: %s", err.Error())
	}
	p := NewXMLParser()
	v, err := p.Parse(context.Background(), data)
	if err != nil {
		fmt.Println(err.Error())
	}
	buf, _ := json.Marshal(v)
	fmt.Println(string(buf))
}

func TestParseFormBpmn(t *testing.T) {
	data, err := util.ReadFile("../../../test_data/form.xml")
	if err != nil {
		t.Errorf("read file failed: %s", err.Error())
	}
	p := NewXMLParser()
	v, err := p.Parse(context.Background(), data)
	if err != nil {
		fmt.Println(err.Error())
	}
	buf, _ := json.Marshal(v)
	fmt.Println(string(buf))
}
