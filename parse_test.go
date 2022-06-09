package pflag_test

import (
	"github.com/rsb/pflag"
	"github.com/stretchr/testify/require"
	"net"
	"reflect"
	"testing"
	"time"
)

func TestParse(t *testing.T) {
	t.Parallel()
	pflag.ResetForTesting(func() { t.Error("bad parse") })
	testParse(pflag.GetCommandLine(), t)
}

func TestParseAll(t *testing.T) {
	t.Parallel()
	pflag.ResetForTesting(func() { t.Error("bad parse") })
	testParseAll(pflag.GetCommandLine(), t)
}

func TestIgnoreUnknownFlags(t *testing.T) {
	t.Parallel()
	pflag.ResetForTesting(func() { t.Error("bad parse") })
	testParseWithUnknownFlags(pflag.GetCommandLine(), t)
}

func TestFlagSetParse(t *testing.T) {
	t.Parallel()
	testParse(pflag.NewFlagSet("test", pflag.ContinueOnError), t)
}

func testParse(f *pflag.FlagSet, t *testing.T) {
	ok := f.Parsed()
	require.False(t, ok, "f.Parse() = true before Parse")

	boolFlag := f.Bool("bool", false, "bool value")
	bool2Flag := f.Bool("bool2", false, "bool2 value")
	bool3Flag := f.Bool("bool3", false, "bool3 value")
	intFlag := f.Int("int", 0, "int value")
	int8Flag := f.Int8("int8", 0, "int value")
	int16Flag := f.Int16("int16", 0, "int value")
	int32Flag := f.Int32("int32", 0, "int value")
	int64Flag := f.Int64("int64", 0, "int64 value")
	uintFlag := f.Uint("uint", 0, "uint value")
	uint8Flag := f.Uint8("uint8", 0, "uint value")
	uint16Flag := f.Uint16("uint16", 0, "uint value")
	uint32Flag := f.Uint32("uint32", 0, "uint value")
	uint64Flag := f.Uint64("uint64", 0, "uint64 value")
	stringFlag := f.String("string", "0", "string value")
	float32Flag := f.Float32("float32", 0, "float32 value")
	float64Flag := f.Float64("float64", 0, "float64 value")
	ipFlag := f.IP("ip", net.ParseIP("127.0.0.1"), "ip value")
	maskFlag := f.IPMask("mask", pflag.ParseIPv4Mask("0.0.0.0"), "mask value")
	durationFlag := f.Duration("duration", 5*time.Second, "time.Duration value")
	optionalIntNoValueFlag := f.Int("optional-int-no-value", 0, "int value")
	f.Lookup("optional-int-no-value").NoOptDefVal = "9"
	optionalIntWithValueFlag := f.Int("optional-int-with-value", 0, "int value")
	f.Lookup("optional-int-no-value").NoOptDefVal = "9"
	extra := "one-extra-argument"
	args := []string{
		"--bool",
		"--bool2=true",
		"--bool3=false",
		"--int=22",
		"--int8=-8",
		"--int16=-16",
		"--int32=-32",
		"--int64=0x23",
		"--uint", "24",
		"--uint8=8",
		"--uint16=16",
		"--uint32=32",
		"--uint64=25",
		"--string=hello",
		"--float32=-172e12",
		"--float64=2718e28",
		"--ip=10.11.12.13",
		"--mask=255.255.255.0",
		"--duration=2m",
		"--optional-int-no-value",
		"--optional-int-with-value=42",
		extra,
	}
	if err := f.Parse(args); err != nil {
		t.Fatal(err)
	}
	if !f.Parsed() {
		t.Error("f.Parse() = false after Parse")
	}
	if *boolFlag != true {
		t.Error("bool flag should be true, is ", *boolFlag)
	}
	if v, err := f.GetBool("bool"); err != nil || v != *boolFlag {
		t.Error("GetBool does not work.")
	}
	if *bool2Flag != true {
		t.Error("bool2 flag should be true, is ", *bool2Flag)
	}
	if *bool3Flag != false {
		t.Error("bool3 flag should be false, is ", *bool2Flag)
	}
	if *intFlag != 22 {
		t.Error("int flag should be 22, is ", *intFlag)
	}
	if v, err := f.GetInt("int"); err != nil || v != *intFlag {
		t.Error("GetInt does not work.")
	}
	if *int8Flag != -8 {
		t.Error("int8 flag should be 0x23, is ", *int8Flag)
	}
	if *int16Flag != -16 {
		t.Error("int16 flag should be -16, is ", *int16Flag)
	}
	if v, err := f.GetInt8("int8"); err != nil || v != *int8Flag {
		t.Error("GetInt8 does not work.")
	}
	if v, err := f.GetInt16("int16"); err != nil || v != *int16Flag {
		t.Error("GetInt16 does not work.")
	}
	if *int32Flag != -32 {
		t.Error("int32 flag should be 0x23, is ", *int32Flag)
	}
	if v, err := f.GetInt32("int32"); err != nil || v != *int32Flag {
		t.Error("GetInt32 does not work.")
	}
	if *int64Flag != 0x23 {
		t.Error("int64 flag should be 0x23, is ", *int64Flag)
	}
	if v, err := f.GetInt64("int64"); err != nil || v != *int64Flag {
		t.Error("GetInt64 does not work.")
	}
	if *uintFlag != 24 {
		t.Error("uint flag should be 24, is ", *uintFlag)
	}
	if v, err := f.GetUint("uint"); err != nil || v != *uintFlag {
		t.Error("GetUint does not work.")
	}
	if *uint8Flag != 8 {
		t.Error("uint8 flag should be 8, is ", *uint8Flag)
	}
	if v, err := f.GetUint8("uint8"); err != nil || v != *uint8Flag {
		t.Error("GetUint8 does not work.")
	}
	if *uint16Flag != 16 {
		t.Error("uint16 flag should be 16, is ", *uint16Flag)
	}
	if v, err := f.GetUint16("uint16"); err != nil || v != *uint16Flag {
		t.Error("GetUint16 does not work.")
	}
	if *uint32Flag != 32 {
		t.Error("uint32 flag should be 32, is ", *uint32Flag)
	}
	if v, err := f.GetUint32("uint32"); err != nil || v != *uint32Flag {
		t.Error("GetUint32 does not work.")
	}
	if *uint64Flag != 25 {
		t.Error("uint64 flag should be 25, is ", *uint64Flag)
	}
	if v, err := f.GetUint64("uint64"); err != nil || v != *uint64Flag {
		t.Error("GetUint64 does not work.")
	}
	if *stringFlag != "hello" {
		t.Error("string flag should be `hello`, is ", *stringFlag)
	}
	if v, err := f.GetString("string"); err != nil || v != *stringFlag {
		t.Error("GetString does not work.")
	}
	if *float32Flag != -172e12 {
		t.Error("float32 flag should be -172e12, is ", *float32Flag)
	}
	if v, err := f.GetFloat32("float32"); err != nil || v != *float32Flag {
		t.Errorf("GetFloat32 returned %v but float32Flag was %v", v, *float32Flag)
	}
	if *float64Flag != 2718e28 {
		t.Error("float64 flag should be 2718e28, is ", *float64Flag)
	}
	if v, err := f.GetFloat64("float64"); err != nil || v != *float64Flag {
		t.Errorf("GetFloat64 returned %v but float64Flag was %v", v, *float64Flag)
	}
	if !(*ipFlag).Equal(net.ParseIP("10.11.12.13")) {
		t.Error("ip flag should be 10.11.12.13, is ", *ipFlag)
	}
	if v, err := f.GetIP("ip"); err != nil || !v.Equal(*ipFlag) {
		t.Errorf("GetIP returned %v but ipFlag was %v", v, *ipFlag)
	}
	if (*maskFlag).String() != pflag.ParseIPv4Mask("255.255.255.0").String() {
		t.Error("mask flag should be 255.255.255.0, is ", (*maskFlag).String())
	}
	if v, err := f.GetIPv4Mask("mask"); err != nil || v.String() != (*maskFlag).String() {
		t.Errorf("GetIP returned %v maskFlag was %v error was %v", v, *maskFlag, err)
	}
	if *durationFlag != 2*time.Minute {
		t.Error("duration flag should be 2m, is ", *durationFlag)
	}
	if v, err := f.GetDuration("duration"); err != nil || v != *durationFlag {
		t.Error("GetDuration does not work.")
	}
	if _, err := f.GetInt("duration"); err == nil {
		t.Error("GetInt parsed a time.Duration?!?!")
	}
	if *optionalIntNoValueFlag != 9 {
		t.Error("optional int flag should be the default value, is ", *optionalIntNoValueFlag)
	}
	if *optionalIntWithValueFlag != 42 {
		t.Error("optional int flag should be 42, is ", *optionalIntWithValueFlag)
	}
	if len(f.Args()) != 1 {
		t.Error("expected one argument, got", len(f.Args()))
	} else if f.Args()[0] != extra {
		t.Errorf("expected argument %q got %q", extra, f.Args()[0])
	}
}

