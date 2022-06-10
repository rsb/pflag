package pflag_test

import (
	"bytes"
	"fmt"
	"github.com/rsb/pflag"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"io/ioutil"
	"reflect"
	"strings"
	"testing"
)

var (
	testBool                     = pflag.Bool("test_bool", false, "bool value")
	testInt                      = pflag.Int("test_int", 0, "int value")
	testInt64                    = pflag.Int64("test_int64", 0, "int64 value")
	testUint                     = pflag.Uint("test_uint", 0, "uint value")
	testUint64                   = pflag.Uint64("test_uint64", 0, "uint64 value")
	testString                   = pflag.String("test_string", "0", "string value")
	testFloat                    = pflag.Float64("test_float64", 0, "float64 value")
	testDuration                 = pflag.Duration("test_duration", 0, "time.Duration value")
	testOptionalInt              = pflag.Int("test_optional_int", 0, "optional int value")
	normalizeFlagNameInvocations = 0
)

func boolString(s string) string {
	if s == "0" {
		return "false"
	}
	return "true"
}

func setup() {
	testBool = pflag.Bool("test_bool", false, "bool value")
	testInt = pflag.Int("test_int", 0, "int value")
	testInt64 = pflag.Int64("test_int64", 0, "int64 value")
	testUint = pflag.Uint("test_uint", 0, "uint value")
	testUint64 = pflag.Uint64("test_uint64", 0, "uint64 value")
	testString = pflag.String("test_string", "0", "string value")
	testFloat = pflag.Float64("test_float64", 0, "float64 value")
	testDuration = pflag.Duration("test_duration", 0, "time.Duration value")
	testOptionalInt = pflag.Int("test_optional_int", 0, "optional int value")
	normalizeFlagNameInvocations = 0
}

func TestEverything(t *testing.T) {

	var err error
	m := make(map[string]*pflag.Flag)
	desired := "0"
	visitor := func(f *pflag.Flag) {
		if len(f.Name) > 5 && f.Name[0:5] == "test_" {
			m[f.Name] = f
			ok := false
			switch {
			case f.Value.String() == desired:
				ok = true
			case f.Name == "test_bool" && f.Value.String() == boolString(desired):
				ok = true
			case f.Name == "test_duration" && f.Value.String() == desired+"s":
				ok = true
			}
			assert.True(t, ok)
		}
	}

	pflag.ResetForTesting(func() {})
	setup()
	pflag.VisitAll(visitor)
	require.Len(t, m, 9, "Visit misses some flags")

	m = make(map[string]*pflag.Flag)
	pflag.Visit(visitor)
	require.Len(t, m, 0, "Visit sees unset flags")

	// Now set all flags
	err = pflag.Set("test_bool", "true")
	require.NoError(t, err)

	err = pflag.Set("test_int", "1")
	require.NoError(t, err)

	err = pflag.Set("test_int64", "1")
	require.NoError(t, err)

	err = pflag.Set("test_uint", "1")
	require.NoError(t, err)

	err = pflag.Set("test_uint64", "1")
	require.NoError(t, err)

	err = pflag.Set("test_string", "1")
	require.NoError(t, err)

	err = pflag.Set("test_float64", "1")
	require.NoError(t, err)

	err = pflag.Set("test_duration", "1s")
	require.NoError(t, err)

	err = pflag.Set("test_optional_int", "1")
	require.NoError(t, err)

	desired = "1"
	pflag.Visit(visitor)
	require.Len(t, m, 9, "Visit fails after set")

	// Now test they're visited in sort order.
	var flagNames []string
	pflag.Visit(func(f *pflag.Flag) { flagNames = append(flagNames, f.Name) })
	// require.True(t, sort.StringsAreSorted(flagNames), "flag names are not sorted: %v", flagNames)
}

func TestUsage(t *testing.T) {
	called := false
	pflag.ResetForTesting(func() { called = true })
	err := pflag.GetCommandLine().Parse([]string{"--x"})
	require.Error(t, err, "parse did not fail for unknown flag")
	require.False(t, called, "did call Usage while using ContinueOnError")
}

