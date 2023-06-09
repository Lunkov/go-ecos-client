package objects

import (
  "testing"
  "github.com/stretchr/testify/assert"

  "github.com/Lunkov/go-ecos-client/utils"
)

func TestMsgPassport(t *testing.T) {
  key, ok := utils.RSAGenerate(2048)
  assert.True(t, ok)
  
  msg := NewPassport()
  msg.DisplayName = "Test User"
  
  msg.IssuerPublicKey, ok = utils.RSASerializePublicKey(&key.PublicKey)
  assert.True(t, ok)
  
  msg.IssuerSign, ok = utils.RSASign(key, msg.HashIssuer())
  assert.True(t, ok)
  
  buf, okb := msg.Serialize()
  assert.True(t, okb) 
  
  msg2 := NewPassport()
  
  assert.True(t, msg2.Deserialize(buf)) 
  
  assert.True(t, utils.RSAVerify(msg.IssuerPublicKey, msg.HashIssuer(), msg.IssuerSign))
  
  assert.Equal(t, "Test User", msg2.DisplayName)
  assert.Equal(t, *msg, *msg2)
  
  benc, bok := msg.SerializeEncrypt("password")
  assert.True(t, bok) 
  msg3 := NewPassport()
  bok = msg3.DeserializeDecrypt("password", benc)
  assert.True(t, bok) 
  assert.Equal(t, *msg, *msg3)
}
