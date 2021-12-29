package repository

import (
	"database/sql"
	"fmt"
	"github.com/chapin666/kitten/model"
	"github.com/chapin666/kitten/pkg/db"
	"github.com/pkg/errors"
	"time"
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

// QueryAllFlowPage 查询流程分页数据
func (f *Flow) QueryAllFlowPage(params model.FlowQueryParam, pageIndex, pageSize uint) (
	int64,
	[]*model.FlowQueryResult,
	error,
) {
	var (
		where = "WHERE deleted=0 AND flag=1"
		args  []interface{}
	)

	if code := params.Code; code != "" {
		where = fmt.Sprintf("%s AND code LIKE ?", where)
		args = append(args, "%"+code+"%")
	}

	if name := params.Name; name != "" {
		where = fmt.Sprintf("%s AND name LIKE ?", where)
		args = append(args, "%"+name+"%")
	}

	if v := params.TypeCode; v != "" {
		where = fmt.Sprintf("%s AND type_code=?", where)
		args = append(args, v)
	}

	if v := params.Status; v > 0 {
		where = fmt.Sprintf("%s AND status=?", where)
		args = append(args, v)
	}

	n, err := f.DB.SelectInt(fmt.Sprintf("SELECT count(*) FROM %s %s", model.FlowTableName, where), args...)
	if err != nil {
		return 0, nil, errors.Wrapf(err, "查询分页数据发生错误")
	} else if n == 0 {
		return 0, nil, nil
	}

	query := fmt.Sprintf("SELECT id,record_id,created,code,name,version FROM %s %s ORDER BY id DESC", model.FlowTableName, where)
	if pageIndex > 0 && pageSize > 0 {
		query = fmt.Sprintf("%s limit %d,%d", query, (pageIndex-1)*pageSize, pageSize)
	}

	var items []*model.FlowQueryResult
	_, err = f.DB.Select(&items, query, args...)
	if err != nil {
		return 0, nil, errors.Wrapf(err, "查询分页数据发生错误")
	}

	return n, items, err
}

// GetFlow 获取流程数据
func (f *Flow) GetFlow(recordID string) (*model.Flow, error) {
	query := fmt.Sprintf("SELECT * FROM %s WHERE deleted=0 AND record_id=? LIMIT 1", model.FlowTableName)

	var flow model.Flow
	err := f.DB.SelectOne(&flow, query, recordID)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, errors.Wrapf(err, "获取流程数据发生错误")
	}
	return &flow, nil
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
	query := fmt.Sprintf(""+
		"SELECT * FROM %s "+
		"WHERE flow_id=? AND code=? AND deleted=0 "+
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

// CheckNodeCandidate 检查节点候选人
func (f *Flow) CheckNodeCandidate(nodeInstanceID, userID string) (bool, error) {
	query := fmt.Sprintf("SELECT "+
		"COUNT(*) "+
		"FROM %s "+
		"WHERE node_instance_id=? AND candidate_id=? AND deleted=0", model.NodeCandidateTableName)

	n, err := f.DB.SelectInt(query, nodeInstanceID, userID)
	if err != nil {
		return false, errors.Wrapf(err, "检查节点候选人发生错误")
	}

	return n > 0, nil
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
	query := fmt.Sprintf("SELECT "+
		"count(*) FROM %s "+
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

// QueryDoneIDs 查询已办理的流程实例ID列表
func (f *Flow) QueryDoneIDs(flowCode, userID string) ([]string, error) {
	query := fmt.Sprintf("SELECT "+
		"record_id "+
		"FROM %s "+
		"WHERE deleted=0 "+
		"AND flow_id IN (SELECT record_id FROM %s WHERE deleted=0 AND flag=1 AND code=?) "+
		"AND record_id IN(SELECT flow_instance_id FROM %s WHERE deleted=0 AND status=2 AND processor=?)",
		model.FlowInstanceTableName, model.FlowTableName, model.NodeInstanceTableName)

	var items []*model.FlowInstance
	_, err := f.DB.Select(&items, query, flowCode, userID)
	if err != nil {
		return nil, errors.Wrapf(err, "查询已办理的流程数据发生错误")
	}

	ids := make([]string, len(items))
	for i, item := range items {
		ids[i] = item.RecordID
	}

	return ids, nil
}

// QueryTodo 查询用户的待办数据
func (f *Flow) QueryTodo(typeCode string, flowCode string, userID string, limit int) ([]*model.FlowTodoResult, error) {
	var args []interface{}
	query := fmt.Sprintf(`SELECT
			ni.record_id,
			ni.flow_instance_id,
			ni.input_data,
			ni.node_id,
			f.data 'form_data',
			f.type_code 'form_type',
			fi.launcher,
			fi.launch_time,
			n.code 'node_code',
			n.name 'node_name',
			fw.name 'flow_name'
		FROM %s ni
			JOIN %s fi ON ni.flow_instance_id = fi.record_id AND fi.deleted = ni.deleted
			LEFT JOIN %s n ON ni.node_id = n.record_id AND n.deleted = ni.deleted
			LEFT JOIN %s f ON n.form_id = f.record_id AND f.deleted = n.deleted
			LEFT JOIN %s fw ON n.flow_id = fw.record_id AND fw.deleted=n.deleted
		WHERE 
			ni.deleted = 0 AND ni.status = 1 AND fi.status = 1 AND 
			ni.record_id IN (SELECT node_instance_id FROM %s WHERE deleted = 0 AND candidate_id = ?)
		`, model.NodeInstanceTableName, model.FlowInstanceTableName, model.NodeTableName,
		model.FormTableName, model.FlowTableName, model.NodeCandidateTableName)

	args = append(args, userID)
	if typeCode != "" {
		query = fmt.Sprintf("%s AND fi.flow_id IN (SELECT record_id FROM %s WHERE deleted=0 AND flag=1 AND type_code=?)", query, model.FlowTableName)
		args = append(args, typeCode)
	} else if flowCode != "" {
		query = fmt.Sprintf("%s AND fi.flow_id IN (SELECT record_id FROM %s WHERE deleted=0 AND flag=1 AND code=?)", query, model.FlowTableName)
		args = append(args, flowCode)
	}
	query = fmt.Sprintf("%s ORDER BY ni.id DESC LIMIT %d", query, limit)

	var items []*model.FlowTodoResult
	_, err := f.DB.Select(&items, query, args...)
	if err != nil {
		return nil, errors.Wrapf(err, "查询用户的待办数据发生错误")
	}
	return items, nil
}

// DeleteFlow 删除流程
func (f *Flow) DeleteFlow(flowID string) error {
	tran, err := f.DB.Begin()
	if err != nil {
		return errors.Wrapf(err, "删除流程开启事物发生错误")
	}

	ctimeUnix := time.Now().Unix()
	_, err = tran.Exec(fmt.Sprintf("UPDATE %s SET deleted=? WHERE deleted=0 AND record_id=?", model.FlowTableName), ctimeUnix, flowID)
	if err != nil {
		_ = tran.Rollback()
		return errors.Wrapf(err, "删除流程发生错误")
	}

	_, err = tran.Exec(fmt.Sprintf("UPDATE %s SET deleted=? WHERE deleted=0 AND source_node_id IN(SELECT record_id FROM %s WHERE deleted=0 AND flow_id=?)", model.NodeRouterTableName, model.NodeTableName), ctimeUnix, flowID)
	if err != nil {
		_ = tran.Rollback()
		return errors.Wrapf(err, "删除流程节点路由发生错误")
	}

	_, err = tran.Exec(fmt.Sprintf("UPDATE %s SET deleted=? WHERE deleted=0 AND node_id IN(SELECT record_id FROM %s WHERE deleted=0 AND flow_id=?)", model.NodeAssignmentTableName, model.NodeTableName), ctimeUnix, flowID)
	if err != nil {
		_ = tran.Rollback()
		return errors.Wrapf(err, "删除流程节点指派发生错误")
	}

	_, err = tran.Exec(fmt.Sprintf("UPDATE %s SET deleted=? WHERE deleted=0 AND node_id IN(SELECT record_id FROM %s WHERE deleted=0 AND flow_id=?)", model.NodePropertyTableName, model.NodeTableName), ctimeUnix, flowID)
	if err != nil {
		_ = tran.Rollback()
		return errors.Wrapf(err, "删除流程节点属性发生错误")
	}

	_, err = tran.Exec(fmt.Sprintf("UPDATE %s SET deleted=? WHERE deleted=0 AND flow_id=?", model.NodeTableName), ctimeUnix, flowID)
	if err != nil {
		_ = tran.Rollback()
		return errors.Wrapf(err, "删除流程节点发生错误")
	}

	_, err = tran.Exec(fmt.Sprintf("UPDATE %s SET deleted=? WHERE deleted=0 AND flow_id=?", model.FormTableName), ctimeUnix, flowID)
	if err != nil {
		_ = tran.Rollback()
		return errors.Wrapf(err, "删除流程表单发生错误")
	}

	err = tran.Commit()
	if err != nil {
		return errors.Wrapf(err, "删除流程提交事物发生错误")
	}
	return nil
}
