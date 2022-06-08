package pflag

import (
	"strconv"
)

// -- count Value
type countValue int

func newCountValue(val int, p *int) *countValue {
	*p = val
	return (*countValue)(p)
}

func (i *countValue) Set(s string) error {
	// "+1" means that no specific value was passed, so increment
	if s == "+1" {
		*i = countValue(*i + 1)
		return nil
	}
	v, err := strconv.ParseInt(s, 0, 0)
	*i = countValue(v)
	return err
}

func (i *countValue) Type() string {
	return "count"
}

func (i *countValue) String() string { return strconv.Itoa(int(*i)) }

func countConv(sval string) (interface{}, error) {
	i, err := strconv.Atoi(sval)
	if err != nil {
		return nil, err
	}
	return i, nil
}