func testParseAll(f *pflag.FlagSet, t *testing.T) {
	if f.Parsed() {
		t.Error("f.Parse() = true before Parse")
	}
	f.BoolP("boola", "a", false, "bool value")
	f.BoolP("boolb", "b", false, "bool2 value")
	f.BoolP("boolc", "c", false, "bool3 value")
	f.BoolP("boold", "d", false, "bool4 value")
	f.StringP("stringa", "s", "0", "string value")
	f.StringP("stringz", "z", "0", "string value")
	f.StringP("stringx", "x", "0", "string value")
	f.StringP("stringy", "y", "0", "string value")
	f.Lookup("stringx").NoOptDefVal = "1"
	args := []string{
		"-ab",
		"-cs=xx",
		"--stringz=something",
		"-d=true",
		"-x",
		"-y",
		"ee",
	}
	want := []string{
		"boola", "true",
		"boolb", "true",
		"boolc", "true",
		"stringa", "xx",
		"stringz", "something",
		"boold", "true",
		"stringx", "1",
		"stringy", "ee",
	}
	got := []string{}
	store := func(flag *pflag.Flag, value string) error {
		got = append(got, flag.Name)
		if len(value) > 0 {
			got = append(got, value)
		}
		return nil
	}
	if err := f.ParseAll(args, store); err != nil {
		t.Errorf("expected no error, got %s", err)
	}
	if !f.Parsed() {
		t.Errorf("f.Parse() = false after Parse")
	}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("f.ParseAll() fail to restore the args")
		t.Errorf("Got:  %v", got)
		t.Errorf("Want: %v", want)
	}
}

