package utils

import (
  "bytes"
  "github.com/Lunkov/go-btcec"
  "encoding/gob"
  "crypto/rand"
  "math/big"
  "crypto/sha512"
  "crypto/ecdsa"
)

const version = byte(0x00)

type ECDSABuf struct {
  // btcec.KoblitzCurve
	BitSize int      // the size of the underlying field
	Name    string   // the canonical name of the curve
  // 
  X, Y []byte
}

func ECDSAPublicKeySerialize(pk *ecdsa.PublicKey) ([]byte, error) {
  ecdsabuf := ECDSABuf{
                       BitSize: pk.Curve.Params().BitSize,
                       Name: "S256",
                       X: pk.X.Bytes(),
                       Y: pk.Y.Bytes(),
                       }
  var buff bytes.Buffer
  encoder := gob.NewEncoder(&buff)
  encoder.Encode(ecdsabuf)
  return buff.Bytes(), nil
}

func ECDSAPublicKeyDeserialize(msg []byte) (*ecdsa.PublicKey, error) {
  var ecdsabuf ECDSABuf
  buf := bytes.NewBuffer(msg)
  decoder := gob.NewDecoder(buf)
  err := decoder.Decode(&ecdsabuf)
  if err != nil {
    return nil, err
  }
  x := big.Int{}
  y := big.Int{}
  x.SetBytes(ecdsabuf.X)
  y.SetBytes(ecdsabuf.Y)
  
  // btcec.KoblitzCurve
  return &ecdsa.PublicKey{Curve: btcec.S256(), X: &x, Y: &y}, nil
}

// 
func ECDSA256VerifyHash512(pubKey []byte, hash []byte, signature []byte) (bool, error) {
  defer func() {
    if r := recover(); r != nil {
    }
  }()
  rawPubKey, err := ECDSAPublicKeyDeserialize(pubKey) //PublicKeyFromBytes(pubKey)
  if err != nil {
    return false, err
  }
  return ecdsa.VerifyASN1(rawPubKey, hash, signature), nil
}

func ECDSA256Sign(pk *ecdsa.PrivateKey, message []byte) ([]byte, error) {
  hashed := sha512.Sum512(message)
  signature, err := ecdsa.SignASN1(rand.Reader, pk, hashed[:])
  return signature, err
}

func ECDSA256VerifySender(address string, pubKey []byte, hash []byte, signature []byte) (bool, error) {
  defer func() {
    if r := recover(); r != nil {
    }
  }()
  rawPubKey, err := ECDSAPublicKeyDeserialize(pubKey)
  if err != nil {
    return false, err
  }
  return ecdsa.VerifyASN1(rawPubKey, hash, signature), nil
}

func ECDSA256SignHash512(pk *ecdsa.PrivateKey, hash []byte) ([]byte, error) {
  signature, err := ecdsa.SignASN1(rand.Reader, pk, hash)
  return signature, err
}
