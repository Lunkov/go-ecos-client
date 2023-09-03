package messages

import (
  "bytes"
  "strings"
  "encoding/gob"
  "encoding/json"
  
  "github.com/Lunkov/go-ecos-client/objects"
)

type OrganizationInfo struct {
  Organization     objects.OrganizationInfo
  CID              string
  TX               objects.Transaction 
}

func NewOrganizationInfo() *OrganizationInfo {
  return &OrganizationInfo{}
}

func (p *OrganizationInfo) GetCID() (string, bool) {
  obj := strings.Split(p.CID, ":")
  if len(obj) < 2 {
    return "", false
  }
  return obj[len(obj) - 1], true
}

func (p *OrganizationInfo) GetType() (string, bool) {
  obj := strings.Split(p.CID, ":")
  if len(obj) < 1 {
    return "", false
  }
  return obj[0], true
}
            
func (p *OrganizationInfo) Serialize() ([]byte, error) {
	var result bytes.Buffer
	encoder := gob.NewEncoder(&result)

	err := encoder.Encode(p)
	if err != nil {
    return nil, err
	}

	return result.Bytes(), nil
}

func (p *OrganizationInfo) Deserialize(data []byte) error {
	decoder := gob.NewDecoder(bytes.NewReader(data))
	return decoder.Decode(p)
}

func (p *OrganizationInfo) ToJSON() ([]byte, error) {
  return json.Marshal(p)
}

func (p *OrganizationInfo) FromJSON(data []byte) error {
  return json.Unmarshal(data, p)
}
