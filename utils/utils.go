package utils

import (
)

// ReverseBytes reverses a byte array
func ReverseBytes(data *[]byte) {
  for i, j := 0, len(*data)-1; i < j; i, j = i+1, j-1 {
    (*data)[i], (*data)[j] = (*data)[j], (*data)[i]
  }
}

func UInt64ToBytes(v uint64) (b []byte) {
  l := 8
  b = make([]byte, l)

  for i := 0; i < l; i++ {
    f := 8 * i
    b[i] = byte(v >> f) 
  }
  return b
}

func UInt32ToBytes(v uint32) (b []byte) {
  l := 4
  b = make([]byte, l)

  for i := 0; i < l; i++ {
    f := 8 * i
    b[i] = byte(v >> f)
  }
  return b
}

func Int64ToBytes(v int64) (b []byte) {
  l := 8
  b = make([]byte, l)

  for i := 0; i < l; i++ {
    f := 8 * i
    b[i] = byte(v >> f)
  }
  return b
}


