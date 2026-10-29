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

const sampleDefaults = `  -help, -h              Show this help message
  -jobs, -j workers      number of workers (default 4)
  -output, -o file       write to file
  -timeout, -t duration  timeout after duration (default 30s)
  -verbose, -v           verbose output
`

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
		if got := original.String(); got != want {
			t.Errorf("help:\n%s\nwant:\n%s", got, want)
		}
	}

	fs.SetOutput(&redirected)
	fs.Usage()
	if redirected.String() != want {
		t.Errorf("Usage did not follow SetOutput: %q", redirected.String())
	}
	redirected.Reset()
	flagx.PrintDefaults(fs)
	if got := redirected.String(); got != sampleDefaults {
		t.Errorf("PrintDefaults = %q, want %q", got, sampleDefaults)
	}
	redirected.Reset()
	flagx.PrintUsage(&explicit, fs, "myapp", "")
	if got, want := explicit.String(), "Usage: myapp\n\nOptions:\n"+sampleDefaults; got != want {
		t.Errorf("PrintUsage = %q, want %q", got, want)
	}
	if fs.Output() != &redirected || redirected.Len() != 0 {
		t.Fatal("PrintUsage changed or wrote to fs.Output")
	}
}

func TestDefaults(t *testing.T) {
	fs := flag.NewFlagSet("test", flag.ContinueOnError)
	fs.String("string", "hello\nworld", "a string")
	fs.Bool("bool", true, "use `mode`")
	fs.Int("int", 9, "an integer")
	fs.Duration("duration", time.Minute, "a duration")
	fs.String("zero-string", "", "empty")
	fs.Bool("zero-bool", false, "off")
	fs.Int("zero-int", 0, "zero")
	fs.Duration("zero-duration", 0, "immediate")
	if err := fs.Set("int", "100"); err != nil {
		t.Fatal(err)
	}
	var out bytes.Buffer
	fs.SetOutput(&out)
	flagx.PrintDefaults(fs)
	for _, want := range []string{
		`a string (default "hello\nworld")`,
		"use mode (default true)",
		"an integer (default 9)",
		"a duration (default 1m0s)",
	} {
		if !strings.Contains(out.String(), want) {
			t.Errorf("missing %q:\n%s", want, out.String())
		}
	}
	for _, line := range strings.Split(out.String(), "\n") {
		if strings.Contains(line, "-zero-") && strings.Contains(line, "(default") {
			t.Errorf("zero default not omitted: %s", line)
		}
		if strings.Contains(line, "-bool ") && strings.Contains(line, "-bool mode") {
			t.Errorf("boolean has a placeholder: %s", line)
		}
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
