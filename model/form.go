package model

// Form 流程表单
type Form struct {
	ID       int64  `db:"id,primarykey,autoincrement" structs:"id" json:"id"`     // 唯一标识(自增ID)
	RecordID string `db:"record_id,size:36" structs:"record_id" json:"record_id"` // 记录内码(uuid)
	FlowID   string `db:"flow_id,size:36" structs:"flow_id" json:"flow_id"`       // 流程内码
	Code     string `db:"code,size:50" structs:"code" json:"code"`                // 表单编号(唯一)
	Name     string `db:"name,size:50" structs:"name" json:"name"`                // 表单名称
	TypeCode string `db:"type_code,size:50" structs:"type_code" json:"type_code"` // 表单类型(URL:表单链接路径 META:表单元数据)
	Data     string `db:"data,size:1024" structs:"data" json:"data"`              // 表单数据
	Created  int64  `db:"created" structs:"created" json:"created"`               // 创建时间戳
	Updated  int64  `db:"updated" structs:"updated" json:"updated"`               // 更新时间戳
	Deleted  int64  `db:"deleted" structs:"deleted" json:"deleted"`               // 删除时间戳
}
