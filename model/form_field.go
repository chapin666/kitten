package model

// FormField 表单字段
type FormField struct {
	ID           int64  `db:"id,primarykey,autoincrement" structs:"id" json:"id"`                  // 唯一标识(自增ID)
	RecordID     string `db:"record_id,size:36" structs:"record_id" json:"record_id"`              // 记录内码(uuid)
	FormID       string `db:"form_id,size:36" structs:"form_id" json:"form_id"`                    // 表单内码
	Code         string `db:"code,size:50" structs:"code" json:"code"`                             // 字段编号
	Label        string `db:"label,size:50" structs:"label" json:"label"`                          // 字段标签
	TypeCode     string `db:"type_code,size:50" structs:"type_code" json:"type_code"`              // 字段类型
	DefaultValue string `db:"default_value,size:100" structs:"default_value" json:"default_value"` // 字段默认值
	Created      int64  `db:"created" structs:"created" json:"created"`                            // 创建时间戳
	Updated      int64  `db:"updated" structs:"updated" json:"updated"`                            // 更新时间戳
	Deleted      int64  `db:"deleted" structs:"deleted" json:"deleted"`                            // 删除时间戳
}
