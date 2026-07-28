package crypto

import (
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"errors"
	"hash"
	"sync"
)

type Service struct {
	key  string
	klen int
}

func NewService(key string) (_ *Service, e error) {
	var klen int
	if klen = len(key); klen == 0 {
		return nil, errors.New("could not initialize crypto subsystem without encryption key")
	}

	return &Service{key, klen}, nil
}

func (m *Service) Encrypt(payload []byte) []byte {
	return m.xorEncDecMethod(payload)
}

func (m *Service) Decrypt(payload []byte) []byte {
	return m.xorEncDecMethod(payload)
}

func (m *Service) xorEncDecMethod(p []byte) []byte {
	if m.klen == 0 {
		return p
	}

	for i := range len(p) {
		p[i] = p[i] ^ m.key[i%m.klen]
	}

	return p
}

func SHA256Hex(value string) string {
	sum := sha256.Sum256([]byte(value))
	return hex.EncodeToString(sum[:])
}

func KeyBytes(value string) []byte {
	if decoded, err := base64.StdEncoding.DecodeString(value); err == nil && len(decoded) == 32 {
		return decoded
	}
	sum := sha256.Sum256([]byte(value))
	return sum[:]
}

// !!!
// !!!
// !!! SEE signer.go !!!

var sha256Pool = sync.Pool{New: func() any {
	return sha256.New()
}}

func acquireSha256Hash() hash.Hash { return sha256Pool.Get().(hash.Hash) }
func releaseSha256Hash(v hash.Hash) {
	v.Reset()
	sha256Pool.Put(v)
}

// SHA256 hasher with minimal allocations
func Sha256Hash(dst, src []byte) (_ []byte, e error) {
	s256 := acquireSha256Hash()
	defer releaseSha256Hash(s256)

	if _, e = s256.Write(src); e != nil {
		return
	}

	return s256.Sum(dst), e
}
