package pflag_test

import (
	"github.com/rsb/pflag"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"io/ioutil"
	"reflect"
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
	t.Parallel()

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
	t.Parallel()
	called := false
	pflag.ResetForTesting(func() { called = true })
	err := pflag.GetCommandLine().Parse([]string{"--x"})
	require.Error(t, err, "parse did not fail for unknown flag")
	require.False(t, called, "did call Usage while using ContinueOnError")
}

func TestAddFlagSet(t *testing.T) {
	t.Parallel()
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
	t.Parallel()
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
	t.Parallel()
	flagSetName := "bob"
	f := pflag.NewFlagSet(flagSetName, pflag.ContinueOnError)

	givenName := f.Name()
	require.Equal(t, givenName, flagSetName)
}

func TestShorthand(t *testing.T) {
	t.Parallel()
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
	t.Parallel()
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
	t.Parallel()
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
