package kitten

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/chapin666/kitten/mapper"
	"github.com/chapin666/kitten/model"
	"github.com/chapin666/kitten/pkg/db"
	"github.com/chapin666/kitten/pkg/parse"
	"github.com/chapin666/kitten/pkg/parse/xml"
	"github.com/chapin666/kitten/pkg/util"
	"github.com/chapin666/kitten/service"
	"github.com/facebookgo/inject"
	"github.com/pkg/errors"
	"log"
	"strconv"
	"time"
)

// Engine .
type Engine struct {
	parser  parse.Parser
	execer  Execer
	flowSvc *service.Flow
}

// 初始化
func New(mysqlDNS string, trace bool) (*Engine, error) {
	var g inject.Graph
	var flowSvc service.Flow

	sqlDB, trace, err := db.NewMySQL(db.SetDSN(mysqlDNS), db.SetTrace(trace))
	if err != nil {
		return nil, err
	}
	dbInstance := db.NewMySQLWithDB(sqlDB, trace)

	if err := g.Provide(&inject.Object{Value: dbInstance}, &inject.Object{Value: &flowSvc}); err != nil {
		return nil, err
	}

	if err := g.Populate(); err != nil {
		return nil, err
	}

	mapper.FlowDBMap(dbInstance)
	if err := dbInstance.CreateTablesIfNotExists(); err != nil {
		return nil, err
	}

	return &Engine{
		parser:  xml.NewXMLParser(),
		execer:  NewQLangExecer(),
		flowSvc: &flowSvc,
	}, nil
}

// 部署
func (e *Engine) Deploy(filePath string) (string, error) {
	// 读取并解析bpmn文件
	data, err := util.ReadFile(filePath)
	if err != nil {
		return "", err
	}
	result, err := e.parser.Parse(context.Background(), data)
	if err != nil {
		return "", err
	}

	// 检查流程是否存在，如果存在则检查版本号是否一致，如果不一致则创建新流程
	oldFlow, err := e.flowSvc.GetFlowByCode(result.FlowID)
	if err != nil {
		return "", err
	}

	// 比较版本
	if oldFlow != nil {
		if result.FlowVersion <= oldFlow.Version {
			return oldFlow.RecordID, nil
		}
	}

	flow := &model.Flow{
		RecordID: util.UUID(),
		Code:     result.FlowID,
		Name:     result.FlowName,
		Version:  result.FlowVersion,
		XML:      string(data),
		Status:   result.FlowStatus,
		Created:  time.Now().Unix(),
	}

	nodeOperating, formOperating := e.parseOperating(flow, result.Nodes)

	// 解析节点表单数据
	for _, node := range result.Nodes {
		// 查找表单ID不为空并且不包含表单字段的节点
		if node.FormResult != nil && node.FormResult.ID != "" && len(node.FormResult.Fields) == 0 {
			// 查找表单ID
			var formID string
			for _, form := range formOperating.FormGroup {
				if form.Code == node.FormResult.ID {
					formID = form.RecordID
					break
				}
			}
			if formID != "" {
				for i, ns := range nodeOperating.NodeGroup {
					if ns.Code == node.NodeID {
						nodeOperating.NodeGroup[i].FormID = formID
						break
					}
				}
			}
		}
	}

	err = e.flowSvc.CreateFlow(flow, nodeOperating, formOperating)
	if err != nil {
		return "", err
	}
	return flow.RecordID, nil
}

