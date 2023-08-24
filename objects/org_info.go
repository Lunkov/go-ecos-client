package objects

import (
  "time"
  "strings"
  
  "bytes"
  "crypto/sha512"
  "encoding/gob"
  
  "github.com/google/uuid"
  
  "github.com/Lunkov/go-ecos-client/utils"
)

// Information about Organization
type OrgInfo struct {
  Id                string

  CID               string
  PrevCID           string

  DisplayName       string
  DisplayNameTr     map[string]string
  
  Description       string
  DescriptionTr     map[string]string
  
  TypeId            uint32

  LogoURL           string          `json:"logo_url"`
  Logo             []byte
  
  URL               string          `json:"url"`
  APIURL          []string          `json:"api_url"`

  EMailInfo         string

  Country           string
  Locality          string
  
  WalletAddress     string
  CreatedAt         time.Time       `json:"created_at"`

  Cert             []byte
  Sign             []byte
}

func NewOrgInfo() *OrgInfo {
  return &OrgInfo{}
}

func (io *OrgInfo) NewID() {
  uid, _ := uuid.NewUUID()
  io.Id = uid.String()
}

func (io *OrgInfo) Hash() []byte {
  sha_512 := sha512.New()
  sha_512.Write([]byte(io.Id + io.Description + io.PrevCID + io.WalletAddress + io.CreatedAt.String()))
  sha_512.Write([]byte(io.DisplayName + io.LogoURL + io.EMailInfo + io.Country + io.Locality))
  sha_512.Write([]byte(strings.Join(io.APIURL, ":")))
  sha_512.Write(io.Cert)
  sha_512.Write(utils.UInt32ToBytes(io.TypeId))
  sha_512.Write(io.Logo)
  return sha_512.Sum(nil)
}

func (oi *OrgInfo) Serialize() []byte {
  var buff bytes.Buffer
  encoder := gob.NewEncoder(&buff)
  encoder.Encode(oi)
  return buff.Bytes()
}

func (oi *OrgInfo) Deserialize(msg []byte) error {
  buf := bytes.NewBuffer(msg)
  decoder := gob.NewDecoder(buf)
  return decoder.Decode(oi)
}

func (oi *OrgInfo) DoSign(cert utils.CertInfo) error {
  crt, err := cert.SerializeCert()
  if err != nil {
    return err
  }
  sign, err2 := cert.Sign(oi.Hash())
  if err2 != nil {
    return err2
  }
  oi.Cert = crt
  oi.Sign = sign
  return nil
}

func (oi *OrgInfo) DoVerify() (error) {
  cert := utils.NewCertInfo()
  
  err := cert.DeserializeCert(oi.Cert)
  if err != nil {
    return err
  }
  
  return cert.Verify(oi.Hash(), oi.Sign)
}
