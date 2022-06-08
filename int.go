package pflag

import (
	"strconv"
)

// -- int Value
type intValue int

func newIntValue(val int, p *int) *intValue {
	*p = val
	return (*intValue)(p)
}

func (i *intValue) Set(s string) error {
	v, err := strconv.ParseInt(s, 0, 64)
	*i = intValue(v)
	return err
}

func (i *intValue) Type() string {
	return "int"
}

func (i *intValue) String() string { return strconv.Itoa(int(*i)) }

func intConv(sval string) (interface{}, error) {
	return strconv.Atoi(sval)
}

// GetInt return the int value of a flag with the given name
func (f *FlagSet) GetInt(name string) (int, error) {
	val, err := f.getFlagType(name, "int", intConv)
	if err != nil {
		return 0, err
	}
	return val.(int), nil
}
