package flagx_test

import (
	"bytes"
	"errors"
	"flag"
	"fmt"
	"io"
	"net/netip"
	"strings"
	"testing"
	"time"

	"github.com/Songmu/flagx"
)

func TestFlagSetRegistration(t *testing.T) {
	negativeIntValues := registrationData.intValue
	negativeIntValues.input, negativeIntValues.want = "-12", "-12"
	positiveDefaultInt64Values := registrationData.int64Value
	positiveDefaultInt64Values.initial = "4"
	tests := []struct {
		name     string
		register func(*flagx.FlagSet, string) func() string
		values   registrationValues
	}{
		{"String", func(fs *flagx.FlagSet, names string) func() string {
			v := fs.String(names, "initial", "an option")
			return func() string { return *v }
		}, registrationData.stringValue},
		{"StringVar", func(fs *flagx.FlagSet, names string) func() string {
			var v string
			fs.StringVar(&v, names, "initial", "an option")
			return func() string { return v }
		}, registrationData.stringValue},
		{"Bool", func(fs *flagx.FlagSet, names string) func() string {
			v := fs.Bool(names, false, "an option")
			return func() string { return fmt.Sprint(*v) }
		}, registrationData.boolValue},
		{"BoolVar", func(fs *flagx.FlagSet, names string) func() string {
			var v bool
			fs.BoolVar(&v, names, false, "an option")
			return func() string { return fmt.Sprint(v) }
		}, registrationData.boolValue},
		{"Int", func(fs *flagx.FlagSet, names string) func() string {
			v := fs.Int(names, 4, "an option")
			return func() string { return fmt.Sprint(*v) }
		}, negativeIntValues},
		{"IntVar", func(fs *flagx.FlagSet, names string) func() string {
			var v int
			fs.IntVar(&v, names, 4, "an option")
			return func() string { return fmt.Sprint(v) }
		}, negativeIntValues},
		{"Int64", func(fs *flagx.FlagSet, names string) func() string {
			v := fs.Int64(names, 4, "an option")
			return func() string { return fmt.Sprint(*v) }
		}, positiveDefaultInt64Values},
		{"Int64Var", func(fs *flagx.FlagSet, names string) func() string {
			var v int64
			fs.Int64Var(&v, names, 4, "an option")
			return func() string { return fmt.Sprint(v) }
		}, positiveDefaultInt64Values},
		{"Uint", func(fs *flagx.FlagSet, names string) func() string {
			v := fs.Uint(names, 4, "an option")
			return func() string { return fmt.Sprint(*v) }
		}, registrationData.uintValue},
		{"UintVar", func(fs *flagx.FlagSet, names string) func() string {
			var v uint
			fs.UintVar(&v, names, 4, "an option")
			return func() string { return fmt.Sprint(v) }
		}, registrationData.uintValue},
		{"Uint64", func(fs *flagx.FlagSet, names string) func() string {
			v := fs.Uint64(names, 4, "an option")
			return func() string { return fmt.Sprint(*v) }
		}, registrationData.uint64Value},
		{"Uint64Var", func(fs *flagx.FlagSet, names string) func() string {
			var v uint64
			fs.Uint64Var(&v, names, 4, "an option")
			return func() string { return fmt.Sprint(v) }
		}, registrationData.uint64Value},
		{"Float64", func(fs *flagx.FlagSet, names string) func() string {
			v := fs.Float64(names, 1.5, "an option")
			return func() string { return fmt.Sprint(*v) }
		}, registrationData.float64Value},
		{"Float64Var", func(fs *flagx.FlagSet, names string) func() string {
			var v float64
			fs.Float64Var(&v, names, 1.5, "an option")
			return func() string { return fmt.Sprint(v) }
		}, registrationData.float64Value},
		{"Duration", func(fs *flagx.FlagSet, names string) func() string {
			v := fs.Duration(names, time.Second, "an option")
			return func() string { return v.String() }
		}, registrationData.durationValue},
		{"DurationVar", func(fs *flagx.FlagSet, names string) func() string {
			var v time.Duration
			fs.DurationVar(&v, names, time.Second, "an option")
			return func() string { return v.String() }
		}, registrationData.durationValue},
		{"TextVar", func(fs *flagx.FlagSet, names string) func() string {
			var v netip.Addr
			fs.TextVar(&v, names, netip.MustParseAddr("127.0.0.1"), "an option")
			return func() string { return v.String() }
		}, registrationData.textValue},
		{"Var", func(fs *flagx.FlagSet, names string) func() string {
			v := otherString("initial")
			fs.Var(&v, names, "an option")
			return v.String
		}, registrationValues{"initial", "changed", "changed"}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			testRegistration(t, tt.values, func(names string) (registrationFlagSet, func() string) {
				fs := flagx.NewFlagSet("test", flag.ContinueOnError)
				fs.SetOutput(io.Discard)
				return fs, tt.register(fs, names)
			})
		})
	}
}