// 解析node和form
func (e *Engine) parseOperating(flow *model.Flow, nodeResults []*parse.NodeResult) (
	*model.NodeOperating,
	*model.FormOperating,
) {
	nodeOperating := &model.NodeOperating{
		NodeGroup: make([]*model.Node, len(nodeResults)),
	}
	formOperating := &model.FormOperating{}

	for i, n := range nodeResults {
		// 节点数据
		node := &model.Node{
			RecordID: util.UUID(),
			FlowID:   flow.RecordID,
			Code:     n.NodeID,
			Name:     n.NodeName,
			TypeCode: n.NodeType.String(),
			OrderNum: strconv.FormatInt(int64(i+10), 10),
			Created:  flow.Created,
		}

		if n.FormResult != nil {
			e.parseFormOperating(formOperating, flow, node, n.FormResult)
		}

		// 当前节点表达式
		for _, exp := range n.CandidateExpressions {
			nodeOperating.AssignmentGroup = append(nodeOperating.AssignmentGroup, &model.NodeAssignment{
				RecordID:   util.UUID(),
				NodeID:     node.RecordID,
				Expression: exp,
				Created:    flow.Created,
			})
		}

		nodeOperating.NodeGroup[i] = node
	}

	var getNodeRecordID = func(nodeCode string) string {
		for _, n := range nodeOperating.NodeGroup {
			if n.Code == nodeCode {
				return n.RecordID
			}
		}
		return ""
	}

	for _, n := range nodeResults {
		// 增加路由
		for _, r := range n.Routers {
			nodeOperating.RouterGroup = append(nodeOperating.RouterGroup, &model.NodeRouter{
				RecordID:     util.UUID(),
				SourceNodeID: getNodeRecordID(n.NodeID),
				TargetNodeID: getNodeRecordID(r.TargetNodeID),
				Expression:   r.Expression,
				Explain:      r.Explain,
				Created:      flow.Created,
			})
		}

		// 增加节点属性
		for _, p := range n.Properties {
			nodeOperating.PropertyGroup = append(nodeOperating.PropertyGroup, &model.NodeProperty{
				RecordID: util.UUID(),
				NodeID:   getNodeRecordID(n.NodeID),
				Name:     p.Name,
				Value:    p.Value,
				Created:  flow.Created,
			})
		}
	}

	return nodeOperating, formOperating
}

// 解析form
func (e *Engine) parseFormOperating(
	formOperating *model.FormOperating,
	flow *model.Flow,
	node *model.Node,
	formResult *parse.NodeFormResult,
) {
	if formResult.ID == "" {
		return
	}

	for _, f := range formOperating.FormGroup {
		if f.Code == formResult.ID {
			node.FormID = f.RecordID
			return
		}
	}

	// form 表单
	form := &model.Form{
		RecordID: util.UUID(),
		FlowID:   flow.RecordID,
		Code:     formResult.ID,
		TypeCode: "META",
		Created:  flow.Created,
	}

	// 解析URL类型
	if fields := formResult.Fields; len(fields) == 2 {
		if fields[0].ID == "type_code" &&
			fields[0].DefaultValue == "URL" &&
			fields[1].ID == "data" {
			form.TypeCode = "URL"
			form.Data = fields[1].DefaultValue
			formOperating.FormGroup = append(formOperating.FormGroup, form)
			node.FormID = form.RecordID
			return
		}
	}

	meta, _ := json.Marshal(formResult.Fields)
	form.Data = string(meta)

	for _, ff := range formResult.Fields {
		// FormField 字段
		field := &model.FormField{
			RecordID:     util.UUID(),
			FormID:       form.RecordID,
			Code:         ff.ID,
			Label:        ff.Label,
			TypeCode:     ff.Type,
			DefaultValue: ff.DefaultValue,
			Created:      flow.Created,
		}

		// 字段值
		for _, item := range ff.Values {
			formOperating.FieldOptionGroup = append(formOperating.FieldOptionGroup, &model.FieldOption{
				RecordID:  util.UUID(),
				FieldID:   field.RecordID,
				ValueID:   item.ID,
				ValueName: item.Name,
				Created:   flow.Created,
			})
		}

		// 属性
		for _, item := range ff.Properties {
			formOperating.FieldPropertyGroup = append(formOperating.FieldPropertyGroup, &model.FieldProperty{
				RecordID: util.UUID(),
				FieldID:  field.RecordID,
				Code:     item.ID,
				Value:    item.Value,
				Created:  flow.Created,
			})
		}

		// 校验配置
		for _, item := range ff.Validations {
			formOperating.FieldValidationGroup = append(formOperating.FieldValidationGroup, &model.FieldValidation{
				RecordID:         util.UUID(),
				FieldID:          field.RecordID,
				ConstraintName:   item.Name,
				ConstraintConfig: item.Config,
				Created:          flow.Created,
			})
		}

		formOperating.FormFieldGroup = append(formOperating.FormFieldGroup, field)
	}

	formOperating.FormGroup = append(formOperating.FormGroup, form)
	node.FormID = form.RecordID
}

