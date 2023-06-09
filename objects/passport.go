package objects

import (
  "bytes"
  "time"
  "crypto/sha512"
  "crypto/rsa"
  "encoding/gob"
  "encoding/json"
  
  "github.com/Lunkov/lib-cipher"
  "github.com/Lunkov/go-hdwallet"
  "github.com/Lunkov/lib-wallets"
  "github.com/Lunkov/go-ecos-client/utils"
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

  CreatedAt     time.Time       `json:"created_at"`
  
  Address             string    `json:"address"`
  PersonPublicKey     []byte    `json:"person_pubkey"`  
  PersonSign          []byte    `json:"person_sign"`

  IssuerName          string    `json:"issuer_name"`
  IssuerPublicKey     []byte    `json:"issuer_pubkey"`
  IssuerSign          []byte    `json:"issuer_sign"`
}

func NewPassport() *Passport {
  return &Passport{}
}

func (p *Passport) HashPerson() []byte {
  sha_512 := sha512.New()
  sha_512.Write([]byte(p.Version + p.IssuerName + p.Address))
  sha_512.Write([]byte(p.DisplayName + p.FirstName + p.MiddleName + p.LastName + p.EMail + p.Phone))
  sha_512.Write([]byte(p.Country + p.Locality + p.Role))
  sha_512.Write(p.PersonPublicKey)
  sha_512.Write(p.Photo)
  return sha_512.Sum(nil)
}

func (p *Passport) HashIssuer() []byte {
  sha_512 := sha512.New()
  sha_512.Write([]byte(p.ID))
  sha_512.Write([]byte(p.CreatedAt.String()))
  sha_512.Write(p.IssuerPublicKey)
  sha_512.Write(p.HashPerson())
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

func (p *Passport) DoSignPerson(wallet wallets.IWallet) bool {
  p.Address = wallet.GetAddress(hdwallet.ECOS)
  PublicKey, ok := utils.ECDSAPublicKeySerialize(wallet.GetECDSAPublicKey())
  if !ok {
    return false
  }
  p.PersonPublicKey = PublicKey
  sign, ok := utils.ECDSA256SignHash512(wallet.GetECDSAPrivateKey(), p.HashPerson())
  if !ok {
    return false
  }
  p.PersonSign = sign
  return true
}

func (p *Passport) DoVerifyPerson() bool {
  return utils.ECDSA256VerifyHash512(p.PersonPublicKey, p.HashPerson(), p.PersonSign)
}

func (p *Passport) DoSignIssuer(pk *rsa.PrivateKey) bool {
  PublicKey, ok := utils.RSASerializePublicKey(&pk.PublicKey)
  if !ok {
    return false
  }
  p.IssuerPublicKey = PublicKey
  sign, ok := utils.RSASign(pk, p.HashIssuer())
  if !ok {
    return false
  }
  p.PersonSign = sign
  return true
}

func (p *Passport) DoVerifyIssuer() bool {
  return utils.RSAVerify(p.IssuerPublicKey, p.HashIssuer(), p.IssuerSign)
}

func (p *Passport) SerializeEncrypt(password string) ([]byte, bool) {
  buf, ok := p.Serialize()
  if !ok {
    return nil, false
  }
  c := cipher.NewSCipher()
  key := c.Password2Key(password)
  enc, okenc := c.AESEncrypt(key, buf)
  if !okenc {
    return nil, false
  }

  return enc, true
}


func (p *Passport) DeserializeDecrypt(password string, buf []byte) (bool) {
  c := cipher.NewSCipher()
  key := c.Password2Key(password)
  dec, ok := c.AESDecrypt(key, buf)
  if !ok {
    return false
  }
  return p.Deserialize(dec)
}
