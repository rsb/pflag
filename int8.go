package pflag

import "strconv"

// -- int8 Value
type int8Value int8

func newInt8Value(val int8, p *int8) *int8Value {
	*p = val
	return (*int8Value)(p)
}

func (i *int8Value) Set(s string) error {
	v, err := strconv.ParseInt(s, 0, 8)
	*i = int8Value(v)
	return err
}

func (i *int8Value) Type() string {
	return "int8"
}

func (i *int8Value) String() string { return strconv.FormatInt(int64(*i), 10) }

func int8Conv(sval string) (interface{}, error) {
	v, err := strconv.ParseInt(sval, 0, 8)
	if err != nil {
		return 0, err
	}
	return int8(v), nil
}

// GetInt8 return the int8 value of a flag with the given name
func (f *FlagSet) GetInt8(name string) (int8, error) {
	val, err := f.getFlagType(name, "int8", int8Conv)
	if err != nil {
		return 0, err
	}
	return val.(int8), nil
}
