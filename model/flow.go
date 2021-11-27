package model

// Flow 流程
type Flow struct {
	ID       int64  `db:"id,primarykey,autoincrement" structs:"id" json:"id"`     // 唯一标识(自增ID)
	RecordID string `db:"record_id,size:36" structs:"record_id" json:"record_id"` // 记录内码(uuid)
	Code     string `db:"code,size:50" structs:"code" json:"code"`                // 流程编号
	Name     string `db:"name,size:50" structs:"name" json:"name"`                // 流程名称
	Version  int64  `db:"version" structs:"version" json:"version"`               // 版本号
	TypeCode string `db:"type_code,size:50" structs:"type_code" json:"type_code"` // 流程类型编号
	XML      string `db:"xml,size:1024" structs:"xml" json:"xml"`                 // XML数据
	Memo     string `db:"memo,size:255" structs:"memo" json:"memo"`               // 流程备注
	Flag     int64  `db:"flag" structs:"flag" json:"flag"`                        // 流程标志(1:主流程 2:子流程)
	ParentID string `db:"parent_id,size:36" structs:"parent_id" json:"parent_id"` // 父级流程内码
	Status   int    `db:"status" structs:"status" json:"status"`                  // 流程状态(1:正常 2:禁用)
	Created  int64  `db:"created" structs:"created" json:"created"`               // 创建时间戳
	Updated  int64  `db:"updated" structs:"updated" json:"updated"`               // 更新时间戳
	Deleted  int64  `db:"deleted" structs:"deleted" json:"deleted"`               // 删除时间戳
}
