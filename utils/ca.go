package utils

import (
  "time"
  "errors"
  "bytes"
  "strings"
  "crypto"
  "crypto/rand"
  "crypto/rsa"
  "crypto/ecdsa"
  "crypto/sha512"
  "crypto/x509"
  "crypto/x509/pkix"
  "encoding/asn1"
  "encoding/pem"
  "gopkg.in/yaml.v3"
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
  DisplayName     string        `yaml:"DisplayName,omitempty"`
  CommonName      string        `yaml:"CommonName,omitempty"`
  Organization    string        `yaml:"Organization,omitempty"`
  Country         string        `yaml:"Country,omitempty"`
  Locality        string        `yaml:"Locality,omitempty"`
  StreetAddress   string        `yaml:"StreetAddress,omitempty"`
  PostalCode      string        `yaml:"PostalCode,omitempty"`
  OrgUnit         string        `yaml:"OrgUnit,omitempty"`
  EMail           string        `yaml:"EMail,omitempty"`

  Bits            int           `yaml:"Bits,omitempty"`
  
  PrivateKey     *rsa.PrivateKey
  Cert           *x509.Certificate
}

var oidEmailAddress = asn1.ObjectIdentifier{1, 2, 840, 113549, 1, 9, 1}
var oidDisplayName = asn1.ObjectIdentifier{2, 16, 76, 1, 3, 3}

func NewCertInfo() *CertInfo {
  return &CertInfo{}
}

func (cai *CertInfo) NewInfo(isCA bool) *x509.Certificate {
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
    EmailAddresses:      []string{cai.EMail},
    NotBefore:             time.Now(),
    NotAfter:              time.Now().AddDate(10, 0, 0),
    IsCA:                  isCA,
    ExtKeyUsage:           []x509.ExtKeyUsage{x509.ExtKeyUsageClientAuth, x509.ExtKeyUsageServerAuth},
    KeyUsage:              x509.KeyUsageKeyEncipherment | x509.KeyUsageDigitalSignature | x509.KeyUsageCertSign,
    BasicConstraintsValid: true,
  }
}

func (cai *CertInfo) СreateNewCA(fileNameCert string, fileNamePriv string, password string) (error) {
  cai.Cert = cai.NewInfo(true)

  caPrivKey, errc := rsa.GenerateKey(rand.Reader, cai.Bits)
  if errc != nil {
    return errors.New("Private Key does not exists")
  }
  cai.PrivateKey = caPrivKey

  buf, err := cai.SerializeCert()
  if err != nil {
    return err
  }

  err = os.WriteFile(fileNameCert, buf, 0640) 
  if err != nil {
    return err
  }

  block := &pem.Block{
      Type:  "RSA PRIVATE KEY",
      Bytes: x509.MarshalPKCS1PrivateKey(caPrivKey),
    }

  if password != "" {
    var err error
    block, err = x509.EncryptPEMBlock(rand.Reader, block.Type, block.Bytes, []byte(password), x509.PEMCipherAES256)
    if err != nil {
      return err
    }
  }

  return os.WriteFile(fileNamePriv, pem.EncodeToMemory(block), 0640) 
}

func (cai *CertInfo) LoadConfig(filename string) error {
  out, err := os.ReadFile(filename) 
  if err != nil {
    return err
  }
  return yaml.Unmarshal(out, &cai)
}

func (cai *CertInfo) Load(fileNameCert string, fileNamePriv string, password string) error {
  content, err := os.ReadFile(fileNameCert)
  if err != nil {
    return err
  }
  
  err = cai.DeserializeCert(content)
  if err != nil {
    return err
  }
  
  content, err = os.ReadFile(fileNamePriv)
  if err != nil {
    return err
  }
  block, _ := pem.Decode(content)
  if block == nil {
    return errors.New("Wrong format")
  }
  enc := x509.IsEncryptedPEMBlock(block)
  b := block.Bytes
  if enc {
    b, err = x509.DecryptPEMBlock(block, []byte(password))
    if err != nil {
      return err
    }
  }
  key, err := x509.ParsePKCS1PrivateKey(b)
  if err != nil {
    return err
  }
  cai.PrivateKey = key
  return nil
}

func (cai *CertInfo) SerializeCert() ([]byte, error) {
  caBytes, err := x509.CreateCertificate(rand.Reader, cai.Cert, cai.Cert, &cai.PrivateKey.PublicKey, cai.PrivateKey)
  if err != nil {
    return nil, err
  }

  caPEM := new(bytes.Buffer)
  pem.Encode(caPEM, &pem.Block{
    Type:  "CERTIFICATE",
    Bytes: caBytes,
  })

  return caPEM.Bytes(), nil
}

