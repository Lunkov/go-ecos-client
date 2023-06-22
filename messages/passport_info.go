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
            
func (p *PassportInfo) Serialize() ([]byte, bool) {
	var result bytes.Buffer
	encoder := gob.NewEncoder(&result)

	err := encoder.Encode(p)
	if err != nil {
    return nil, false
	}

	return result.Bytes(), true
}

func (p *PassportInfo) Deserialize(data []byte) bool {
	decoder := gob.NewDecoder(bytes.NewReader(data))
	err := decoder.Decode(p)
	if err != nil {
		return false
	}

	return true
}

func (p *PassportInfo) ToJSON() ([]byte, bool) {
  jsonAnswer, err := json.Marshal(p)
  if err != nil {
    return jsonAnswer, false
  }
	return jsonAnswer, true
}

func (p *PassportInfo) FromJSON(data []byte) bool {
  if err := json.Unmarshal(data, p); err != nil {
    return false
  }
  return true
}
