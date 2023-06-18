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
  
  tx.Transactions = append(tx.Transactions, WalletTransaction{Timestamp: 4, DirectionCoins: DirectionInput,  Amount: 10000})
  tx.Transactions = append(tx.Transactions, WalletTransaction{Timestamp: 1, DirectionCoins: DirectionInput,  Amount: 10000})
  tx.Transactions = append(tx.Transactions, WalletTransaction{Timestamp: 3, DirectionCoins: DirectionOutput, Amount: 10000})
  tx.Transactions = append(tx.Transactions, WalletTransaction{Timestamp: 2, DirectionCoins: DirectionInput,  Amount: 10000})
  tx.Transactions = append(tx.Transactions, WalletTransaction{Timestamp: 5, DirectionCoins: DirectionInput,  Amount: 10000})
  tx.Transactions = append(tx.Transactions, WalletTransaction{Timestamp: 7, DirectionCoins: DirectionOutput, Amount: 10000})
  tx.Transactions = append(tx.Transactions, WalletTransaction{Timestamp: 6, DirectionCoins: DirectionInput,  Amount: 10000})
  tx.Transactions = append(tx.Transactions, WalletTransaction{Timestamp: 8, DirectionCoins: DirectionOutput, Amount: 10000})
  
  in, out := tx.GetStat()
  assert.Equal(t, uint64(50000), in) 
  assert.Equal(t, uint64(30000), out) 
  
  tx.RecalcBalance(uint64(5000))
  assert.Equal(t, []WalletTransaction{
                         WalletTransaction{Timestamp: 1, DirectionCoins: DirectionInput,   Amount: 10000,  Balance: 15000},
                         WalletTransaction{Timestamp: 2, DirectionCoins: DirectionInput,   Amount: 10000,  Balance: 25000},
                         WalletTransaction{Timestamp: 3, DirectionCoins: DirectionOutput,  Amount: 10000,  Balance: 15000},
                         WalletTransaction{Timestamp: 4, DirectionCoins: DirectionInput,   Amount: 10000,  Balance: 25000},
                         WalletTransaction{Timestamp: 5, DirectionCoins: DirectionInput,   Amount: 10000,  Balance: 35000},
                         WalletTransaction{Timestamp: 6, DirectionCoins: DirectionInput,   Amount: 10000,  Balance: 45000},
                         WalletTransaction{Timestamp: 7, DirectionCoins: DirectionOutput,  Amount: 10000,  Balance: 35000},
                         WalletTransaction{Timestamp: 8, DirectionCoins: DirectionOutput,  Amount: 10000,  Balance: 25000},
                    }, tx.Transactions) 
  
  tx2 := NewWalletTransactionsEmpty()

  buf, ok := tx.Serialize()
  assert.True(t, ok) 
  assert.True(t, tx2.Deserialize(buf)) 
  
  
  assert.Equal(t, []WalletTransaction{
                         WalletTransaction{Timestamp: 1, DirectionCoins: DirectionInput,   Amount: 10000,  Balance: 15000},
                         WalletTransaction{Timestamp: 2, DirectionCoins: DirectionInput,   Amount: 10000,  Balance: 25000},
                         WalletTransaction{Timestamp: 3, DirectionCoins: DirectionOutput,  Amount: 10000,  Balance: 15000},
                         WalletTransaction{Timestamp: 4, DirectionCoins: DirectionInput,   Amount: 10000,  Balance: 25000},
                         WalletTransaction{Timestamp: 5, DirectionCoins: DirectionInput,   Amount: 10000,  Balance: 35000},
                         WalletTransaction{Timestamp: 6, DirectionCoins: DirectionInput,   Amount: 10000,  Balance: 45000},
                         WalletTransaction{Timestamp: 7, DirectionCoins: DirectionOutput,  Amount: 10000,  Balance: 35000},
                         WalletTransaction{Timestamp: 8, DirectionCoins: DirectionOutput,  Amount: 10000,  Balance: 25000},
                    }, tx2.Transactions) 

  tx3 := NewWalletTransactionsBalance(w1.GetAddress(hdwallet.ECOS))
  tx3.Transactions = append(tx3.Transactions, WalletTransaction{Timestamp: 16, DirectionCoins: DirectionOutput, Amount: 5000})
  tx3.Transactions = append(tx3.Transactions, WalletTransaction{Timestamp: 18, DirectionCoins: DirectionOutput, Amount: 10000})
  tx3.RecalcBalance(uint64(25000))

  tx2.Append(tx3) 
  
  assert.Equal(t, []WalletTransaction{
                         WalletTransaction{Timestamp: 1,  DirectionCoins: DirectionInput,   Amount: 10000,  Balance: 15000},
                         WalletTransaction{Timestamp: 2,  DirectionCoins: DirectionInput,   Amount: 10000,  Balance: 25000},
                         WalletTransaction{Timestamp: 3,  DirectionCoins: DirectionOutput,  Amount: 10000,  Balance: 15000},
                         WalletTransaction{Timestamp: 4,  DirectionCoins: DirectionInput,   Amount: 10000,  Balance: 25000},
                         WalletTransaction{Timestamp: 5,  DirectionCoins: DirectionInput,   Amount: 10000,  Balance: 35000},
                         WalletTransaction{Timestamp: 6,  DirectionCoins: DirectionInput,   Amount: 10000,  Balance: 45000},
                         WalletTransaction{Timestamp: 7,  DirectionCoins: DirectionOutput,  Amount: 10000,  Balance: 35000},
                         WalletTransaction{Timestamp: 8,  DirectionCoins: DirectionOutput,  Amount: 10000,  Balance: 25000},
                         WalletTransaction{Timestamp: 16, DirectionCoins: DirectionOutput,  Amount: 5000,   Balance: 20000},
                         WalletTransaction{Timestamp: 18, DirectionCoins: DirectionOutput,  Amount: 10000,  Balance: 10000},
                    }, tx2.Transactions)

  tx4 := NewWalletTransactionsBalance(w1.GetAddress(hdwallet.ECOS))
  assert.False(t, tx4.Have("0x000123456789")) 
  tx4.Transactions = append(tx4.Transactions, WalletTransaction{Timestamp: 20, DirectionNFT: DirectionInput,  CIDNFT: "0x000123456789", Amount: 0})
  assert.True(t, tx4.Have("0x000123456789")) 
  tx4.Transactions = append(tx4.Transactions, WalletTransaction{Timestamp: 21, DirectionNFT: DirectionOutput, CIDNFT: "0x000123456789", Amount: 0})
  assert.False(t, tx4.Have("0x000123456789")) 
  tx4.RecalcBalance(uint64(10000))
  
  tx2.Append(tx4)
  tx2.SortByDate() 

  assert.Equal(t, []WalletTransaction{
                         WalletTransaction{Timestamp: 1,  DirectionCoins: DirectionInput,   Amount: 10000,  Balance: 15000},
                         WalletTransaction{Timestamp: 2,  DirectionCoins: DirectionInput,   Amount: 10000,  Balance: 25000},
                         WalletTransaction{Timestamp: 3,  DirectionCoins: DirectionOutput,  Amount: 10000,  Balance: 15000},
                         WalletTransaction{Timestamp: 4,  DirectionCoins: DirectionInput,   Amount: 10000,  Balance: 25000},
                         WalletTransaction{Timestamp: 5,  DirectionCoins: DirectionInput,   Amount: 10000,  Balance: 35000},
                         WalletTransaction{Timestamp: 6,  DirectionCoins: DirectionInput,   Amount: 10000,  Balance: 45000},
                         WalletTransaction{Timestamp: 7,  DirectionCoins: DirectionOutput,  Amount: 10000,              Balance: 35000},
                         WalletTransaction{Timestamp: 8,  DirectionCoins: DirectionOutput,  Amount: 10000,              Balance: 25000},
                         WalletTransaction{Timestamp: 16, DirectionCoins: DirectionOutput,  Amount: 5000,               Balance: 20000},
                         WalletTransaction{Timestamp: 18, DirectionCoins: DirectionOutput,  Amount: 10000,              Balance: 10000},
                         WalletTransaction{Timestamp: 20, DirectionNFT:   DirectionInput,   CIDNFT: "0x000123456789",   Balance: 10000},
                         WalletTransaction{Timestamp: 21, DirectionNFT:   DirectionOutput,  CIDNFT: "0x000123456789",   Balance: 10000},
                    }, tx2.Transactions)

  assert.Equal(t, WalletTransactions{
                   Address: "0x5f7ae710cED588D42E863E9b55C7c51e56869963",
                   Transactions: []WalletTransaction{
                         WalletTransaction{Timestamp: 20, DirectionNFT:   DirectionInput,   CIDNFT: "0x000123456789",   Balance: 10000},
                         WalletTransaction{Timestamp: 21, DirectionNFT:   DirectionOutput,  CIDNFT: "0x000123456789",   Balance: 10000},
                    }}, *tx2.FilterNFT("0x000123456789"))

}
