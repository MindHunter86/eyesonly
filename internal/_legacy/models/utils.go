package models

import (
	"crypto/sha256"
	"fmt"
	"hash"
	"sync"
)

var sha256Pool = sync.Pool{New: func() any {
	return sha256.New()
}}

func acquireSha256Hash() hash.Hash { return sha256Pool.Get().(hash.Hash) }
func releaseSha256Hash(v hash.Hash) {
	v.Reset()
	sha256Pool.Put(v)
}

// SHA256 hasher with minimal allocations
func hashSha256(dst, src []byte) (_ []byte, e error) {
	s256 := acquireSha256Hash()
	defer releaseSha256Hash(s256)

	if _, e = s256.Write(src); e != nil {
		return
	}

	return s256.Sum(dst), e
}

var accPasswordPool = sync.Pool{New: func() any {
	return &accountPassword{
		s256: makeBytesWithCap(sha256.Size),
	}
}}

func makeBytesWithCap(c int) []byte {
	return make([]byte, 0, c)
}

// todo - performance tests, please...
func acquirePooledType[V any]() (v V) {
	switch d := any(v).(type) {
	case *accountPassword:
		return accPasswordPool.Get().(V)
	default:
		panic(fmt.Sprintf("undefined type in acquirePooledType %s, %T", d, d))
	}
}

func releasePooledType(v any) {
	switch any(v).(type) {
	case *accountPassword:
		v.(*accountPassword).release()
	}
}
