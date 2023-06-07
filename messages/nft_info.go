package messages

import (
  "bytes"
  "crypto/sha512"
  "encoding/gob"
  "encoding/json"
)

type NFTInfo struct {
  Address       string          `json:"address"`

  Type          string          `json:"type"`

  CID           string          `json:"cid"`

}

func NewNFTInfo() *NFTInfo {
  return &NFTInfo{}
}

func (p *NFTInfo) Hash() []byte {
  sha_512 := sha512.New()
  sha_512.Write([]byte(p.Address + p.Type + p.CID))
  return sha_512.Sum(nil)
}

func (p *NFTInfo) Serialize() ([]byte, bool) {
	var result bytes.Buffer
	encoder := gob.NewEncoder(&result)

	err := encoder.Encode(p)
	if err != nil {
    return nil, false
	}

	return result.Bytes(), true
}

func (p *NFTInfo) Deserialize(data []byte) bool {
	decoder := gob.NewDecoder(bytes.NewReader(data))
	err := decoder.Decode(p)
	if err != nil {
		return false
	}

	return true
}

func (p *NFTInfo) ToJSON() ([]byte, bool) {
  jsonAnswer, err := json.Marshal(p)
  if err != nil {
    return jsonAnswer, false
  }
	return jsonAnswer, true
}

func (p *NFTInfo) FromJSON(data []byte) bool {
  if err := json.Unmarshal(data, p); err != nil {
    return false
  }
  return true
}
