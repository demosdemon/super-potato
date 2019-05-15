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
