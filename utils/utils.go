package utils

import (
  "strconv"
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


func UInt64Format(n uint64) string {
  in := strconv.FormatUint(n, 10)
  numOfCommas := (len(in) - 1) / 3

  out := make([]byte, len(in) + numOfCommas)

  for i, j, k := len(in)-1, len(out)-1, 0; ; i, j = i-1, j-1 {
    out[j] = in[i]
    if i == 0 {
      return string(out)
    }
    if k++; k == 3 {
      j, k = j-1, 0
      out[j] = ','
    }
  }
}

