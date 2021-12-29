package kitten

import (
	"context"
	"encoding/json"
	"github.com/chapin666/kitten/model"
	"os"
	"testing"
)

var (
	client   *Engine
	err      error
	mysqlDNS = "root@tcp(127.0.0.1:3306)/flow_test?charset=utf8"
)

func TestMain(m *testing.M) {
	client, err = New(mysqlDNS, true)
	if err != nil {
		panic(err)
	}
	os.Exit(m.Run())
}

func TestDeploy(t *testing.T) {
	result, err := client.Deploy("./test_data/leave.xml")
	if err != nil {
		t.Errorf("deploy flow define failed: %s", err.Error())
	}
	t.Log(result)
}

func TestQueryAll(t *testing.T) {
	params := model.FlowQueryParam{
	}
	total, result, err := client.QueryAllFlowPage(params, 1, 10)
	if err != nil {
		t.Errorf("query all flow page failed: %s", err.Error())
	}

	t.Logf("total=%d", total)
	for _, item := range result {
		t.Logf("%#v", item)
	}
}

func TestGetFlow(t *testing.T) {
	flow, err := client.GetFlow("7099d58c-8df2-4781-bace-8a439d3bae9c")
	if err != nil {
		t.Errorf("get flow  failed: %s", err.Error())
	}

	t.Logf("%#v", flow)
}

func TestStartFlow(t *testing.T) {
	flowCode := "process_leave_test"
	nodeCode := "node_start"
	userID := "F001"
	input, _ := json.Marshal(map[string]interface{}{
		"day": 1,
		"bzr": "F002",
	})
	result, err := client.StartFlow(context.Background(), flowCode, nodeCode, userID, input)
	if err != nil {
		t.Errorf("start flow define failed: %s", err.Error())
	}
	t.Log(result)
}

func TestQueryTodoFlows(t *testing.T) {
	flowCode := "process_leave_test"
	userID := "F002"
	limit := 100
	todos, err := client.QueryTodoFlows(flowCode, userID, limit)
	if err != nil {
		t.Fatalf("query flow failed: %s", err.Error())
	}

	for _, todo := range todos {
		t.Logf("%#v", todo)
	}
}

func TestQueryNodeCandidates(t *testing.T) {
	nodeInstanceID := "164f4a70-6d60-4447-b332-bfa8af875676"
	userIDs, err := client.QueryNodeCandidates(nodeInstanceID)
	if err != nil {
		t.Errorf("query node candidate failed: %s", err.Error())
	}
	t.Log(userIDs)
}

func TestHandleFlow(t *testing.T) {
	nodeInstanceID := "164f4a70-6d60-4447-b332-bfa8af875676"
	userID := "F002"
	input, _ := json.Marshal(map[string]interface{}{
		"action": "pass",
	})
	result, err := client.HandleFlow(context.Background(), nodeInstanceID, userID, input)
	if err != nil {
		t.Errorf("hanle flow failed: %s", err.Error())
	}
	t.Log(result)
}

func TestQueryDoneFlowIDs(t *testing.T) {
	flowCode := "process_leave_test"
	userID := "T002"
	ids, err := client.QueryDoneFlowIDs(flowCode, userID)
	if err != nil {
		t.Errorf("query done flow ids failed: %s", err.Error())
	}
	t.Log(ids)
}

func TestStopFlowInstance(t *testing.T) {
	nodeInstanceID := "4c66bea5-01fa-463f-8da5-bedb290e419e"
	err := client.StopFlowInstance(nodeInstanceID, func(instance *model.FlowInstance) bool {
		return true
	})
	if err != nil {
		t.Errorf("hanle flow failed: %s", err.Error())
	}
}
