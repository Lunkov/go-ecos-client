package engine

import (
  "github.com/jinzhu/gorm"
)

type Balance struct {
  Balance              uint64      `gorm:"column:balance"`
  UnconfirmedBalance   uint64      `gorm:"column:unconfirmed_balance"`
  TotalReceived        uint64      `gorm:"column:total_received"`
  TotalSent            uint64      `gorm:"column:total_sent"`
}

func NewBalance() *Balance {
  return &Balance{}
}

func (c *Balance) GetLastData(address string, trdb *gorm.DB) (bool) {
  var trx TransactionDB
  result := trdb.Raw("SELECT address, sum(Vin), sum(Vout) FROM transactions WHERE address = ? GROUP BY address", address).Scan(&trx)
  if result.Error != nil {
    return false
  }
  c.Balance = trx.Vin - trx.Vout
  c.UnconfirmedBalance = trx.Vin - trx.Vout
  c.TotalReceived = trx.Vin
  c.TotalSent = trx.Vout
  return true
}