func (cai *CertInfo) DeserializeCert(buf []byte) (error) {
  block, _ := pem.Decode(buf)
  if block == nil {
    return errors.New("Wrong format")
  }
  cert, err2 := x509.ParseCertificate(block.Bytes)
  if err2 != nil {
    return err2
  }
  cai.Cert = cert
  return nil
}

func (cai *CertInfo) Sign(message []byte) ([]byte, error) {
  if cai.PrivateKey == nil {
    return nil, errors.New("Private Key does not exists")
  }
  hash := sha512.Sum512(message)
  signature, err := rsa.SignPKCS1v15(rand.Reader, cai.PrivateKey, crypto.SHA512, hash[:])
  if err != nil {
    return nil, err
  }
  return signature, nil 
}

func (cai *CertInfo) Verify(message []byte, signature []byte) (error) {
  if cai.Cert == nil {
    return errors.New("Certificate does not exists")
  }
  hash := sha512.Sum512(message)
  var err error
  switch pub := cai.Cert.PublicKey.(type) {
	case *rsa.PublicKey:
    err = rsa.VerifyPKCS1v15(pub, crypto.SHA512, hash[:], signature)

	case *ecdsa.PublicKey:
		err = cai.Cert.CheckSignature(x509.ECDSAWithSHA512, nil, signature)

	default:
		return errors.New("Wrong type of key")
	}
  
  if err != nil {
    return err
  }

  return nil
}

func (cai *CertInfo) SignFile(filename string) error {
  content, err := os.ReadFile(filename)
  if err != nil {
    return err
  }
  sign, errs := cai.Sign(content)
  if errs != nil {
    return errs
  }
  filens := filename + ".sign"
  return os.WriteFile(filens, sign, 0640) 
}

func (cai *CertInfo) VerifyFile(filename string) error {
  content, err := os.ReadFile(filename)
  if err != nil {
    return err
  }
  filesign := filename + ".sign"
  sign, errs := os.ReadFile(filesign)
  if errs != nil {
    return err
  }
  return cai.Verify(content, sign)
}


func (cai *CertInfo) CreateSubCert(subCert *CertInfo) ([]byte, error) {
  certPrivKey, errg := rsa.GenerateKey(rand.Reader, subCert.Bits)
  if errg != nil {
    return nil, errg
  }
  
  cn := strings.Split(subCert.EMail, "@")
  if len(cn) != 2 {
    return nil, errors.New("Wrong email: " + subCert.EMail)
  }

  subCert.Cert = &x509.Certificate{
    SerialNumber: big.NewInt(1658),
    Subject: pkix.Name{
      CommonName:         cn[0],
      Country:            []string{subCert.Country},
      Province:           []string{""},
      Locality:           []string{subCert.Locality},
      Organization:       []string{cai.Organization},
      OrganizationalUnit: []string{subCert.OrgUnit},
      ExtraNames: []pkix.AttributeTypeAndValue{
            {
                Type:  oidEmailAddress, 
                Value: asn1.RawValue{
                    Tag:   asn1.TagIA5String, 
                    Bytes: []byte(subCert.EMail),
                },
            },
            {
                Type:  oidDisplayName, 
                Value: asn1.RawValue{
                    Tag:   asn1.TagIA5String, 
                    Bytes: []byte(subCert.DisplayName),
                },
            },
        },
    },
    IsCA:         false,
    EmailAddresses: []string{subCert.EMail},
    //IPAddresses:  []net.IP{net.IPv4(127, 0, 0, 1), net.IPv6loopback},
    NotBefore:    time.Now(),
    NotAfter:     time.Now().AddDate(10, 0, 0),
    SubjectKeyId: []byte{1, 2, 3, 4, 6},
    ExtKeyUsage:  []x509.ExtKeyUsage{x509.ExtKeyUsageClientAuth, x509.ExtKeyUsageServerAuth},
    KeyUsage:     x509.KeyUsageDigitalSignature,
  }
  subCert.PrivateKey = certPrivKey
  
  certBytes, err := x509.CreateCertificate(rand.Reader, subCert.Cert, cai.Cert, &subCert.PrivateKey.PublicKey, cai.PrivateKey)
  if err != nil {
    return nil, err
  }

  caPEM := new(bytes.Buffer)
  pem.Encode(caPEM, &pem.Block{
    Type:  "CERTIFICATE",
    Bytes: certBytes,
  })

  return caPEM.Bytes(), nil
}
