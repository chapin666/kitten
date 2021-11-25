package xml

import (
	"context"
	"encoding/json"
	"fmt"
	"testing"
)

const (
	basic = `<?xml version="1.0" encoding="UTF-8"?>
	<bpmn:definitions xmlns:bpmn="http://www.omg.org/spec/BPMN/20100524/MODEL" xmlns:bpmndi="http://www.omg.org/spec/BPMN/20100524/DI" xmlns:di="http://www.omg.org/spec/DD/20100524/DI" xmlns:dc="http://www.omg.org/spec/DD/20100524/DC" xmlns:xsi="http://www.w3.org/2001/XMLSchema-instance" id="Definitions_1" targetNamespace="http://bpmn.io/schema/bpmn" exporter="Camunda Modeler" exporterVersion="1.11.3">
	  <bpmn:process id="Process_1" name="很牛逼的流程一" isExecutable="true">
		<bpmn:startEvent id="StartEvent_1" name="节点一">
		  <bpmn:outgoing>fdf</bpmn:outgoing>
		</bpmn:startEvent>
		<bpmn:userTask id="Task_0yrbi3a">
		  <bpmn:incoming>fdf</bpmn:incoming>
		  <bpmn:outgoing>SequenceFlow_1r953n5</bpmn:outgoing>
		</bpmn:userTask>
		<bpmn:endEvent id="EndEvent_1lxbe5s" name="结束节点">
		  <bpmn:documentation>结束节点的文档</bpmn:documentation>
		  <bpmn:incoming>SequenceFlow_1r953n5</bpmn:incoming>
		</bpmn:endEvent>
		<bpmn:sequenceFlow id="fdf" name="连线一" sourceRef="StartEvent_1" targetRef="Task_0yrbi3a">
		  <bpmn:documentation>连线一的文档</bpmn:documentation>
		</bpmn:sequenceFlow>
		<bpmn:sequenceFlow id="SequenceFlow_1r953n5" name="连线二" sourceRef="Task_0yrbi3a" targetRef="EndEvent_1lxbe5s">
		  <bpmn:conditionExpression xsi:type="bpmn:tFormalExpression">这是一段测试用的expression</bpmn:conditionExpression>
		</bpmn:sequenceFlow>
	  </bpmn:process>
	  <bpmndi:BPMNDiagram id="BPMNDiagram_1">
		<bpmndi:BPMNPlane id="BPMNPlane_1" bpmnElement="Process_1">
		  <bpmndi:BPMNShape id="_BPMNShape_StartEvent_2" bpmnElement="StartEvent_1">
			<dc:Bounds x="346" y="447" width="36" height="36" />
			<bpmndi:BPMNLabel>
			  <dc:Bounds x="348" y="483" width="35" height="13" />
			</bpmndi:BPMNLabel>
		  </bpmndi:BPMNShape>
		  <bpmndi:BPMNShape id="UserTask_1xub0ug_di" bpmnElement="Task_0yrbi3a">
			<dc:Bounds x="520" y="425" width="100" height="80" />
		  </bpmndi:BPMNShape>
		  <bpmndi:BPMNShape id="EndEvent_1lxbe5s_di" bpmnElement="EndEvent_1lxbe5s">
			<dc:Bounds x="819" y="447" width="36" height="36" />
			<bpmndi:BPMNLabel>
			  <dc:Bounds x="814" y="487" width="46" height="13" />
			</bpmndi:BPMNLabel>
		  </bpmndi:BPMNShape>
		  <bpmndi:BPMNEdge id="SequenceFlow_0fghwzx_di" bpmnElement="fdf">
			<di:waypoint xsi:type="dc:Point" x="382" y="465" />
			<di:waypoint xsi:type="dc:Point" x="450" y="465" />
			<di:waypoint xsi:type="dc:Point" x="450" y="465" />
			<di:waypoint xsi:type="dc:Point" x="520" y="465" />
			<bpmndi:BPMNLabel>
			  <dc:Bounds x="448" y="459" width="35" height="12" />
			</bpmndi:BPMNLabel>
		  </bpmndi:BPMNEdge>
		  <bpmndi:BPMNEdge id="SequenceFlow_1r953n5_di" bpmnElement="SequenceFlow_1r953n5">
			<di:waypoint xsi:type="dc:Point" x="620" y="465" />
			<di:waypoint xsi:type="dc:Point" x="721" y="465" />
			<di:waypoint xsi:type="dc:Point" x="721" y="465" />
			<di:waypoint xsi:type="dc:Point" x="819" y="465" />
			<bpmndi:BPMNLabel>
			  <dc:Bounds x="719" y="459" width="35" height="12" />
			</bpmndi:BPMNLabel>
		  </bpmndi:BPMNEdge>
		</bpmndi:BPMNPlane>
	  </bpmndi:BPMNDiagram>
	</bpmn:definitions>
`
	form = `
	<?xml version="1.0" encoding="UTF-8"?>
<bpmn:definitions xmlns:bpmn="http://www.omg.org/spec/BPMN/20100524/MODEL" xmlns:bpmndi="http://www.omg.org/spec/BPMN/20100524/DI" xmlns:di="http://www.omg.org/spec/DD/20100524/DI" xmlns:dc="http://www.omg.org/spec/DD/20100524/DC" xmlns:camunda="http://camunda.org/schema/1.0/bpmn" xmlns:xsi="http://www.w3.org/2001/XMLSchema-instance" id="Definitions_1" targetNamespace="http://bpmn.io/schema/bpmn" exporter="Camunda Modeler" exporterVersion="1.11.3">
  <bpmn:process id="Process_1" isExecutable="true">
    <bpmn:startEvent id="StartEvent_1">
      <bpmn:outgoing>SequenceFlow_1hpjl43</bpmn:outgoing>
    </bpmn:startEvent>
    <bpmn:userTask id="Task_1b1eha8" camunda:formKey="key_test">
      <bpmn:extensionElements>
        <camunda:formData>
          <camunda:formField id="id_1" label="label_1" type="string" defaultValue="default_1">
            <camunda:properties>
              <camunda:property id="Property_1" value="Property_1_value" />
              <camunda:property id="Property_2" value="Property_2_value" />
            </camunda:properties>
            <camunda:validation>
              <camunda:constraint name="constraint_1" config="constraint_1_config" />
              <camunda:constraint name="constraint_2" config="constraint_2_config" />
            </camunda:validation>
          </camunda:formField>
          <camunda:formField id="FormField_2" label="FormField_2_label" type="long" />
          <camunda:formField id="FormField_3" type="date" />
          <camunda:formField id="FormField_4" label="FormField_4_label" type="enum" defaultValue="FormField_4_label_default">
            <camunda:value id="Value_1" name="Value_1_value" />
            <camunda:value id="Value_2" name="Value_2_value" />
          </camunda:formField>
          <camunda:formField id="FormField_5" type="boolean" />
        </camunda:formData>
      </bpmn:extensionElements>
      <bpmn:incoming>SequenceFlow_1hpjl43</bpmn:incoming>
      <bpmn:outgoing>SequenceFlow_1oe5xbn</bpmn:outgoing>
    </bpmn:userTask>
    <bpmn:endEvent id="EndEvent_0h4y9rm">
      <bpmn:incoming>SequenceFlow_1oe5xbn</bpmn:incoming>
    </bpmn:endEvent>
    <bpmn:sequenceFlow id="SequenceFlow_1hpjl43" sourceRef="StartEvent_1" targetRef="Task_1b1eha8" />
    <bpmn:sequenceFlow id="SequenceFlow_1oe5xbn" sourceRef="Task_1b1eha8" targetRef="EndEvent_0h4y9rm" />
  </bpmn:process>
  <bpmndi:BPMNDiagram id="BPMNDiagram_1">
    <bpmndi:BPMNPlane id="BPMNPlane_1" bpmnElement="Process_1">
      <bpmndi:BPMNShape id="_BPMNShape_StartEvent_2" bpmnElement="StartEvent_1">
        <dc:Bounds x="272" y="368" width="36" height="36" />
        <bpmndi:BPMNLabel>
          <dc:Bounds x="245" y="404" width="90" height="20" />
        </bpmndi:BPMNLabel>
      </bpmndi:BPMNShape>
      <bpmndi:BPMNShape id="UserTask_1n2a22u_di" bpmnElement="Task_1b1eha8">
        <dc:Bounds x="406" y="346" width="100" height="80" />
      </bpmndi:BPMNShape>
      <bpmndi:BPMNShape id="EndEvent_0h4y9rm_di" bpmnElement="EndEvent_0h4y9rm">
        <dc:Bounds x="638" y="368" width="36" height="36" />
        <bpmndi:BPMNLabel>
          <dc:Bounds x="656" y="408" width="0" height="12" />
        </bpmndi:BPMNLabel>
      </bpmndi:BPMNShape>
      <bpmndi:BPMNEdge id="SequenceFlow_1hpjl43_di" bpmnElement="SequenceFlow_1hpjl43">
        <di:waypoint xsi:type="dc:Point" x="308" y="386" />
        <di:waypoint xsi:type="dc:Point" x="406" y="386" />
        <bpmndi:BPMNLabel>
          <dc:Bounds x="357" y="365" width="0" height="12" />
        </bpmndi:BPMNLabel>
      </bpmndi:BPMNEdge>
      <bpmndi:BPMNEdge id="SequenceFlow_1oe5xbn_di" bpmnElement="SequenceFlow_1oe5xbn">
        <di:waypoint xsi:type="dc:Point" x="506" y="386" />
        <di:waypoint xsi:type="dc:Point" x="638" y="386" />
        <bpmndi:BPMNLabel>
          <dc:Bounds x="572" y="365" width="0" height="12" />
        </bpmndi:BPMNLabel>
      </bpmndi:BPMNEdge>
    </bpmndi:BPMNPlane>
  </bpmndi:BPMNDiagram>
</bpmn:definitions>
`
)

func TestParseBasicBpmn(t *testing.T) {
	p := NewXMLParser()
	v, err := p.Parse(context.Background(), []byte(basic))
	if err != nil {
		fmt.Println(err.Error())
	}
	buf, _ := json.Marshal(v)
	fmt.Println(string(buf))
}

func TestParseFormBpmn(t *testing.T) {
	p := NewXMLParser()
	v, err := p.Parse(context.Background(), []byte(form))
	if err != nil {
		fmt.Println(err.Error())
	}
	buf, _ := json.Marshal(v)
	fmt.Println(string(buf))
}
