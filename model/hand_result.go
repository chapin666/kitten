package model

import "encoding/json"

// NextNode 下一节点
type NextNode struct {
	Node         *Node         // 节点信息
	CandidateIDs []string      // 节点候选人
	NodeInstance *NodeInstance // 节点实例
}

// HandleResult 处理结果
type HandleResult struct {
	IsEnd        bool          `json:"is_end"`        // 是否结束
	NextNodes    []*NextNode   `json:"next_nodes"`    // 下一处理节点
	FlowInstance *FlowInstance `json:"flow_instance"` // 流程实例
}

func (r *HandleResult) String() string {
	b, _ := json.Marshal(r)
	return string(b)
}

