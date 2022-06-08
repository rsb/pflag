package pflag

import (
	"fmt"
	"strconv"
	"strings"
)

// -- intSlice Value
type intSliceValue struct {
	value   *[]int
	changed bool
}

func newIntSliceValue(val []int, p *[]int) *intSliceValue {
	isv := new(intSliceValue)
	isv.value = p
	*isv.value = val
	return isv
}

func (s *intSliceValue) Set(val string) error {
	ss := strings.Split(val, ",")
	out := make([]int, len(ss))
	for i, d := range ss {
		var err error
		out[i], err = strconv.Atoi(d)
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

func (s *intSliceValue) Type() string {
	return "intSlice"
}

func (s *intSliceValue) String() string {
	out := make([]string, len(*s.value))
	for i, d := range *s.value {
		out[i] = fmt.Sprintf("%d", d)
	}
	return "[" + strings.Join(out, ",") + "]"
}

func (s *intSliceValue) Append(val string) error {
	i, err := strconv.Atoi(val)
	if err != nil {
		return err
	}
	*s.value = append(*s.value, i)
	return nil
}

func (s *intSliceValue) Replace(val []string) error {
	out := make([]int, len(val))
	for i, d := range val {
		var err error
		out[i], err = strconv.Atoi(d)
		if err != nil {
			return err
		}
	}
	*s.value = out
	return nil
}

func (s *intSliceValue) GetSlice() []string {
	out := make([]string, len(*s.value))
	for i, d := range *s.value {
		out[i] = strconv.Itoa(d)
	}
	return out
}

func intSliceConv(val string) (interface{}, error) {
	val = strings.Trim(val, "[]")
	// Empty string would cause a slice with one (empty) entry
	if len(val) == 0 {
		return []int{}, nil
	}
	ss := strings.Split(val, ",")
	out := make([]int, len(ss))
	for i, d := range ss {
		var err error
		out[i], err = strconv.Atoi(d)
		if err != nil {
			return nil, err
		}

	}
	return out, nil
}
