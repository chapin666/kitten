package goworkflow

import (
	"goworkflow/pkg/db"
	"testing"
)

func TestInit(t *testing.T) {
	dnsOption := db.SetDSN("root@tcp(127.0.0.1:3306)/flow_test?charset=utf8")
	traceOption := db.SetTrace(true)
	Init(dnsOption, traceOption)
}
