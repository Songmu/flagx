package flagx_test

import (
	"bytes"
	"errors"
	"flag"
	"fmt"
	"io"
	"net/netip"
	"reflect"
	"strings"
	"testing"
	"time"

	"github.com/Songmu/flagx"
)

type listValue []string

func (v *listValue) String() string { return strings.Join(*v, ",") }
func (v *listValue) Set(s string) error {
	*v = append(*v, s)
	return nil
}

func TestRegistration(t *testing.T) {
	tests := []struct {
		name     string
		register func(*flag.FlagSet, string) func() string
		initial  string
		input    string
		want     string
	}{
		{"string", func(fs *flag.FlagSet, names string) func() string {
			var v string
			flagx.StringVar(fs, &v, names, "initial", "an option")
			return func() string { return v }
		}, "initial", "changed", "changed"},
		{"bool", func(fs *flag.FlagSet, names string) func() string {
			var v bool
			flagx.BoolVar(fs, &v, names, false, "an option")
			return func() string { return fmt.Sprint(v) }
		}, "false", "true", "true"},
		{"int", func(fs *flag.FlagSet, names string) func() string {
			var v int
			flagx.IntVar(fs, &v, names, 4, "an option")
			return func() string { return fmt.Sprint(v) }
		}, "4", "12", "12"},
		{"duration", func(fs *flag.FlagSet, names string) func() string {
			var v time.Duration
			flagx.DurationVar(fs, &v, names, time.Second, "an option")
			return func() string { return v.String() }
		}, "1s", "2m", "2m0s"},
		{"int64", func(fs *flag.FlagSet, names string) func() string {
			var v int64
			flagx.Int64Var(fs, &v, names, -4, "an option")
			return func() string { return fmt.Sprint(v) }
		}, "-4", "-1099511627776", "-1099511627776"},
		{"uint", func(fs *flag.FlagSet, names string) func() string {
			var v uint
			flagx.UintVar(fs, &v, names, 4, "an option")
			return func() string { return fmt.Sprint(v) }
		}, "4", "0xff", "255"},
		{"uint64", func(fs *flag.FlagSet, names string) func() string {
			var v uint64
			flagx.Uint64Var(fs, &v, names, 4, "an option")
			return func() string { return fmt.Sprint(v) }
		}, "4", "18446744073709551615", "18446744073709551615"},
		{"float64", func(fs *flag.FlagSet, names string) func() string {
			var v float64
			flagx.Float64Var(fs, &v, names, 1.5, "an option")
			return func() string { return fmt.Sprint(v) }
		}, "1.5", "-2.5e2", "-250"},
		{"text", func(fs *flag.FlagSet, names string) func() string {
			var v netip.Addr
			flagx.TextVar(fs, &v, names, netip.MustParseAddr("127.0.0.1"), "an option")
			return func() string { return v.String() }
		}, "127.0.0.1", "::1", "::1"},
		{"String", func(fs *flag.FlagSet, names string) func() string {
			v := flagx.String(fs, names, "initial", "an option")
			return func() string { return *v }
		}, "initial", "changed", "changed"},
		{"Bool", func(fs *flag.FlagSet, names string) func() string {
			v := flagx.Bool(fs, names, false, "an option")
			return func() string { return fmt.Sprint(*v) }
		}, "false", "true", "true"},
		{"Int", func(fs *flag.FlagSet, names string) func() string {
			v := flagx.Int(fs, names, 4, "an option")
			return func() string { return fmt.Sprint(*v) }
		}, "4", "12", "12"},
		{"Int64", func(fs *flag.FlagSet, names string) func() string {
			v := flagx.Int64(fs, names, -4, "an option")
			return func() string { return fmt.Sprint(*v) }
		}, "-4", "-1099511627776", "-1099511627776"},
		{"Uint", func(fs *flag.FlagSet, names string) func() string {
			v := flagx.Uint(fs, names, 4, "an option")
			return func() string { return fmt.Sprint(*v) }
		}, "4", "0xff", "255"},
		{"Uint64", func(fs *flag.FlagSet, names string) func() string {
			v := flagx.Uint64(fs, names, 4, "an option")
			return func() string { return fmt.Sprint(*v) }
		}, "4", "18446744073709551615", "18446744073709551615"},
		{"Float64", func(fs *flag.FlagSet, names string) func() string {
			v := flagx.Float64(fs, names, 1.5, "an option")
			return func() string { return fmt.Sprint(*v) }
		}, "1.5", "-2.5e2", "-250"},
		{"Duration", func(fs *flag.FlagSet, names string) func() string {
			v := flagx.Duration(fs, names, time.Second, "an option")
			return func() string { return v.String() }
		}, "1s", "2m", "2m0s"},
		{"value", func(fs *flag.FlagSet, names string) func() string {
			v := listValue{"initial"}
			flagx.Var(fs, &v, names, "an option")
			return v.String
		}, "initial", "changed", "initial,changed"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			for _, name := range []string{"-option", "--option", "-o", "--o", "--other"} {
				t.Run(name, func(t *testing.T) {
					fs := flag.NewFlagSet("test", flag.ContinueOnError)
					fs.SetOutput(io.Discard)
					value := tt.register(fs, "o|other|option")
					if got := value(); got != tt.initial {
						t.Fatalf("initial value = %q, want %q", got, tt.initial)
					}
					if err := fs.Parse([]string{name + "=" + tt.input}); err != nil {
						t.Fatal(err)
					}
					if got := value(); got != tt.want {
						t.Errorf("value = %q, want %q", got, tt.want)
					}
					for _, alias := range []string{"option", "o", "other"} {
						if got := fs.Lookup(alias).DefValue; got != tt.initial {
							t.Errorf("%s default = %q, want %q", alias, got, tt.initial)
						}
					}
					if fs.Lookup("o").Value.String() != fs.Lookup("option").Value.String() ||
						fs.Lookup("other").Value.String() != fs.Lookup("option").Value.String() {
						t.Error("aliases do not share the parsed value")
					}
				})
			}
			t.Run("single name", func(t *testing.T) {
				fs := flag.NewFlagSet("test", flag.ContinueOnError)
				tt.register(fs, "option")
				var names []string
				fs.VisitAll(func(f *flag.Flag) { names = append(names, f.Name) })
				if len(names) != 1 || names[0] != "option" {
					t.Errorf("registered names = %v, want [option]", names)
				}
			})
		})
	}
}

