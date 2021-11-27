package model

// FieldOption 字段选项
type FieldOption struct {
	ID        int64  `db:"id,primarykey,autoincrement" structs:"id" json:"id"`         // 唯一标识(自增ID)
	RecordID  string `db:"record_id,size:36" structs:"record_id" json:"record_id"`     // 记录内码(uuid)
	FieldID   string `db:"field_id,size:36" structs:"field_id" json:"field_id"`        // 字段内码
	ValueID   string `db:"value_id,size:50" structs:"value_id" json:"value_id"`        // 选项值ID
	ValueName string `db:"value_name,size:100" structs:"value_name" json:"value_name"` // 选项值名称
	Created   int64  `db:"created" structs:"created" json:"created"`                   // 创建时间戳
	Updated   int64  `db:"updated" structs:"updated" json:"updated"`                   // 更新时间戳
	Deleted   int64  `db:"deleted" structs:"deleted" json:"deleted"`                   // 删除时间戳
}
