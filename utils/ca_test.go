package utils

import (
  "testing"
  "github.com/stretchr/testify/assert"
)

func TestCA(t *testing.T) {
  ca := NewCertInfo()
  ca.Bits = 2048
  pwd := "12345"
  errc := ca.СreateNewCA("./test/test.cert", "./test/test.priv", pwd)
  assert.Nil(t, errc)
  
  ca.Load("./test/test.cert", "./test/test.priv", pwd)
  
  sig, err := ca.Sign([]byte("1234567890"))
  
  assert.Nil(t, err)
  //assert.Equal(t, "", fmt.Sprintf("%x", sig))

  err = ca.Verify([]byte("1234567890"), sig)
  assert.Nil(t, err)
  
  bufCert, berr := ca.SerializeCert()
  assert.Nil(t, berr)
  
  ca2 := NewCertInfo()
  
  err = ca2.DeserializeCert(bufCert)
  assert.Nil(t, err)
  
  err = ca2.Verify([]byte("1234567890"), sig)
  assert.Nil(t, err)
  
  
  ca1sub := NewCertInfo()
  ca1sub.Bits = 2048
  ca1sub.EMail = "user@myorg.ru"
  
  buf1sub, errs := ca.CreateSubCert(ca1sub)
  assert.Nil(t, errs)
  sig, err = ca1sub.Sign([]byte("1234567890"))

  ca2sub := NewCertInfo()
  
  err = ca2sub.DeserializeCert(buf1sub)
  assert.Nil(t, err)
  err = ca2sub.Verify([]byte("1234567890"), sig)
  assert.Nil(t, err)

}

