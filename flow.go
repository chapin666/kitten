package goworkflow

import (
	"context"
	"encoding/json"
	"goworkflow/model"
	"goworkflow/pkg/db"
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
	qlangExecer := NewQLangExecer()
	e, err := new(Engine).Init(xmlParser, qlangExecer, dbInstance, trace)
	if err != nil {
		panic(err)
	}
	engine = e
}

// Deploy 部署流程定义
// filePath: 流程文件
func Deploy(filePath string) (string, error) {
	return engine.Deploy(filePath)
}

// StartFlow 启动流程
// ctx context
// flowCode 流程编号
// nodeCode 开始节点编号
// userID 发起人
// input 输入数据
func StartFlow(
	ctx context.Context,
	flowCode string,
	nodeCode string,
	userID string,
	input interface{},
) (*model.HandleResult, error) {
	inputData, err := json.Marshal(input)
	if err != nil {
		return nil, err
	}
	return engine.StartFlow(ctx, flowCode, nodeCode, userID, inputData)
}