func TestRegistrationParseErrors(t *testing.T) {
	tests := []struct {
		name     string
		register func(*flag.FlagSet)
		input    string
	}{
		{"Float64", func(fs *flag.FlagSet) { flagx.Float64(fs, "option|o", 1.5, "an option") }, "invalid"},
		{"Int64", func(fs *flag.FlagSet) { flagx.Int64(fs, "option|o", 4, "an option") }, "9223372036854775808"},
		{"Uint", func(fs *flag.FlagSet) { flagx.Uint(fs, "option|o", 4, "an option") }, "-1"},
		{"Uint64", func(fs *flag.FlagSet) { flagx.Uint64(fs, "option|o", 4, "an option") }, "18446744073709551616"},
		{"TextVar", func(fs *flag.FlagSet) {
			var v netip.Addr
			flagx.TextVar(fs, &v, "option|o", netip.MustParseAddr("127.0.0.1"), "an option")
		}, "invalid"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fs := flag.NewFlagSet("test", flag.ContinueOnError)
			fs.SetOutput(io.Discard)
			tt.register(fs)
			for _, name := range []string{"option", "o"} {
				if err := fs.Parse([]string{"--" + name + "=" + tt.input}); err == nil {
					t.Errorf("Parse(%s=%s) succeeded, want error", name, tt.input)
				}
			}
		})
	}
}

func TestCallbackRegistration(t *testing.T) {
	for _, tt := range []struct {
		name     string
		register func(*flag.FlagSet, string, string, func(string) error)
		args     []string
		want     []string
		help     string
	}{
		{"Func", flagx.Func, []string{"--option=first", "-o", "second", "--option=third"}, []string{"first", "second", "third"}, "  -option, -o value  callback\n"},
		{"BoolFunc", flagx.BoolFunc, []string{"--option", "-o=false", "--option=true"}, []string{"true", "false", "true"}, "  -option, -o  callback\n"},
	} {
		t.Run(tt.name, func(t *testing.T) {
			fs := flag.NewFlagSet("test", flag.ContinueOnError)
			var out bytes.Buffer
			fs.SetOutput(&out)
			var got []string
			rejected := errors.New("rejected value")
			tt.register(fs, "option|o", "callback", func(value string) error {
				if value == "reject" {
					return rejected
				}
				got = append(got, value)
				return nil
			})
			if len(got) != 0 {
				t.Fatal("callback invoked during registration")
			}
			if err := fs.Parse(tt.args); err != nil {
				t.Fatal(err)
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("callback values = %v, want %v", got, tt.want)
			}
			flagx.PrintDefaults(fs)
			if !strings.Contains(out.String(), tt.help) || strings.Count(out.String(), "callback") != 1 {
				t.Errorf("callback aliases not grouped correctly:\n%s", out.String())
			}
			for _, name := range []string{"option", "o"} {
				err := fs.Parse([]string{"--" + name + "=reject", "--option=unreachable"})
				if err == nil || !strings.Contains(err.Error(), rejected.Error()) {
					t.Errorf("Parse(%s) = %v, want callback error", name, err)
				}
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("parsing continued after callback error: %v", got)
			}
		})
	}
}

func TestBoolWithoutArgument(t *testing.T) {
	fs := flag.NewFlagSet("test", flag.ContinueOnError)
	var verbose bool
	flagx.BoolVar(fs, &verbose, "verbose|v", false, "verbose output")
	if err := fs.Parse([]string{"-v", "--verbose=false", "--verbose"}); err != nil {
		t.Fatal(err)
	}
	if !verbose {
		t.Fatal("boolean flag did not become true")
	}
}

func TestInvalidNames(t *testing.T) {
	tests := []struct {
		names string
		want  string
	}{
		{"", `flagx: empty flag name in ""`},
		{"|o", `flagx: empty flag name in "|o"`},
		{"option|", `flagx: empty flag name in "option|"`},
		{"option||o", `flagx: empty flag name in "option||o"`},
		{"a\n|", `flagx: empty flag name in "a\n|"`},
		{"option|option", ""},
		{"-option", ""},
		{"option=value", ""},
	}
	for _, tt := range tests {
		t.Run(tt.names, func(t *testing.T) {
			fs := flag.NewFlagSet("test", flag.ContinueOnError)
			fs.SetOutput(io.Discard)
			defer func() {
				got := recover()
				if got == nil {
					t.Errorf("registration of %q did not panic", tt.names)
				} else if tt.want != "" && fmt.Sprint(got) != tt.want {
					t.Errorf("panic = %q, want %q", got, tt.want)
				}
			}()
			var value string
			flagx.StringVar(fs, &value, tt.names, "", "an option")
		})
	}
}
