package pflag

import "strconv"

// -- float64 Value
type float64Value float64

func newFloat64Value(val float64, p *float64) *float64Value {
	*p = val
	return (*float64Value)(p)
}

func (f *float64Value) Set(s string) error {
	v, err := strconv.ParseFloat(s, 64)
	*f = float64Value(v)
	return err
}

func (f *float64Value) Type() string {
	return "float64"
}

func (f *float64Value) String() string { return strconv.FormatFloat(float64(*f), 'g', -1, 64) }

func float64Conv(sval string) (interface{}, error) {
	return strconv.ParseFloat(sval, 64)
}

// GetFloat64 return the float64 value of a flag with the given name
func (f *FlagSet) GetFloat64(name string) (float64, error) {
	val, err := f.getFlagType(name, "float64", float64Conv)
	if err != nil {
		return 0, err
	}
	return val.(float64), nil
}
