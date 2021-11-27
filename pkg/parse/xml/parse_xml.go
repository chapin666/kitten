package xml

import (
	"context"
	"strconv"

	"goworkflow/pkg/parse"
	"goworkflow/pkg/types"
	"goworkflow/pkg/util"

	"github.com/beevik/etree"
)

type xmlParser struct {
}

// NewXMLParser xml解析器
func NewXMLParser() parse.Parser {
	return &xmlParser{}
}

func (p *xmlParser) Parse(ctx context.Context, content []byte) (*parse.ParseResult, error) {
	result := &parse.ParseResult{
		FlowStatus: 2,
	}
	var err error

	doc := etree.NewDocument()
	if err = doc.ReadFromBytes(content); err != nil {
		panic(err)
	}

	// definitions
	root := doc.SelectElement("definitions")

	// process
	//id：流程定义ID，代表该流程的唯一性，启动该流程时需要使用该ID
	//isExecutable：表示该流程是否可执行，其值有true和false，默认为true
	//name：流程名称
	//type：流程类型
	//isClosed：流程是否已关闭,关闭不能执行
	//versionTag：版本号
	process := root.SelectElement("process")

	// process id property
	if id := process.SelectAttr("id"); id != nil {
		result.FlowID = id.Value
	}
	// process name property
	if name := process.SelectAttr("name"); name != nil {
		result.FlowName = name.Value
	}
	// process isExecutable property
	if v := process.SelectAttr("isExecutable"); v != nil {
		b, _ := strconv.ParseBool(v.Value)
		if b {
			result.FlowStatus = 1
		}
	}
	// process versionTag property
	if version := process.SelectAttr("versionTag"); version != nil {
		result.FlowVersion, err = util.StringToInt(version.Value)
		if err != nil {
			return nil, err
		}
	}


	// 解析节点
	// 定义一个用于辅助的 map，由节点 id 映射到 NodeResult
	nodeMap := make(map[string]*parse.NodeResult)
	// 遍历找到所有的节点，因为是解析一个树，所以先解析节点，再解析sequenceFlow部分
	for _, element := range process.ChildElements() {
		if element.Tag == "documentation" ||
			element.Tag == "extensionElements" ||
			element.Tag == "sequenceFlow" {
			continue
		}
		node, _ := p.ParseNode(element)
		var nodeResult parse.NodeResult
		nodeResult.NodeID = node.Code
		nodeResult.NodeName = node.Name
		nodeResult.NodeType, err = types.GetNodeTypeByName(node.Type)
		if err != nil {
			return nil, err
		}
		nodeResult.CandidateExpressions = node.CandidateUsers
		nodeResult.FormResult = node.FormResult
		nodeResult.Properties = node.Properties
		nodeMap[nodeResult.NodeID] = &nodeResult
		// 如果节点是一个路由的话，需要特殊处理
	}

	// 解析sequenceFlow
	// 解析sequenceFlow部分时，nodeMap里面应该已经有对应的nodeId了
	for _, element := range process.ChildElements() {
		if element.Tag == "sequenceFlow" {
			sFlow, _ := p.ParseSequenceFlow(element)
			var routerResult parse.RouterResult
			routerResult.Expression = sFlow.Expression
			routerResult.Explain = sFlow.Explain
			routerResult.TargetNodeID = sFlow.TargetRef
			if nodeResult, exist := nodeMap[sFlow.SourceRef]; exist {
				nodeResult.Routers = append(nodeResult.Routers, &routerResult)
			}
		}
	}

	for _, nodeResult := range nodeMap {
		result.Nodes = append(result.Nodes, nodeResult)
	}
	return result, nil
}

func (p *xmlParser) ParseNode(element *etree.Element) (*nodeInfo, error) {
	var node nodeInfo

	node.Type = element.Tag
	if node.Type == "endEvent" {
		for _, e := range element.ChildElements() {
			if e.Tag == "terminateEventDefinition" {
				node.Type = "terminateEvent"
			}
		}
	}
	if name := element.SelectAttr("name"); name != nil {
		node.Name = name.Value
	}
	if id := element.SelectAttr("id"); id != nil {
		node.Code = id.Value
	}
	if candidateUsers := element.SelectAttr("candidateUsers"); candidateUsers != nil {
		node.CandidateUsers = []string{candidateUsers.Value}
	}

	nodeFormResult := new(parse.NodeFormResult)
	if formKey := element.SelectAttr("formKey"); formKey != nil {
		nodeFormResult.ID = formKey.Value
	}

	if extensionElements := element.SelectElement("extensionElements"); extensionElements != nil {
		if formData := extensionElements.SelectElement("formData"); formData != nil {
			form, err := p.ParseFormData(formData)
			if err != nil {
				return nil, err
			}
			if form != nil {
				nodeFormResult.Fields = form.Fields
			}
		}

		if propertyData := extensionElements.SelectElement("properties"); propertyData != nil {
			// 解析节点属性
			for _, p := range propertyData.SelectElements("property") {
				var item parse.PropertyResult
				if name := p.SelectAttr("name"); name != nil {
					item.Name = name.Value
				}
				if value := p.SelectAttr("value"); value != nil {
					item.Value = value.Value
				}
				if item.Name != "" {
					node.Properties = append(node.Properties, &item)
				}
			}
		}
	}
	node.FormResult = nodeFormResult

	return &node, nil
}

