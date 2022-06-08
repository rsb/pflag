package pflag

import (
	"time"
)

// -- time.Duration Value
type durationValue time.Duration

func newDurationValue(val time.Duration, p *time.Duration) *durationValue {
	*p = val
	return (*durationValue)(p)
}

func (d *durationValue) Set(s string) error {
	v, err := time.ParseDuration(s)
	*d = durationValue(v)
	return err
}

func (d *durationValue) Type() string {
	return "duration"
}

func (d *durationValue) String() string { return (*time.Duration)(d).String() }

func durationConv(sval string) (interface{}, error) {
	return time.ParseDuration(sval)
}
