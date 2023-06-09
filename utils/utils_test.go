package utils

import (
  "testing"
  "github.com/stretchr/testify/assert"
)

func TestUtils(t *testing.T) {
  
  ui64buf := Int64ToBytes(1234567890)
  assert.Equal(t, []byte{0xd2, 0x2, 0x96, 0x49, 0x0, 0x0, 0x0, 0x0}, ui64buf)
  ReverseBytes(&ui64buf)
  assert.Equal(t, []byte{0x0, 0x0, 0x0, 0x0, 0x49, 0x96, 0x2, 0xd2}, ui64buf)
  assert.Equal(t, []byte{0xd2, 0x2, 0x96, 0x49, 0x0, 0x0, 0x0, 0x0}, UInt64ToBytes(1234567890))
  assert.Equal(t, []byte{0xd2, 0x2, 0x96, 0x49}, UInt32ToBytes(1234567890))
  
  assert.Equal(t, "1,234,567,890", UInt64Format(1234567890))
}
