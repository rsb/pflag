package pflag

import (
	"bytes"
	goflag "flag"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/rsb/failure"
)

type parseFunc func(flag *Flag, value string) error

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

// GetNormalizeFunc returns the previously set NormalizeFunc of a function which
// does no translation, if not set previously.
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

// Visit visits the flags in lexicographical order or
// in primordial order if f.SortFlags is false, calling fn for each.
// It visits only those flags that have been set.
func (f *FlagSet) Visit(fn func(*Flag)) {
	if len(f.actual) == 0 {
		return
	}

	var flags []*Flag
	if f.SortFlags {
		if len(f.actual) != len(f.sortedActual) {
			f.sortedActual = sortFlags(f.actual)
		}
		flags = f.sortedActual
	} else {
		flags = f.orderedActual
	}

	for _, flag := range flags {
		fn(flag)
	}
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
		_, _ = fmt.Fprintf(
			f.Output(),
			"Flag --%s has been deprecated, %s\n", flag.Name, flag.Deprecated,
		)
	}

	return nil
}

// SetAnnotation allows one to set arbitrary annotations on a flag in
// the FlagSet. This is sometimes used to generate additional bash
// completion information.
func (f *FlagSet) SetAnnotation(name, key string, values []string) error {
	normalName := f.normalizeFlagName(name)

	flag := f.lookup(normalName)
	if flag == nil {
		return failure.NotFound("flag (%s), does not exist", name)
	}

	if flag.Annotations == nil {
		flag.Annotations = map[string][]string{}
	}
	flag.Annotations[key] = values

	return nil
}

// Changed returns true if the flag was explicitly set during Parse()
// otherwise it will be false
func (f *FlagSet) Changed(name string) bool {
	flag := f.Lookup(name)
	// if a flag doesn't exist, it was not changed
	if flag == nil {
		return false
	}

	return flag.Changed
}

// PrintDefaults prints, to standard error unless configured otherwise,
// the values of all defined flags in a set.
func (f *FlagSet) PrintDefaults() {
	usages := f.FlagUsages()
	_, _ = fmt.Fprintf(f.Output(), usages)
}

// FlagUsagesWrapped returns a string containing the usage information
// for all flags in the FlagSet. Wrapped to `cols` columns (0 for no wrapping)
func (f *FlagSet) FlagUsagesWrapped(cols int) string {
	buf := new(bytes.Buffer)

	lines := make([]string, 0, len(f.formal))

	maxLen := 0

	f.VisitAll(func(flag *Flag) {
		if flag.Hidden {
			return
		}

		line := ""
		if flag.Short != "" && flag.ShortDeprecated == "" {
			line = fmt.Sprintf(" %s, --%s", flag.Short, flag.Name)
		} else {
			line = fmt.Sprintf("     --%s", flag.Name)
		}

		varName, usage := UnquoteUsage(flag)
		if varName != "" {
			line += " " + varName
		}

		if flag.NoOptDefVal != "" {
			switch flag.Value.Type() {
			case "string":
				line += fmt.Sprintf("[=\"%s\"]", flag.NoOptDefVal)
			case "bool":
				if flag.NoOptDefVal != "true" {
					line += fmt.Sprintf("[=%s]", flag.NoOptDefVal)
				}
			case "count":
				if flag.NoOptDefVal != "+1" {
					line += fmt.Sprintf("[=%s]", flag.NoOptDefVal)
				}
			default:
				line += fmt.Sprintf("[=%s]", flag.NoOptDefVal)
			}
		}

		// This special character will be replaced with spacing once the
		// correct alignment is calculated
		line += "\x00"
		if len(line) > maxLen {
			maxLen = len(line)
		}

		line += usage
		if !flag.defaultIsZeroValue() {
			if flag.Value.Type() == "string" {
				line += fmt.Sprintf(" (default %q)", flag.Default)
			} else {
				line += fmt.Sprintf(" (default %s)", flag.Default)
			}

			if len(flag.Deprecated) != 0 {
				line += fmt.Sprintf(" (DEPRECATED: %s)", flag.Deprecated)
			}

			lines = append(lines, line)
		}
	})

	for _, line := range lines {
		sIndex := strings.Index(line, "\x00")
		spacing := strings.Repeat(" ", maxLen-sIndex)
		// maxlen + 2 comes from + 1 for the \x00 and + 1
		// for the (deliberate) off-by-one in maxlen-sIndex
		_, _ = fmt.Fprintln(buf, line[:sIndex], spacing, wrap(maxLen+2, cols, line[sIndex+1:]))
	}

	return buf.String()
}

