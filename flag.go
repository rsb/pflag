package pflag

import (
	goflag "flag"
	"io"
	"sort"
)

// NormalizedName is a flag name that has been normalized according to rules
// for the FlagSet (e.g. making '-' and '_' equivalent).
type NormalizedName string

// ParseErrorsWhitelist defines the parsing errors that can be ignored
type ParseErrorsWhitelist struct {
	// UnknownFlags will ignore unknown flags errors and continue parsing rest of the flags
	UnknownFlags bool
}

// FlagSet represents a collection of defined flags.
type FlagSet struct {
	// Usage is the function called when an error occurs while parsing flags.
	// The field is a function (not a method) that may be changed to point to
	// a custom error handler
	Usage func()

	// SortFlags is used to indicate, if user wants to have sorted flags in
	// help/usage message
	SortFlags bool

	// ParseErrorsWhitelist is used to configure a whitelist of errors
	ParseErrorsWhitelist ParseErrorsWhitelist

	name              string
	parsed            bool
	actual            map[NormalizedName]*Flag
	orderedActual     []*Flag
	sortedActual      []*Flag
	formal            map[NormalizedName]*Flag
	orderedFormal     []*Flag
	sortedFormal      []*Flag
	shorthands        map[byte]*Flag
	args              []string  // arguments after flags
	argsLenAtDash     int       // len(args) when a '--' was located when parsing, or -1 if no --
	output            io.Writer // nil means stderr; use Output() accessor
	interspersed      bool      // allow interspersed option/non-option args
	normalizeNameFunc func(f *FlagSet, name string) NormalizedName

	addedGoFlagSets []*goflag.FlagSet
}

// Flag represents the state of a command line flag.
type Flag struct {
	Name            string              // name as it appears on the command line
	Short           string              // one-letter abbreviated flag
	Usage           string              // help message
	Value           Value               // value as set
	Default         string              // default value (as text); for usage message
	Changed         bool                // if the user changed the value (or if left to default)
	NoOptDefVal     string              // default value (as text); if the flag in on the command line without options
	Deprecated      string              // if this flag is deprecated, this string is the new or now thing to use
	Hidden          bool                // allow flags to be hidden from help/usage text
	ShortDeprecated string              // if the shorthand of this flag is deprecated, this string is the new or now thing to use
	Annotations     map[string][]string // used for bash autocomplete code
}

// Value is the interface to the dynamic value stored in a flag.
// (The default value is represented as a string.)
type Value interface {
	String() string
	Set(string) error
	Type() string
}

// SliceValue is a secondary interface to all flags which hold a list
// of values.  This allows full control over the value of list flags,
// and avoids complicated marshalling and unmarshalling to csv.
type SliceValue interface {
	// Append adds the specified value to the end of the flag value list.
	Append(string) error
	// Replace will fully overwrite any data currently in the flag value list.
	Replace([]string) error
	// GetSlice returns the flag value list as an array of strings.
	GetSlice() []string
}

// sortFlags returns the flags as a slice in lexicographical sorted order.
func sortFlags(flags map[NormalizedName]*Flag) []*Flag {
	list := make(sort.StringSlice, len(flags))
	i := 0
	for k := range flags {
		list[i] = string(k)
		i++
	}
	list.Sort()
	result := make([]*Flag, len(list))
	for i, name := range list {
		result[i] = flags[NormalizedName(name)]
	}
	return result
}
