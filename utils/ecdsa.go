package utils

import (
  "bytes"
  //"github.com/btcsuite/btcd/btcec"
  "github.com/Lunkov/go-btcec"
  "encoding/gob"
  "crypto/rand"
  "math/big"
  "crypto/sha512"
  "crypto/ecdsa"
  "github.com/golang/glog"
)

const version = byte(0x00)

type ECDSABuf struct {
  // btcec.KoblitzCurve
	BitSize int      // the size of the underlying field
	Name    string   // the canonical name of the curve
  // 
  X, Y []byte
}

func ECDSAPublicKeySerialize(pk *ecdsa.PublicKey) ([]byte, bool) {
  ecdsabuf := ECDSABuf{
                       BitSize: pk.Curve.Params().BitSize,
                       Name: "S256",
                       X: pk.X.Bytes(),
                       Y: pk.Y.Bytes(),
                       }
  var buff bytes.Buffer
  encoder := gob.NewEncoder(&buff)
  encoder.Encode(ecdsabuf)
  return buff.Bytes(), true
}

func ECDSAPublicKeyDeserialize(msg []byte) (*ecdsa.PublicKey, bool) {
  var ecdsabuf ECDSABuf
  buf := bytes.NewBuffer(msg)
  decoder := gob.NewDecoder(buf)
  err := decoder.Decode(&ecdsabuf)
  if err != nil {
    glog.Errorf("ERR: gob.Decode %v", err)
    return nil, false
  }
  x := big.Int{}
  y := big.Int{}
  x.SetBytes(ecdsabuf.X)
  y.SetBytes(ecdsabuf.Y)
  
  // btcec.KoblitzCurve
  return &ecdsa.PublicKey{Curve: btcec.S256(), X: &x, Y: &y}, true
}

// 
func ECDSA256VerifyHash512(pubKey []byte, hash []byte, signature []byte) bool {
  defer func() {
    if r := recover(); r != nil {
      glog.Errorf("ERR: ECDSA signature verification error: %v", r)
    }
  }()
  rawPubKey, ok := ECDSAPublicKeyDeserialize(pubKey) //PublicKeyFromBytes(pubKey)
  if !ok {
    return false
  }
  return ecdsa.VerifyASN1(rawPubKey, hash, signature)
}

func ECDSA256Sign(pk *ecdsa.PrivateKey, message []byte) ([]byte, bool) {
  hashed := sha512.Sum512(message)
  signature, err := ecdsa.SignASN1(rand.Reader, pk, hashed[:])
  return signature, err == nil 
}

func ECDSA256VerifySender(address string, pubKey []byte, hash []byte, signature []byte) bool {
  defer func() {
    if r := recover(); r != nil {
      glog.Errorf("ERR: ECDSA signature verification error: %v", r)
    }
  }()
  rawPubKey, ok := ECDSAPublicKeyDeserialize(pubKey)
  if !ok {
    glog.Errorf("ERR: ECDSAPublicKeyDeserialize %v", rawPubKey)
    return false
  }
  return ecdsa.VerifyASN1(rawPubKey, hash, signature)
}

func ECDSA256SignHash512(pk *ecdsa.PrivateKey, hash []byte) ([]byte, bool) {
  signature, err := ecdsa.SignASN1(rand.Reader, pk, hash)
  return signature, err == nil 
}
