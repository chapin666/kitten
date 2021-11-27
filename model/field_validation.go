package model

// FieldValidation 字段校验
type FieldValidation struct {
	ID               int64  `db:"id,primarykey,autoincrement" structs:"id" json:"id"`                              // 唯一标识(自增ID)
	RecordID         string `db:"record_id,size:36" structs:"record_id" json:"record_id"`                          // 记录内码(uuid)
	FieldID          string `db:"field_id,size:36" structs:"field_id" json:"field_id"`                             // 字段内码
	ConstraintName   string `db:"constraint_name,size:50" structs:"constraint_name" json:"constraint_name"`        // 约束名称
	ConstraintConfig string `db:"constraint_config,size:100" structs:"constraint_config" json:"constraint_config"` // 约束配置
	Created          int64  `db:"created" structs:"created" json:"created"`                                        // 创建时间戳
	Updated          int64  `db:"updated" structs:"updated" json:"updated"`                                        // 更新时间戳
	Deleted          int64  `db:"deleted" structs:"deleted" json:"deleted"`                                        // 删除时间戳
}
