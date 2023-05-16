package objects

import (
  "strings"
  "bytes"
  "encoding/gob"
)

/*
 * 
 * 
 */
type ServiceFunc func(string) (string, bool, string)

type CallFunc struct {
  Urls       []string                   `db:"urls"           json:"urls,omitempty"            yaml:"urls"`
  Func         ServiceFunc
}

type Service struct {
  CODE         string                   `db:"code"           json:"code,omitempty"            yaml:"code"             unique:"true"`
  /*
   * Types of service:
   * - func
   * - grpc
   * - rest
  */
  Type         string                   `db:"type"           json:"type,omitempty"            yaml:"type"`
  Version      string                   `db:"version"        json:"version,omitempty"         yaml:"version"`
  Name         string                   `db:"name"           json:"name,omitempty"            yaml:"name"`
  Description  string                   `db:"description"    json:"description,omitempty"     yaml:"description"`
  Disabled     bool                     `db:"disabled"       json:"disabled,omitempty"        yaml:"disabled"`
  MainCall     CallFunc                 `db:"main"           json:"main,omitempty"            yaml:"main"`
  TestCall     CallFunc                 `db:"test"           json:"test,omitempty"            yaml:"test"`
}

type Services struct {
  a map[string]Service
}

func NewServices() *Services {
  return &Services{
                 a: make(map[string]Service),
               }
}

func (s *Services) FileExtension() string {
  return ".service"
}

func (s *Services) Count() int64 {
  return int64(len(s.a))
}

func (s *Services) Append(info *Service) {
  info.CODE = strings.ToLower(info.CODE)
  s.a[info.CODE] = *info
}

func (s *Services) GetByCODE(code string) (*Service) {
  item, ok := s.a[code]
  if ok {
    return &item
  }
  return nil
}

func (s *Services) ExistsByCODE(code string) bool {
  _, ok := s.a[code]
  return ok
}

func (s *Services) GetList() []Service {
  res := make([]Service, 0)
  for _, item := range s.a {
    res = append(res, item)
  }
  return res
}

func (s *Services) Serialize() []byte {
  var buff bytes.Buffer
  encoder := gob.NewEncoder(&buff)
  encoder.Encode(s)
  return buff.Bytes()
}

func (s *Services) Deserialize(msg []byte) bool {
  buf := bytes.NewBuffer(msg)
  decoder := gob.NewDecoder(buf)
  err := decoder.Decode(s)
  if err != nil {
    return false
  }
  return true
}