func TestAddFlagSet(t *testing.T) {
	oldSet := pflag.NewFlagSet("old", pflag.ContinueOnError)
	newSet := pflag.NewFlagSet("new", pflag.ContinueOnError)

	oldSet.String("flag1", "flag1", "flag1")
	oldSet.String("flag2", "flag2", "flag2")

	newSet.String("flag2", "flag2", "flag2")
	newSet.String("flag3", "flag3", "flag3")

	f2 := newSet.Lookup("flag2")
	require.NotNil(t, f2)
	f3 := newSet.Lookup("flag3")
	require.NotNil(t, f3)

	oldSet.AddFlagSet(newSet)

	f2Old := oldSet.Lookup("flag2")
	require.NotNil(t, f2Old)
	assert.Equal(t, f2, f2Old)

	f3Old := oldSet.Lookup("flag3")
	require.NotNil(t, f3Old)
	assert.Equal(t, f3, f3Old)
}

func TestAnnotation(t *testing.T) {
	f := pflag.NewFlagSet("shorthand", pflag.ContinueOnError)

	err := f.SetAnnotation("missing-flag", "key", nil)
	require.Error(t, err, "Expected error setting annotation on non-existent flag")

	f.StringP("stringa", "a", "", "string value")
	err = f.SetAnnotation("stringa", "key", nil)
	require.NoError(t, err, "f.SetAnnotation is not expected to fail with nil value")

	annotation := f.Lookup("stringa").Annotations["key"]
	require.Nil(t, annotation, "Not expecting to find annotation")

	f.StringP("stringb", "b", "", "string2 value")

	err = f.SetAnnotation("stringb", "key", []string{"value1"})
	require.NoError(t, err, "f.SetAnnotation is not expected to fail")

	annotation = f.Lookup("stringb").Annotations["key"]
	require.True(t, reflect.DeepEqual(annotation, []string{"value1"}))

	err = f.SetAnnotation("stringb", "key", []string{"value2"})
	require.NoError(t, err, "f.SetAnnotation is not expected to fail")

	annotation = f.Lookup("stringb").Annotations["key"]
	require.True(t, reflect.DeepEqual(annotation, []string{"value2"}))
}

func TestName(t *testing.T) {
	flagSetName := "bob"
	f := pflag.NewFlagSet(flagSetName, pflag.ContinueOnError)

	givenName := f.Name()
	require.Equal(t, givenName, flagSetName)
}

func TestShorthand(t *testing.T) {
	f := pflag.NewFlagSet("shorthand", pflag.ContinueOnError)
	if f.Parsed() {
		t.Error("f.Parse() = true before Parse")
	}
	boolaFlag := f.BoolP("boola", "a", false, "bool value")
	boolbFlag := f.BoolP("boolb", "b", false, "bool2 value")
	boolcFlag := f.BoolP("boolc", "c", false, "bool3 value")
	booldFlag := f.BoolP("boold", "d", false, "bool4 value")
	stringaFlag := f.StringP("stringa", "s", "0", "string value")
	stringzFlag := f.StringP("stringz", "z", "0", "string value")
	extra := "interspersed-argument"
	notaflag := "--i-look-like-a-flag"
	args := []string{
		"-ab",
		extra,
		"-cs",
		"hello",
		"-z=something",
		"-d=true",
		"--",
		notaflag,
	}
	f.SetOutput(ioutil.Discard)

	err := f.Parse(args)
	require.NoError(t, err)
	require.True(t, f.Parsed(), "f.Parse() should not be false after Parse")

	require.True(t, *boolaFlag)
	require.True(t, *boolbFlag)
	require.True(t, *boolcFlag)
	require.True(t, *booldFlag)
	require.Equal(t, "hello", *stringaFlag)
	require.Equal(t, "something", *stringzFlag)

	resultArgs := f.Args()
	require.Len(t, resultArgs, 2)
	require.Equal(t, 1, f.ArgsLenAtDash())
	require.Equal(t, extra, resultArgs[0])
	require.Equal(t, notaflag, resultArgs[1])
}

func TestShorthandLookup(t *testing.T) {
	f := pflag.NewFlagSet("shorthand", pflag.ContinueOnError)
	if f.Parsed() {
		t.Error("f.Parse() = true before Parse")
	}
	f.BoolP("boola", "a", false, "bool value")
	f.BoolP("boolb", "b", false, "bool2 value")
	args := []string{
		"-ab",
	}
	f.SetOutput(ioutil.Discard)

	err := f.Parse(args)
	require.NoError(t, err)

	require.True(t, f.Parsed(), "f.Parse() = false after Parse")

	flag := f.ShortLookup("a")
	require.NotNil(t, flag)
	require.Equal(t, "boola", flag.Name)

	flag = f.ShortLookup("")
	require.Nil(t, flag)
	defer func() {
		recover()
	}()

	flag = f.ShortLookup("ab")
	// should NEVER get here. lookup should panic. defer'd func should recover it.
	t.Errorf("f.ShorthandLookup(\"ab\") did not panic")
}

