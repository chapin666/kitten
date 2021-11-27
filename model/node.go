package model


// Node 流程节点
type Node struct {
	ID       int64  `db:"id,primarykey,autoincrement" structs:"id" json:"id"`     // 唯一标识(自增ID)
	RecordID string `db:"record_id,size:36" structs:"record_id" json:"record_id"` // 记录内码(uuid)
	FlowID   string `db:"flow_id,size:36" structs:"flow_id" json:"flow_id"`       // 流程内码
	Code     string `db:"code,size:50" structs:"code" json:"code"`                // 节点编号
	Name     string `db:"name,size:50" structs:"name" json:"name"`                // 节点名称
	TypeCode string `db:"type_code,size:50" structs:"type_code" json:"type_code"` // 节点类型编号
	OrderNum string `db:"order_num,size:10" structs:"order_num" json:"order_num"` // 排序值
	FormID   string `db:"form_id,size:36" structs:"form_id" json:"form_id"`       // 表单内码
	Created  int64  `db:"created" structs:"created" json:"created"`               // 创建时间戳
	Updated  int64  `db:"updated" structs:"updated" json:"updated"`               // 更新时间戳
	Deleted  int64  `db:"deleted" structs:"deleted" json:"deleted"`               // 删除时间戳
}

