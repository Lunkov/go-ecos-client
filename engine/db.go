package engine

import (
  "fmt"
  "github.com/jinzhu/gorm"
  _ "github.com/lib/pq"

  "github.com/golang/glog"
  
  "github.com/Lunkov/go-ecos-client/messages"
)


type DBO struct {
  connectStr       string
  handle          *gorm.DB
}

func NewDBO() *DBO {
  return &DBO{  }
}

func (dbo *DBO) GetHandle() *gorm.DB {
  return dbo.handle
}

func (dbo *DBO) Conn(dbHost, dbPort, dbUser, dbPwd, dbName string) bool {
  // create a connection string
  dbo.connectStr = fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
                                   dbHost, dbPort, dbUser, dbPwd, dbName)
  db, err := gorm.Open("postgres", dbo.connectStr)
  if err != nil {
    glog.Errorf("ERR: DB: Connect '%s:%s/%s' by '%s'", dbHost, dbPort, dbName, dbUser)
    return false
  }
  if glog.V(9) {
    glog.Infof("DBG: DB: CONNECT (%s)", dbo.connectStr)
  }
  dbo.handle = db
  return true
}

func (c *DBO) Migrate() {
  var balance messages.Balance
  var tx TransactionDB
  var block BlockDB
  c.handle.Table("blocks").AutoMigrate(block)
  c.handle.Table("balances").AutoMigrate(balance)
  c.handle.Table("transactions").AutoMigrate(tx)
}

func (c *DBO) DropTables() {
  c.handle.Exec("DROP TABLE balances")
  c.handle.Exec("DROP TABLE transactions")
}
