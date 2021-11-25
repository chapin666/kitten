package goworkflow

import (
	"database/sql"
	"github.com/facebookgo/inject"
	"goworkflow/mapper"
	"goworkflow/pkg/db"
	"goworkflow/pkg/expression"
	"goworkflow/pkg/parse"
	"goworkflow/service"
)

// Engine .
type Engine struct {
	parser  parse.Parser
	execer  expression.Execer
	flowSvc *service.Flow
}

// Init .
func (e *Engine) Init(parser parse.Parser, execer expression.Execer, sqlDB *sql.DB, trace bool) (*Engine, error) {
	var g inject.Graph
	var flowSvc service.Flow

	dbInstance := db.NewMySQLWithDB(sqlDB, trace)

	err := g.Provide(&inject.Object{Value: dbInstance}, &inject.Object{Value: &flowSvc})
	if err != nil {
		return e, err
	}

	err = g.Populate()
	if err != nil {
		return e, err
	}

	mapper.FlowDBMap(dbInstance)
	err = dbInstance.CreateTablesIfNotExists()
	if err != nil {
		return e, err
	}

	e.parser = parser
	e.execer = execer
	e.flowSvc = &flowSvc

	return e, nil
}
