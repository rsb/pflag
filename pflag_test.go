package pflag_test

import (
	"github.com/rsb/pflag"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"reflect"
	"sort"
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
	require.True(t, sort.StringsAreSorted(flagNames), "flag names are not sorted: %v", flagNames)
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
