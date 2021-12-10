package kitten

import (
	"encoding/json"
	"fmt"
	"github.com/chapin666/kitten/pkg/expression"
)

// Execer 表达式执行器
type Execer interface {
	// 执行表达式返回布尔类型的值
	ExecReturnBool(exp, params []byte) (bool, error)

	// 执行表达式返回字符串切片类型的值
	ExecReturnStringSlice(exp, params []byte) ([]string, error)
}

type execer struct {}

// NewQLangExecer 创建基于qlang的表达式执行器
func NewQLangExecer() Execer {
	return &execer{}
}


func (*execer) ExecReturnBool(exp, params []byte) (bool, error) {
	var m map[string]interface{}
	err := json.Unmarshal(params, &m)
	if err != nil {
		return false, err
	}
	return expression.ExecParamBool(string(exp), m)
}

func (*execer) ExecReturnStringSlice(exp, params []byte) ([]string, error) {
	var m map[string]interface{}
	err := json.Unmarshal(params, &m)
	if err != nil {
		return nil, err
	}

	fmt.Println(string(exp))
	fmt.Println(m)

	return expression.ExecParamSliceStr(string(exp), m)
}


