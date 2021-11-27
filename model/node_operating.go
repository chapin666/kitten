package model

// NodeOperating 节点操作
type NodeOperating struct {
	NodeGroup       []*Node
	RouterGroup     []*NodeRouter
	AssignmentGroup []*NodeAssignment
	PropertyGroup   []*NodeProperty
}

// All 获取所有节点操作的组
func (a *NodeOperating) All() []interface{} {
	var group []interface{}

	for _, item := range a.NodeGroup {
		group = append(group, item)
	}
	for _, item := range a.RouterGroup {
		group = append(group, item)
	}
	for _, item := range a.AssignmentGroup {
		group = append(group, item)
	}
	for _, item := range a.PropertyGroup {
		group = append(group, item)
	}

	return group
}