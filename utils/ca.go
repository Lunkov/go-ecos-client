package utils

import (
  "time"
  "bytes"
  "crypto"
  "crypto/rand"
  "crypto/rsa"
  "crypto/ecdsa"
  "crypto/sha512"
  "crypto/x509"
  "crypto/x509/pkix"
  
  "encoding/pem"
  
  "os"
  "math/big"
)

/*
 * https://opensource.com/article/22/9/dynamically-update-tls-certificates-golang-server-no-downtime
 * 
openssl genrsa -out localhost.key 2048
openssl req -new -key localhost.key -out localhost.csr
openssl x509 -req -in localhost.csr -CA rootCA.crt -CAkey rootCA.key -CAcreateserial -out localhost.crt -days 500 -sha256
* 
* https://shaneutt.com/blog/golang-ca-and-signed-cert-go/
* https://gist.github.com/shaneutt/5e1995295cff6721c89a71d13a71c251
*/

type CertInfo struct {
  CommonName      string        `yaml:"CommonName,omitempty"`
  Organization    string        `yaml:"Organization,omitempty"`
  Country         string        `yaml:"Country,omitempty"`
  Locality        string        `yaml:"Locality,omitempty"`
  StreetAddress   string        `yaml:"StreetAddress,omitempty"`
  PostalCode      string        `yaml:"PostalCode,omitempty"`

  Bits            int           `yaml:"Bits,omitempty"`
  
  PrivateKey     *rsa.PrivateKey
  Cert           *x509.Certificate
}

func NewCertInfo() *CertInfo {
  return &CertInfo{}
}

func (cai *CertInfo) NewInfo() *x509.Certificate {
  return &x509.Certificate{
    SerialNumber: big.NewInt(2023),
    Subject: pkix.Name{
      Organization:  []string{cai.Organization},
      Country:       []string{cai.Country},
      Province:      []string{""},
      Locality:      []string{cai.Locality},
      StreetAddress: []string{cai.StreetAddress},
      PostalCode:    []string{cai.PostalCode},
    },
    NotBefore:             time.Now(),
    NotAfter:              time.Now().AddDate(10, 0, 0),
    IsCA:                  true,
    ExtKeyUsage:           []x509.ExtKeyUsage{x509.ExtKeyUsageClientAuth, x509.ExtKeyUsageServerAuth},
    KeyUsage:              x509.KeyUsageKeyEncipherment | x509.KeyUsageDigitalSignature | x509.KeyUsageCertSign,
    BasicConstraintsValid: true,
  }
}

func (cai *CertInfo) СreateNewCA(fileNameCert string, fileNamePriv string, password string) (bool) {
  cai.Cert = cai.NewInfo()

  caPrivKey, errc := rsa.GenerateKey(rand.Reader, cai.Bits)
  if errc != nil {
    return false
  }
  cai.PrivateKey = caPrivKey

  buf, ok := cai.SerializeCert()
  if !ok {
    return false
  }

  err := os.WriteFile(fileNameCert, buf, 0640) 
  if err != nil {
    return false
  }

  block := &pem.Block{
      Type:  "RSA PRIVATE KEY",
      Bytes: x509.MarshalPKCS1PrivateKey(caPrivKey),
    }

  if password != "" {
    var err error
    block, err = x509.EncryptPEMBlock(rand.Reader, block.Type, block.Bytes, []byte(password), x509.PEMCipherAES256)
    if err != nil {
      return false
    }
  }

  err = os.WriteFile(fileNamePriv, pem.EncodeToMemory(block), 0640) 
  if err != nil {
    return false
  }
  return true
}

func (cai *CertInfo) Load(fileNameCert string, fileNamePriv string, password string) bool {
  content, err := os.ReadFile(fileNameCert)
  if err != nil {
    return false
  }
  
  if !cai.DeserializeCert(content) {
    return false
  }
  
  content, err = os.ReadFile(fileNamePriv)
  if err != nil {
    return false
  }
  block, _ := pem.Decode(content)
  if block == nil {
    return false
  }
  enc := x509.IsEncryptedPEMBlock(block)
  b := block.Bytes
  if enc {
    b, err = x509.DecryptPEMBlock(block, []byte(password))
    if err != nil {
      return false
    }
  }
  key, err := x509.ParsePKCS1PrivateKey(b)
  if err != nil {
    return false
  }
  cai.PrivateKey = key
  return true
}

func (cai *CertInfo) SerializeCert() ([]byte, bool) {
  caBytes, err := x509.CreateCertificate(rand.Reader, cai.Cert, cai.Cert, &cai.PrivateKey.PublicKey, cai.PrivateKey)
  if err != nil {
    return nil, false
  }

  caPEM := new(bytes.Buffer)
  pem.Encode(caPEM, &pem.Block{
    Type:  "CERTIFICATE",
    Bytes: caBytes,
  })

  return caPEM.Bytes(), true
}

func (cai *CertInfo) DeserializeCert(buf []byte) (bool) {
  block, _ := pem.Decode(buf)
  if block == nil {
    return false
  }
  cert, err2 := x509.ParseCertificate(block.Bytes)
  if err2 != nil {
    return false
  }
  cai.Cert = cert
  return true
}

func (cai *CertInfo) Sign(message []byte) ([]byte, bool) {
  if cai.PrivateKey == nil {
    return nil, false
  }
  hash := sha512.Sum512(message)
  signature, err := rsa.SignPKCS1v15(rand.Reader, cai.PrivateKey, crypto.SHA512, hash[:])
  if err != nil {
    return nil, false
  }
  return signature, true 
}

func (cai *CertInfo) Verify(message []byte, signature []byte) (bool) {
  if cai.Cert == nil {
    return false
  }
  hash := sha512.Sum512(message)
  var err error
  switch pub := cai.Cert.PublicKey.(type) {
	case *rsa.PublicKey:
    err = rsa.VerifyPKCS1v15(pub, crypto.SHA512, hash[:], signature)

	case *ecdsa.PublicKey:
		err = cai.Cert.CheckSignature(x509.ECDSAWithSHA512, nil, signature)

	default:
		return false
	}
  
  if err != nil {
    return false
  }

  return true
}

func (cai *CertInfo) SignFile(filename string) bool {
  content, err := os.ReadFile(filename)
  if err != nil {
    return false
  }
  sign, ok := cai.Sign(content)
  if !ok {
    return false
  }
  filens := filename + ".sign"
  err = os.WriteFile(filens, sign, 0640) 
  if err != nil {
    return false
  }
  return true
}

func (cai *CertInfo) VerifyFile(filename string) bool {
  content, err := os.ReadFile(filename)
  if err != nil {
    return false
  }
  filesign := filename + ".sign"
  sign, errs := os.ReadFile(filesign)
  if errs != nil {
    return false
  }
  return cai.Verify(content, sign)
}
