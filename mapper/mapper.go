package mapper

import (
	"kitten/model"
	"kitten/pkg/db"
)

// FlowDBMap 注册流程相关的数据库映射
func FlowDBMap(dbInstance *db.DB) {
	dbInstance.AddTableWithName(model.Flow{}, model.FlowTableName)
	dbInstance.AddTableWithName(model.Node{}, model.NodeTableName)
	dbInstance.AddTableWithName(model.NodeRouter{}, model.NodeRouterTableName)
	dbInstance.AddTableWithName(model.NodeAssignment{}, model.NodeAssignmentTableName)
	dbInstance.AddTableWithName(model.FlowInstance{}, model.FlowInstanceTableName)
	dbInstance.AddTableWithName(model.NodeInstance{}, model.NodeInstanceTableName)
	dbInstance.AddTableWithName(model.NodeTiming{}, model.NodeTimingTableName)
	dbInstance.AddTableWithName(model.NodeCandidate{}, model.NodeCandidateTableName)
	dbInstance.AddTableWithName(model.Form{}, model.FormTableName)
	dbInstance.AddTableWithName(model.FormField{}, model.FormFieldTableName)
	dbInstance.AddTableWithName(model.FieldOption{}, model.FieldOptionTableName)
	dbInstance.AddTableWithName(model.FieldProperty{}, model.FieldPropertyTableName)
	dbInstance.AddTableWithName(model.FieldValidation{}, model.FieldValidationTableName)
	dbInstance.AddTableWithName(model.NodeProperty{}, model.NodePropertyTableName)
}
