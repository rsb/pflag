package pflag

// -- string Value
type stringValue string

func newStringValue(val string, p *string) *stringValue {
	*p = val
	return (*stringValue)(p)
}

func (s *stringValue) Set(val string) error {
	*s = stringValue(val)
	return nil
}
func (s *stringValue) Type() string {
	return "string"
}

func (s *stringValue) String() string { return string(*s) }

func stringConv(sval string) (interface{}, error) {
	return sval, nil
}

// GetString return the string value of a flag with the given name
func (f *FlagSet) GetString(name string) (string, error) {
	val, err := f.getFlagType(name, "string", stringConv)
	if err != nil {
		return "", err
	}
	return val.(string), nil
}