// SaveFlow 保存流程
func (e *Engine) SaveFlow(data []byte) (string, error) {

	result, err := e.parser.Parse(context.Background(), data)
	if err != nil {
		return "", err
	}

	// 检查流程是否存在，如果存在则检查版本号是否一致，如果不一致则创建新流程
	oldFlow, err := e.flowSvc.GetFlowByCode(result.FlowID)
	if err != nil {
		return "", err
	} else if oldFlow != nil {
		if result.FlowVersion <= oldFlow.Version {
			return oldFlow.RecordID, nil
		}
	}

	flow := &model.Flow{
		RecordID: util.UUID(),
		Code:     result.FlowID,
		Name:     result.FlowName,
		Version:  result.FlowVersion,
		XML:      string(data),
		Status:   result.FlowStatus,
		Created:  time.Now().Unix(),
	}
	nodeOperating, formOperating := e.parseOperating(flow, result.Nodes)

	// 解析节点表单数据
	for _, node := range result.Nodes {
		// 查找表单ID不为空并且不包含表单字段的节点
		if node.FormResult != nil && node.FormResult.ID != "" && len(node.FormResult.Fields) == 0 {
			// 查找表单ID
			var formID string
			for _, form := range formOperating.FormGroup {
				if form.Code == node.FormResult.ID {
					formID = form.RecordID
					break
				}
			}
			if formID != "" {
				for i, ns := range nodeOperating.NodeGroup {
					if ns.Code == node.NodeID {
						nodeOperating.NodeGroup[i].FormID = formID
						break
					}
				}
			}
		}
	}

	err = e.flowSvc.CreateFlow(flow, nodeOperating, formOperating)
	if err != nil {
		return "", err
	}
	return flow.RecordID, nil
}

// 启动流程
func (e *Engine) StartFlow(
	ctx context.Context,
	flowCode string,
	nodeCode string,
	userID string,
	inputData []byte,
) (*model.HandleResult, error) {
	nodeInstance, err := e.flowSvc.LaunchFlowInstance(flowCode, nodeCode, userID, inputData)
	if err != nil {
		return nil, err
	}
	if nodeInstance == nil {
		return nil, errors.New("未找到流程信息")
	}

	return e.nextFlowHandle(ctx, nodeInstance.RecordID, userID, inputData)
}

func (e *Engine) nextFlowHandle(
	ctx context.Context,
	nodeInstanceID string,
	userID string,
	inputData []byte,
) (*model.HandleResult, error) {
	var result model.HandleResult

	var onNextNode = OnNextNodeOption(func(
		node *model.Node,
		nodeInstance *model.NodeInstance,
		nodeCandidates []*model.NodeCandidate,
	) {
		var cIDs []string
		for _, nc := range nodeCandidates {
			cIDs = append(cIDs, nc.CandidateID)
		}
		result.NextNodes = append(result.NextNodes, &model.NextNode{
			Node:         node,
			NodeInstance: nodeInstance,
			CandidateIDs: cIDs,
		})
	})

	var onFlowEnd = OnFlowEndOption(func(_ *model.FlowInstance) {
		result.IsEnd = true
	})

	nr, err := new(NodeRouter).Init(ctx, e, nodeInstanceID, inputData, onNextNode, onFlowEnd)
	if err != nil {
		return nil, err
	}

	err = nr.Next(userID)
	if err != nil {
		return nil, err
	}
	result.FlowInstance = nr.GetFlowInstance()

	if !result.IsEnd {
		for _, item := range result.NextNodes {
			prop, err := e.flowSvc.GetNodeProperty(item.Node.RecordID)
			if err != nil {
				return nil, err
			}

			// 检查节点是否设定定时器，如果设定则加入定时
			if v := prop["timing"]; v != "" {
				expired, err := strconv.Atoi(v)
				if err == nil && expired > 0 {
					nt := &model.NodeTiming{
						NodeInstanceID: item.NodeInstance.RecordID,
						Processor:      item.CandidateIDs[0],
						Input:          prop["timing_input"],
						ExpiredAt:      time.Now().Add(time.Duration(expired) * time.Minute).Unix(),
						Created:        time.Now().Unix(),
					}

					if flag, ok := FromFlagContext(ctx); ok {
						nt.Flag = flag
					}

					err = e.flowSvc.CreateNodeTiming(nt)
					if err != nil {
						e.errorf("%+v", err)
					}
				}
			}
		}
	}

	return &result, nil
}

