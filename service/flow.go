package service

import (
	"fmt"
	"github.com/chapin666/kitten/pkg/util"
	"github.com/chapin666/kitten/repository"
	"github.com/chapin666/kitten/model"
	"sync"
	"time"
)

// Flow 流程管理
type Flow struct {
	sync.RWMutex
	FlowModel *repository.Flow `inject:""`
}

// CreateFlow 创建流程数据
func (f *Flow) CreateFlow(flow *model.Flow, nodes *model.NodeOperating, forms *model.FormOperating) error {
	if flow.Flag == 0 {
		flow.Flag = 1
	}
	return f.FlowModel.CreateFlow(flow, nodes, forms)
}


// QueryAllFlowPage 查询流程分页数据
func (f *Flow) QueryAllFlowPage(params model.FlowQueryParam, pageIndex, pageSize uint) (int64, []*model.FlowQueryResult, error) {
	return f.FlowModel.QueryAllFlowPage(params, pageIndex, pageSize)
}

// GetFlowByCode 根据编号查询流程数据
func (f *Flow) GetFlowByCode(code string) (*model.Flow, error) {
	return f.FlowModel.GetFlowByCode(code)
}

// LaunchFlowInstance 发起流程实例
func (f *Flow) LaunchFlowInstance(flowCode, nodeCode, launcher string, inputData []byte) (*model.NodeInstance, error) {

	// 根据工作流id查询数据库
	flow, err := f.FlowModel.GetFlowByCode(flowCode)
	if err != nil {
		return nil, err
	}
	if flow == nil {
		return nil, nil
	}

	// 根据nodeCode获取node
	node, err := f.FlowModel.GetNodeByCode(flow.RecordID, nodeCode)
	if err != nil {
		return nil, err
	}
	if node == nil {
		return nil, nil
	}

	// 创建flow实例
	flowInstance := &model.FlowInstance{
		RecordID:   util.UUID(),
		FlowID:     flow.RecordID,
		Launcher:   launcher,
		LaunchTime: time.Now().Unix(),
		Status:     1,
		Created:    time.Now().Unix(),
	}
	// 创建node实例
	nodeInstance := &model.NodeInstance{
		RecordID:       util.UUID(),
		FlowInstanceID: flowInstance.RecordID,
		NodeID:         node.RecordID,
		InputData:      string(inputData),
		Status:         1,
		Created:        flowInstance.Created,
	}

	err = f.FlowModel.CreateFlowInstance(flowInstance, nodeInstance)
	if err != nil {
		return nil, err
	}

	return nodeInstance, nil
}


// GetNode 获取流程节点
func (f *Flow) GetNode(recordID string) (*model.Node, error) {
	return f.FlowModel.GetNode(recordID)
}

// GetFlowInstance 获取流程实例
func (f *Flow) GetFlowInstance(recordID string) (*model.FlowInstance, error) {
	return f.FlowModel.GetFlowInstance(recordID)
}

// GetNodeInstance 获取流程节点实例
func (f *Flow) GetNodeInstance(recordID string) (*model.NodeInstance, error) {
	return f.FlowModel.GetNodeInstance(recordID)
}


// QueryNodeCandidates 查询节点候选人
func (f *Flow) QueryNodeCandidates(nodeInstanceID string) ([]*model.NodeCandidate, error) {
	return f.FlowModel.QueryNodeCandidates(nodeInstanceID)
}


// CheckNodeCandidate 检查节点候选人
func (f *Flow) CheckNodeCandidate(nodeInstanceID, userID string) (bool, error) {
	return f.FlowModel.CheckNodeCandidate(nodeInstanceID, userID)
}


// DoneNodeInstance 完成节点实例
func (f *Flow) DoneNodeInstance(nodeInstanceID, processor string, outData []byte) error {
	// 加锁保证节点实例的处理过程
	f.Lock()
	defer f.Unlock()

	nodeInstance, err := f.FlowModel.GetNodeInstance(nodeInstanceID)
	if err != nil {
		return err
	}

	if nodeInstance == nil || nodeInstance.Status == 2 {
		return fmt.Errorf("无效的处理节点")
	}

	info := map[string]interface{}{
		"processor":    processor,
		"process_time": time.Now().Unix(),
		"out_data":     string(outData),
		"status":       2,
		"updated":      time.Now().Unix(),
	}
	return f.FlowModel.UpdateNodeInstance(nodeInstanceID, info)
}


// CheckFlowInstanceTodo 检查流程实例待办事项
func (f *Flow) CheckFlowInstanceTodo(flowInstanceID string) (bool, error) {
	return f.FlowModel.CheckFlowInstanceTodo(flowInstanceID)
}


// DoneFlowInstance 完成流程实例
func (f *Flow) DoneFlowInstance(flowInstanceID string) error {
	info := map[string]interface{}{
		"status": 9,
	}
	return f.FlowModel.UpdateFlowInstance(flowInstanceID, info)
}

// QueryNodeRouters 查询节点路由
func (f *Flow) QueryNodeRouters(sourceNodeID string) ([]*model.NodeRouter, error) {
	return f.FlowModel.QueryNodeRouters(sourceNodeID)
}

// QueryNodeAssignments 查询节点指派
func (f *Flow) QueryNodeAssignments(nodeID string) ([]*model.NodeAssignment, error) {
	return f.FlowModel.QueryNodeAssignments(nodeID)
}

// CreateNodeInstance 创建节点实例
func (f *Flow) CreateNodeInstance(flowInstanceID, nodeID string, inputData []byte, candidates []string) (string, error) {
	nodeInstance := &model.NodeInstance{
		RecordID:       util.UUID(),
		FlowInstanceID: flowInstanceID,
		NodeID:         nodeID,
		InputData:      string(inputData),
		Status:         1,
		Created:        time.Now().Unix(),
	}

	var nodeCandidates []*model.NodeCandidate
	for _, c := range candidates {
		nodeCandidates = append(nodeCandidates, &model.NodeCandidate{
			RecordID:       util.UUID(),
			NodeInstanceID: nodeInstance.RecordID,
			CandidateID:    c,
			Created:        nodeInstance.Created,
		})
	}

	err := f.FlowModel.CreateNodeInstance(nodeInstance, nodeCandidates)
	if err != nil {
		return "", err
	}

	return nodeInstance.RecordID, nil
}

// GetNodeProperty 获取节点属性
func (f *Flow) GetNodeProperty(nodeID string) (map[string]string, error) {
	items, err := f.FlowModel.QueryNodeProperty(nodeID)
	if err != nil {
		return nil, err
	}

	data := make(map[string]string)
	for _, item := range items {
		data[item.Name] = item.Value
	}
	return data, nil
}

// CreateNodeTiming 创建定时节点
func (f *Flow) CreateNodeTiming(item *model.NodeTiming) error {
	item.ID = 0
	return f.FlowModel.CreateNodeTiming(item)
}

// QueryTodo 查询用户的待办节点实例数据
func (f *Flow) QueryTodo(typeCode string, flowCode string, userID string, limit int) ([]*model.FlowTodoResult, error) {
	return f.FlowModel.QueryTodo(typeCode, flowCode, userID, limit)
}

// QueryDoneIDs 查询已办理的流程实例ID列表
func (f *Flow) QueryDoneIDs(flowCode, userID string) ([]string, error) {
	return f.FlowModel.QueryDoneIDs(flowCode, userID)
}

// StopFlowInstance 停止流程实例
func (f *Flow) StopFlowInstance(flowInstanceID string) error {
	info := map[string]interface{}{
		"status": 9,
	}
	return f.FlowModel.UpdateFlowInstance(flowInstanceID, info)
}