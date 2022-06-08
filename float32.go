package pflag

import (
	"strconv"
)

// -- float32 Value
type float32Value float32

func newFloat32Value(val float32, p *float32) *float32Value {
	*p = val
	return (*float32Value)(p)
}

func (f *float32Value) Set(s string) error {
	v, err := strconv.ParseFloat(s, 32)
	*f = float32Value(v)
	return err
}

func (f *float32Value) Type() string {
	return "float32"
}

func (f *float32Value) String() string { return strconv.FormatFloat(float64(*f), 'g', -1, 32) }

func float32Conv(sval string) (interface{}, error) {
	v, err := strconv.ParseFloat(sval, 32)
	if err != nil {
		return 0, err
	}
	return float32(v), nil
}
