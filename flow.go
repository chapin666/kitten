package kitten

import (
	"context"
	"encoding/json"
	"kitten/model"
	"kitten/pkg/db"
	"kitten/pkg/parse/xml"
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

// HandleFlow处理流程节点
// nodeInstanceID 节点实例内码
// userID 处理人
// input 输入数据
func HandleFlow(
	ctx context.Context,
	nodeInstanceID string,
	userID string,
	input interface{},
) (*model.HandleResult, error) {
	inputData, err := json.Marshal(input)
	if err != nil {
		return nil, err
	}
	return engine.HandleFlow(ctx, nodeInstanceID, userID, inputData)
}

// QueryTodoFlows 查询流程待办数据
// flowCode 流程编号
// userID 待办人
func QueryTodoFlows(flowCode string, userID string, limit int) ([]*model.FlowTodoResult, error) {
	return engine.QueryTodoFlows(flowCode, userID, limit)
}

// QueryNodeCandidates 查询节点实例的候选人ID列表
func QueryNodeCandidates(nodeInstanceID string) ([]string, error) {
	return engine.QueryNodeCandidates(nodeInstanceID)
}


// QueryDoneFlowIDs 查询已办理的流程实例ID列表
func QueryDoneFlowIDs(flowCode, userID string) ([]string, error) {
	return engine.QueryDoneFlowIDs(flowCode, userID)
}

// StopFlowInstance 停止流程实例
func StopFlowInstance(flowInstanceID string, allowStop func(*model.FlowInstance) bool) error {
	return engine.StopFlowInstance(flowInstanceID, allowStop)
}

