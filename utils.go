package clientecos

import (
  "fmt"
  "bytes"
  //"errors"
  "encoding/gob"
  "crypto/rand"
  "math/big"
  "crypto/sha512"
  "crypto/ecdsa"
  "crypto/elliptic"
)

type ECDSABuf struct {
  X, Y []byte
}

func Serialize(pk *ecdsa.PublicKey) ([]byte, bool) {
  fmt.Printf("\nrawPubKey Serialize %x\n", pk)
  ecdsabuf := ECDSABuf{X: pk.X.Bytes(), Y: pk.Y.Bytes()}
  var buff bytes.Buffer
  encoder := gob.NewEncoder(&buff)
  encoder.Encode(ecdsabuf)
  return buff.Bytes(), true
}

func Deserialize(msg []byte) (*ecdsa.PublicKey, bool) {
  var ecdsabuf ECDSABuf
  curve := elliptic.P521()
  buf := bytes.NewBuffer(msg)
  decoder := gob.NewDecoder(buf)
  err := decoder.Decode(&ecdsabuf)
  if err != nil {
    fmt.Printf("\ndecoder.Decode %v\n", err)
    return nil, false
  }
  x := big.Int{}
  y := big.Int{}
  x.SetBytes(ecdsabuf.X)
  y.SetBytes(ecdsabuf.Y)
  return &ecdsa.PublicKey{Curve: curve, X: &x, Y: &y}, true
}

func PublicKeyToBytes(key *ecdsa.PrivateKey) ([]byte, bool) {
  /*
  // Marshal the curve parameters to a byte slice
  curveBytes := elliptic.Marshal(pubKey.Curve, pubKey.X, pubKey.Y)

  // Serialize the public key coordinates to a byte slice
  keyLen := (pubKey.Curve.Params().BitSize + 7) / 8
  xBytes := pubKey.X.Bytes()
  yBytes := pubKey.Y.Bytes()
  if len(xBytes) < keyLen {
       pad := make([]byte, keyLen-len(xBytes))
       xBytes = append(pad, xBytes...)
  }
  if len(yBytes) < keyLen {
       pad := make([]byte, keyLen-len(yBytes))
       yBytes = append(pad, yBytes...)
  }
  pubKeyBytes := append(xBytes, yBytes...)

  // Combine the curve parameters and public key coordinates into a single byte slice
  bytes := append(curveBytes, pubKeyBytes...)
  */
  //pubKeyBytes := elliptic.Marshal(key.Curve, key.PublicKey.X, key.PublicKey.Y)
  return elliptic.Marshal(key.Curve, key.PublicKey.X, key.PublicKey.Y), true
}

/*
func PublicKeyFromBytes(pubKey []byte) (*ecdsa.PublicKey, bool) {
	curve := elliptic.P256()

  x := big.Int{}
  y := big.Int{}
  keyLen := len(pubKey)
  x.SetBytes(pubKey[:(keyLen / 2)])
  y.SetBytes(pubKey[(keyLen / 2):])

  return &ecdsa.PublicKey{Curve: curve, X: &x, Y: &y}, true
}*/

func PublicKeyFromBytes(pubKeyBytes []byte) (*ecdsa.PublicKey, bool) {
  // Unmarshal the curve parameters from the byte slice
  /*curve := elliptic.P256() // Replace with the appropriate curve for your use case
  curveParams := curve.Params()
  curveLen := (curveParams.BitSize + 7) / 8
  if len(bytes) < curveLen {
    // fmt.Errorf("invalid byte slice length")
    return nil, false
  }
  //curveBytes := bytes[:curveLen]
  //x, y := elliptic.Unmarshal(curve, curveBytes)

  // Set the public key coordinates from the remaining bytes
  keyLen := (curveParams.BitSize + 7) / 8
  if len(bytes) != curveLen+2*keyLen {
    // fmt.Errorf("invalid byte slice length")
    return nil, false
  }
  xBytes := bytes[curveLen : curveLen+keyLen]
  yBytes := bytes[curveLen+keyLen:]
  pubKey := &ecdsa.PublicKey{
      Curve: curve,
      X:     new(big.Int).SetBytes(xBytes),
      Y:     new(big.Int).SetBytes(yBytes),
  }*/
  curve := elliptic.P256()
  deserializedX, deserializedY := elliptic.Unmarshal(curve, pubKeyBytes)
  //deserializedKey := ecdsa.PublicKey{Curve: key.Curve, X: deserializedX, Y: deserializedY}
  return &ecdsa.PublicKey{Curve: curve, X: deserializedX, Y: deserializedY}, true
}

// 
func ECDSA256Verify(pubKey []byte, message []byte, signature []byte) bool {
  /*
  r := big.Int{}
  s := big.Int{}
  sigLen := len(signature)
  r.SetBytes(signature[:(sigLen / 2)])
  s.SetBytes(signature[(sigLen / 2):])
  */
  defer func() {
		if r := recover(); r != nil {
			fmt.Printf("\nECDSA signature verification error: %v\n", r)
		}
	}()
  rawPubKey, ok := Deserialize(pubKey) //PublicKeyFromBytes(pubKey)
  if !ok {
    return false
  }
  hashed := sha512.Sum512(message)
  //return ecdsa.Verify(rawPubKey, msg, &r, &s)
  fmt.Printf("\nHASH VERIFY %x\n", hashed)
  fmt.Printf("\nrawPubKey VERIFY %x\n", rawPubKey)
  fmt.Printf("\nsignature VERIFY %x\n", signature)
  return ecdsa.VerifyASN1(rawPubKey, hashed[:], signature)
}

func ECDSA256Sign(pk *ecdsa.PrivateKey, message []byte) ([]byte, bool) {
  hashed := sha512.Sum512(message)
  fmt.Printf("\nHASH SIGN %x\n", hashed)
  signature, err := ecdsa.SignASN1(rand.Reader, pk, hashed[:])
  return signature, err == nil 
}

