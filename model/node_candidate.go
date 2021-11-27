package model


// NodeCandidate 节点候选人
type NodeCandidate struct {
	ID             int64  `db:"id,primarykey,autoincrement" structs:"id" json:"id"`                          // 唯一标识(自增ID)
	RecordID       string `db:"record_id,size:36" structs:"record_id" json:"record_id"`                      // 记录内码(uuid)
	NodeInstanceID string `db:"node_instance_id,size:36" structs:"node_instance_id" json:"node_instance_id"` // 节点实例内码
	CandidateID    string `db:"candidate_id,size:36" structs:"candidate_id" json:"candidate_id"`             // 候选人ID(根据节点指派表达式生成)
	Created        int64  `db:"created" structs:"created" json:"created"`                                    // 创建时间戳
	Updated        int64  `db:"updated" structs:"updated" json:"updated"`                                    // 更新时间戳
	Deleted        int64  `db:"deleted" structs:"deleted" json:"deleted"`                                    // 删除时间戳
}