func TestFlagSetCallbackRegistration(t *testing.T) {
	for _, tt := range []struct {
		name     string
		register func(*flagx.FlagSet) func() string
		input    string
		want     string
	}{
		{"Func", func(fs *flagx.FlagSet) func() string {
			v := "initial"
			fs.Func("option|o", "an option", func(s string) error { v = s; return nil })
			return func() string { return v }
		}, "changed", "changed"},
		{"BoolFunc", func(fs *flagx.FlagSet) func() string {
			v := "initial"
			fs.BoolFunc("option|o", "an option", func(s string) error { v = s; return nil })
			return func() string { return v }
		}, "false", "false"},
	} {
		t.Run(tt.name, func(t *testing.T) {
			for _, name := range []string{"--option", "-o"} {
				fs := flagx.NewFlagSet("test", flag.ContinueOnError)
				value := tt.register(fs)
				if got := value(); got != "initial" {
					t.Fatalf("initial value = %q, want %q", got, "initial")
				}
				if fs.Lookup("option") == nil || fs.Lookup("o") == nil || fs.Lookup("option|o") != nil {
					t.Fatal("names were not registered as aliases")
				}
				if err := fs.Parse([]string{name + "=" + tt.input}); err != nil {
					t.Fatal(err)
				}
				if got := value(); got != tt.want {
					t.Errorf("value after %s = %q, want %q", name, got, tt.want)
				}
			}
		})
	}
}

func TestFlagSetUsage(t *testing.T) {
	for _, tt := range []struct {
		name     string
		new      func() *flagx.FlagSet
		embedded bool
	}{
		{"constructor", func() *flagx.FlagSet { return flagx.NewFlagSet("test", flag.ContinueOnError) }, true},
		{"zero", func() *flagx.FlagSet { return new(flagx.FlagSet) }, false},
		{"initialized", func() *flagx.FlagSet {
			fs := new(flagx.FlagSet)
			fs.Init("test", flag.ContinueOnError)
			return fs
		}, true},
		{"renamed", func() *flagx.FlagSet {
			fs := flagx.NewFlagSet("old", flag.ContinueOnError)
			fs.Init("test", flag.ContinueOnError)
			return fs
		}, true},
		{"nil usage", func() *flagx.FlagSet {
			fs := flagx.NewFlagSet("test", flag.ContinueOnError)
			fs.Usage = nil
			return fs
		}, false},
		{"nil usage initialized", func() *flagx.FlagSet {
			fs := flagx.NewFlagSet("test", flag.ContinueOnError)
			fs.Usage = nil
			fs.Init(fs.Name(), fs.ErrorHandling())
			return fs
		}, true},
	} {
		t.Run(tt.name, func(t *testing.T) {
			fs := tt.new()
			if got := fs.Usage != nil; got != tt.embedded {
				t.Fatalf("Usage installed before parsing = %t, want %t", got, tt.embedded)
			}
			parse := fs.Parse
			if tt.embedded {
				parse = fs.FlagSet.Parse
			}
			fs.String("output|o", "", "write to `file`")
			var original, redirected bytes.Buffer
			fs.SetOutput(&original)
			const defaults = "  -output, -o file  write to file\n  -help, -h         Show this help message\n"
			want := "Usage: " + fs.Name() + "\n\nOptions:\n" + defaults
			for _, name := range []string{"-h", "--help"} {
				original.Reset()
				if err := parse([]string{name}); !errors.Is(err, flag.ErrHelp) {
					t.Fatalf("Parse(%s) = %v, want flag.ErrHelp", name, err)
				}
				assertUsageOutput(t, "help "+name, original.String(), want)
			}

			original.Reset()
			fs.SetOutput(&redirected)
			fs.Usage()
			assertUsageOutput(t, "Usage after SetOutput", redirected.String(), want)
			assertUsageOutput(t, "original output after SetOutput", original.String(), "")

			redirected.Reset()
			fs.PrintDefaults()
			assertUsageOutput(t, "PrintDefaults", redirected.String(), defaults)

			redirected.Reset()
			if err := parse([]string{"--unknown"}); err == nil {
				t.Fatal("unknown flag accepted")
			}
			assertUsageOutput(t, "unknown flag error and help", redirected.String(), "flag provided but not defined: -unknown\n"+want)
		})
	}
}