// FlagUsages returns a string containing the usage information for
// all flags in the FlagSet
func (f *FlagSet) FlagUsages() string {
	return f.FlagUsagesWrapped(0)
}

// NFlag returns the number of flags that have been set.
func (f *FlagSet) NFlag() int { return len(f.actual) }

// NArg is the number of arguments remaining after flags have been processed.
func (f *FlagSet) NArg() int { return len(f.args) }

// Arg returns the i'th argument.  Arg(0) is the first remaining argument
// after flags have been processed.
func (f *FlagSet) Arg(i int) string {
	if i < 0 || i >= len(f.args) {
		return ""
	}
	return f.args[i]
}

// Args returns the non-flag arguments.
func (f *FlagSet) Args() []string { return f.args }

// Var defines a flag with the specified name and usage string. The type and
// value of the flag are represented by the first argument, of type Value, which
// typically holds a user-defined implementation of Value. For instance, the
// caller could create a flag that turns a comma-separated string into a slice
// of strings by giving the slice the methods of Value; in particular, Set would
// decompose the comma-separated string into the slice.
func (f *FlagSet) Var(value Value, name string, usage string) {
	f.VarP(value, name, "", usage)
}

// VarPF is like VarP, but returns the flag created
func (f *FlagSet) VarPF(value Value, name, short, usage string) *Flag {
	// Remember the default value as a string; it won't change.
	flag := &Flag{
		Name:    name,
		Short:   short,
		Usage:   usage,
		Value:   value,
		Default: value.String(),
	}

	f.AddFlag(flag)
	return flag
}

// VarP is like Var, but accepts a shorthand letter that can be used after a single dash.
func (f *FlagSet) VarP(value Value, name, short, usage string) {
	_ = f.VarPF(value, name, short, usage)
}

// AddFlag will add the flag to the FlagSet
func (f *FlagSet) AddFlag(flag *Flag) {
	normalizedFlagName := f.normalizeFlagName(flag.Name)

	_, alreadyThere := f.formal[normalizedFlagName]
	if alreadyThere {
		msg := fmt.Sprintf("%s flag redefined: %s", f.name, flag.Name)
		_, _ = fmt.Fprintln(f.Output(), msg)
		panic(msg) // Happens only if flags are declared with identical names
	}
	if f.formal == nil {
		f.formal = make(map[NormalizedName]*Flag)
	}

	flag.Name = string(normalizedFlagName)
	f.formal[normalizedFlagName] = flag
	f.orderedFormal = append(f.orderedFormal, flag)

	if flag.Short == "" {
		return
	}
	if len(flag.Short) > 1 {
		msg := fmt.Sprintf("%q shorthand is more than one ASCII character", flag.Short)
		_, _ = fmt.Fprintf(f.Output(), msg)
		panic(msg)
	}
	if f.shorts == nil {
		f.shorts = make(map[byte]*Flag)
	}
	c := flag.Short[0]
	used, alreadyThere := f.shorts[c]
	if alreadyThere {
		msg := fmt.Sprintf("unable to redefine %q shorthand in %q flagset: it's already used for %q flag", c, f.name, used.Name)
		_, _ = fmt.Fprintf(f.Output(), msg)
		panic(msg)
	}
	f.shorts[c] = flag
}

