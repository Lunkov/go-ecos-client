package utils

import (
  "bytes"
  "encoding/gob"
  "crypto"
  "crypto/rand"
  "math/big"
  "crypto/sha512"
  "crypto/rsa"
  "github.com/golang/glog"
)

// const version = byte(0x00)

type RSABuf struct {
  N []byte
  E int
}

func RSASerialize(pk *rsa.PublicKey) ([]byte, bool) {
  rsabuf := RSABuf{N: pk.N.Bytes(), E: pk.E}
  var buff bytes.Buffer
  encoder := gob.NewEncoder(&buff)
  encoder.Encode(rsabuf)
  return buff.Bytes(), true
}

func RSADeserialize(msg []byte) (*rsa.PublicKey, bool) {
  var rsabuf RSABuf
  buf := bytes.NewBuffer(msg)
  decoder := gob.NewDecoder(buf)
  err := decoder.Decode(&rsabuf)
  if err != nil {
    glog.Errorf("ERR: gob.Decode %v", err)
    return nil, false
  }
  n := big.Int{}
  n.SetBytes(rsabuf.N)
  return &rsa.PublicKey{E: rsabuf.E, N: &n}, true
}

func RSAVerify(pubKey []byte, message []byte, signature []byte) bool {
  defer func() {
    if r := recover(); r != nil {
      glog.Errorf("ERR: ECDSA signature verification error: %v", r)
    }
  }()
  rawPubKey, ok := RSADeserialize(pubKey) //PublicKeyFromBytes(pubKey)
  if !ok {
    return false
  }
  hashed := sha512.Sum512(message)
  return rsa.VerifyPKCS1v15(rawPubKey, crypto.SHA512, hashed[:], signature) == nil
}

func RSASign(pk *rsa.PrivateKey, message []byte) ([]byte, bool) {
  hashed := sha512.Sum512(message)
  signature, err := rsa.SignPKCS1v15(rand.Reader, pk, crypto.SHA512, hashed[:])
  return signature, err == nil 
}

