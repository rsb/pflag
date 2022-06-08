package pflag

import (
	"strconv"
)

// -- int16 Value
type int16Value int16

func newInt16Value(val int16, p *int16) *int16Value {
	*p = val
	return (*int16Value)(p)
}

func (i *int16Value) Set(s string) error {
	v, err := strconv.ParseInt(s, 0, 16)
	*i = int16Value(v)
	return err
}

func (i *int16Value) Type() string {
	return "int16"
}

func (i *int16Value) String() string { return strconv.FormatInt(int64(*i), 10) }

func int16Conv(sval string) (interface{}, error) {
	v, err := strconv.ParseInt(sval, 0, 16)
	if err != nil {
		return 0, err
	}
	return int16(v), nil
}

// GetInt16 returns the int16 value of a flag with the given name
func (f *FlagSet) GetInt16(name string) (int16, error) {
	val, err := f.getFlagType(name, "int16", int16Conv)
	if err != nil {
		return 0, err
	}
	return val.(int16), nil
}
