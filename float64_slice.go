package pflag

import (
	"fmt"
	"strconv"
	"strings"
)

// -- float64Slice Value
type float64SliceValue struct {
	value   *[]float64
	changed bool
}

func newFloat64SliceValue(val []float64, p *[]float64) *float64SliceValue {
	isv := new(float64SliceValue)
	isv.value = p
	*isv.value = val
	return isv
}

func (s *float64SliceValue) Set(val string) error {
	ss := strings.Split(val, ",")
	out := make([]float64, len(ss))
	for i, d := range ss {
		var err error
		out[i], err = strconv.ParseFloat(d, 64)
		if err != nil {
			return err
		}

	}
	if !s.changed {
		*s.value = out
	} else {
		*s.value = append(*s.value, out...)
	}
	s.changed = true
	return nil
}

func (s *float64SliceValue) Type() string {
	return "float64Slice"
}

func (s *float64SliceValue) String() string {
	out := make([]string, len(*s.value))
	for i, d := range *s.value {
		out[i] = fmt.Sprintf("%f", d)
	}
	return "[" + strings.Join(out, ",") + "]"
}

func (s *float64SliceValue) fromString(val string) (float64, error) {
	return strconv.ParseFloat(val, 64)
}

func (s *float64SliceValue) toString(val float64) string {
	return fmt.Sprintf("%f", val)
}

func (s *float64SliceValue) Append(val string) error {
	i, err := s.fromString(val)
	if err != nil {
		return err
	}
	*s.value = append(*s.value, i)
	return nil
}

func (s *float64SliceValue) Replace(val []string) error {
	out := make([]float64, len(val))
	for i, d := range val {
		var err error
		out[i], err = s.fromString(d)
		if err != nil {
			return err
		}
	}
	*s.value = out
	return nil
}

func (s *float64SliceValue) GetSlice() []string {
	out := make([]string, len(*s.value))
	for i, d := range *s.value {
		out[i] = s.toString(d)
	}
	return out
}

func float64SliceConv(val string) (interface{}, error) {
	val = strings.Trim(val, "[]")
	// Empty string would cause a slice with one (empty) entry
	if len(val) == 0 {
		return []float64{}, nil
	}
	ss := strings.Split(val, ",")
	out := make([]float64, len(ss))
	for i, d := range ss {
		var err error
		out[i], err = strconv.ParseFloat(d, 64)
		if err != nil {
			return nil, err
		}

	}
	return out, nil
}
