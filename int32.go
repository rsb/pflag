package pflag

import "strconv"

// -- int32 Value
type int32Value int32

func newInt32Value(val int32, p *int32) *int32Value {
	*p = val
	return (*int32Value)(p)
}

func (i *int32Value) Set(s string) error {
	v, err := strconv.ParseInt(s, 0, 32)
	*i = int32Value(v)
	return err
}

func (i *int32Value) Type() string {
	return "int32"
}

func (i *int32Value) String() string { return strconv.FormatInt(int64(*i), 10) }

func int32Conv(sval string) (interface{}, error) {
	v, err := strconv.ParseInt(sval, 0, 32)
	if err != nil {
		return 0, err
	}
	return int32(v), nil
}

// GetInt32 return the int32 value of a flag with the given name
func (f *FlagSet) GetInt32(name string) (int32, error) {
	val, err := f.getFlagType(name, "int32", int32Conv)
	if err != nil {
		return 0, err
	}
	return val.(int32), nil
}
