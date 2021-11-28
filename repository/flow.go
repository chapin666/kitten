package repository

import (
	"database/sql"
	"fmt"
	"github.com/pkg/errors"
	"goworkflow/pkg/db"
	"goworkflow/model"
)

// Flow 流程管理
type Flow struct {
	DB *db.DB `inject:""`
}

// CreateFlow 创建流程数据
func (f *Flow) CreateFlow(flow *model.Flow, nodes *model.NodeOperating, forms *model.FormOperating) error {
	tran, err := f.DB.Begin()
	if err != nil {
		return errors.Wrapf(err, "创建流程基础数据开启事物发生错误")
	}

	// 写入数据到flow表
	err = tran.Insert(flow)
	if err != nil {
		_ = tran.Rollback()
		return errors.Wrapf(err, "插入流程数据发生错误")
	}

	// 写入数据到 NodeGroup、RouterGroup、AssignmentGroup、PropertyGroup
	if list := nodes.All(); len(list) > 0 {
		err = tran.Insert(list...)
		if err != nil {
			_ = tran.Rollback()
			return errors.Wrapf(err, "插入节点数据发生错误")
		}
	}

	// 写入数据到 FormGroup、FormFieldGroup、FieldOptionGroup、FieldPropertyGroup、FieldValidationGroup
	if list := forms.All(); len(list) > 0 {
		err = tran.Insert(list...)
		if err != nil {
			_ = tran.Rollback()
			return errors.Wrapf(err, "插入表单数据发生错误")
		}
	}

	err = tran.Commit()
	if err != nil {
		return errors.Wrapf(err, "创建流程基础数据提交事物发生错误")
	}
	return nil
}


// GetFlowByCode 根据编号查询流程数据
func (f *Flow) GetFlowByCode(code string) (*model.Flow, error) {
	query := fmt.Sprintf("SELECT * FROM %s WHERE flag=1 AND status=1 AND code=? AND deleted=0 "+
		"ORDER BY version DESC LIMIT 1", model.FlowTableName)

	var flow model.Flow
	err := f.DB.SelectOne(&flow, query, code)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, errors.Wrapf(err, "根据编号查询流程数据发生错误")
	}

	return &flow, nil
}

// GetNodeByCode 根据节点编号获取流程节点
func (f *Flow) GetNodeByCode(flowID, nodeCode string) (*model.Node, error) {
	query := fmt.Sprintf("" +
		"SELECT * FROM %s " +
		"WHERE flow_id=? AND code=? AND deleted=0 " +
		"ORDER BY order_num LIMIT 1", model.NodeTableName)

	var node model.Node
	err := f.DB.SelectOne(&node, query, flowID, nodeCode)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, errors.Wrapf(err, "根据节点编号获取流程节点发生错误")
	}

	return &node, nil
}

// GetNode 获取流程节点
func (f *Flow) GetNode(recordID string) (*model.Node, error) {
	query := fmt.Sprintf("SELECT * FROM %s WHERE record_id=? AND deleted=0", model.NodeTableName)

	var item model.Node
	err := f.DB.SelectOne(&item, query, recordID)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, errors.Wrapf(err, "获取流程节点发生错误")
	}

	return &item, nil
}

// QueryNodeCandidates 查询节点候选人
func (f *Flow) QueryNodeCandidates(nodeInstanceID string) ([]*model.NodeCandidate, error) {
	query := fmt.Sprintf("SELECT * FROM %s WHERE node_instance_id=? AND deleted=0", model.NodeCandidateTableName)

	var items []*model.NodeCandidate
	_, err := f.DB.Select(&items, query, nodeInstanceID)
	if err != nil {
		return nil, errors.Wrapf(err, "查询节点候选人发生错误")
	}

	return items, nil
}

// QueryNodeProperty 查询节点属性
func (f *Flow) QueryNodeProperty(nodeID string) ([]*model.NodeProperty, error) {
	query := fmt.Sprintf("SELECT * FROM %s WHERE node_id=? AND deleted=0", model.NodePropertyTableName)

	var items []*model.NodeProperty
	_, err := f.DB.Select(&items, query, nodeID)
	if err != nil {
		return nil, errors.Wrapf(err, "查询节点属性发生错误")
	}

	return items, nil
}

// QueryNodeRouters 查询节点路由
func (f *Flow) QueryNodeRouters(sourceNodeID string) ([]*model.NodeRouter, error) {
	query := fmt.Sprintf("SELECT * FROM %s WHERE source_node_id=? AND deleted=0", model.NodeRouterTableName)

	var items []*model.NodeRouter
	_, err := f.DB.Select(&items, query, sourceNodeID)
	if err != nil {
		return nil, errors.Wrapf(err, "查询节点路由发生错误")
	}

	return items, nil
}

// QueryNodeAssignments 查询节点指派
func (f *Flow) QueryNodeAssignments(nodeID string) ([]*model.NodeAssignment, error) {
	query := fmt.Sprintf("SELECT * FROM %s WHERE node_id=? AND deleted=0", model.NodeAssignmentTableName)

	var items []*model.NodeAssignment
	_, err := f.DB.Select(&items, query, nodeID)
	if err != nil {
		return nil, errors.Wrapf(err, "查询节点指派发生错误")
	}

	return items, nil
}


