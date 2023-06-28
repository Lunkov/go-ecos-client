package utils

import (
  "testing"
  "github.com/stretchr/testify/assert"
)

func TestCA(t *testing.T) {
  ca := NewCertInfo()
  ca.Bits = 2048
  pwd := "12345"
  okc := ca.СreateNewCA("./test/test.cert", "./test/test.priv", pwd)
  assert.True(t, okc)
  
  ca.Load("./test/test.cert", "./test/test.priv", pwd)
  
  sig, ok := ca.Sign([]byte("1234567890"))
  
  assert.True(t, ok)
  //assert.Equal(t, "", fmt.Sprintf("%x", sig))

  ok = ca.Verify([]byte("1234567890"), sig)
  assert.True(t, ok)
  
  bufCert, bok := ca.SerializeCert()
  assert.True(t, bok)
  
  ca2 := NewCertInfo()
  
  ok = ca2.DeserializeCert(bufCert)
  assert.True(t, ok)
  
  ok = ca2.Verify([]byte("1234567890"), sig)
  assert.True(t, ok)
  
}

