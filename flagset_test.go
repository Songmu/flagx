package flagx_test

import (
	"bytes"
	"errors"
	"flag"
	"fmt"
	"net/netip"
	"strings"
	"testing"
	"time"

	"github.com/Songmu/flagx"
)

func TestFlagSetRegistration(t *testing.T) {
	tests := []struct {
		name     string
		register func(*flagx.FlagSet) func() string
		initial  string
		input    string
		want     string
	}{
		{"String", func(fs *flagx.FlagSet) func() string {
			v := fs.String("option|o", "initial", "an option")
			return func() string { return *v }
		}, "initial", "changed", "changed"},
		{"StringVar", func(fs *flagx.FlagSet) func() string {
			var v string
			fs.StringVar(&v, "option|o", "initial", "an option")
			return func() string { return v }
		}, "initial", "changed", "changed"},
		{"Bool", func(fs *flagx.FlagSet) func() string {
			v := fs.Bool("option|o", false, "an option")
			return func() string { return fmt.Sprint(*v) }
		}, "false", "true", "true"},
		{"BoolVar", func(fs *flagx.FlagSet) func() string {
			var v bool
			fs.BoolVar(&v, "option|o", false, "an option")
			return func() string { return fmt.Sprint(v) }
		}, "false", "true", "true"},
		{"Int", func(fs *flagx.FlagSet) func() string {
			v := fs.Int("option|o", 4, "an option")
			return func() string { return fmt.Sprint(*v) }
		}, "4", "-12", "-12"},
		{"IntVar", func(fs *flagx.FlagSet) func() string {
			var v int
			fs.IntVar(&v, "option|o", 4, "an option")
			return func() string { return fmt.Sprint(v) }
		}, "4", "-12", "-12"},
		{"Int64", func(fs *flagx.FlagSet) func() string {
			v := fs.Int64("option|o", 4, "an option")
			return func() string { return fmt.Sprint(*v) }
		}, "4", "-1099511627776", "-1099511627776"},
		{"Int64Var", func(fs *flagx.FlagSet) func() string {
			var v int64
			fs.Int64Var(&v, "option|o", 4, "an option")
			return func() string { return fmt.Sprint(v) }
		}, "4", "-1099511627776", "-1099511627776"},
		{"Uint", func(fs *flagx.FlagSet) func() string {
			v := fs.Uint("option|o", 4, "an option")
			return func() string { return fmt.Sprint(*v) }
		}, "4", "0xff", "255"},
		{"UintVar", func(fs *flagx.FlagSet) func() string {
			var v uint
			fs.UintVar(&v, "option|o", 4, "an option")
			return func() string { return fmt.Sprint(v) }
		}, "4", "0xff", "255"},
		{"Uint64", func(fs *flagx.FlagSet) func() string {
			v := fs.Uint64("option|o", 4, "an option")
			return func() string { return fmt.Sprint(*v) }
		}, "4", "18446744073709551615", "18446744073709551615"},
		{"Uint64Var", func(fs *flagx.FlagSet) func() string {
			var v uint64
			fs.Uint64Var(&v, "option|o", 4, "an option")
			return func() string { return fmt.Sprint(v) }
		}, "4", "18446744073709551615", "18446744073709551615"},
		{"Float64", func(fs *flagx.FlagSet) func() string {
			v := fs.Float64("option|o", 1.5, "an option")
			return func() string { return fmt.Sprint(*v) }
		}, "1.5", "-2.5e2", "-250"},
		{"Float64Var", func(fs *flagx.FlagSet) func() string {
			var v float64
			fs.Float64Var(&v, "option|o", 1.5, "an option")
			return func() string { return fmt.Sprint(v) }
		}, "1.5", "-2.5e2", "-250"},
		{"Duration", func(fs *flagx.FlagSet) func() string {
			v := fs.Duration("option|o", time.Second, "an option")
			return func() string { return v.String() }
		}, "1s", "2m", "2m0s"},
		{"DurationVar", func(fs *flagx.FlagSet) func() string {
			var v time.Duration
			fs.DurationVar(&v, "option|o", time.Second, "an option")
			return func() string { return v.String() }
		}, "1s", "2m", "2m0s"},
		{"TextVar", func(fs *flagx.FlagSet) func() string {
			var v netip.Addr
			fs.TextVar(&v, "option|o", netip.MustParseAddr("127.0.0.1"), "an option")
			return func() string { return v.String() }
		}, "127.0.0.1", "::1", "::1"},
		{"Var", func(fs *flagx.FlagSet) func() string {
			v := otherString("initial")
			fs.Var(&v, "option|o", "an option")
			return v.String
		}, "initial", "changed", "changed"},
		{"Func", func(fs *flagx.FlagSet) func() string {
			v := "initial"
			fs.Func("option|o", "an option", func(s string) error { v = s; return nil })
			return func() string { return v }
		}, "initial", "changed", "changed"},
		{"BoolFunc", func(fs *flagx.FlagSet) func() string {
			v := "initial"
			fs.BoolFunc("option|o", "an option", func(s string) error { v = s; return nil })
			return func() string { return v }
		}, "initial", "false", "false"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			for _, name := range []string{"--option", "-o"} {
				fs := flagx.NewFlagSet("test", flag.ContinueOnError)
				value := tt.register(fs)
				if got := value(); got != tt.initial {
					t.Fatalf("initial value = %q, want %q", got, tt.initial)
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
		name string
		new  func() *flagx.FlagSet
	}{
		{"constructor", func() *flagx.FlagSet { return flagx.NewFlagSet("test", flag.ContinueOnError) }},
		{"zero", func() *flagx.FlagSet { return new(flagx.FlagSet) }},
		{"initialized", func() *flagx.FlagSet {
			fs := new(flagx.FlagSet)
			fs.Init("test", flag.ContinueOnError)
			return fs
		}},
		{"renamed", func() *flagx.FlagSet {
			fs := flagx.NewFlagSet("old", flag.ContinueOnError)
			fs.Init("test", flag.ContinueOnError)
			return fs
		}},
		{"nil usage", func() *flagx.FlagSet {
			fs := flagx.NewFlagSet("test", flag.ContinueOnError)
			fs.Usage = nil
			return fs
		}},
	} {
		t.Run(tt.name, func(t *testing.T) {
			fs := tt.new()
			fs.String("output|o", "", "write to `file`")
			var original, redirected bytes.Buffer
			fs.SetOutput(&original)
			const defaults = "  -output, -o file  write to file\n  -help, -h         Show this help message\n"
			want := "Usage: " + fs.Name() + "\n\nOptions:\n" + defaults
			for _, name := range []string{"-h", "--help"} {
				original.Reset()
				if err := fs.Parse([]string{name}); !errors.Is(err, flag.ErrHelp) {
					t.Fatalf("Parse(%s) = %v, want flag.ErrHelp", name, err)
				}
				if got := original.String(); got != want {
					t.Errorf("help = %q, want %q", got, want)
				}
			}
			original.Reset()
			fs.SetOutput(&redirected)
			fs.Usage()
			if got := redirected.String(); got != want || original.Len() != 0 {
				t.Errorf("Usage did not follow SetOutput: %q", got)
			}
			redirected.Reset()
			fs.PrintDefaults()
			if got := redirected.String(); got != defaults {
				t.Errorf("PrintDefaults = %q, want %q", got, defaults)
			}
			redirected.Reset()
			if err := fs.Parse([]string{"--unknown"}); err == nil {
				t.Fatal("unknown flag accepted")
			}
			if got := redirected.String(); got != "flag provided but not defined: -unknown\n"+want {
				t.Errorf("error help = %q", got)
			}
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
