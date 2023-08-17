package utils

import (
  "bytes"
  "encoding/gob"
  "crypto"
  "crypto/rand"
  "math/big"
  "crypto/sha512"
  "crypto/rsa"
)

// const version = byte(0x00)

type RSABuf struct {
  N []byte
  E int
}

func RSAGenerate(bits int) (*rsa.PrivateKey, bool) {
  caPrivKey, errc := rsa.GenerateKey(rand.Reader, bits)
  if errc != nil {
    return nil, false
  }
  return caPrivKey, true
}

func RSASerializePublicKey(pk *rsa.PublicKey) ([]byte, bool) {
  rsabuf := RSABuf{N: pk.N.Bytes(), E: pk.E}
  var buff bytes.Buffer
  encoder := gob.NewEncoder(&buff)
  encoder.Encode(rsabuf)
  return buff.Bytes(), true
}

func RSADeserializePublicKey(msg []byte) (*rsa.PublicKey, bool) {
  var rsabuf RSABuf
  buf := bytes.NewBuffer(msg)
  decoder := gob.NewDecoder(buf)
  err := decoder.Decode(&rsabuf)
  if err != nil {
    return nil, false
  }
  n := big.Int{}
  n.SetBytes(rsabuf.N)
  return &rsa.PublicKey{E: rsabuf.E, N: &n}, true
}

func RSAVerify(pubKey []byte, message []byte, signature []byte) bool {
  defer func() {
    if r := recover(); r != nil {
    }
  }()
  rawPubKey, ok := RSADeserializePublicKey(pubKey) //PublicKeyFromBytes(pubKey)
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

