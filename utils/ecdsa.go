package utils

import (
  "bytes"
  "encoding/gob"
  "crypto/rand"
  "math/big"
  "crypto/sha512"
  "crypto/ecdsa"
  "crypto/elliptic"
  "github.com/golang/glog"
)

const version = byte(0x00)

type ECDSABuf struct {
  X, Y []byte
}

func ECDSASerialize(pk *ecdsa.PublicKey) ([]byte, bool) {
  ecdsabuf := ECDSABuf{X: pk.X.Bytes(), Y: pk.Y.Bytes()}
  var buff bytes.Buffer
  encoder := gob.NewEncoder(&buff)
  encoder.Encode(ecdsabuf)
  return buff.Bytes(), true
}

func ECDSADeserialize(msg []byte) (*ecdsa.PublicKey, bool) {
  var ecdsabuf ECDSABuf
  curve := elliptic.P521()
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
  return &ecdsa.PublicKey{Curve: curve, X: &x, Y: &y}, true
}

func ECDSA256PublicKeyFromBytes(pubKeyBytes []byte) (*ecdsa.PublicKey, bool) {
  curve := elliptic.P256()
  deserializedX, deserializedY := elliptic.Unmarshal(curve, pubKeyBytes)
  return &ecdsa.PublicKey{Curve: curve, X: deserializedX, Y: deserializedY}, true
}

// 
func ECDSA256Verify(pubKey []byte, message []byte, signature []byte) bool {
  defer func() {
    if r := recover(); r != nil {
      glog.Errorf("ERR: ECDSA signature verification error: %v", r)
    }
  }()
  rawPubKey, ok := ECDSADeserialize(pubKey) //PublicKeyFromBytes(pubKey)
  if !ok {
    return false
  }
  hashed := sha512.Sum512(message)
  return ecdsa.VerifyASN1(rawPubKey, hashed[:], signature)
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
  rawPubKey, ok := ECDSADeserialize(pubKey)
  if !ok {
    return false
  }
  if address != PubkeyToAddress(rawPubKey).Hex() {
    return false
  }
  return ecdsa.VerifyASN1(rawPubKey, hash, signature)
}

func ECDSA256SignHash512(pk *ecdsa.PrivateKey, hash []byte) ([]byte, bool) {
  signature, err := ecdsa.SignASN1(rand.Reader, pk, hash)
  return signature, err == nil 
}

/*
func ECDSAEncrypt(publicKey *ecdsa.PublicKey, message []byte) ([]byte, bool) {
    k, err := rand.Int(rand.Reader, publicKey.Params().N) // Generate a random number k
    if err != nil {
      glog.Errorf("ERR: ECDSA Encrypt: %v", err)
      return nil, false
    }
    // x, y := publicKey.Curve.ScalarBaseMult(k.Bytes()) // Calculate the point (x,y) on the elliptic curve
    x, y := publicKey.Curve.ScalarBaseMult(k.Bytes()) // Calculate the point (x,y) on the elliptic curve
    s := new(big.Int).SetBytes(message) // Convert the message to a big integer
    s.Mul(s, x) // Multiply the message by x
    s.Mod(s, publicKey.Params().N) // Take the modulus of the result
    yBytes := elliptic.Marshal(publicKey.Curve, x, y) // Convert y to a byte array
    //yBytes := elliptic.Marshal(publicKey.Curve, x, y) // Convert y to a byte array
    return append(yBytes, s.Bytes()...), true // Return the encrypted message
    // return s.Bytes(), true // Return the encrypted message
}

func ECDSADecrypt(privateKey *ecdsa.PrivateKey, encryptedMessage []byte) ([]byte, bool) {
    yBytes := encryptedMessage[:len(encryptedMessage)/2] // Extract the y component from the encrypted message
    sBytes := encryptedMessage[len(encryptedMessage)/2:] // Extract the s component from the encrypted message
    y, _ := elliptic.Unmarshal(privateKey.Curve, yBytes) // Convert y back to a point on the elliptic curve
    // x, _ := privateKey.Curve.ScalarMult(y, privateKey.D.Bytes()) // Calculate the point (x,y) on the elliptic curve
    x, _ := privateKey.Curve.ScalarMult(y, y, privateKey.D.Bytes()) // Calculate the point (x,y) on the elliptic curve
    s := new(big.Int).SetBytes(sBytes) // Convert s to a big integer
    // s := new(big.Int).SetBytes(encryptedMessage) // Convert s to a big integer
    s.ModInverse(s, privateKey.Params().N) // Calculate the modular inverse of s
    s.Mul(s, x) // Multiply the result by x
    s.Mod(s, privateKey.Params().N) // Take the modulus of the result
    return s.Bytes(), true // Return the decrypted message
}
*/
