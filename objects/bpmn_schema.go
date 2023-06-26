package objects

import (
  "fmt"
  "os"
  "encoding/xml"
)

type PropertyInfo struct {
  Name   string `json:"name"   yaml:"name"   xml:"name,attr"`
  Value  string `json:"value"  yaml:"value"  xml:"value,attr"`
}

type inputOutputParameter struct {
  Name   string `json:"name"   yaml:"name"   xml:"name,attr"`
  Value  string `json:"value"  yaml:"value"  xml:",chardata"`
}

type ExtentionsInfo struct {
  Input  []inputOutputParameter   `json:"inputParameter"  yaml:"inputParameter"  xml:"inputOutput>inputParameter"`
}

type CollaborationInfo struct {
  Id          string           `json:"id"           yaml:"id"          xml:"id,attr"`
  Properies   []PropertyInfo   `json:"properties"   yaml:"properties"  xml:"properties>property"`
}

type StartEventInfo struct {
  Id            string           `json:"id"                       yaml:"id"                       xml:"id,attr"`
  Outgoing      []string         `json:"outgoing"                 yaml:"outgoing"                 xml:"outgoing"`
  Properies     []PropertyInfo   `json:"properties"               yaml:"properties"               xml:"properties>property"`
  MessageEvent  []MessageInfo    `json:"messageEventDefinition"   yaml:"messageEventDefinition"   xml:"messageEventDefinition"`
  Timer         TimeInfo         `json:"timerEventDefinition"     yaml:"timerEventDefinition"     xml:"timerEventDefinition"`
}

type EndEventInfo struct {
  Id          string            `json:"id"          yaml:"id"          xml:"id,attr"`
  Incoming    []string          `json:"incoming"    yaml:"incoming"    xml:"incoming"`
  Properies   []PropertyInfo    `json:"properties"  yaml:"properties"  xml:"properties>property"`
}

type ExclusiveGatewayInfo struct {
  Id            string           `json:"id"           yaml:"id"           xml:"id,attr"`
  Name          string           `json:"name"         yaml:"name"         xml:"name,attr"`
  Incoming      []string         `json:"incoming"     yaml:"incoming"     xml:"incoming"`
  Outgoing      []string         `json:"outgoing"     yaml:"outgoing"     xml:"outgoing"`
  Properies     []PropertyInfo   `json:"properties"   yaml:"properties"   xml:"properties>property"`
}

type ParallelGatewayInfo struct {
  Id             string           `json:"id"           yaml:"id"           xml:"id,attr"`
  Name           string           `json:"name"         yaml:"name"         xml:"name,attr"`
  Incoming       []string         `json:"incoming"     yaml:"incoming"     xml:"incoming"`
  Outgoing       []string         `json:"outgoing"     yaml:"outgoing"     xml:"outgoing"`
  Properies      []PropertyInfo   `json:"properties"   yaml:"properties"   xml:"properties>property"`
}

type TaskInfo struct {
  Id                 string           `json:"id"               yaml:"id"                xml:"id,attr"`
  Name               string           `json:"name"             yaml:"name"              xml:"name,attr"`
  Type               string           `json:"type"             yaml:"type"              xml:"type,attr"`
  Class              string           `json:"class"            yaml:"class"             xml:"class,attr"`
  Topic              string           `json:"topic"            yaml:"topic"             xml:"topic,attr"`
  CandidateGroups    string           `json:"candidateGroups"  yaml:"candidateGroups"   xml:"candidateGroups,attr"`
  FormKey            string           `json:"formKey"          yaml:"formKey"           xml:"formKey,attr"`
  Incoming           []string         `json:"incoming"         yaml:"incoming"          xml:"incoming"`
  Outgoing           []string         `json:"outgoing"         yaml:"outgoing"          xml:"outgoing"`
  Properies          []PropertyInfo   `json:"properties"       yaml:"properties"        xml:"properties>property"`
  Extentions         ExtentionsInfo `json:"extensionElements"       yaml:"extensionElements"        xml:"extensionElements"`
}

