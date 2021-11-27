package goworkflow

import (
	"goworkflow/pkg/db"
	"os"
	"testing"
)

var (
	dnsOption = db.SetDSN("root@tcp(127.0.0.1:3306)/flow_test?charset=utf8")
	traceOption = db.SetTrace(true)
)

func TestMain(m *testing.M) {
	Init(dnsOption, traceOption)
	os.Exit(m.Run())
}

func TestDeploy(t *testing.T) {
	result, err := Deploy("./test_data/leave.xml")
	if err != nil {
		t.Errorf("deploy flow define failed: %s", err.Error())
	}
	t.Log(result)
}