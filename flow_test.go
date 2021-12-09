package kitten

import (
	"context"
	"kitten/model"
	"kitten/pkg/db"
	"os"
	"testing"
)

var (
	dnsOption   = db.SetDSN("root@tcp(127.0.0.1:3306)/flow_test?charset=utf8")
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

func TestStartFlow(t *testing.T) {
	var flowCode = "process_leave_test"
	var nodeCode = "node_start"
	var userID = "F001"
	var input = map[string]interface{}{
		"day": 1,
		"bzr": "F002",
	}
	result, err := StartFlow(context.Background(), flowCode, nodeCode, userID, input)
	if err != nil {
		t.Errorf("start flow define failed: %s", err.Error())
	}
	t.Log(result)
}

func TestQueryTodoFlows(t *testing.T) {
	flowCode := "process_leave_test"
	userID := "F002"
	limit := 100
	todos, err := QueryTodoFlows(flowCode, userID, limit)
	if err != nil {
		t.Fatalf("query flow failed: %s", err.Error())
	}

	for _, todo := range todos {
		t.Logf("%#v", todo)
	}
}

func TestQueryNodeCandidates(t *testing.T) {
	nodeInstanceID := "164f4a70-6d60-4447-b332-bfa8af875676"
	userIDs, err := QueryNodeCandidates(nodeInstanceID)
	if err != nil {
		t.Errorf("query node candidate failed: %s", err.Error())
	}
	t.Log(userIDs)
}

func TestHandleFlow(t *testing.T) {
	nodeInstanceID := "164f4a70-6d60-4447-b332-bfa8af875676"
	userID := "F002"
	input := map[string]interface{}{
		"action": "pass",
	}
	result, err := HandleFlow(context.Background(), nodeInstanceID, userID, input)
	if err != nil {
		t.Errorf("hanle flow failed: %s", err.Error())
	}
	t.Log(result)
}

func TestQueryDoneFlowIDs(t *testing.T) {
	flowCode := "process_leave_test"
	userID := "T002"
	ids, err := QueryDoneFlowIDs(flowCode, userID)
	if err != nil {
		t.Errorf("query done flow ids failed: %s", err.Error())
	}
	t.Log(ids)
}


func TestStopFlowInstance(t *testing.T) {
	nodeInstanceID := "4c66bea5-01fa-463f-8da5-bedb290e419e"
	err := StopFlowInstance(nodeInstanceID, func(instance *model.FlowInstance) bool {
		return true
	})
	if err != nil {
		t.Errorf("hanle flow failed: %s", err.Error())
	}
}
