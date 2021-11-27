package model

// NodeRouter 节点路由
type NodeRouter struct {
	ID              int64  `db:"id,primarykey,autoincrement" structs:"id" json:"id"`                     // 唯一标识(自增ID)
	RecordID        string `db:"record_id,size:36" structs:"record_id" json:"record_id"`                 // 记录内码(uuid)
	SourceNodeID    string `db:"source_node_id,size:36" structs:"source_node_id" json:"source_node_id"`  // 源节点内码
	TargetNodeID    string `db:"target_node_id,size:36" structs:"target_node_id" json:"target_node_id"`  // 目标节点内码
	Expression      string `db:"expression,size:1024" structs:"expression" json:"expression"`            // 条件表达式(使用qlang作为表达式脚本语言(返回值bool))
	Explain         string `db:"explain,size:255" structs:"explain" json:"explain"`                      // 说明
	IsDefaultTarget int64  `db:"is_default_target" structs:"is_default_target" json:"is_default_target"` // 是否是默认节点(1:是 2:否)
	Created         int64  `db:"created" structs:"created" json:"created"`                               // 创建时间戳
	Updated         int64  `db:"updated" structs:"updated" json:"updated"`                               // 更新时间戳
	Deleted         int64  `db:"deleted" structs:"deleted" json:"deleted"`                               // 删除时间戳
}
