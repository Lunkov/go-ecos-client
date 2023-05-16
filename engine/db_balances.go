package main

import (
  "sync"
  "github.com/Lunkov/lib-messages"
)

type Balances struct {
  db                 *DBO
	items               map[string]*Balance
	mu                  sync.RWMutex
}

func NewBalances(db *DBO) *Balances {
  return &Balances{ db: db, items: make(map[string]*Balance) }
}

func (c *Balances) MintCoin(address string, value uint64) {
  c.mu.Lock()
  item, okf := c.items[address]
  if !okf {
		c.items[address] = &Balance{
                                 Balance: value,
                                 UnconfirmedBalance: value,
                                 TotalReceived: value,
                                 TotalSent: 0,
                                 }
	} else {
    item.Balance += value
    item.UnconfirmedBalance += value
    item.TotalReceived += value
  }
  c.mu.Unlock()
}

func (c *Balances) Get(address string) (*messages.Balance, bool)  {
  c.mu.RLock()
  item, ok := c.items[address]
  defer c.mu.RUnlock()
  if ok {
    return &messages.Balance{
               Address: address,
               Balance: item.Balance,
               UnconfirmedBalance: item.UnconfirmedBalance,
               TotalReceived: item.TotalReceived,
               TotalSent: item.TotalSent,
      }, true
  }
  
  b := NewBalance()
  b.GetLastData(address, c.db.GetHandle())
  return &messages.Balance{
           Address: address,
           Balance: b.Balance,
           UnconfirmedBalance: b.UnconfirmedBalance,
           TotalReceived: b.TotalReceived,
           TotalSent: b.TotalSent,
         }, true
}

func (c *Balances) Transfer(addressFrom string, addressTo string, value uint64) (bool)  {
  c.mu.Lock()
  defer c.mu.Unlock()
  itemFrom, okf := c.items[addressFrom]
  if !okf {
		return false
	}
  if itemFrom.Balance < value {
		return false
	}
  itemTo, okt := c.items[addressTo]
  if !okt {
    c.items[addressTo] = &Balance{
                               Balance: value,
                               UnconfirmedBalance: value,
                               TotalReceived: value,
                               TotalSent: 0,
                             }
	} else {
    itemTo.Balance += value
    itemTo.UnconfirmedBalance += value
    itemTo.TotalReceived += value
  }
  itemFrom.Balance -= value
  itemFrom.UnconfirmedBalance -= value
  itemFrom.TotalSent += value
	return true
}