type ConditionExpressionInfo struct {
  Expression      string           `json:"expression"   yaml:"expression"    xml:",chardata"`
  Type            string           `json:"type"         yaml:"type"          xml:"type,attr"`
  Language        string           `json:"language"     yaml:"language"      xml:"language,attr"`
  Properies       []PropertyInfo   `json:"properties"   yaml:"properties"    xml:"properties>property"`
}

type SequenceFlowInfo struct {
  Id                   string                      `json:"id"                    yaml:"id"                    xml:"id,attr"`
  SourceRef            string                      `json:"sourceRef"             yaml:"sourceRef"             xml:"sourceRef,attr"`
  TargetRef            string                      `json:"targetRef"             yaml:"targetRef"             xml:"targetRef,attr"`
  ConditionExpression  []ConditionExpressionInfo   `json:"conditionExpression"   yaml:"conditionExpression"   xml:"conditionExpression"`
  Properies            []PropertyInfo              `json:"properties"            yaml:"properties"            xml:"properties>property"`
}

type ProcessInfo struct {
  Id                 string                   `json:"id"        yaml:"id"                  xml:"id,attr"`
  Name               string                   `json:"name"      yaml:"name"                xml:"name,attr"`
  Version            string                   `json:"version"   yaml:"version"             xml:"versionTag,attr"`
  StartEvent         []StartEventInfo         `json:"-"         yaml:"startEvent"          xml:"startEvent"`
  EndEvent           []EndEventInfo           `json:"-"         yaml:"endEvent"            xml:"endEvent"`
  ExclusiveGateway   []ExclusiveGatewayInfo   `json:"-"         yaml:"exclusiveGateway"    xml:"exclusiveGateway"`
  ParallelGateway    []ParallelGatewayInfo    `json:"-"         yaml:"parallelGateway"     xml:"parallelGateway"`
  UserTask           []TaskInfo               `json:"-"         yaml:"startEvent"          xml:"userTask"`
  ServiceTask        []TaskInfo               `json:"-"         yaml:"serviceTask"         xml:"serviceTask"`
  SendTask           []TaskInfo               `json:"-"         yaml:"sendTask"            xml:"sendTask"`
  SequenceFlow       []SequenceFlowInfo       `json:"-"         yaml:"sequenceFlow"        xml:"sequenceFlow"`
  Properies          []PropertyInfo           `json:"-"         yaml:"properties"          xml:"properties>property"`
}

type MessageInfo struct {
  Id                 string                   `json:"id"            yaml:"id"              xml:"id,attr"`
  Name               string                   `json:"name"          yaml:"name"            xml:"name,attr"`
  MessageRef         string                   `json:"messageRef"    yaml:"messageRef"      xml:"messageRef,attr"`
}

type TimeInfo struct {
  Type            string                   `json:"type"           yaml:"type"           xml:"type,attr"`
  ValueCyrcle     string                   `json:"timeCycle"      yaml:"timeCycle"      xml:"timeCycle"`
  ValueDuration   string                   `json:"timeDuration"   yaml:"timeDuration"   xml:"timeDuration"`
}

type Schema struct {
  Id              string                `json:"id"             yaml:"id"             xml:"id,attr"`  
  Collaboration   []CollaborationInfo   `json:"collaboration"  yaml:"collaboration"  xml:"collaboration"`
  Process         []ProcessInfo         `json:"process"        yaml:"process"        xml:"process"`
  Message         []MessageInfo         `json:"message"        yaml:"message"        xml:"message"`
  SequenceFlow    []SequenceFlowInfo    `json:"sequenceflow"   yaml:"sequenceflow"   xml:"sequenceFlow"`
}