// CreateFlowInstance 创建流程实例
func (f *Flow) CreateFlowInstance(flowInstance *model.FlowInstance, nodeInstances ...*model.NodeInstance) error {
	tran, err := f.DB.Begin()
	if err != nil {
		return errors.Wrapf(err, "创建流程实例开启事物发生错误")
	}

	err = tran.Insert(flowInstance)
	if err != nil {
		err = tran.Rollback()
		if err != nil {
			return errors.Wrapf(err, "创建流程实例回滚事物发生错误")
		}
		return errors.Wrapf(err, "插入流程实例数据发生错误")
	}

	for _, n := range nodeInstances {
		err = tran.Insert(n)
		if err != nil {
			err = tran.Rollback()
			if err != nil {
				return errors.Wrapf(err, "创建流程实例回滚事物发生错误")
			}
			return errors.Wrapf(err, "插入流程节点实例数据发生错误")
		}
	}

	err = tran.Commit()
	if err != nil {
		return errors.Wrapf(err, "创建流程实例提交事物发生错误")
	}
	return nil
}

// UpdateFlowInstance 更新流程实例信息
func (f *Flow) UpdateFlowInstance(recordID string, info map[string]interface{}) error {
	_, err := f.DB.UpdateByPK(model.FlowInstanceTableName, db.M{"record_id": recordID}, info)
	if err != nil {
		return errors.Wrapf(err, "更新流程实例信息发生错误")
	}
	return nil
}

// CheckFlowInstanceTodo 检查流程实例待办事项
func (f *Flow) CheckFlowInstanceTodo(flowInstanceID string) (bool, error) {
	query := fmt.Sprintf("SELECT " +
		"count(*) FROM %s " +
		"WHERE status=1 AND flow_instance_id=? AND deleted=0", model.NodeInstanceTableName)
	n, err := f.DB.SelectInt(query, flowInstanceID)
	if err != nil {
		return false, errors.Wrapf(err, "检查流程待办事项发生错误")
	}
	return n > 0, nil
}


// CreateNodeInstance 创建流程节点实例
func (f *Flow) CreateNodeInstance(nodeInstance *model.NodeInstance, nodeCandidates []*model.NodeCandidate) error {
	tran, err := f.DB.Begin()
	if err != nil {
		return errors.Wrapf(err, "创建流程节点实例开启事物发生错误")
	}

	err = tran.Insert(nodeInstance)
	if err != nil {
		err = tran.Rollback()
		if err != nil {
			return errors.Wrapf(err, "创建流程节点实例回滚事物发生错误")
		}
		return errors.Wrapf(err, "插入流程节点实例数据发生错误")
	}

	for _, c := range nodeCandidates {
		err = tran.Insert(c)
		if err != nil {
			err = tran.Rollback()
			if err != nil {
				return errors.Wrapf(err, "创建流程节点候选人回滚事物发生错误")
			}
			return errors.Wrapf(err, "插入流程节点候选人数据发生错误")
		}
	}

	err = tran.Commit()
	if err != nil {
		return errors.Wrapf(err, "创建流程节点实例提交事物发生错误")
	}
	return nil
}

// UpdateNodeInstance 更新节点实例信息
func (f *Flow) UpdateNodeInstance(recordID string, info map[string]interface{}) error {
	_, err := f.DB.UpdateByPK(model.NodeInstanceTableName, db.M{"record_id": recordID}, info)
	if err != nil {
		return errors.Wrapf(err, "更新节点实例信息发生错误")
	}
	return nil
}


// GetFlowInstance 获取流程实例
func (f *Flow) GetFlowInstance(recordID string) (*model.FlowInstance, error) {
	query := fmt.Sprintf("SELECT * FROM %s WHERE record_id=? AND deleted=0 LIMIT 1", model.FlowInstanceTableName)

	var item model.FlowInstance
	err := f.DB.SelectOne(&item, query, recordID)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, errors.Wrapf(err, "获取流程实例发生错误")
	}

	return &item, nil
}

// GetNodeInstance 获取流程节点实例
func (f *Flow) GetNodeInstance(recordID string) (*model.NodeInstance, error) {
	query := fmt.Sprintf("SELECT * FROM %s WHERE record_id=? AND deleted=0 LIMIT 1", model.NodeInstanceTableName)

	var item model.NodeInstance
	err := f.DB.SelectOne(&item, query, recordID)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, errors.Wrapf(err, "获取流程节点实例发生错误")
	}

	return &item, nil
}

// CreateNodeTiming 创建定时节点
func (f *Flow) CreateNodeTiming(item *model.NodeTiming) error {
	err := f.DB.Insert(item)
	if err != nil {
		return errors.Wrapf(err, "创建节点定时发生错误")
	}
	return nil
}

// UpdateNodeTiming 更新定时节点
func (f *Flow) UpdateNodeTiming(nodeInstanceID string, info map[string]interface{}) error {
	_, err := f.DB.UpdateByPK(model.NodeTimingTableName, db.M{"node_instance_id": nodeInstanceID}, db.M(info))
	if err != nil {
		return errors.Wrapf(err, "更新节点定时发生错误")
	}
	return nil
}
