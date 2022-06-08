package pflag

import (
	"strconv"
)

// -- uint64 Value
type uint64Value uint64

func newUint64Value(val uint64, p *uint64) *uint64Value {
	*p = val
	return (*uint64Value)(p)
}

func (i *uint64Value) Set(s string) error {
	v, err := strconv.ParseUint(s, 0, 64)
	*i = uint64Value(v)
	return err
}

func (i *uint64Value) Type() string {
	return "uint64"
}

func (i *uint64Value) String() string { return strconv.FormatUint(uint64(*i), 10) }

func uint64Conv(sval string) (interface{}, error) {
	v, err := strconv.ParseUint(sval, 0, 64)
	if err != nil {
		return 0, err
	}
	return uint64(v), nil
}
