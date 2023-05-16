package engine

import (
  "testing"
  "github.com/stretchr/testify/assert"
  
  "github.com/Lunkov/lib-messages"
)

func TestBalances(t *testing.T) {
	
  b1 := NewBalances(nil)
  
  b1.MintCoin("a1", 10000)

  res1, ok1 := b1.Get("b2")
  assert.False(t, ok1)
  assert.Nil(t, res1)

  res2, ok2 := b1.Get("a1")
  assert.True(t, ok2)
  assert.Equal(t, &messages.Balance{Address: "a1", Balance: 10000, TotalReceived: 10000, UnconfirmedBalance: 10000}, res2)
  
  okt := b1.Transfer("b1", "b2", 100)
  assert.False(t, okt)
  okt = b1.Transfer("a1", "b2", 100000)
  assert.False(t, okt)
  okt = b1.Transfer("a1", "b2", 1000)
  assert.True(t, okt)

  res3, ok3 := b1.Get("a1")
  assert.True(t, ok3)
  assert.Equal(t, &messages.Balance{Address: "a1", Balance: 9000, TotalReceived: 10000, UnconfirmedBalance: 9000, TotalSent: 1000}, res3)
  
  res4, ok4 := b1.Get("b2")
  assert.True(t, ok4)
  assert.Equal(t, &messages.Balance{Address: "b2", Balance: 1000, TotalReceived: 1000, UnconfirmedBalance: 1000}, res4)

  okt = b1.Transfer("b2", "c1", 700)
  assert.True(t, okt)

  res4, ok4 = b1.Get("b2")
  assert.True(t, ok4)
  assert.Equal(t, &messages.Balance{Address: "b2", Balance: 300, TotalReceived: 1000, UnconfirmedBalance: 300, TotalSent: 700}, res4)


  okt = b1.Transfer("b2", "c2", 100)
  assert.True(t, okt)

  res4, ok4 = b1.Get("b2")
  assert.True(t, ok4)
  assert.Equal(t, &messages.Balance{Address: "b2", Balance: 200, TotalReceived: 1000, UnconfirmedBalance: 200, TotalSent: 800}, res4)

}
