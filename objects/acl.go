package objects

import (
  "os"
  "time"
  "strings"
  "bytes"
  "io/ioutil"
  "encoding/gob"
  "encoding/hex"
  "github.com/google/uuid"
  "github.com/golang/glog"
  "github.com/Lunkov/lib-cipher"
)

type StorageFiles struct {
  path         string
  masterKey    cipher.IACipher
}

type ACLFile struct {
  ownerKeyId     []byte
  userKeyId      []byte
  createAt       time.Time
  key            []byte
}

type ACLStorage struct {
  acl            ACLFile
  sign           []byte
}

func NewStorageFiles() *StorageFiles {
  return &StorageFiles{}
}

func (f *StorageFiles) Open(path string, masterKey cipher.IACipher) {
  f.path = path
  f.masterKey = masterKey
}

func (f *StorageFiles) GetFileName(fileId string, mkdir bool) string {
  p := f.Chunks(fileId, 2)
  if mkdir {
    p1 := p[:len(p)-1]
    pathDir := strings.Join(p1, "/")
    os.MkdirAll(f.path + "/" + pathDir, os.ModePerm)
  }
  return f.path + "/" + strings.Join(p, "/")
}

func (f *StorageFiles) GetACLFileName(fileId string, userKey cipher.IACipher, mkdir bool) string {
  p := f.Chunks(fileId, 2)
  if mkdir {
    p1 := p[:len(p)-1]
    pathDir := strings.Join(p1, "/")
    os.MkdirAll(f.path + "/" + pathDir, os.ModePerm)
  }
  return f.path + "/" + strings.Join(p, "/") + ".acl." + hex.EncodeToString(userKey.GetID())
}

func (f *StorageFiles) CreateFile(pubKeyId string, data []byte) (string, bool) {
  // Generate FileId
  uid, _ := uuid.NewUUID()
  fileId := uid.String()
  
  // Generate New Symmetric Key
  sk := cipher.NewSCipher()
  key, ok := sk.AESCreateKey()
  if !ok {
    return fileId, false
  }
  
  // Encrypt Symmetric Key By MasterKey
  acl := ACLStorage{
           acl: ACLFile{
                      ownerKeyId: f.masterKey.GetID(),
                      userKeyId: f.masterKey.GetID(),
                      createAt: time.Now(),
                  },
        }
  acl.acl.key, ok = f.masterKey.EncryptWithPublicKey(key)
  if !ok {
    return fileId, false
  }
  acl.sign, ok = f.masterKey.Sign(key)
  if !ok {
    return fileId, false
  }
  
  // Get FileName
  filename := f.GetFileName(fileId, true)
  
  // Save New Symmetric Key
  var buff bytes.Buffer
  encoder := gob.NewEncoder(&buff)
  encoder.Encode(acl)
  aclfilename := f.GetACLFileName(fileId, f.masterKey, false)
  err := ioutil.WriteFile(aclfilename, buff.Bytes(), 0640) // just pass the file name
  if err != nil {
    glog.Errorf("ERR: SaveKey Write(%s): %v", aclfilename, err)
    return "", false
  }
  
  // Encrypt and Save File
  cf := cipher.NewCFile()
  return filename, cf.SaveFile(filename, key, data)
}

func (f *StorageFiles) AddShareFile(fileId string, userKey cipher.IACipher) (bool) {
  acluserfilename := f.GetACLFileName(fileId, userKey, false)
  
  // Read ACL master file
  var userACL ACLStorage
  
  skey, ok := f.LoadACL(fileId)
  if !ok {
    return false
  }
  
  // Encrypt Symmetric Key by Asymmetric User Key
  userACL.acl.ownerKeyId = f.masterKey.GetID()
  userACL.acl.userKeyId = userKey.GetID()
  userACL.acl.createAt = time.Now()
  userACL.acl.key, ok = userKey.EncryptWithPublicKey(skey)
  if !ok {
    return false
  }
  userACL.sign, ok = f.masterKey.Sign(skey)
  if !ok {
    return false
  }

  // Save User ACL
  var buff bytes.Buffer
  encoder := gob.NewEncoder(&buff)
  encoder.Encode(userACL)
  aclfilename := f.GetACLFileName(fileId, userKey, false)
  err := ioutil.WriteFile(acluserfilename, buff.Bytes(), 0640) // just pass the file name
  if err != nil {
    glog.Errorf("ERR: SaveACL Write(%s): %v", aclfilename, err)
    return false
  }
  return true
}

func (f *StorageFiles) RemoveShareFile(fileId string, userKey cipher.IACipher) (bool) {
  acluserfilename := f.GetACLFileName(fileId, userKey, false)

  // Remove File
  err := os.Remove(acluserfilename)
  if err != nil {
    glog.Errorf("ERR: RemoveShareFile(%s): %v", acluserfilename, err)
    return false
  }
  return true
}

func (f *StorageFiles) LoadACL(fileId string) ([]byte, bool) {
  aclmasterfilename := f.GetACLFileName(fileId, f.masterKey, false)

  // Read ACL master file
  var masterACL ACLStorage
  data, err := ioutil.ReadFile(aclmasterfilename) 
  if err != nil {
    glog.Errorf("ERR: LoadMasterACL (%s) err='%v'", aclmasterfilename, err)
    return nil, false
  }
  
  // Decrypt Symmetric Key by Asymmetric Master Key
  buf := bytes.NewBuffer(data)
  decoder := gob.NewDecoder(buf)
  err = decoder.Decode(masterACL)
  if err != nil {
    glog.Errorf("ERR: gob.Decoder('%s'): GOB: %v", aclmasterfilename, err)
    return nil, false
  }
  key, ok := f.masterKey.DecryptWithPrivateKey(masterACL.acl.key)
  if !ok {
    return nil, false
  }
  return key, true
}

func (f *StorageFiles) SaveFile(fileId string, data []byte) (bool) {
  filename := f.GetFileName(fileId, true)
  // Load User ACL
  skey, ok := f.LoadACL(fileId)
  if !ok {
    return false
  }
  
  // Save File
  cf := cipher.NewCFile()
  return cf.SaveFile(filename, skey, data)
}

func (f *StorageFiles) LoadFile(fileId string)  ([]byte, bool) {
  filename := f.GetFileName(fileId, true)
  // Load User ACL
  skey, ok := f.LoadACL(fileId)
  if !ok {
    return nil, false
  }
  
  // Load File
  cf := cipher.NewCFile()
  return cf.LoadFile(filename, skey)
}

func (f *StorageFiles) Chunks(s string, chunkSize int) []string {
  if len(s) == 0 {
    return nil
  }
  if chunkSize >= len(s) {
    return []string{s}
  }
  var chunks []string = make([]string, 0, (len(s)-1)/chunkSize+1)
  currentLen := 0
  currentStart := 0
  for i := range s {
    if currentLen == chunkSize {
      chunks = append(chunks, s[currentStart:i])
      currentLen = 0
      currentStart = i
    }
    currentLen++
  }
  chunks = append(chunks, s[currentStart:])
  return chunks
}
