package kitten

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/chapin666/kitten/model"
	"github.com/chapin666/kitten/pkg/types"
)

// 定义错误
var (
	ErrNotFound = errors.New("未找到流程相关的信息")
)

// NextNodeHandle 定义下一节点处理函数
type NextNodeHandle func(*model.Node, *model.NodeInstance, []*model.NodeCandidate)

// EndHandle 定义流程结束处理函数
type EndHandle func(*model.FlowInstance)

type nodeRouterOptions struct {
	autoStart  bool
	onNextNode NextNodeHandle
	onFlowEnd  EndHandle
}

// NodeRouterOption 节点路由配置
type NodeRouterOption func(*nodeRouterOptions)

// OnNextNodeOption 注册下一节点处理事件配置
func OnNextNodeOption(fn NextNodeHandle) NodeRouterOption {
	return func(o *nodeRouterOptions) {
		o.onNextNode = fn
	}
}

// OnFlowEndOption 注册流程结束事件
func OnFlowEndOption(fn EndHandle) NodeRouterOption {
	return func(o *nodeRouterOptions) {
		o.onFlowEnd = fn
	}
}

// NodeRouter 节点路由
type NodeRouter struct {
	ctx          context.Context
	parent       *NodeRouter
	engine       *Engine
	node         *model.Node
	inputData    []byte
	opts         *nodeRouterOptions
	flowInstance *model.FlowInstance
	nodeInstance *model.NodeInstance
	stop         bool
}

// Init 初始化节点路由
func (r *NodeRouter) Init(
	ctx context.Context,
	engine *Engine,
	nodeInstanceID string,
	inputData []byte,
	options ...NodeRouterOption,
) (*NodeRouter, error) {
	opts := &nodeRouterOptions{
		autoStart: true,
	}
	for _, opt := range options {
		opt(opts)
	}

	r.ctx = ctx
	r.opts = opts
	r.inputData = inputData
	r.engine = engine

	// nodeInstance
	nodeInstance, err := r.engine.flowSvc.GetNodeInstance(nodeInstanceID)
	if err != nil {
		return nil, err
	}
	if nodeInstance == nil {
		return nil, ErrNotFound
	}
	r.nodeInstance = nodeInstance

	// flowInstance
	flowInstance, err := r.engine.flowSvc.GetFlowInstance(nodeInstance.FlowInstanceID)
	if err != nil {
		return nil, err
	}
	if flowInstance == nil {
		return nil, ErrNotFound
	}
	r.flowInstance = flowInstance

	// 获取node
	node, err := r.engine.flowSvc.GetNode(nodeInstance.NodeID)
	if err != nil {
		return nil, err
	}

	if node == nil {
		return nil, ErrNotFound
	}
	r.node = node

	return r, nil
}

// GetFlowInstance 获取流程实例
func (r *NodeRouter) GetFlowInstance() *model.FlowInstance {
	return r.flowInstance
}

// Next 流向下一节点
func (r *NodeRouter) Next(processor string) error {
	nodeType, err := types.GetNodeTypeByName(r.node.TypeCode)
	if err != nil {
		return err
	}

	if nodeType == types.UserTask && r.parent != nil {
		pNodeType, err := types.GetNodeTypeByName(r.parent.node.TypeCode)
		if err != nil {
			return err
		}

		// 不是开始事件也不是自动开始
		if !(pNodeType == types.StartEvent && r.parent.opts.autoStart) {
			// 通知下一节点实例事件
			if fn := r.opts.onNextNode; fn != nil {
				candidates, err := r.engine.flowSvc.QueryNodeCandidates(r.nodeInstance.RecordID)
				if err != nil {
					return err
				}
				fn(r.node, r.nodeInstance, candidates)
			}
			return nil
		}
	}

	// 完成当前节点
	err = r.engine.flowSvc.DoneNodeInstance(r.nodeInstance.RecordID, processor, r.inputData)
	if err != nil {
		return err
	}

	// 如果当前节点是人工任务，检查下一节点是否是并行网关，如果是则检查还未完成的待办事项，如果有则停止流转
	if nodeType == types.UserTask && r.parent == nil {
		ok, err := r.checkNextNodeType(types.ParallelGateway)
		if err != nil {
			return err
		}

		if ok {
			exists, err := r.engine.flowSvc.CheckFlowInstanceTodo(r.flowInstance.RecordID)
			if err != nil {
				return err
			}
			if exists {
				return nil
			}
		}
	}

	// 如果是结束事件或终止事件，则停止流转
	if nodeType == types.EndEvent || nodeType == types.TerminateEvent {
		isEnd := false

		// 如果是结束事件，则检查还未完成的待办事项，如果没有则结束流程并通知结束事件
		if nodeType == types.EndEvent {
			exists, err := r.engine.flowSvc.CheckFlowInstanceTodo(r.flowInstance.RecordID)
			if err != nil {
				return err
			}

			if !exists {
				isEnd = true
			}
		}

		// 如果是终止事件，则结束流程并通知结束事件
		if nodeType == types.TerminateEvent {
			isEnd = true
		}

		if isEnd {
			// 流程实例结束处理
			err = r.engine.flowSvc.DoneFlowInstance(r.flowInstance.RecordID)
			if err != nil {
				return err
			}

			r.stop = true
			if fn := r.opts.onFlowEnd; fn != nil {
				fn(r.flowInstance)
			}
		}
		return nil
	}

	// 增加下一节点
	nodeInstanceIDs, err := r.addNextNodeInstances()
	if err != nil {
		return err
	}

	for _, instanceID := range nodeInstanceIDs {
		nextRouter, err := r.next(instanceID, processor)
		if err != nil {
			return err
		}
		if nextRouter.stop {
			break
		}
	}

	return nil
}

