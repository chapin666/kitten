package model

// NodeInstance 节点实例表
type NodeInstance struct {
	ID             int64  `db:"id,primarykey,autoincrement" structs:"id" json:"id"`                          // 唯一标识(自增ID)
	RecordID       string `db:"record_id,size:36" structs:"record_id" json:"record_id"`                      // 记录内码(uuid)
	FlowInstanceID string `db:"flow_instance_id,size:36" structs:"flow_instance_id" json:"flow_instance_id"` // 流程实例内码
	NodeID         string `db:"node_id,size:36" structs:"node_id" json:"node_id"`                            // 节点内码
	Processor      string `db:"processor,size:36" structs:"processor" json:"processor"`                      // 处理人
	ProcessTime    int64  `db:"process_time" structs:"process_time" json:"process_time"`                     // 处理时间(秒时间戳)
	InputData      string `db:"input_data,size:1024" structs:"input_data" json:"input_data"`                 // 输入数据
	OutData        string `db:"out_data,size:1024" structs:"out_data" json:"out_data"`                       // 输出数据
	Status         int64  `db:"status" structs:"status" json:"status"`                                       // 处理状态(1:待处理 2:已完成)
	Created        int64  `db:"created" structs:"created" json:"created"`                                    // 创建时间戳
	Updated        int64  `db:"updated" structs:"updated" json:"updated"`                                    // 更新时间戳
	Deleted        int64  `db:"deleted" structs:"deleted" json:"deleted"`                                    // 删除时间戳
}
