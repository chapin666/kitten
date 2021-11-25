package repository

import "goworkflow/pkg/db"

// Flow 流程管理
type Flow struct {
	DB *db.DB `inject:""`
}

