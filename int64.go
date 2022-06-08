package pflag

import (
	"strconv"
)

// -- int64 Value
type int64Value int64

func newInt64Value(val int64, p *int64) *int64Value {
	*p = val
	return (*int64Value)(p)
}

func (i *int64Value) Set(s string) error {
	v, err := strconv.ParseInt(s, 0, 64)
	*i = int64Value(v)
	return err
}

func (i *int64Value) Type() string {
	return "int64"
}

func (i *int64Value) String() string { return strconv.FormatInt(int64(*i), 10) }

func int64Conv(sval string) (interface{}, error) {
	return strconv.ParseInt(sval, 0, 64)
}

// GetInt64 return the int64 value of a flag with the given name
func (f *FlagSet) GetInt64(name string) (int64, error) {
	val, err := f.getFlagType(name, "int64", int64Conv)
	if err != nil {
		return 0, err
	}
	return val.(int64), nil
}
