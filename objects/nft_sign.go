package objects

import (
  "bytes"
  "time"
  "crypto/sha512"
  "crypto/rsa"
  "encoding/gob"
  "encoding/json"
  
  "github.com/Lunkov/go-ecos-client/utils"
)

type SignInfo struct {
  Version       string          `json:"version"`
  CID           string          `json:"cid"`
  Type          string          `json:"type"`
  CreatedAt     time.Time       `json:"created_at"`
  
  IssuerName          string    `json:"issuer_name"`
  IssuerPublicKey     []byte    `json:"issuer_pubkey"`
  IssuerSignType      string    `json:"issuer_sign_type"`
  IssuerSign          []byte    `json:"issuer_sign"`
}

func NewSignInfo() *SignInfo {
  return &SignInfo{}
}

func (si *SignInfo) Hash(data []byte) []byte {
  sha_512 := sha512.New()
  sha_512.Write([]byte(si.Version + si.CID + si.Type + si.IssuerName + si.IssuerSignType))
  sha_512.Write([]byte(si.CreatedAt.String()))
  sha_512.Write(data)
  sha_512.Write(si.IssuerPublicKey)
  return sha_512.Sum(nil)
}

func (si *SignInfo) Serialize() ([]byte, error) {
	var result bytes.Buffer
	encoder := gob.NewEncoder(&result)

	err := encoder.Encode(si)
	if err != nil {
    return nil, err
	}

	return result.Bytes(), nil
}

func (si *SignInfo) Deserialize(data []byte) error {
	decoder := gob.NewDecoder(bytes.NewReader(data))
	err := decoder.Decode(si)
	if err != nil {
		return err
	}

	return nil
}

func (si *SignInfo) ToJSON() ([]byte, error) {
  jsonAnswer, err := json.Marshal(si)
  if err != nil {
    return jsonAnswer, err
  }
	return jsonAnswer, nil
}

func (si *SignInfo) FromJSON(data []byte) error {
  if err := json.Unmarshal(data, si); err != nil {
    return err
  }
  return nil
}

func (si *SignInfo) DoSignIssuer(data []byte, pk *rsa.PrivateKey) error {
  PublicKey, err := utils.RSASerializePublicKey(&pk.PublicKey)
  if err != nil {
    return err
  }
  si.IssuerPublicKey = PublicKey
  sign, errs := utils.RSASign(pk, si.Hash(data))
  if errs != nil {
    return errs
  }
  si.IssuerSign = sign
  return nil
}

func (si *SignInfo) DoVerifyIssuer(data []byte, ) (error) {
  return utils.RSAVerify(si.IssuerPublicKey, si.Hash(data), si.IssuerSign)
}
