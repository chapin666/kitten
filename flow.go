package goworkflow

import (
	"goworkflow/pkg/db"
	"goworkflow/pkg/expression"
	"goworkflow/pkg/parse/xml"
)

var (
	engine *Engine
)

// Init 初始化流程配置
func Init(opts ...db.Option) {
	dbInstance, trace, err := db.NewMySQL(opts...)
	if err != nil {
		panic(err)
	}

	xmlParser := xml.NewXMLParser()
	qlangExecer := expression.CreateExecer("")
	e, err := new(Engine).Init(xmlParser, qlangExecer, dbInstance, trace)
	if err != nil {
		panic(err)
	}
	engine = e
}