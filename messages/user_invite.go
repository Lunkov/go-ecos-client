package messages

import (
  "time"
  "bytes"
  "crypto/sha512"
  "encoding/gob"
  "encoding/base64"
)

// Invite User to Organization
type DBUserInvite struct {
  Version       string          `json:"version"     gorm:"column:version"`
  InviteToken   string          `json:"token"       gorm:"column:invite_token"`
  Login         string          `json:"login"       gorm:"column:login"`
  Role          string          `json:"role"        gorm:"column:role"`
  EMail         string          `json:"email"       gorm:"index:idx_email,unique;column:email"`
  UpdatedAt     time.Time       `json:"updated_at"  gorm:"column:updated_at"`
  OrgSign     []byte            `json:"org_sign"    gorm:"column:org_sign"`
}

type ReqUserInviteFromOrg struct {
  Version       string          `json:"version"`
  NodeUrl       string          `json:"url"`
  InviteToken   string          `json:"token"`
  Login         string          `json:"login"`
  EMail         string          `json:"email"`
  Role          string          `json:"role"`
  UpdatedAt     time.Time       `json:"updated_at"`
  Sign          []byte          `json:"sign"`
}

type ReqUserRegisterToOrg struct {
  Version       string          `json:"version"`
  NodeUrl       string          `json:"url"`
  InviteToken   string          `json:"token"`
  Login         string          `json:"login"`
  EMail         string          `json:"email"`

  DisplayName   string          `json:"displayName"`
  FirstName     string          `json:"first_name"`
  MiddleName    string          `json:"middle_name"`
  LastName      string          `json:"last_name"`

  Country       string          `json:"country"`
  Locality      string          `json:"locality"`
  OrgUnit       string          `json:"unit"`
  Role          string          `json:"role"`

  PubKey        []byte          `json:"pubkey"`
  Sign          []byte          `json:"sign"`
}

type UserRegister struct {
  Version       string          `json:"version"`

  Cert          []byte          `json:"cert"`
  PrivKey       []byte          `json:"privkey"`
  Sign          []byte          `json:"sign"`
}

/*
 * DBUserInvite
 * 
 * 
 */
 
func NewDBUserInvite() *DBUserInvite {
  return &DBUserInvite{Version: "1"}
}

func (i *DBUserInvite) Init(token string, login string, email string, role string) {
  i.InviteToken = token
  i.Login = login
  i.EMail = email
  i.Role = role
  i.UpdatedAt = time.Now()
}

func (i *DBUserInvite) Hash() []byte {
  sha_512 := sha512.New()
  sha_512.Write([]byte(i.InviteToken + i.Version + i.Login + i.EMail + i.Role + i.UpdatedAt.String()))
  return sha_512.Sum(nil)
}

/*
 * ReqUserInviteFromOrg
 * 
 * 
 */

func NewReqUserInviteFromOrg() (*ReqUserInviteFromOrg) {
  return &ReqUserInviteFromOrg{Version: "1"}
}

func (i *ReqUserInviteFromOrg) Init(nodeUrl string, dbui *DBUserInvite) {
  i.NodeUrl = nodeUrl
  i.InviteToken = dbui.InviteToken
  i.Login = dbui.Login 
  i.EMail = dbui.EMail
  i.Role = dbui.Role
  i.UpdatedAt = dbui.UpdatedAt
}

func (i *ReqUserInviteFromOrg) Hash() []byte {
  sha_512 := sha512.New()
  sha_512.Write([]byte(i.InviteToken + i.Version + i.NodeUrl + i.Login + i.EMail + i.Role + i.UpdatedAt.String()))
  return sha_512.Sum(nil)
}

func (i *ReqUserInviteFromOrg) Pack() string {
  var buff bytes.Buffer
  encoder := gob.NewEncoder(&buff)
  encoder.Encode(i)
  return base64.StdEncoding.EncodeToString(buff.Bytes())
}

func (i *ReqUserInviteFromOrg) Unpack(msg string) bool {
  input, err := base64.StdEncoding.DecodeString(msg)
  if err != nil {
    return false
  }
  buf := bytes.NewBuffer(input)
  decoder := gob.NewDecoder(buf)
  err = decoder.Decode(i)
  if err != nil {
    return false
  }
  return true
}

/*
 * ReqUserRegisterToOrg
 * 
 * 
 */

func NewReqUserRegisterToOrg() (*ReqUserRegisterToOrg) {
  return &ReqUserRegisterToOrg{Version: "1"}
}

func (i *ReqUserRegisterToOrg) Hash() []byte {
  sha_512 := sha512.New()
  sha_512.Write([]byte(i.InviteToken + i.Version + i.NodeUrl + i.Login + i.EMail + i.DisplayName + 
                       i.FirstName + i.MiddleName + i.LastName + i.Role +
                       i.Country + i.Locality))
  sha_512.Write(i.PubKey)
  return sha_512.Sum(nil)
}

/*
 * UserRegister
 * 
 */
 
func NewUserRegister() (*UserRegister) {
  return &UserRegister{Version: "1"}
}


