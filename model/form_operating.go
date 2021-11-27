package model

// FormOperating 表单操作
type FormOperating struct {
	FormGroup            []*Form
	FormFieldGroup       []*FormField
	FieldOptionGroup     []*FieldOption
	FieldPropertyGroup   []*FieldProperty
	FieldValidationGroup []*FieldValidation
}

// All 获取所有表单操作的组
func (a *FormOperating) All() []interface{} {
	var group []interface{}

	for _, item := range a.FormGroup {
		group = append(group, item)
	}
	for _, item := range a.FormFieldGroup {
		group = append(group, item)
	}
	for _, item := range a.FieldOptionGroup {
		group = append(group, item)
	}
	for _, item := range a.FieldPropertyGroup {
		group = append(group, item)
	}
	for _, item := range a.FieldValidationGroup {
		group = append(group, item)
	}

	return group
}

