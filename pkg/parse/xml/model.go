package xml

import "kitten/pkg/parse"

type nodeInfo struct {
	ProcessCode    string
	Type           string
	Code           string
	Name           string
	CandidateUsers []string
	Properties     []*parse.PropertyResult
	FormResult     *parse.NodeFormResult
}

type sequenceFlow struct {
	ProcessCode string
	XMLName     string
	Code        string
	SourceRef   string
	TargetRef   string
	Explain     string
	Expression  string
}

