package pflag

import (
	"strconv"
)

// -- uint32 value
type uint32Value uint32

func newUint32Value(val uint32, p *uint32) *uint32Value {
	*p = val
	return (*uint32Value)(p)
}

func (i *uint32Value) Set(s string) error {
	v, err := strconv.ParseUint(s, 0, 32)
	*i = uint32Value(v)
	return err
}

func (i *uint32Value) Type() string {
	return "uint32"
}

func (i *uint32Value) String() string { return strconv.FormatUint(uint64(*i), 10) }

func uint32Conv(sval string) (interface{}, error) {
	v, err := strconv.ParseUint(sval, 0, 32)
	if err != nil {
		return 0, err
	}
	return uint32(v), nil
}

// GetUint32 return the uint32 value of a flag with the given name
func (f *FlagSet) GetUint32(name string) (uint32, error) {
	val, err := f.getFlagType(name, "uint32", uint32Conv)
	if err != nil {
		return 0, err
	}
	return val.(uint32), nil
}
