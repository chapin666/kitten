package repository

import (
	"database/sql"
	"fmt"
	"github.com/pkg/errors"
	"goworkflow/pkg/db"
	"goworkflow/model"
)

// Flow 流程管理
type Flow struct {
	DB *db.DB `inject:""`
}

// CreateFlow 创建流程数据
func (f *Flow) CreateFlow(flow *model.Flow, nodes *model.NodeOperating, forms *model.FormOperating) error {
	tran, err := f.DB.Begin()
	if err != nil {
		return errors.Wrapf(err, "创建流程基础数据开启事物发生错误")
	}

	// 写入数据到flow表
	err = tran.Insert(flow)
	if err != nil {
		_ = tran.Rollback()
		return errors.Wrapf(err, "插入流程数据发生错误")
	}

	// 写入数据到 NodeGroup、RouterGroup、AssignmentGroup、PropertyGroup
	if list := nodes.All(); len(list) > 0 {
		err = tran.Insert(list...)
		if err != nil {
			_ = tran.Rollback()
			return errors.Wrapf(err, "插入节点数据发生错误")
		}
	}

	// 写入数据到 FormGroup、FormFieldGroup、FieldOptionGroup、FieldPropertyGroup、FieldValidationGroup
	if list := forms.All(); len(list) > 0 {
		err = tran.Insert(list...)
		if err != nil {
			_ = tran.Rollback()
			return errors.Wrapf(err, "插入表单数据发生错误")
		}
	}

	err = tran.Commit()
	if err != nil {
		return errors.Wrapf(err, "创建流程基础数据提交事物发生错误")
	}
	return nil
}


// GetFlowByCode 根据编号查询流程数据
func (f *Flow) GetFlowByCode(code string) (*model.Flow, error) {
	query := fmt.Sprintf("SELECT * FROM %s WHERE deleted=0 AND flag=1 AND status=1 AND code=? "+
		"ORDER BY version DESC LIMIT 1", model.FlowTableName)

	var flow model.Flow
	err := f.DB.SelectOne(&flow, query, code)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, errors.Wrapf(err, "根据编号查询流程数据发生错误")
	}

	return &flow, nil
}


