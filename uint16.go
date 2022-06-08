package pflag

import (
	"strconv"
)

// -- uint16 value
type uint16Value uint16

func newUint16Value(val uint16, p *uint16) *uint16Value {
	*p = val
	return (*uint16Value)(p)
}

func (i *uint16Value) Set(s string) error {
	v, err := strconv.ParseUint(s, 0, 16)
	*i = uint16Value(v)
	return err
}

func (i *uint16Value) Type() string {
	return "uint16"
}

func (i *uint16Value) String() string { return strconv.FormatUint(uint64(*i), 10) }

func uint16Conv(sval string) (interface{}, error) {
	v, err := strconv.ParseUint(sval, 0, 16)
	if err != nil {
		return 0, err
	}
	return uint16(v), nil
}