func testParseWithUnknownFlags(f *pflag.FlagSet, t *testing.T) {
	if f.Parsed() {
		t.Error("f.Parse() = true before Parse")
	}
	f.ParseErrorsWhitelist.UnknownFlags = true

	f.BoolP("boola", "a", false, "bool value")
	f.BoolP("boolb", "b", false, "bool2 value")
	f.BoolP("boolc", "c", false, "bool3 value")
	f.BoolP("boold", "d", false, "bool4 value")
	f.BoolP("boole", "e", false, "bool4 value")
	f.StringP("stringa", "s", "0", "string value")
	f.StringP("stringz", "z", "0", "string value")
	f.StringP("stringx", "x", "0", "string value")
	f.StringP("stringy", "y", "0", "string value")
	f.StringP("stringo", "o", "0", "string value")
	f.Lookup("stringx").NoOptDefVal = "1"
	args := []string{
		"-ab",
		"-cs=xx",
		"--stringz=something",
		"--unknown1",
		"unknown1Value",
		"-d=true",
		"-x",
		"--unknown2=unknown2Value",
		"-u=unknown3Value",
		"-p",
		"unknown4Value",
		"-q", // another unknown with bool value
		"-y",
		"ee",
		"--unknown7=unknown7value",
		"--stringo=ovalue",
		"--unknown8=unknown8value",
		"--boole",
		"--unknown6",
		"",
		"-uuuuu",
		"",
		"--unknown10",
		"--unknown11",
	}
	want := []string{
		"boola", "true",
		"boolb", "true",
		"boolc", "true",
		"stringa", "xx",
		"stringz", "something",
		"boold", "true",
		"stringx", "1",
		"stringy", "ee",
		"stringo", "ovalue",
		"boole", "true",
	}
	got := []string{}
	store := func(flag *pflag.Flag, value string) error {
		got = append(got, flag.Name)
		if len(value) > 0 {
			got = append(got, value)
		}
		return nil
	}
	if err := f.ParseAll(args, store); err != nil {
		t.Errorf("expected no error, got %s", err)
	}
	if !f.Parsed() {
		t.Errorf("f.Parse() = false after Parse")
	}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("f.ParseAll() fail to restore the args")
		t.Errorf("Got:  %v", got)
		t.Errorf("Want: %v", want)
	}
}
