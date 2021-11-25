package service

import (
	"goworkflow/repository"
	"sync"
)

// Flow 流程管理
type Flow struct {
	sync.RWMutex
	FlowModel *repository.Flow `inject:""`
}