func TestChangedHelper(t *testing.T) {
	f := pflag.NewFlagSet("changedtest", pflag.ContinueOnError)
	f.Bool("changed", false, "changed bool")
	f.Bool("settrue", true, "true to true")
	f.Bool("setfalse", false, "false to false")
	f.Bool("unchanged", false, "unchanged bool")

	args := []string{"--changed", "--settrue", "--setfalse=false"}

	err := f.Parse(args)
	require.NoError(t, err)
	require.True(t, f.Parsed())
	require.True(t, f.Changed("changed"))
	require.True(t, f.Changed("settrue"))
	require.True(t, f.Changed("setfalse"))
	require.False(t, f.Changed("unchanged"))
	require.False(t, f.Changed("invalid"))

	require.Equal(t, -1, f.ArgsLenAtDash())
}

func testWordSepNormalizedNames(args []string, t *testing.T) {
	f := pflag.NewFlagSet("normalized", pflag.ContinueOnError)
	if f.Parsed() {
		t.Error("f.Parse() = true before Parse")
	}
	withDashFlag := f.Bool("with-dash-flag", false, "bool value")
	// Set this after some flags have been added and before others.
	f.SetNormalizeFunc(wordSepNormalizeFunc)
	withUnderFlag := f.Bool("with_under_flag", false, "bool value")
	withBothFlag := f.Bool("with-both_flag", false, "bool value")

	err := f.Parse(args)
	require.NoError(t, err)
	require.True(t, f.Parsed())
	require.True(t, *withDashFlag)
	require.True(t, *withUnderFlag)
	require.True(t, *withBothFlag)
}

func replaceSeparators(name string, from []string, to string) string {
	result := name
	for _, sep := range from {
		result = strings.Replace(result, sep, to, -1)
	}
	// Type convert to indicate normalization has been done.
	return result
}

func TestWordSepNormalizedNames(t *testing.T) {
	args := []string{
		"--with-dash-flag",
		"--with-under-flag",
		"--with-both-flag",
	}
	testWordSepNormalizedNames(args, t)

	args = []string{
		"--with_dash_flag",
		"--with_under_flag",
		"--with_both_flag",
	}
	testWordSepNormalizedNames(args, t)

	args = []string{
		"--with-dash_flag",
		"--with-under_flag",
		"--with-both_flag",
	}
	testWordSepNormalizedNames(args, t)
}

func TestCustomNormalizedNames(t *testing.T) {
	f := pflag.NewFlagSet("normalized", pflag.ContinueOnError)
	require.False(t, f.Parsed())

	validFlag := f.Bool("valid-flag", false, "bool value")
	f.SetNormalizeFunc(aliasAndWordSepFlagNames)
	someOtherFlag := f.Bool("some-other-flag", false, "bool value")

	args := []string{"--old_valid_flag", "--some-other_flag"}

	err := f.Parse(args)
	require.NoError(t, err)
	require.True(t, *validFlag)
	require.True(t, *someOtherFlag)
}

// Every flag we add, the name (displayed also in usage) should normalized
func TestNormalizationFuncShouldChangeFlagName(t *testing.T) {
	// Test normalization after addition
	f := pflag.NewFlagSet("normalized", pflag.ContinueOnError)

	f.Bool("valid_flag", false, "bool value")
	result := f.Lookup("valid_flag")
	require.Equal(t, "valid_flag", result.Name)

	f.SetNormalizeFunc(wordSepNormalizeFunc)
	result = f.Lookup("valid_flag")
	require.Equal(t, "valid.flag", result.Name)

	// Test normalization before addition
	f = pflag.NewFlagSet("normalized", pflag.ContinueOnError)
	f.SetNormalizeFunc(wordSepNormalizeFunc)
	f.Bool("valid_flag", false, "bool value")

	result = f.Lookup("valid_flag")
	require.Equal(t, "valid.flag", result.Name)
}

