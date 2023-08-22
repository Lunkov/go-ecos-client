package objects

import (
  "testing"
  "github.com/stretchr/testify/assert"

  "github.com/Lunkov/go-ecos-client/utils"
)

func TestMsgPassport(t *testing.T) {
  key, err := utils.RSAGenerate(2048)
  assert.Nil(t, err)
  
  msg := NewPassport()
  msg.DisplayName = "Test User"
  
  msg.IssuerPublicKey, err = utils.RSASerializePublicKey(&key.PublicKey)
  assert.Nil(t, err)
  
  msg.IssuerSign, err = utils.RSASign(key, msg.HashIssuer())
  assert.Nil(t, err)
  
  buf, okb := msg.Serialize()
  assert.Nil(t, okb) 
  
  msg2 := NewPassport()
  
  assert.Nil(t, msg2.Deserialize(buf)) 
  
  verr := utils.RSAVerify(msg.IssuerPublicKey, msg.HashIssuer(), msg.IssuerSign)
  assert.Nil(t, verr)
  
  assert.Equal(t, "Test User", msg2.DisplayName)
  assert.Equal(t, *msg, *msg2)
  
  benc, berr := msg.SerializeEncrypt("password")
  assert.Nil(t, berr) 
  msg3 := NewPassport()
  berr = msg3.DeserializeDecrypt("password", benc)
  assert.Nil(t, berr) 
  assert.Equal(t, *msg, *msg3)
}