// AddFlagSet adds one FlagSet to another. If a flag is already present in f
// the flag from newSet will be ignored.
func (f *FlagSet) AddFlagSet(newSet *FlagSet) {
	if newSet == nil {
		return
	}
	newSet.VisitAll(func(flag *Flag) {
		if f.Lookup(flag.Name) == nil {
			f.AddFlag(flag)
		}
	})
}

// Parse parses flag definitions from the argument list, which should not
// include the command name.  Must be called after all flags in the FlagSet
// are defined and before flags are accessed by the program.
// The return value will be ErrHelp if -help was set but not defined.
func (f *FlagSet) Parse(arguments []string) error {
	if f.addedGoFlagSets != nil {
		for _, goFlagSet := range f.addedGoFlagSets {
			_ = goFlagSet.Parse(nil)
		}
	}
	f.parsed = true

	if len(arguments) < 0 {
		return nil
	}

	f.args = make([]string, 0, len(arguments))

	set := func(flag *Flag, value string) error {
		return f.Set(flag.Name, value)
	}

	err := f.parseArgs(arguments, set)
	if err != nil {
		switch f.errorHandling {
		case ContinueOnError:
			return err
		case ExitOnError:
			fmt.Println(err)
			os.Exit(2)
		case PanicOnError:
			panic(err)
		}
	}
	return nil
}

// ParseAll parses flag definitions from the argument list, which should not
// include the command name. The arguments for fn are flag and value. Must be
// called after all flags in the FlagSet are defined and before flags are
// accessed by the program. The return value will be ErrHelp if -help was set
// but not defined.
func (f *FlagSet) ParseAll(arguments []string, fn func(flag *Flag, value string) error) error {
	f.parsed = true
	f.args = make([]string, 0, len(arguments))

	err := f.parseArgs(arguments, fn)
	if err != nil {
		switch f.errorHandling {
		case ContinueOnError:
			return err
		case ExitOnError:
			os.Exit(2)
		case PanicOnError:
			panic(err)
		}
	}
	return nil
}

// Parsed reports whether f.Parse has been called.
func (f *FlagSet) Parsed() bool {
	return f.parsed
}

// lookup returns the Flag structure of the named flag, returning nil
// if none exists.
func (f *FlagSet) lookup(name NormalizedName) *Flag {
	return f.formal[name]
}

func (f *FlagSet) normalizeFlagName(name string) NormalizedName {
	return f.GetNormalizeFunc()(f, name)
}

// func to return a given type for a given flag name
func (f *FlagSet) getFlagType(name string, ftype string, convFunc func(sval string) (interface{}, error)) (interface{}, error) {
	flag := f.Lookup(name)
	if flag == nil {
		err := failure.InvalidState("flag accessed but not defined: %s", name)
		return nil, err
	}

	if flag.Value.Type() != ftype {
		err := failure.InvalidState("trying to get %s value of flag of type %s", ftype, flag.Value.Type())
		return nil, err
	}

	sval := flag.Value.String()
	result, err := convFunc(sval)
	if err != nil {
		return nil, err
	}
	return result, nil
}

// failf prints to standard error a formatted error and usage message and
// returns the error.
func (f *FlagSet) failf(format string, a ...interface{}) error {
	err := fmt.Errorf(format, a...)
	if f.errorHandling != ContinueOnError {
		_, _ = fmt.Fprintln(f.Output(), err)
		f.usage()
	}
	return err
}

// usage calls the Usage method for the flag set, or the usage function if
// the flag set is CommandLine.
func (f *FlagSet) usage() {
	if f == CommandLine {
		// Usage()
	} else if f.Usage == nil {
		defaultUsage(f)
	} else {
		f.Usage()
	}
}

