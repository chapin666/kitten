package model

// FieldProperty 字段属性
type FieldProperty struct {
	ID       int64  `db:"id,primarykey,autoincrement" structs:"id" json:"id"`     // 唯一标识(自增ID)
	RecordID string `db:"record_id,size:36" structs:"record_id" json:"record_id"` // 记录内码(uuid)
	FieldID  string `db:"field_id,size:36" structs:"field_id" json:"field_id"`    // 字段内码
	Code     string `db:"code,size:50" structs:"code" json:"code"`                // 属性编号
	Value    string `db:"value,size:100" structs:"value" json:"value"`            // 属性值
	Created  int64  `db:"created" structs:"created" json:"created"`               // 创建时间戳
	Updated  int64  `db:"updated" structs:"updated" json:"updated"`               // 更新时间戳
	Deleted  int64  `db:"deleted" structs:"deleted" json:"deleted"`               // 删除时间戳
}
