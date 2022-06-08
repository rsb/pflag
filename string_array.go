package pflag

// -- stringArray Value
type stringArrayValue struct {
	value   *[]string
	changed bool
}

func newStringArrayValue(val []string, p *[]string) *stringArrayValue {
	ssv := new(stringArrayValue)
	ssv.value = p
	*ssv.value = val
	return ssv
}

func (s *stringArrayValue) Set(val string) error {
	if !s.changed {
		*s.value = []string{val}
		s.changed = true
	} else {
		*s.value = append(*s.value, val)
	}
	return nil
}

func (s *stringArrayValue) Append(val string) error {
	*s.value = append(*s.value, val)
	return nil
}

func (s *stringArrayValue) Replace(val []string) error {
	out := make([]string, len(val))
	for i, d := range val {
		out[i] = d
	}
	*s.value = out
	return nil
}

func (s *stringArrayValue) GetSlice() []string {
	out := make([]string, len(*s.value))
	for i, d := range *s.value {
		out[i] = d
	}
	return out
}

func (s *stringArrayValue) Type() string {
	return "stringArray"
}

func (s *stringArrayValue) String() string {
	str, _ := writeAsCSV(*s.value)
	return "[" + str + "]"
}

func stringArrayConv(sval string) (interface{}, error) {
	sval = sval[1 : len(sval)-1]
	// An empty string would cause a array with one (empty) string
	if len(sval) == 0 {
		return []string{}, nil
	}
	return readAsCSV(sval)
}
