package objects

import (
  "testing"
  "github.com/stretchr/testify/assert"

  "github.com/Lunkov/go-hdwallet"
  "github.com/Lunkov/lib-wallets"
)

func TestMsgWalletTx(t *testing.T) {
  w1 := wallets.NewWallet(wallets.TypeWalletHD)
  w1.Create(&map[string]string{"mnemonic": "chase oil pigeon elegant ketchup whip frozen beauty unknown brass amount slender pony pottery attitude flavor rifle primary beach sign glue oven crazy lottery"})

  tx := NewWalletTransactionsBalance(w1.GetAddress(hdwallet.ECOS))
  
  tx.Transactions = append(tx.Transactions, WalletTransaction{IdTx: []byte{0x04}, Timestamp: 4, DirectionCoins: DirectionInput,  Amount: 10000})
  tx.Transactions = append(tx.Transactions, WalletTransaction{IdTx: []byte{0x01}, Timestamp: 1, DirectionCoins: DirectionInput,  Amount: 10000})
  tx.Transactions = append(tx.Transactions, WalletTransaction{IdTx: []byte{0x03}, Timestamp: 3, DirectionCoins: DirectionOutput, Amount: 10000})
  tx.Transactions = append(tx.Transactions, WalletTransaction{IdTx: []byte{0x02}, Timestamp: 2, DirectionCoins: DirectionInput,  Amount: 10000})
  tx.Transactions = append(tx.Transactions, WalletTransaction{IdTx: []byte{0x05}, Timestamp: 5, DirectionCoins: DirectionInput,  Amount: 10000})
  tx.Transactions = append(tx.Transactions, WalletTransaction{IdTx: []byte{0x07}, Timestamp: 7, DirectionCoins: DirectionOutput, Amount: 10000})
  tx.Transactions = append(tx.Transactions, WalletTransaction{IdTx: []byte{0x06}, Timestamp: 6, DirectionCoins: DirectionInput,  Amount: 10000})
  tx.Transactions = append(tx.Transactions, WalletTransaction{IdTx: []byte{0x08}, Timestamp: 8, DirectionCoins: DirectionOutput, Amount: 10000})
  
  in, out := tx.GetStat()
  assert.Equal(t, uint64(50000), in) 
  assert.Equal(t, uint64(30000), out) 
  
  tx.RecalcBalance(uint64(5000))
  assert.Equal(t, []WalletTransaction{
                         WalletTransaction{IdTx: []byte{0x01}, Timestamp: 1, DirectionCoins: DirectionInput,   Amount: 10000,  Balance: 15000},
                         WalletTransaction{IdTx: []byte{0x02}, Timestamp: 2, DirectionCoins: DirectionInput,   Amount: 10000,  Balance: 25000},
                         WalletTransaction{IdTx: []byte{0x03}, Timestamp: 3, DirectionCoins: DirectionOutput,  Amount: 10000,  Balance: 15000},
                         WalletTransaction{IdTx: []byte{0x04}, Timestamp: 4, DirectionCoins: DirectionInput,   Amount: 10000,  Balance: 25000},
                         WalletTransaction{IdTx: []byte{0x05}, Timestamp: 5, DirectionCoins: DirectionInput,   Amount: 10000,  Balance: 35000},
                         WalletTransaction{IdTx: []byte{0x06}, Timestamp: 6, DirectionCoins: DirectionInput,   Amount: 10000,  Balance: 45000},
                         WalletTransaction{IdTx: []byte{0x07}, Timestamp: 7, DirectionCoins: DirectionOutput,  Amount: 10000,  Balance: 35000},
                         WalletTransaction{IdTx: []byte{0x08}, Timestamp: 8, DirectionCoins: DirectionOutput,  Amount: 10000,  Balance: 25000},
                    }, tx.Transactions) 
  
  tx2 := NewWalletTransactionsEmpty()

  buf, err := tx.Serialize()
  assert.Nil(t, err) 
  assert.Nil(t, tx2.Deserialize(buf)) 
  
  
  assert.Equal(t, []WalletTransaction{
                         WalletTransaction{IdTx: []byte{0x01}, Timestamp: 1, DirectionCoins: DirectionInput,   Amount: 10000,  Balance: 15000},
                         WalletTransaction{IdTx: []byte{0x02}, Timestamp: 2, DirectionCoins: DirectionInput,   Amount: 10000,  Balance: 25000},
                         WalletTransaction{IdTx: []byte{0x03}, Timestamp: 3, DirectionCoins: DirectionOutput,  Amount: 10000,  Balance: 15000},
                         WalletTransaction{IdTx: []byte{0x04}, Timestamp: 4, DirectionCoins: DirectionInput,   Amount: 10000,  Balance: 25000},
                         WalletTransaction{IdTx: []byte{0x05}, Timestamp: 5, DirectionCoins: DirectionInput,   Amount: 10000,  Balance: 35000},
                         WalletTransaction{IdTx: []byte{0x06}, Timestamp: 6, DirectionCoins: DirectionInput,   Amount: 10000,  Balance: 45000},
                         WalletTransaction{IdTx: []byte{0x07}, Timestamp: 7, DirectionCoins: DirectionOutput,  Amount: 10000,  Balance: 35000},
                         WalletTransaction{IdTx: []byte{0x08}, Timestamp: 8, DirectionCoins: DirectionOutput,  Amount: 10000,  Balance: 25000},
                    }, tx2.Transactions) 

  tx3 := NewWalletTransactionsBalance(w1.GetAddress(hdwallet.ECOS))
  tx3.Transactions = append(tx3.Transactions, WalletTransaction{IdTx: []byte{0x16}, Timestamp: 16, DirectionCoins: DirectionOutput, Amount: 5000})
  tx3.Transactions = append(tx3.Transactions, WalletTransaction{IdTx: []byte{0x18}, Timestamp: 18, DirectionCoins: DirectionOutput, Amount: 10000})
  tx3.RecalcBalance(uint64(25000))

  tx2.Append(tx3) 
  
  assert.Equal(t, []WalletTransaction{
                         WalletTransaction{IdTx: []byte{0x01}, Timestamp: 1,  DirectionCoins: DirectionInput,   Amount: 10000,  Balance: 15000},
                         WalletTransaction{IdTx: []byte{0x02}, Timestamp: 2,  DirectionCoins: DirectionInput,   Amount: 10000,  Balance: 25000},
                         WalletTransaction{IdTx: []byte{0x03}, Timestamp: 3,  DirectionCoins: DirectionOutput,  Amount: 10000,  Balance: 15000},
                         WalletTransaction{IdTx: []byte{0x04}, Timestamp: 4,  DirectionCoins: DirectionInput,   Amount: 10000,  Balance: 25000},
                         WalletTransaction{IdTx: []byte{0x05}, Timestamp: 5,  DirectionCoins: DirectionInput,   Amount: 10000,  Balance: 35000},
                         WalletTransaction{IdTx: []byte{0x06}, Timestamp: 6,  DirectionCoins: DirectionInput,   Amount: 10000,  Balance: 45000},
                         WalletTransaction{IdTx: []byte{0x07}, Timestamp: 7,  DirectionCoins: DirectionOutput,  Amount: 10000,  Balance: 35000},
                         WalletTransaction{IdTx: []byte{0x08}, Timestamp: 8,  DirectionCoins: DirectionOutput,  Amount: 10000,  Balance: 25000},
                         WalletTransaction{IdTx: []byte{0x16}, Timestamp: 16, DirectionCoins: DirectionOutput,  Amount: 5000,   Balance: 20000},
                         WalletTransaction{IdTx: []byte{0x18}, Timestamp: 18, DirectionCoins: DirectionOutput,  Amount: 10000,  Balance: 10000},
                    }, tx2.Transactions)

  tx4 := NewWalletTransactionsBalance(w1.GetAddress(hdwallet.ECOS))
  assert.False(t, tx4.Have("0x000123456789")) 
  tx4.Transactions = append(tx4.Transactions, WalletTransaction{IdTx: []byte{0x20}, Timestamp: 20, DirectionNFT: DirectionInput,  CIDNFT: "0x000123456789", Amount: 0})
  assert.True(t, tx4.Have("0x000123456789")) 
  tx4.Transactions = append(tx4.Transactions, WalletTransaction{IdTx: []byte{0x21}, Timestamp: 21, DirectionNFT: DirectionOutput, CIDNFT: "0x000123456789", Amount: 0})
  assert.False(t, tx4.Have("0x000123456789")) 
  tx4.RecalcBalance(uint64(10000))
  
  tx2.Append(tx4)
  tx2.SortByDate() 

  assert.Equal(t, []WalletTransaction{
                         WalletTransaction{IdTx: []byte{0x01}, Timestamp: 1,  DirectionCoins: DirectionInput,   Amount: 10000,  Balance: 15000},
                         WalletTransaction{IdTx: []byte{0x02}, Timestamp: 2,  DirectionCoins: DirectionInput,   Amount: 10000,  Balance: 25000},
                         WalletTransaction{IdTx: []byte{0x03}, Timestamp: 3,  DirectionCoins: DirectionOutput,  Amount: 10000,  Balance: 15000},
                         WalletTransaction{IdTx: []byte{0x04}, Timestamp: 4,  DirectionCoins: DirectionInput,   Amount: 10000,  Balance: 25000},
                         WalletTransaction{IdTx: []byte{0x05}, Timestamp: 5,  DirectionCoins: DirectionInput,   Amount: 10000,  Balance: 35000},
                         WalletTransaction{IdTx: []byte{0x06}, Timestamp: 6,  DirectionCoins: DirectionInput,   Amount: 10000,  Balance: 45000},
                         WalletTransaction{IdTx: []byte{0x07}, Timestamp: 7,  DirectionCoins: DirectionOutput,  Amount: 10000,              Balance: 35000},
                         WalletTransaction{IdTx: []byte{0x08}, Timestamp: 8,  DirectionCoins: DirectionOutput,  Amount: 10000,              Balance: 25000},
                         WalletTransaction{IdTx: []byte{0x16}, Timestamp: 16, DirectionCoins: DirectionOutput,  Amount: 5000,               Balance: 20000},
                         WalletTransaction{IdTx: []byte{0x18}, Timestamp: 18, DirectionCoins: DirectionOutput,  Amount: 10000,              Balance: 10000},
                         WalletTransaction{IdTx: []byte{0x20}, Timestamp: 20, DirectionNFT:   DirectionInput,   CIDNFT: "0x000123456789",   Balance: 10000},
                         WalletTransaction{IdTx: []byte{0x21}, Timestamp: 21, DirectionNFT:   DirectionOutput,  CIDNFT: "0x000123456789",   Balance: 10000},
                    }, tx2.Transactions)

  assert.Equal(t, WalletTransactions{
                   Address: "0x5f7ae710cED588D42E863E9b55C7c51e56869963",
                   Transactions: []WalletTransaction{
                         WalletTransaction{IdTx: []byte{0x20}, Timestamp: 20, DirectionNFT:   DirectionInput,   CIDNFT: "0x000123456789",   Balance: 10000},
                         WalletTransaction{IdTx: []byte{0x21}, Timestamp: 21, DirectionNFT:   DirectionOutput,  CIDNFT: "0x000123456789",   Balance: 10000},
                    }}, *tx2.FilterNFT("0x000123456789"))

  assert.Equal(t, []TXInput{}, tx2.FindCoins(20000))
  assert.Equal(t, []TXInput{TXInput{Txid:[]uint8{0x6}, Address:"0x5f7ae710cED588D42E863E9b55C7c51e56869963", CIDNFT:"", Vout:uint64(500)}}, tx2.FindCoins(500))

  tx5 := NewWalletTransactionsBalance(w1.GetAddress(hdwallet.ECOS))
  tx5.Transactions = append(tx5.Transactions, WalletTransaction{IdTx: []byte{0x26}, Timestamp: 26, DirectionCoins: DirectionInput, Amount: 500})
  tx5.Transactions = append(tx5.Transactions, WalletTransaction{IdTx: []byte{0x28}, Timestamp: 28, DirectionCoins: DirectionInput, Amount: 1000})
  tx5.RecalcBalance(uint64(10000))

  tx2.Append(tx5) 

  assert.Equal(t, []TXInput{
                             TXInput{Txid:[]uint8{0x6}, Address:"0x5f7ae710cED588D42E863E9b55C7c51e56869963", CIDNFT:"", Vout:uint64(10000)},
                             TXInput{Txid:[]uint8{0x26}, Address:"0x5f7ae710cED588D42E863E9b55C7c51e56869963", CIDNFT:"", Vout:uint64(250)},
                           }, tx2.FindCoins(10250))

  assert.Equal(t, []TXInput{
                             TXInput{Txid:[]uint8{0x6}, Address:"0x5f7ae710cED588D42E863E9b55C7c51e56869963", CIDNFT:"", Vout:uint64(10000)},
                             TXInput{Txid:[]uint8{0x26}, Address:"0x5f7ae710cED588D42E863E9b55C7c51e56869963", CIDNFT:"", Vout:uint64(500)},
                             TXInput{Txid:[]uint8{0x28}, Address:"0x5f7ae710cED588D42E863E9b55C7c51e56869963", CIDNFT:"", Vout:uint64(250)},
                           }, tx2.FindCoins(10750))

}
