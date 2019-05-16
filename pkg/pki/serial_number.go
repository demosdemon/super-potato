package pki

import (
	"fmt"
	"math/big"
	"strings"
)

type SerialNumber struct {
	*big.Int
}

func NewSerialNumber(b []byte) SerialNumber {
	rv := big.Int{}
	rv.SetBytes(b)
	return SerialNumber{&rv}
}

func (n SerialNumber) String() string {
	b := n.Bytes()
	rv := make([]string, len(b))
	for idx, v := range b {
		rv[idx] = fmt.Sprintf("%02x", v)
	}
	return strings.Join(rv, ":")
}

func (n SerialNumber) MarshalText() ([]byte, error) {
	return []byte(n.String()), nil
}

func (n *SerialNumber) UnmarshalText(text []byte) error {
	nibbles := strings.Split(string(text), ":")
	bigint := make([]byte, len(nibbles))
	for idx, nib := range nibbles {
		_, _ = fmt.Sscanf(nib, "%02x", &bigint[idx])
	}
	n.Int = new(big.Int)
	n.SetBytes(bigint)
	return nil
}
