package messages

import (

)

/*
 * System
 * 
 */

type System struct {
  Id           string
  Name         string
  Type         string
  RedirectURL  string
  ImageURL     string
}
 
func NewSystem() *System {
  return &System{}
}


/*
 * Systems
 * 
 */

type Systems struct {
  Systems    []System
}
 
func NewSystems() *Systems {
  return &Systems{Systems: make([]System, 0)}
}


// User
type UserACL struct {
  Version       string          `json:"version"     gorm:"column:version"`
  UserId        string          `json:"user_id"     gorm:"column:user_id"`
  SystemId      string          `json:"system_id"   gorm:"column:system_id"`
  Role          string          `json:"role"        gorm:"column:role"`
}

/*
 * UserACL
 * 
 */
 
func NewUserACL() *UserACL {
  return &UserACL{Version: "1"}
}

// User
type UserACLs struct {
  Version       string          `json:"version"     gorm:"column:version"`
  ACLs        []UserACL         `json:"acls"        gorm:"column:acls"`
}
