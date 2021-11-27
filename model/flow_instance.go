package model

// FlowInstance 流程实例
type FlowInstance struct {
	ID         int64  `db:"id,primarykey,autoincrement" structs:"id" json:"id"`     // 唯一标识(自增ID)
	RecordID   string `db:"record_id,size:36" structs:"record_id" json:"record_id"` // 记录内码(uuid)
	FlowID     string `db:"flow_id,size:36" structs:"flow_id" json:"flow_id"`       // 流程内码
	Status     int64  `db:"status" structs:"status" json:"status"`                  // 流程状态(0:未开始 1:进行中 2:暂停 3:已停止 9:已完成)
	Launcher   string `db:"launcher,size:36" structs:"launcher" json:"launcher"`    // 发起人
	LaunchTime int64  `db:"launch_time" structs:"launch_time" json:"launch_time"`   // 发起时间
	Created    int64  `db:"created" structs:"created" json:"created"`               // 创建时间戳
	Updated    int64  `db:"updated" structs:"updated" json:"updated"`               // 更新时间戳
	Deleted    int64  `db:"deleted" structs:"deleted" json:"deleted"`               // 删除时间戳
}