func (p *xmlParser) ParseSequenceFlow(element *etree.Element) (*sequenceFlow, error) {
	hasExpression := false
	var seq sequenceFlow
	seq.XMLName = element.Tag
	seq.Code = element.SelectAttr("id").Value
	seq.SourceRef = element.SelectAttr("sourceRef").Value
	seq.TargetRef = element.SelectAttr("targetRef").Value
	for _, childEle := range element.ChildElements() {
		if childEle.Tag == "documentation" {
			seq.Explain = childEle.Text()
		} else if childEle.Tag == "conditionExpression" {
			seq.Expression = childEle.Text()
			hasExpression = true
		}
	}
	if !hasExpression {
		seq.Expression = ""
	}
	return &seq, nil
}

func (p *xmlParser) ParseFormData(element *etree.Element) (*parse.NodeFormResult, error) {
	var formResult = &parse.NodeFormResult{}
	if id := element.SelectAttr("id"); id != nil {
		formResult.ID = id.Value
	}

	if fieldList := element.SelectElements("formField"); fieldList != nil {
		for _, item := range fieldList {
			var field = &parse.FormFieldResult{}
			var err error
			if properties := item.SelectElement("properties"); properties != nil {
				field.Properties, err = p.ParseProperties(properties)
				if err != nil {
					return nil, err
				}
			}
			if validations := item.SelectElement("validation"); validations != nil {
				field.Validations, err = p.ParseValidations(validations)
				if err != nil {
					return nil, err
				}
			}
			if nodeType := item.SelectAttr("type"); nodeType != nil {
				field.Type = nodeType.Value
				if field.Type == "enum" {
					field.Values, err = p.ParseEnumValues(item)
					if err != nil {
						return nil, err
					}
				}
			}
			if id := item.SelectAttr("id"); id != nil {
				field.ID = id.Value
			}
			if label := item.SelectAttr("label"); label != nil {
				field.Label = label.Value
			}
			if defaultValue := item.SelectAttr("defaultValue"); defaultValue != nil {
				field.DefaultValue = defaultValue.Value
			}
			formResult.Fields = append(formResult.Fields, field)
		}
	}

	return formResult, nil
}

func (p *xmlParser) ParseProperties(element *etree.Element) ([]*parse.FieldProperty, error) {
	var properties = make([]*parse.FieldProperty, 0)
	if propertyList := element.SelectElements("property"); propertyList != nil {
		for _, item := range propertyList {
			var property = &parse.FieldProperty{}
			if id := item.SelectAttr("id"); id != nil {
				property.ID = id.Value
			}
			if value := item.SelectAttr("value"); value != nil {
				property.Value = value.Value
			}
			properties = append(properties, property)
		}
	}
	return properties, nil
}

func (p *xmlParser) ParseValidations(element *etree.Element) ([]*parse.FieldValidation, error) {
	var validations = make([]*parse.FieldValidation, 0)
	if validationList := element.SelectElements("constraint"); validationList != nil {
		for _, item := range validationList {
			var validation = &parse.FieldValidation{}
			if name := item.SelectAttr("name"); name != nil {
				validation.Name = name.Value
			}
			if config := item.SelectAttr("config"); config != nil {
				validation.Config = config.Value
			}
			validations = append(validations, validation)
		}
	}
	return validations, nil
}

func (p *xmlParser) ParseEnumValues(element *etree.Element) ([]*parse.FieldOption, error) {
	var options = make([]*parse.FieldOption, 0)
	if optionList := element.SelectElements("value"); optionList != nil {
		for _, item := range optionList {
			var option = &parse.FieldOption{}
			if id := item.SelectAttr("id"); id != nil {
				option.ID = id.Value
			}
			if name := item.SelectAttr("name"); name != nil {
				option.Name = name.Value
			}
			options = append(options, option)
		}
	}
	return options, nil
}