func TestFlagSetCustomUsage(t *testing.T) {
	fs := flagx.NewFlagSet("old", flag.ContinueOnError)
	var out, explicit bytes.Buffer
	fs.SetOutput(&out)
	fs.String("output|o", "", "write to `file`")
	fs.Usage = func() {
		fs.PrintUsage(&explicit, "test [options] <input>", "Process input files.")
	}
	fs.Init("test", flag.ContinueOnError)
	if err := fs.Parse([]string{"--help"}); !errors.Is(err, flag.ErrHelp) {
		t.Fatalf("Parse = %v, want flag.ErrHelp", err)
	}
	want := "Usage: test [options] <input>\n\nProcess input files.\n\nOptions:\n" +
		"  -output, -o file  write to file\n  -help, -h         Show this help message\n"
	if got := explicit.String(); got != want {
		t.Errorf("custom Usage = %q, want %q", got, want)
	}
	if fs.Output() != &out || out.Len() != 0 {
		t.Fatal("PrintUsage changed or wrote to fs.Output")
	}
}

func TestFlagSetInterop(t *testing.T) {
	fs := flagx.NewFlagSet("test", flag.PanicOnError)
	verbose := fs.Bool("verbose|v", false, "verbose output")
	jobs := flagx.Int(&fs.FlagSet, "jobs|j", 4, "workers")
	raw := fs.FlagSet.String("raw", "", "standard flag")
	if err := fs.Parse([]string{"-v", "-j=8", "--raw=value", "input", "--verbose=false"}); err != nil {
		t.Fatal(err)
	}
	if !*verbose || *jobs != 8 || *raw != "value" {
		t.Fatalf("mixed registrations parsed as %v, %d, %q", *verbose, *jobs, *raw)
	}
	if got := fs.Args(); len(got) != 2 || got[0] != "input" || got[1] != "--verbose=false" {
		t.Errorf("parsing did not stop at positional argument: %v", got)
	}
	var out bytes.Buffer
	fs.SetOutput(&out)
	fs.PrintDefaults()
	help := out.String()
	if a, b, c := strings.Index(help, "-verbose, -v"), strings.Index(help, "-jobs, -j"), strings.Index(help, "-raw"); a < 0 || b <= a || c <= b {
		t.Errorf("mixed registration order not preserved:\n%s", help)
	}
	defer func() {
		if recovered := recover(); recovered == nil {
			t.Error("PanicOnError did not panic")
		}
	}()
	fs.Parse([]string{"--unknown"})
}
