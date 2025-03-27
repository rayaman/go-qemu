package types

import (
	"math/big"
	"strings"

	"github.com/dustin/go-humanize"
)

type Size string

var (
	KB  = big.NewInt(1000)
	KIB = big.NewInt(1024)
	MB  = big.NewInt(1000000)
	MIB = big.NewInt(1048576)
	GB  = big.NewInt(1000000000)
	GIB = big.NewInt(1073741824)
	TB  = big.NewInt(1000000000000)
	TIB = big.NewInt(1099511627776)
)

// Takes a base and a multiplier and returns a human readable size
func GetSize(base *big.Int, mul int64) Size {
	b := humanize.BigBytes(big.NewInt(0).Mul(base, big.NewInt(mul)))
	return Size(strings.ReplaceAll(strings.ReplaceAll(strings.ReplaceAll(b, ".0", ""), " ", ""), "B", ""))
}
