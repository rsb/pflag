package pflag

import (
	goflag "flag"
	"fmt"
	"io"
	"os"

	"github.com/rsb/failure"
)

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
	shorts            map[byte]*Flag
	args              []string // arguments after flags
	argsLenAtDash     int      // len(args) when a '--' was located when parsing, or -1 if no --
	errorHandling     ErrorHandling
	output            io.Writer // nil means stderr; use Output() accessor
	interspersed      bool      // allow interspersed option/non-option args
	normalizeNameFunc func(f *FlagSet, name string) NormalizedName

	addedGoFlagSets []*goflag.FlagSet
}

// NewFlagSet returns a new, empty flag set with the specified name,
// error handling property and SortFlags set to true.
func NewFlagSet(name string, errorHandling ErrorHandling) *FlagSet {
	f := &FlagSet{
		name:          name,
		errorHandling: errorHandling,
		argsLenAtDash: -1,
		interspersed:  true,
		SortFlags:     true,
	}
	return f
}

// SetInterspersed sets whether to support interspersed option/non-option arguments.
func (f *FlagSet) SetInterspersed(interspersed bool) {
	f.interspersed = interspersed
}

// Init sets the name and error handling property for a flag set.
// By default, the zero FlagSet uses an empty name and the
// ContinueOnError error handling policy.
func (f *FlagSet) Init(name string, errorHandling ErrorHandling) {
	f.name = name
	f.errorHandling = errorHandling
	f.argsLenAtDash = -1
}

// SetNormalizeFunc allows you to add a function which can translate flag names.
// Flags added to the FlagSet will be translated and then when anything tries to
// look up the flag that will also be translated. So it would be possible to create
// a flag named "getURL" and have it translated to "geturl". A user could then pass
// "--getUrl" which may also be translated to "geturl" and everything will work.
func (f *FlagSet) SetNormalizeFunc(n func(f *FlagSet, name string) NormalizedName) {
	f.normalizeNameFunc = n
	f.sortedFormal = f.sortedFormal[:0]
	for fName, flag := range f.formal {
		nName := f.normalizeFlagName(flag.Name)
		if fName == nName {
			continue
		}

		flag.Name = string(nName)
		delete(f.formal, fName)

		f.formal[nName] = flag
		if _, set := f.actual[fName]; set {
			delete(f.actual, fName)
			f.actual[nName] = flag
		}
	}
}

func (f *FlagSet) GetNormalizeFunc() func(f *FlagSet, name string) NormalizedName {
	if f.normalizeNameFunc != nil {
		return f.normalizeNameFunc
	}

	return func(f *FlagSet, name string) NormalizedName { return NormalizedName(name) }
}

// Output returns the destination for usage and error messages. os.Stderr is
// returned if output was not set or was set to nil
func (f *FlagSet) Output() io.Writer {
	if f.output == nil {
		return os.Stderr
	}
	return f.output
}

// SetOutput sets the destination for usage and error messages.
// if output is nil, os.Stderr is used.
func (f *FlagSet) SetOutput(o io.Writer) {
	f.output = o
}

// Name returns the name of the flag set.
func (f *FlagSet) Name() string {
	return f.name
}

// VisitAll visits the flags in lexicographical order of in primordial
// order if f.SortFlags is false, calling fn for each.
// It visits all flags, even those not set.
func (f *FlagSet) VisitAll(fn func(*Flag)) {
	if len(f.formal) == 0 {
		return
	}

	var flags []*Flag
	if f.SortFlags {
		if len(f.formal) != len(f.sortedFormal) {
			f.sortedFormal = sortFlags(f.formal)
		}
		flags = f.sortedFormal
	} else {
		flags = f.orderedFormal
	}

	for _, flag := range flags {
		fn(flag)
	}
}

// HasFlags return a bool to indicate if the FlagSet has any flags defined
func (f *FlagSet) HasFlags() bool {
	return len(f.formal) > 0
}

