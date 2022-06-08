package pflag

import (
	"bytes"
	"encoding/csv"
	"strings"
)

// -- stringSlice Value
type stringSliceValue struct {
	value   *[]string
	changed bool
}

func newStringSliceValue(val []string, p *[]string) *stringSliceValue {
	ssv := new(stringSliceValue)
	ssv.value = p
	*ssv.value = val
	return ssv
}

func readAsCSV(val string) ([]string, error) {
	if val == "" {
		return []string{}, nil
	}
	stringReader := strings.NewReader(val)
	csvReader := csv.NewReader(stringReader)
	return csvReader.Read()
}

func writeAsCSV(vals []string) (string, error) {
	b := &bytes.Buffer{}
	w := csv.NewWriter(b)
	err := w.Write(vals)
	if err != nil {
		return "", err
	}
	w.Flush()
	return strings.TrimSuffix(b.String(), "\n"), nil
}

func (s *stringSliceValue) Set(val string) error {
	v, err := readAsCSV(val)
	if err != nil {
		return err
	}
	if !s.changed {
		*s.value = v
	} else {
		*s.value = append(*s.value, v...)
	}
	s.changed = true
	return nil
}

func (s *stringSliceValue) Type() string {
	return "stringSlice"
}

func (s *stringSliceValue) String() string {
	str, _ := writeAsCSV(*s.value)
	return "[" + str + "]"
}

func (s *stringSliceValue) Append(val string) error {
	*s.value = append(*s.value, val)
	return nil
}

func (s *stringSliceValue) Replace(val []string) error {
	*s.value = val
	return nil
}

func (s *stringSliceValue) GetSlice() []string {
	return *s.value
}

func stringSliceConv(sval string) (interface{}, error) {
	sval = sval[1 : len(sval)-1]
	// An empty string would cause a slice with one (empty) string
	if len(sval) == 0 {
		return []string{}, nil
	}
	return readAsCSV(sval)
}

// GetStringSlice return the []string value of a flag with the given name
func (f *FlagSet) GetStringSlice(name string) ([]string, error) {
	val, err := f.getFlagType(name, "stringSlice", stringSliceConv)
	if err != nil {
		return []string{}, err
	}
	return val.([]string), nil
}
