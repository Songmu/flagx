package flagx_test

import (
	"bytes"
	"errors"
	"flag"
	"strings"
	"testing"
	"time"

	"github.com/Songmu/flagx"
)

func sampleFlags() *flag.FlagSet {
	fs := flag.NewFlagSet("myapp", flag.ContinueOnError)
	var output string
	var verbose bool
	var jobs int
	var timeout time.Duration
	flagx.StringVar(fs, &output, "output|o", "", "write to `file`")
	flagx.BoolVar(fs, &verbose, "verbose|v", false, "verbose output")
	flagx.IntVar(fs, &jobs, "jobs|j", 4, "number of `workers`")
	flagx.DurationVar(fs, &timeout, "timeout|t", 30*time.Second, "timeout after `duration`")
	return fs
}

const sampleDefaults = `  -output, -o file       write to file
  -verbose, -v           verbose output
  -jobs, -j workers      number of workers (default 4)
  -timeout, -t duration  timeout after duration (default 30s)
  -help, -h              Show this help message
`

func assertUsageOutput(t *testing.T, label, got, want string) {
	t.Helper()
	if got != want {
		t.Errorf("%s:\n%s\nwant:\n%s\ngot (quoted): %q\nwant (quoted): %q", label, got, want, got, want)
	}
}

