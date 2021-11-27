package service

import (
	"goworkflow/repository"
	"goworkflow/model"
	"sync"
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

// GetFlowByCode 根据编号查询流程数据
func (f *Flow) GetFlowByCode(code string) (*model.Flow, error) {
	return f.FlowModel.GetFlowByCode(code)
}