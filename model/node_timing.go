package model

// NodeTiming 节点定时
type NodeTiming struct {
	ID             int64  `db:"id,primarykey,autoincrement" structs:"id" json:"id"`                  // 唯一标识(自增ID)
	NodeInstanceID string `db:"node_instance_id" structs:"node_instance_id" json:"node_instance_id"` // 节点实例ID
	Flag           string `db:"flag" structs:"flag" json:"flag"`                                     // 标志
	Processor      string `db:"processor,size:36" structs:"processor" json:"processor"`              // 处理人
	Input          string `db:"input,size:1024" structs:"input" json:"input"`                        // 输入数据
	ExpiredAt      int64  `db:"expired_at" structs:"expired_at" json:"expired_at"`                   // 过期时间戳
	Created        int64  `db:"created" structs:"created" json:"created"`                            // 创建时间戳
	Deleted        int64  `db:"deleted" structs:"deleted" json:"deleted"`                            // 删除时间戳
}