func TestUsage(t *testing.T) {
	fs := sampleFlags()
	var original, redirected, explicit bytes.Buffer
	fs.SetOutput(&original)
	fs.Usage = flagx.Usage(fs, "myapp [options] <input>", "Process input files.")
	want := "Usage: myapp [options] <input>\n\nProcess input files.\n\nOptions:\n" + sampleDefaults
	for _, name := range []string{"-h", "--help"} {
		original.Reset()
		if err := fs.Parse([]string{name}); !errors.Is(err, flag.ErrHelp) {
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
	flagx.PrintDefaults(fs)
	assertUsageOutput(t, "PrintDefaults", redirected.String(), sampleDefaults)

	redirected.Reset()
	flagx.PrintUsage(&explicit, fs, "myapp", "")
	assertUsageOutput(t, "PrintUsage without description", explicit.String(), "Usage: myapp\n\nOptions:\n"+sampleDefaults)
	assertUsageOutput(t, "fs.Output after PrintUsage", redirected.String(), "")
	if fs.Output() != &redirected {
		t.Fatal("PrintUsage changed fs.Output")
	}
}

func TestDefaults(t *testing.T) {
	for _, tt := range []struct {
		name      string
		register  func(*flag.FlagSet)
		set       string
		wantLabel string
		wantUsage string
	}{
		{"string", func(fs *flag.FlagSet) {
			fs.String("string", "hello\nworld", "a string")
		}, "", "-string string", `a string (default "hello\nworld")`},
		{"bool", func(fs *flag.FlagSet) {
			fs.Bool("bool", true, "use `mode`")
		}, "", "-bool", "use mode (default true)"},
		{"int", func(fs *flag.FlagSet) {
			fs.Int("int", 9, "an integer")
		}, "100", "-int int", "an integer (default 9)"},
		{"duration", func(fs *flag.FlagSet) {
			fs.Duration("duration", time.Minute, "a duration")
		}, "", "-duration duration", "a duration (default 1m0s)"},
		{"zero-string", func(fs *flag.FlagSet) {
			fs.String("zero-string", "", "empty")
		}, "", "-zero-string string", "empty"},
		{"zero-bool", func(fs *flag.FlagSet) {
			fs.Bool("zero-bool", false, "off")
		}, "", "-zero-bool", "off"},
		{"zero-int", func(fs *flag.FlagSet) {
			fs.Int("zero-int", 0, "zero")
		}, "", "-zero-int int", "zero"},
		{"zero-duration", func(fs *flag.FlagSet) {
			fs.Duration("zero-duration", 0, "immediate")
		}, "", "-zero-duration duration", "immediate"},
	} {
		t.Run(tt.name, func(t *testing.T) {
			fs := flag.NewFlagSet("test", flag.ContinueOnError)
			tt.register(fs)
			if tt.set != "" {
				if err := fs.Set(tt.name, tt.set); err != nil {
					t.Fatal(err)
				}
			}
			var out bytes.Buffer
			fs.SetOutput(&out)
			flagx.PrintDefaults(fs)
			// Check option semantics independently of column padding and built-in help.
			index := strings.Index(out.String(), "  -"+tt.name+" ")
			if index < 0 {
				t.Fatalf("option %q missing:\n%s", tt.name, out.String())
			}
			line, _, _ := strings.Cut(out.String()[index:], "\n")
			label, usage, _ := strings.Cut(strings.TrimSpace(line), "  ")
			assertUsageOutput(t, "option label", label, tt.wantLabel)
			assertUsageOutput(t, "option usage", strings.TrimSpace(usage), tt.wantUsage)
		})
	}
}

type mapValue map[string]string

func (v mapValue) String() string     { return v["value"] }
func (v mapValue) Set(s string) error { v["value"] = s; return nil }

type fragileValue struct{ value *string }

func (v *fragileValue) String() string     { return *v.value }
func (v *fragileValue) Set(s string) error { *v.value = s; return nil }

type otherString string

func (v *otherString) String() string     { return string(*v) }
func (v *otherString) Set(s string) error { *v = otherString(s); return nil }

func TestAliasIdentity(t *testing.T) {
	fs := flag.NewFlagSet("test", flag.ContinueOnError)
	var shared, separate string
	fs.StringVar(&shared, "zulu", "", "shared")
	fs.StringVar(&shared, "z", "", "shared")
	fs.StringVar(&shared, "alpha", "", "shared")
	fs.StringVar(&shared, "a", "", "shared")
	fs.StringVar(&separate, "separate", "", "separate")
	// A different Value type at the same address must not be merged.
	fs.Var((*otherString)(&shared), "other", "other type")
	custom := listValue{"first"}
	flagx.Var(fs, &custom, "custom|c", "custom list")
	nonComparable := mapValue{"value": "entry"}
	fs.Var(nonComparable, "map-one", "map one")
	fs.Var(nonComparable, "map-two", "map two")
	initial := "fragile"
	fs.Var(&fragileValue{&initial}, "fragile", "fragile value")

	var out bytes.Buffer
	fs.SetOutput(&out)
	flagx.PrintDefaults(fs)
	for _, want := range []string{
		"-alpha, -zulu, -a, -z string", "-separate string", "-other value",
		"-custom, -c value", "custom list (default first)",
		"-map-one value", "-map-two value", "fragile value (default fragile)",
	} {
		if !strings.Contains(out.String(), want) {
			t.Errorf("missing %q:\n%s", want, out.String())
		}
	}
	if strings.Count(out.String(), "  -") != 8 {
		t.Errorf("unexpected alias grouping:\n%s", out.String())
	}
}

func TestExistingHelp(t *testing.T) {
	for _, names := range [][]string{nil, {"h"}, {"help"}, {"help", "h"}} {
		t.Run(strings.Join(names, ","), func(t *testing.T) {
			fs := flag.NewFlagSet("test", flag.ContinueOnError)
			for _, name := range names {
				fs.Bool(name, false, "custom "+name)
			}
			var out bytes.Buffer
			fs.SetOutput(&out)
			flagx.PrintDefaults(fs)
			if strings.Count(out.String(), "-help") != 1 {
				t.Errorf("help missing or duplicated:\n%s", out.String())
			}
			withoutHelp := strings.ReplaceAll(out.String(), "-help", "")
			if strings.Count(withoutHelp, "-h") != 1 {
				t.Errorf("h missing or duplicated:\n%s", out.String())
			}
			for _, name := range names {
				if !strings.Contains(out.String(), "custom "+name) {
					t.Errorf("custom help was replaced:\n%s", out.String())
				}
			}
			var registered int
			fs.VisitAll(func(*flag.Flag) { registered++ })
			if registered != len(names) {
				t.Fatal("printing help registered flags")
			}
		})
	}
}

func TestLongAndMultilineUsage(t *testing.T) {
	fs := flag.NewFlagSet("test", flag.ContinueOnError)
	fs.String("an-option-name-that-is-too-long", "", "first line\nsecond line")
	fs.Bool("v", false, "verbose")
	var out bytes.Buffer
	fs.SetOutput(&out)
	flagx.PrintDefaults(fs)
	indent := strings.Repeat(" ", 32)
	if want := "  -an-option-name-that-is-too-long string\n" + indent + "first line\n" + indent + "second line\n"; !strings.Contains(out.String(), want) {
		t.Errorf("long/multiline layout:\n%s", out.String())
	}
	if want := "  -v" + strings.Repeat(" ", 28) + "verbose\n"; !strings.Contains(out.String(), want) {
		t.Errorf("short-only layout:\n%s", out.String())
	}
}

func TestSpecifiedNameOrder(t *testing.T) {
	for _, names := range []string{"o|output", "bb|aa", "z|a", "output|out|o"} {
		t.Run(names, func(t *testing.T) {
			t.Parallel()
			fs := flag.NewFlagSet("test", flag.ContinueOnError)
			var value string
			flagx.StringVar(fs, &value, names, "", "an option")
			parts := strings.Split(names, "|")
			want := "-" + strings.Join(parts, ", -") + " string"
			var out bytes.Buffer
			fs.SetOutput(&out)
			flagx.PrintDefaults(fs)
			if !strings.Contains(out.String(), want) {
				t.Errorf("missing specified order %q:\n%s", want, out.String())
			}
			for _, name := range parts {
				f := fs.Lookup(name)
				if placeholder, _ := flag.UnquoteUsage(f); placeholder != "string" {
					t.Errorf("standard placeholder = %q, want string", placeholder)
				}
				if _, ok := f.Value.(flag.Getter); !ok {
					t.Errorf("%s lost flag.Getter", name)
				}
			}
		})
	}
}

func TestDeclaredAndInferredAliases(t *testing.T) {
	fs := flag.NewFlagSet("test", flag.ContinueOnError)
	var value string
	fs.StringVar(&value, "direct", "", "direct flag")
	flagx.StringVar(fs, &value, "z|zulu|alpha", "", "declared aliases")
	flagx.StringVar(fs, &value, "extra|e", "", "more aliases")
	custom := mapValue{"value": "entry"}
	flagx.Var(fs, custom, "named-map|m", "named map")
	fs.Bool("middle", false, "middle option")
	var out bytes.Buffer
	fs.SetOutput(&out)
	flagx.PrintDefaults(fs)
	want := "-z, -zulu, -alpha, -extra, -e, -direct string"
	if !strings.HasPrefix(out.String(), "  "+want) {
		t.Errorf("option sorting or alias order:\n%s", out.String())
	}
	if !strings.Contains(out.String(), "-named-map, -m value") {
		t.Errorf("non-pointer declared aliases were not combined:\n%s", out.String())
	}
	if strings.Count(out.String(), "named map (default entry)") != 1 {
		t.Errorf("non-pointer aliases duplicated:\n%s", out.String())
	}
}

func TestOptionRegistrationOrder(t *testing.T) {
	fs := flag.NewFlagSet("test", flag.ContinueOnError)
	fs.Bool("beta", false, "direct option registered first")
	var first bool
	fs.BoolVar(&first, "aardvark", false, "direct alias")
	flagx.BoolVar(fs, &first, "zulu|z", false, "first registered option")
	fs.Bool("middle", false, "direct option")
	var second bool
	flagx.BoolVar(fs, &second, "alpha|a", false, "second registered option")
	var out bytes.Buffer
	fs.SetOutput(&out)
	flagx.PrintDefaults(fs)
	previous := -1
	for _, name := range []string{"-zulu, -z, -aardvark", "-alpha, -a", "-beta", "-help, -h", "-middle"} {
		index := strings.Index(out.String(), "  "+name)
		if index <= previous {
			t.Fatalf("option %q missing or out of order:\n%s", name, out.String())
		}
		previous = index
	}
}

type zeroSizeValue struct{}

func (*zeroSizeValue) String() string   { return "" }
func (*zeroSizeValue) Set(string) error { return nil }

func TestZeroSizeValues(t *testing.T) {
	fs := flag.NewFlagSet("test", flag.ContinueOnError)
	value := &zeroSizeValue{}
	fs.Var(value, "first", "first action")
	fs.Var(value, "second", "second action")
	flagx.Var(fs, value, "declared|d", "declared action")
	var out bytes.Buffer
	fs.SetOutput(&out)
	flagx.PrintDefaults(fs)
	for _, want := range []string{
		"-first value", "-second value", "-declared, -d value",
		"first action", "second action", "declared action",
	} {
		if !strings.Contains(out.String(), want) {
			t.Errorf("missing %q:\n%s", want, out.String())
		}
	}
}
