package goworkflow

import (
	"context"
	"database/sql"
	"encoding/json"
	"github.com/facebookgo/inject"
	"goworkflow/mapper"
	"goworkflow/model"
	"goworkflow/pkg/db"
	"goworkflow/pkg/expression"
	"goworkflow/pkg/parse"
	"goworkflow/pkg/util"
	"goworkflow/service"
	"strconv"
	"time"
)

// Engine .
type Engine struct {
	parser  parse.Parser
	execer  expression.Execer
	flowSvc *service.Flow
}

// 初始化
func (e *Engine) Init(parser parse.Parser, execer expression.Execer, sqlDB *sql.DB, trace bool) (*Engine, error) {
	var g inject.Graph
	var flowSvc service.Flow

	dbInstance := db.NewMySQLWithDB(sqlDB, trace)

	err := g.Provide(&inject.Object{Value: dbInstance}, &inject.Object{Value: &flowSvc})
	if err != nil {
		return e, err
	}

	err = g.Populate()
	if err != nil {
		return e, err
	}

	mapper.FlowDBMap(dbInstance)
	err = dbInstance.CreateTablesIfNotExists()
	if err != nil {
		return e, err
	}

	e.parser = parser
	e.execer = execer
	e.flowSvc = &flowSvc

	return e, nil
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

// 创建节点操作
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