// HandleFlow 处理流程节点
// nodeInstanceID 节点实例内码
// userID 处理人
// inputData 输入数据
func (e *Engine) HandleFlow(
	ctx context.Context,
	nodeInstanceID,
	userID string,
	inputData []byte,
) (*model.HandleResult, error) {
	// 检查是否是节点候选人
	exists, err := e.flowSvc.CheckNodeCandidate(nodeInstanceID, userID)
	if err != nil {
		return nil, err
	}
	if !exists {
		return nil, fmt.Errorf("无效的节点处理人")
	}

	nodeInstance, err := e.flowSvc.GetNodeInstance(nodeInstanceID)
	if err != nil {
		return nil, err
	}
	if nodeInstance == nil || nodeInstance.Status != 1 {
		return nil, fmt.Errorf("无效的处理节点")
	}

	return e.nextFlowHandle(ctx, nodeInstanceID, userID, inputData)
}

// QueryAllFlowPage 查询流程分页数据
func (e *Engine) QueryAllFlowPage(params model.FlowQueryParam, pageIndex, pageSize uint) (
	int64,
	[]*model.FlowQueryResult,
	error,
) {
	return e.flowSvc.QueryAllFlowPage(params, pageIndex, pageSize)
}

// GetFlow 获取流程数据
func (e *Engine) GetFlow(recordID string) (*model.Flow, error) {
	return e.flowSvc.GetFlow(recordID)
}

// QueryTodoFlows 查询流程待办数据
// flowCode 流程编号
// userID 待办人
func (e *Engine) QueryTodoFlows(flowCode string, userID string, limit int) ([]*model.FlowTodoResult, error) {
	return e.flowSvc.QueryTodo("", flowCode, userID, limit)
}

// QueryNodeCandidates 查询节点实例的候选人ID列表
func (e *Engine) QueryNodeCandidates(nodeInstanceID string) ([]string, error) {
	candidates, err := e.flowSvc.QueryNodeCandidates(nodeInstanceID)
	if err != nil {
		return nil, err
	}

	ids := make([]string, len(candidates))

	for i, c := range candidates {
		ids[i] = c.CandidateID
	}

	return ids, nil
}

// QueryDoneFlowIDs 查询已办理的流程实例ID列表
func (e *Engine) QueryDoneFlowIDs(flowCode, userID string) ([]string, error) {
	return e.flowSvc.QueryDoneIDs(flowCode, userID)
}

// StopFlowInstance 停止流程实例
func (e *Engine) StopFlowInstance(flowInstanceID string, allowStop func(*model.FlowInstance) bool) error {
	flowInstance, err := e.flowSvc.GetFlowInstance(flowInstanceID)
	if err != nil {
		return err
	}

	if allowStop != nil && !allowStop(flowInstance) {
		return errors.New("不允许停止流程")
	}
	return e.flowSvc.StopFlowInstance(flowInstanceID)
}

// DeleteFlow 删除流程
func (e *Engine) DeleteFlow(flowID string) error {
	return e.flowSvc.DeleteFlow(flowID)
}


func (e *Engine) errorf(format string, args ...interface{}) {
	log.Printf(format, args...)
}