// Related to https://github.com/spf13/cobra/issues/521.
func TestNormalizationSharedFlags(t *testing.T) {
	f := pflag.NewFlagSet("set f", pflag.ContinueOnError)
	g := pflag.NewFlagSet("set g", pflag.ContinueOnError)

	nfunc := wordSepNormalizeFunc
	testName := "valid_flag"
	normName := nfunc(nil, testName)
	require.NotEqual(t, testName, string(normName))

	f.Bool(testName, false, "bool value")
	g.AddFlagSet(f)

	f.SetNormalizeFunc(nfunc)
	g.SetNormalizeFunc(nfunc)

	require.Len(t, f.Formal(), 1, "Normalizing flags should not result in duplication")

	require.Equal(t, string(normName), f.OrderedFormal()[0].Name)
	for k := range f.Formal() {
		require.Equal(t, "valid.flag", string(k))
	}

	if !reflect.DeepEqual(f.Formal(), g.Formal()) || !reflect.DeepEqual(f.OrderedFormal(), g.OrderedFormal()) {
		t.Error("Two flag sets sharing the same flags should stay consistent after being normalized. Original set:", f.Formal(), "Duplicate set:", g.Formal())
	}
}

func TestNormalizationSetFlags(t *testing.T) {
	f := pflag.NewFlagSet("normalized", pflag.ContinueOnError)
	nfunc := wordSepNormalizeFunc
	testName := "valid_flag"
	normName := nfunc(nil, testName)
	require.NotEqual(t, testName, string(normName))

	f.Bool(testName, false, "bool value")
	err := f.Set(testName, "true")
	require.NoError(t, err)
	f.SetNormalizeFunc(nfunc)
	require.Len(t, f.Formal(), 1)

	require.Equal(t, f.OrderedFormal()[0].Name, string(normName))

	for k := range f.Formal() {
		require.Equal(t, "valid.flag", string(k))
	}

	if !reflect.DeepEqual(f.Formal(), f.Actual()) {
		t.Error("The map of set flags should get normalized. Formal:", f.Formal(), "Actual:", f.Actual())
	}
}

// Declare a user-defined flag type.
type flagVar []string

func (f *flagVar) String() string {
	return fmt.Sprint([]string(*f))
}

func (f *flagVar) Set(value string) error {
	*f = append(*f, value)
	return nil
}

func (f *flagVar) Type() string {
	return "flagVar"
}

func TestUserDefined(t *testing.T) {
	var flags pflag.FlagSet
	flags.Init("test", pflag.ContinueOnError)
	var v flagVar
	flags.VarP(&v, "v", "v", "usage")

	err := flags.Parse([]string{"--v=1", "-v2", "-v", "3"})
	require.NoError(t, err)
	require.Len(t, v, 3)

	expect := "[1 2 3]"
	require.Equal(t, expect, v.String())
}

func TestSetOutput(t *testing.T) {
	var flags pflag.FlagSet
	var buf bytes.Buffer
	flags.SetOutput(&buf)
	flags.Init("test", pflag.ContinueOnErrorWithWarn)
	err := flags.Parse([]string{"--unknown"})
	require.Error(t, err)
	require.Contains(t, buf.String(), "--unknown")
}

func TestOutput(t *testing.T) {
	var flags pflag.FlagSet
	var buf bytes.Buffer
	expect := "an example string"
	flags.SetOutput(&buf)
	_, _ = fmt.Fprint(flags.Output(), expect)
	if out := buf.String(); !strings.Contains(out, expect) {
		t.Errorf("expected output %q; got %q", expect, out)
	}
}

func wordSepNormalizeFunc(f *pflag.FlagSet, name string) pflag.NormalizedName {
	seps := []string{"-", "_"}
	name = replaceSeparators(name, seps, ".")
	normalizeFlagNameInvocations++

	return pflag.NormalizedName(name)
}

func aliasAndWordSepFlagNames(f *pflag.FlagSet, name string) pflag.NormalizedName {
	seps := []string{"-", "_"}

	oldName := replaceSeparators("old-valid_flag", seps, ".")
	newName := replaceSeparators("valid-flag", seps, ".")

	name = replaceSeparators(name, seps, ".")
	switch name {
	case oldName:
		name = newName
	}

	return pflag.NormalizedName(name)
}
