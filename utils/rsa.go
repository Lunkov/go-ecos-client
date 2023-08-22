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

func RSAGenerate(bits int) (*rsa.PrivateKey, error) {
  caPrivKey, err := rsa.GenerateKey(rand.Reader, bits)
  if err != nil {
    return nil, err
  }
  return caPrivKey, nil
}

func RSASerializePublicKey(pk *rsa.PublicKey) ([]byte, error) {
  rsabuf := RSABuf{N: pk.N.Bytes(), E: pk.E}
  var buff bytes.Buffer
  encoder := gob.NewEncoder(&buff)
  encoder.Encode(rsabuf)
  return buff.Bytes(), nil
}

func RSADeserializePublicKey(msg []byte) (*rsa.PublicKey, error) {
  var rsabuf RSABuf
  buf := bytes.NewBuffer(msg)
  decoder := gob.NewDecoder(buf)
  err := decoder.Decode(&rsabuf)
  if err != nil {
    return nil, err
  }
  n := big.Int{}
  n.SetBytes(rsabuf.N)
  return &rsa.PublicKey{E: rsabuf.E, N: &n}, nil
}

func RSAVerify(pubKey []byte, message []byte, signature []byte) (error) {
  defer func() {
    if r := recover(); r != nil {
    }
  }()
  rawPubKey, err := RSADeserializePublicKey(pubKey) //PublicKeyFromBytes(pubKey)
  if err != nil {
    return err
  }
  hashed := sha512.Sum512(message)
  return rsa.VerifyPKCS1v15(rawPubKey, crypto.SHA512, hashed[:], signature)
}

func RSASign(pk *rsa.PrivateKey, message []byte) ([]byte, error) {
  hashed := sha512.Sum512(message)
  signature, err := rsa.SignPKCS1v15(rand.Reader, pk, crypto.SHA512, hashed[:])
  return signature, err 
}

