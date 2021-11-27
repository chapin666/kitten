package model

// 定义表名
const (
	FlowTableName            = "f_flow"             // 流程表
	NodeTableName            = "f_node"             // 流程节点表
	NodeRouterTableName      = "f_node_router"      // 节点路由
	NodeAssignmentTableName  = "f_node_assignment"  // 节点指派
	NodePropertyTableName    = "f_node_property"    // 节点属性
	FlowInstanceTableName    = "f_flow_instance"    // 流程实例
	NodeInstanceTableName    = "f_node_instance"    // 节点实例
	NodeTimingTableName      = "f_node_timing"      // 节点定时
	NodeCandidateTableName   = "f_node_candidate"   // 节点候选人
	FormTableName            = "f_form"             // 流程表单
	FormFieldTableName       = "f_form_field"       // 流程表单字段
	FieldOptionTableName     = "f_field_option"     // 流程表单字段选项
	FieldPropertyTableName   = "f_field_property"   // 流程表单字段属性
	FieldValidationTableName = "f_field_validation" // 流程表单字段校验
)