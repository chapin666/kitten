package model


// FlowQueryResult 流程查询结果
type FlowQueryResult struct {
	ID       int64  `db:"id,primarykey,autoincrement" structs:"id" json:"id"`     // 唯一标识(自增ID)
	RecordID string `db:"record_id,size:36" structs:"record_id" json:"record_id"` // 记录内码(uuid)
	Code     string `db:"code,size:50" structs:"code" json:"code"`                // 流程编号(唯一)
	Name     string `db:"name,size:50" structs:"name" json:"name"`                // 流程名称
	Version  int64  `db:"version" structs:"version" json:"version"`               // 版本号
	TypeCode string `db:"type_code,size:50" structs:"type_code" json:"type_code"` // 流程类型编号
	Status   int    `db:"status" structs:"status" json:"status"`                  // 流程状态(1:正常 2:禁用)
	Created  int64  `db:"created" structs:"created" json:"created"`               // 创建时间戳
	Memo     string `db:"memo,size:255" structs:"memo" json:"memo"`               // 流程备注
}

// FlowQueryParam 流程查询参数
type FlowQueryParam struct {
	Code     string // 流程编号
	Name     string // 流程名称
	TypeCode string // 流程类型编号
	Status   int    // 流程状态(1:正常 2:禁用)
}