// HasAvailableFlags returns a bool to indicate if the FlagSet has any flags
// that are not hidden.
func (f *FlagSet) HasAvailableFlags() bool {
	for _, flag := range f.formal {
		if !flag.Hidden {
			return true
		}
	}
	return false
}

// Lookup returns the Flag structure of the named flag, returning nil
// if none exists
func (f *FlagSet) Lookup(name string) *Flag {
	return f.lookup(f.normalizeFlagName(name))
}

// ShortLookup returns the Flag structure of the short-handed flag,
// returning nil if none exists.
// It panics, if len(name) > 1.
func (f *FlagSet) ShortLookup(name string) *Flag {
	if name == "" {
		return nil
	}
	if len(name) > 1 {
		msg := fmt.Sprintf(
			"can not look up short flag which is more than one ASCII character: %q",
			name,
		)
		_, _ = fmt.Fprintf(f.Output(), msg)
		panic(msg)
	}
	c := name[0]
	return f.shorts[c]
}

// ArgsLenAtDash will return the length of f.Args at the moment when a -- was
// found during arg parsing. This allows your program to know which args were
// before the -- and which came after.
func (f *FlagSet) ArgsLenAtDash() int {
	return f.argsLenAtDash
}

// MarkDeprecated indicated that a flag is deprecated in your program. It will
// continue to function but will not show up in help or usage messages. Using
// this flag will also print the given usage.
func (f *FlagSet) MarkDeprecated(name, usage string) error {
	flag := f.Lookup(name)
	if flag == nil {
		return failure.NotFound("flag (%s), does not exist", name)
	}

	if usage == "" {
		return failure.InvalidParam("usage is empty, deprecated msg for (%s) must be set", name)
	}

	flag.Deprecated = usage
	flag.Hidden = true
	return nil
}

// MarkShortDeprecated will mark the shorthand of a flag deprecated in your
// program. It will continue to function but will not show up in help or
// usage messages. Using this flag will also print the given usage
func (f *FlagSet) MarkShortDeprecated(name, usage string) error {
	flag := f.Lookup(name)
	if flag == nil {
		return failure.NotFound("flag (%s), does not exist", name)
	}

	if usage == "" {
		return failure.InvalidParam("usage is empty, deprecated msg for (%s) must be set", name)
	}

	flag.ShortDeprecated = usage
	return nil
}

// MarkHidden sets the flag to 'hidden' in your program.  It will continue to
// function but will not show up in help or usage messages.
func (f *FlagSet) MarkHidden(name string) error {
	flag := f.Lookup(name)
	if flag == nil {
		return failure.NotFound("flag (%s), does not exist", name)
	}

	flag.Hidden = true
	return nil
}

// Set sets the value of the named flag
func (f *FlagSet) Set(name, value string) error {
	normalName := f.normalizeFlagName(name)
	flag := f.lookup(normalName)
	if flag == nil {
		return failure.NotFound("flag (%s), does not exist", name)
	}

	if err := flag.Value.Set(value); err != nil {
		var flagName string
		if flag.Short != "" && flag.ShortDeprecated == "" {
			flagName = fmt.Sprintf("-%s, --%s", flag.Short, flag.Name)
		} else {
			flagName = fmt.Sprintf("--%s", flag.Name)
		}
		return failure.InvalidParam("invalid argument %q for %q flag: %v", value, flagName, err)
	}

	if !flag.Changed {
		if f.actual == nil {
			f.actual = make(map[NormalizedName]*Flag)
		}
		f.actual[normalName] = flag
		f.orderedActual = append(f.orderedActual, flag)
		flag.Changed = true
	}

	if flag.Deprecated != "" {
		fmt.Fprintf(f.Output(), "Flag --%s has been deprecated, %s\n", flag.Name, flag.Deprecated)
	}

	return nil
}

// lookup returns the Flag structure of the named flag, returning nil
// if none exists.
func (f *FlagSet) lookup(name NormalizedName) *Flag {
	return f.formal[name]
}

func (f *FlagSet) normalizeFlagName(name string) NormalizedName {
	return f.GetNormalizeFunc()(f, name)
}