// Simple test schema
func (b *Schema) Validate(s *Services) (bool, string) {
  for _, pr := range b.Process {
    for _, se := range pr.StartEvent {
      if(len(se.Outgoing) < 1) {
        return false, fmt.Sprintf("ERR: StartEvent(%s.%s) no has outgoing", pr.Id, se.Id)
      }
    }
    for _, eg := range pr.ExclusiveGateway {
      if(len(eg.Outgoing) < 1) {
        return false, fmt.Sprintf("ERR: ExclusiveGateway(%s.%s) no has outgoing", pr.Id, eg.Id)
      }
    }
    for _, eg := range pr.ParallelGateway {
      if(len(eg.Outgoing) < 1) {
        return false, fmt.Sprintf("ERR: ParallelGateway(%s.%s) no has outgoing", pr.Id, eg.Id)
      }
    }
    for _, st := range pr.ServiceTask {
      if(len(st.Outgoing) < 1) {
        return false, fmt.Sprintf("ERR: ServiceTask(%s.%s) no has outgoing", pr.Id, st.Id)
      }
      if s == nil || (s != nil && !s.ExistsByCODE(st.Topic)) {
        return false, fmt.Sprintf("ERR: ServiceTask(%s.%s) not found service '%s'", pr.Id, st.Id, st.Topic)
      }
    }
    for _, st := range pr.SendTask {
      if(len(st.Outgoing) < 1) {
        return false, fmt.Sprintf("ERR: SendTask(%s.%s) no has outgoing", pr.Id, st.Id)
      }
      if s == nil || (s != nil && !s.ExistsByCODE(st.Topic)) {
        return false, fmt.Sprintf("ERR: SendTask(%s.%s) not found service '%s'", pr.Id, st.Id, st.Topic)
      }
    }
    for _, ut := range pr.UserTask {
      if(len(ut.Outgoing) < 1) {
        return false, fmt.Sprintf("ERR: UserTask(%s.%s) no has outgoing", pr.Id, ut.Id)
      }
    }
  }
  return true, "OK"
}

func contains(a []string, x string) bool {
  for _, n := range a {
    if x == n {
      return true
    }
  }
  return false
}

func (b *Schema) findNextUserTasks(SequenceFlow string) []TaskInfo {
  var userTasks []TaskInfo
  for _, pr := range b.Process {
    for _, ut := range pr.UserTask {
      if(contains(ut.Incoming, SequenceFlow)) {
        userTasks = append(userTasks, ut)
      }
    }
  }
  return userTasks
}

func (b *Schema) findNextSendTasks(SequenceFlow string) []TaskInfo {
  var sendTasks []TaskInfo
  for _, pr := range b.Process {
    for _, st := range pr.SendTask {
      if(contains(st.Incoming, SequenceFlow)) {
        sendTasks = append(sendTasks, st)
      }
    }
  }
  return sendTasks
}

func (b *Schema) findNextServiceTasks(SequenceFlow string) []TaskInfo {
  var serviceTasks []TaskInfo
  for _, pr := range b.Process {
    for _, st := range pr.ServiceTask {
      if(contains(st.Incoming, SequenceFlow)) {
        serviceTasks = append(serviceTasks, st)
      }
    }
  }
  return serviceTasks
}

func (b *Schema) isFinish(SequenceFlow string) bool {
  for _, pr := range b.Process {
    for _, ee := range pr.EndEvent {
      if(contains(ee.Incoming, SequenceFlow)) {
        return true
      }
    }
  }
  return false
}

func (b *Schema) LoadXML(filename string) (bool, string) {
  xmlFile, err := os.ReadFile(filename)
  if err != nil {
    return false, fmt.Sprintf("ERR: ReadFile(%s). %v", filename, err)
  }
  err = xml.Unmarshal(xmlFile, b)
  if err != nil {
    return false, fmt.Sprintf("ERR: unmarshalFile(%s). %v", filename, err)
  }
  return true, "OK"
}

