package pflag

import (
	"strconv"
)

// -- uint Value
type uintValue uint

func newUintValue(val uint, p *uint) *uintValue {
	*p = val
	return (*uintValue)(p)
}

func (i *uintValue) Set(s string) error {
	v, err := strconv.ParseUint(s, 0, 64)
	*i = uintValue(v)
	return err
}

func (i *uintValue) Type() string {
	return "uint"
}

func (i *uintValue) String() string { return strconv.FormatUint(uint64(*i), 10) }

func uintConv(sval string) (interface{}, error) {
	v, err := strconv.ParseUint(sval, 0, 0)
	if err != nil {
		return 0, err
	}
	return uint(v), nil
}
