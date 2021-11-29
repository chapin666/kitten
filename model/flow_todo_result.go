package model

// FlowTodoResult 流程待办结果
type FlowTodoResult struct {
	RecordID       string  `db:"record_id" structs:"record_id" json:"record_id"`                      // 节点实例内码
	FlowInstanceID string  `db:"flow_instance_id" structs:"flow_instance_id" json:"flow_instance_id"` // 流程实例内码
	FlowName       string  `db:"flow_name" structs:"flow_name" json:"flow_name"`                      // 流程名称
	NodeID         string  `db:"node_id" structs:"node_id" json:"node_id"`                            // 节点内码
	NodeCode       string  `db:"node_code" structs:"node_code" json:"node_code"`                      // 节点编号
	NodeName       string  `db:"node_name" structs:"node_name" json:"node_name"`                      // 节点名称
	InputData      string  `db:"input_data" structs:"input_data" json:"input_data"`                   // 输入数据
	Launcher       string  `db:"launcher" structs:"launcher" json:"launcher"`                         // 发起人
	LaunchTime     int64   `db:"launch_time" structs:"launch_time" json:"launch_time"`                // 发起时间
	FormType       *string `db:"form_type" structs:"form_type" json:"form_type"`                      // 表单类型
	FormData       *string `db:"form_data" structs:"form_data" json:"form_data"`                      // 表单数据
}
