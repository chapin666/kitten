package model

// NodeAssignment 节点指派
type NodeAssignment struct {
	ID         int64  `db:"id,primarykey,autoincrement" structs:"id" json:"id"`          // 唯一标识(自增ID)
	RecordID   string `db:"record_id,size:36" structs:"record_id" json:"record_id"`      // 记录内码(uuid)
	NodeID     string `db:"node_id,size:36" structs:"node_id" json:"node_id"`            // 节点内码
	Expression string `db:"expression,size:1024" structs:"expression" json:"expression"` // 执行表达式(基于qlang可提供多种内置函数支持，支持SQL查询)
	Created    int64  `db:"created" structs:"created" json:"created"`                    // 创建时间戳
	Updated    int64  `db:"updated" structs:"updated" json:"updated"`                    // 更新时间戳
	Deleted    int64  `db:"deleted" structs:"deleted" json:"deleted"`                    // 删除时间戳
}
