package model

// NodeProperty 节点属性
type NodeProperty struct {
	ID       int64  `db:"id,primarykey,autoincrement" structs:"id" json:"id"`     // 唯一标识(自增ID)
	RecordID string `db:"record_id,size:36" structs:"record_id" json:"record_id"` // 记录内码(uuid)
	NodeID   string `db:"node_id,size:36" structs:"node_id" json:"node_id"`       // 节点内码
	Name     string `db:"name,size:50" structs:"name" json:"name"`                // 属性名称
	Value    string `db:"value,size:255" structs:"value" json:"value"`            // 属性值
	Created  int64  `db:"created" structs:"created" json:"created"`               // 创建时间戳
	Updated  int64  `db:"updated" structs:"updated" json:"updated"`               // 更新时间戳
	Deleted  int64  `db:"deleted" structs:"deleted" json:"deleted"`               // 删除时间戳
}
