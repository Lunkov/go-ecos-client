package objects

import (
  "bytes"
  "time"
  "crypto/sha512"
  "encoding/gob"
  "encoding/json"
)

type Passport struct {
  Version       string          `json:"version"`

  ID            string          `json:"id"`

  DisplayName   string          `json:"displayName"`
  FirstName     string          `json:"first_name"`
  MiddleName    string          `json:"middle_name"`
  LastName      string          `json:"last_name"`

  Country       string          `json:"country"`
  Locality      string          `json:"locality"`
  Role          string          `json:"role"        gorm:"column:role"`
  EMail         string          `json:"email"       gorm:"index:idx_email,unique;column:email"`

  Phone         string          `json:"phone"`

  Photo         []byte          `json:"photo"`

  UpdatedAt     time.Time       `json:"updated_at"  gorm:"column:updated_at"`
  
  Address             string          `json:"address"`
  PersonPublicKey     []byte          `json:"person_pubkey"   gorm:"column:public_key"`  
  PersonSign          []byte          `json:"person_sign"     gorm:"column:sign"`

  IssuerName          string          `json:"issuer_name"`
  IssuerPubKey        []byte          `json:"issuer_pubkey"`
  IssuerSign          []byte          `json:"issuer_sign"`
}

func NewPassport() *Passport {
  return &Passport{}
}

func (p *Passport) Hash() []byte {
  sha_512 := sha512.New()
  sha_512.Write([]byte(p.Version + p.ID + p.IssuerName + p.Address))
  sha_512.Write([]byte(p.DisplayName + p.FirstName + p.MiddleName + p.LastName + p.EMail + p.Phone + p.UpdatedAt.String()))
  sha_512.Write([]byte(p.Country + p.Locality + p.Role))
  sha_512.Write(p.Photo)
  return sha_512.Sum(nil)
}

func (p *Passport) Serialize() ([]byte, bool) {
	var result bytes.Buffer
	encoder := gob.NewEncoder(&result)

	err := encoder.Encode(p)
	if err != nil {
    return nil, false
	}

	return result.Bytes(), true
}

func (p *Passport) Deserialize(data []byte) bool {
	decoder := gob.NewDecoder(bytes.NewReader(data))
	err := decoder.Decode(p)
	if err != nil {
		return false
	}

	return true
}

func (p *Passport) ToJSON() ([]byte, bool) {
  jsonAnswer, err := json.Marshal(p)
  if err != nil {
    return jsonAnswer, false
  }
	return jsonAnswer, true
}

func (p *Passport) FromJSON(data []byte) bool {
  if err := json.Unmarshal(data, p); err != nil {
    return false
  }
  return true
}
