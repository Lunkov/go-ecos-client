package objects

/*
import (
  "flag"
  "testing"
  "github.com/stretchr/testify/assert"
  "github.com/golang/glog"
)

func TestCheckBPMN(t *testing.T) {
  flag.Set("alsologtostderr", "true")
  flag.Set("log_dir", ".")
  flag.Set("v", "9")
  flag.Parse()

  glog.Info("Logging configured: TestCheckBPMN")

  var b Schema
  ok, msg := b.LoadXML("etc.test/wrong_diagram_1.bpmn")
  assert.Equal(t, true, ok)
  assert.Equal(t, "OK", msg)
  
  ok, msg = b.Validate(nil)
  assert.Equal(t, false, ok)
  assert.Equal(t, "ERR: SendTask(Process_Wrong_1.TASK_WITHOUT_OUTPUT) no has outgoing", msg)

  b = Schema{}
  ok, msg = b.LoadXML("etc.test/wrong_diagram_2.bpmn")
  assert.Equal(t, true, ok)
  assert.Equal(t, "OK", msg)
  
  ok, msg = b.Validate(nil)
  assert.Equal(t, false, ok)
  assert.Equal(t, "ERR: StartEvent(Process_Wrong_2.StartEvent_WITHOUT_OUTPUT) no has outgoing", msg)

  srvc := NewServices()
  //srvc.LoadFromFiles("etc.test")

  b = Schema{}
  ok, msg = b.LoadXML("etc.test/new_user.bpmn")
  assert.Equal(t, true, ok)
  assert.Equal(t, "OK", msg)

  ok, msg = b.Validate(nil)
  assert.Equal(t, false, ok)
  assert.Equal(t, "ERR: ServiceTask(Process_NewUser.Create_Email) not found service 'srv-report'", msg)
  
  ok, msg = b.Validate(srvc)
  assert.Equal(t, true, ok)
  assert.Equal(t, "OK", msg)
  
  assert.Equal(t, false, b.findNextServiceTasks("StartEvent_NewUser"))
  
  assert.Equal(t, false, b.isFinish("Event_FINISH_1"))
  assert.Equal(t, true, b.isFinish("Flow_to_the_end"))

}

*/
