package messages

import (
  "bytes"
  "strings"
  "encoding/gob"
  "encoding/json"
  
  "github.com/Lunkov/go-ecos-client/objects"
)

type PassportInfo struct {
  Passport     objects.Passport
  CodePhrase   string
  CID          string
  TX           objects.Transaction 
}

func NewPassportInfo() *PassportInfo {
  return &PassportInfo{}
}

func (p *PassportInfo) GetCID() (string, bool) {
  obj := strings.Split(p.CID, ":")
  if len(obj) < 2 {
    return "", false
  }
  return obj[len(obj) - 1], true
}

func (p *PassportInfo) GetType() (string, bool) {
  obj := strings.Split(p.CID, ":")
  if len(obj) < 1 {
    return "", false
  }
  return obj[0], true
}
            
func (p *PassportInfo) Serialize() ([]byte, error) {
	var result bytes.Buffer
	encoder := gob.NewEncoder(&result)

	err := encoder.Encode(p)
	if err != nil {
    return nil, err
	}

	return result.Bytes(), nil
}

func (p *PassportInfo) Deserialize(data []byte) error {
	decoder := gob.NewDecoder(bytes.NewReader(data))
	return decoder.Decode(p)
}

func (p *PassportInfo) ToJSON() ([]byte, error) {
  return json.Marshal(p)
}

func (p *PassportInfo) FromJSON(data []byte) error {
  return json.Unmarshal(data, p)
}