// 增加下一处理节点实例
func (r *NodeRouter) addNextNodeInstances() ([]string, error) {
	routers, err := r.engine.flowSvc.QueryNodeRouters(r.node.RecordID)
	if err != nil {
		return nil, err
	}
	if len(routers) == 0 {
		return nil, nil
	}

	var nodeInstanceIDs []string
	for _, routerItem := range routers {
		if routerItem.Expression != "" {
			allow, err := r.engine.execer.ExecReturnBool([]byte(routerItem.Expression), r.getExpData())
			if err != nil {
				return nil, err
			}
			if !allow {
				continue
			}
		}

		// 查询指派人表达式
		assigns, err := r.engine.flowSvc.QueryNodeAssignments(routerItem.TargetNodeID)
		if err != nil {
			return nil, err
		}

		var candidates []string
		for _, assign := range assigns {
			ss, err := r.engine.execer.ExecReturnStringSlice([]byte(assign.Expression), r.getExpData())
			if err != nil {
				return nil, err
			}
			candidates = append(candidates, ss...)

		}

		instanceID, err := r.engine.flowSvc.CreateNodeInstance(
			r.flowInstance.RecordID,
			routerItem.TargetNodeID,
			r.inputData,
			candidates,
		)
		if err != nil {
			return nil, err
		}
		nodeInstanceIDs = append(nodeInstanceIDs, instanceID)
	}
	return nodeInstanceIDs, nil
}

func (r *NodeRouter) next(nodeInstanceID, processor string) (*NodeRouter, error) {
	nextRouter, err := new(NodeRouter).Init(r.ctx, r.engine, nodeInstanceID, r.inputData)
	if err != nil {
		return nil, err
	}
	nextRouter.opts = r.opts
	nextRouter.parent = r

	err = nextRouter.Next(processor)
	if err != nil {
		return nil, err
	}
	return nextRouter, nil
}

// 检查下一节点类型
func (r *NodeRouter) checkNextNodeType(t types.NodeType) (bool, error) {
	routers, err := r.engine.flowSvc.QueryNodeRouters(r.node.RecordID)
	if err != nil {
		return false, err
	}
	if len(routers) == 0 {
		return false, nil
	}

	for _, routerItem := range routers {
		if routerItem.Expression != "" {
			allow, err := r.engine.execer.ExecReturnBool([]byte(routerItem.Expression), r.getExpData())
			if err != nil {
				return false, err
			}

			if !allow {
				continue
			}
		}

		node, err := r.engine.flowSvc.GetNode(routerItem.TargetNodeID)
		if err != nil {
			return false, err
		}

		if node == nil {
			return false, nil
		}

		if node.TypeCode == t.String() {
			return true, nil
		}
	}

	return false, nil
}

// 获取表达式数据
func (r *NodeRouter) getExpData() []byte {
	var input map[string]interface{}
	json.Unmarshal(r.inputData, &input)

	expData := map[string]interface{}{
		"input": input,
		"flow":  r.flowInstance,
		"node":  r.nodeInstance,
	}
	b, _ := json.Marshal(expData)
	return b
}