func (f *FlagSet) parseLongArg(s string, args []string, fn parseFunc) (a []string, err error) {
	a = args
	name := s[2:]
	if len(name) == 0 || name[0] == '-' || name[0] == '=' {
		err = f.failf("bad flag syntax: %s", s)
		return
	}

	split := strings.SplitN(name, "=", 2)
	name = split[0]
	flag, exists := f.formal[f.normalizeFlagName(name)]

	if !exists {
		switch {
		case name == "help":
			f.usage()
			return a, ErrHelp
		case f.ParseErrorsWhitelist.UnknownFlags:
			// --unknown=unknownval arg ...
			// we do not want to lose arg in this case
			if len(split) >= 2 {
				return a, nil
			}

			return stripUnknownFlagValue(a), nil
		default:
			err = f.failf("unknown flag: --%s", name)
			return
		}
	}

	var value string
	if len(split) == 2 {
		// '--flag=arg'
		value = split[1]
	} else if flag.NoOptDefVal != "" {
		// '--flag' (arg was optional)
		value = flag.NoOptDefVal
	} else if len(a) > 0 {
		// '--flag arg'
		value = a[0]
		a = a[1:]
	} else {
		// '--flag' (arg was required)
		err = f.failf("flag needs an argument: %s", s)
		return
	}

	err = fn(flag, value)
	if err != nil {
		_ = f.failf(err.Error())
	}
	return
}

func (f *FlagSet) parseSingleShortArg(shorthands string, args []string, fn parseFunc) (outShorts string, outArgs []string, err error) {
	outArgs = args

	if strings.HasPrefix(shorthands, "test.") {
		return
	}

	outShorts = shorthands[1:]
	c := shorthands[0]

	flag, exists := f.shorts[c]
	if !exists {
		switch {
		case c == 'h':
			f.usage()
			err = ErrHelp
			return
		case f.ParseErrorsWhitelist.UnknownFlags:
			// '-f=arg arg ...'
			// we do not want to lose arg in this case
			if len(shorthands) > 2 && shorthands[1] == '=' {
				outShorts = ""
				return
			}

			outArgs = stripUnknownFlagValue(outArgs)
			return
		default:
			err = f.failf("unknown shorthand flag: %q in -%s", c, shorthands)
			return
		}
	}

	var value string
	if len(shorthands) > 2 && shorthands[1] == '=' {
		// '-f=arg'
		value = shorthands[2:]
		outShorts = ""
	} else if flag.NoOptDefVal != "" {
		// '-f' (arg was optional)
		value = flag.NoOptDefVal
	} else if len(shorthands) > 1 {
		// '-farg'
		value = shorthands[1:]
		outShorts = ""
	} else if len(args) > 0 {
		// '-f arg'
		value = args[0]
		outArgs = args[1:]
	} else {
		// '-f' (arg was required)
		err = f.failf("flag needs an argument: %q in -%s", c, shorthands)
		return
	}

	if flag.ShortDeprecated != "" {
		_, _ = fmt.Fprintf(f.Output(), "Flag shorthand -%s has been deprecated, %s\n", flag.Short, flag.ShortDeprecated)
	}

	err = fn(flag, value)
	if err != nil {
		_ = f.failf(err.Error())
	}
	return
}

func (f *FlagSet) parseShortArg(s string, args []string, fn parseFunc) (a []string, err error) {
	a = args
	shorthands := s[1:]

	// "shorthands" can be a series of shorthand letters of flags (e.g. "-vvv").
	for len(shorthands) > 0 {
		shorthands, a, err = f.parseSingleShortArg(shorthands, args, fn)
		if err != nil {
			return
		}
	}

	return
}

func (f *FlagSet) parseArgs(args []string, fn parseFunc) (err error) {
	for len(args) > 0 {
		s := args[0]
		args = args[1:]
		if len(s) == 0 || s[0] != '-' || len(s) == 1 {
			if !f.interspersed {
				f.args = append(f.args, s)
				f.args = append(f.args, args...)
				return nil
			}
			f.args = append(f.args, s)
			continue
		}

		if s[1] == '-' {
			if len(s) == 2 { // "--" terminates the flags
				f.argsLenAtDash = len(f.args)
				f.args = append(f.args, args...)
				break
			}
			args, err = f.parseLongArg(s, args, fn)
		} else {
			args, err = f.parseShortArg(s, args, fn)
		}
		if err != nil {
			return
		}
	}
	return
}
