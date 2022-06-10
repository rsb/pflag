package pflag_test

import (
	"bytes"
	"github.com/rsb/pflag"
	"github.com/stretchr/testify/require"
	"strconv"
	"testing"
)

// This value can be a boolean ("true", "false") or "maybe"
type triStateValue int

const (
	triStateFalse triStateValue = 0
	triStateTrue  triStateValue = 1
	triStateMaybe triStateValue = 2
)

const strTriStateMaybe = "maybe"

func (v *triStateValue) IsBoolFlag() bool {
	return true
}

func (v *triStateValue) Get() interface{} {
	return triStateValue(*v)
}

func (v *triStateValue) Set(s string) error {
	if s == strTriStateMaybe {
		*v = triStateMaybe
		return nil
	}
	boolVal, err := strconv.ParseBool(s)
	if boolVal {
		*v = triStateTrue
	} else {
		*v = triStateFalse
	}
	return err
}

func (v *triStateValue) String() string {
	if *v == triStateMaybe {
		return strTriStateMaybe
	}
	return strconv.FormatBool(*v == triStateTrue)
}

// The type of the flag as required by the pflag.Value interface
func (v *triStateValue) Type() string {
	return "version"
}

func setUpFlagSet(tristate *triStateValue) *pflag.FlagSet {
	f := pflag.NewFlagSet("test", pflag.ContinueOnError)
	*tristate = triStateFalse
	flag := f.VarPF(tristate, "tristate", "t", "tristate value (true, maybe or false)")
	flag.NoOptDefVal = "true"
	return f
}

func TestExplicitTrue(t *testing.T) {
	var tristate triStateValue
	f := setUpFlagSet(&tristate)
	err := f.Parse([]string{"--tristate=true"})
	require.NoError(t, err)
	require.Equal(t, triStateTrue, tristate)
}

func TestImplicitTrue(t *testing.T) {
	var tristate triStateValue
	f := setUpFlagSet(&tristate)
	err := f.Parse([]string{"--tristate"})
	require.NoError(t, err)
	require.Equal(t, triStateTrue, tristate)
}

func TestShortFlag(t *testing.T) {
	var tristate triStateValue
	f := setUpFlagSet(&tristate)
	err := f.Parse([]string{"-t"})
	require.NoError(t, err)
	require.Equal(t, triStateTrue, tristate)
}

func TestShortFlagExtraArgument(t *testing.T) {
	var tristate triStateValue
	f := setUpFlagSet(&tristate)
	// The"maybe"turns into an arg, since short boolean options will only do true/false
	err := f.Parse([]string{"-t", "maybe"})
	require.NoError(t, err)
	require.Equal(t, triStateTrue, tristate)

	args := f.Args()
	require.Len(t, args, 1)
	require.Equal(t, "maybe", args[0])
}

func TestExplicitMaybe(t *testing.T) {
	var tristate triStateValue
	f := setUpFlagSet(&tristate)
	err := f.Parse([]string{"--tristate=maybe"})
	require.NoError(t, err)
	require.Equal(t, triStateMaybe, tristate)
}

func TestExplicitFalse(t *testing.T) {
	var tristate triStateValue
	f := setUpFlagSet(&tristate)
	err := f.Parse([]string{"--tristate=false"})
	require.NoError(t, err)
	require.Equal(t, triStateFalse, tristate)
}

func TestImplicitFalse(t *testing.T) {
	var tristate triStateValue
	f := setUpFlagSet(&tristate)
	err := f.Parse([]string{})
	require.NoError(t, err)
	require.Equal(t, triStateFalse, tristate)
}

func TestInvalidValue(t *testing.T) {
	var tristate triStateValue
	f := setUpFlagSet(&tristate)
	var buf bytes.Buffer
	f.SetOutput(&buf)
	err := f.Parse([]string{"--tristate=invalid"})
	require.Error(t, err)
}

func TestBoolP(t *testing.T) {
	b := pflag.BoolP("bool", "b", false, "bool value in CommandLine")
	c := pflag.BoolP("c", "c", false, "other bool value")
	args := []string{"--bool"}

	err := pflag.CommandLine.Parse(args)
	require.NoError(t, err)
	require.True(t, *b)
	require.False(t, *c)
}
