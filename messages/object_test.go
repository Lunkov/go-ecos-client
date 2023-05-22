package messages

import (
  "testing"
  "github.com/stretchr/testify/assert"
)

func TestMessageObject(t *testing.T) {
	
  obj := NewReqActionObject()
  obj.Init(1, "user", []byte{})
  
  obj.Hash()
  
  buf := obj.Serialize()

  assert.Equal(t, 188, len(buf))
}
