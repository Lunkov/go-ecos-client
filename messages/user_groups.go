package messages

import (
  "github.com/google/uuid"
)

/*
 * UserGroup
 * 
 */
 
type UserInGroup struct {
  Id             string
  Name           string
  Type           string
  TypeKey        string
  Key          []byte
}

func NewUserInGroup() (*UserInGroup) {
  return &UserInGroup{}
}

/*
 * UserGroup
 * 
 */
 
type UserGroup struct {
  Id             string
  Name           string
  Type           string
  Users        []*UserInGroup
}

func NewUserGroup() (*UserGroup) {
  return &UserGroup{Users: make([]*UserInGroup, 0)}
}

func (t *UserGroup) AddUser(name string, tp string, tpkey string, key []byte)  {
  tu := NewUserInGroup()
  uid, _ := uuid.NewUUID()
  tu.Id = uid.String()
  tu.Name = name
  tu.Type = tp
  tu.TypeKey = tpkey
  tu.Key = key
  t.Users = append(t.Users, tu)
}

/*
 * UserGroups
 * 
 */
 
type UserGroups []*UserGroup

func NewUserGroups() UserGroups {
  return make([]*UserGroup, 0)
}
