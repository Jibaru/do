package utils

import (
	"crypto/rand"
	"fmt"
	"io"
)

type UuidFactory interface {
	New() string
}

type randomUuidFactory struct{}

type fixedUuidFactory struct {
	uuid string
}

func NewRandomUuidFactory() UuidFactory {
	return &randomUuidFactory{}
}

func NewFixedUuidFactory(uuid string) UuidFactory {
	return &fixedUuidFactory{uuid: uuid}
}

func (f *randomUuidFactory) New() string {
	uuid := make([]byte, 16)
	_, _ = io.ReadFull(rand.Reader, uuid)

	// Version 4 (pseudo-random) UUID
	uuid[6] = (uuid[6] & 0x0f) | 0x40 // Set version to 4
	uuid[8] = (uuid[8] & 0x3f) | 0x80 // Set variant to RFC 4122

	return fmt.Sprintf("%08x-%04x-%04x-%04x-%12x", uuid[0:4], uuid[4:6], uuid[6:8], uuid[8:10], uuid[10:])
}

func (f *fixedUuidFactory) New() string {
	return f.uuid
}
