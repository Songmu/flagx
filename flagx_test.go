package flagx_test

import (
	"flag"
	"fmt"
	"io"
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
		{"string", func(fs *flag.FlagSet, short string) func() string {
			var v string
			flagx.StringVar(fs, &v, "option", short, "initial", "an option")
			return func() string { return v }
		}, "initial", "changed", "changed"},
		{"bool", func(fs *flag.FlagSet, short string) func() string {
			var v bool
			flagx.BoolVar(fs, &v, "option", short, false, "an option")
			return func() string { return fmt.Sprint(v) }
		}, "false", "true", "true"},
		{"int", func(fs *flag.FlagSet, short string) func() string {
			var v int
			flagx.IntVar(fs, &v, "option", short, 4, "an option")
			return func() string { return fmt.Sprint(v) }
		}, "4", "12", "12"},
		{"duration", func(fs *flag.FlagSet, short string) func() string {
			var v time.Duration
			flagx.DurationVar(fs, &v, "option", short, time.Second, "an option")
			return func() string { return v.String() }
		}, "1s", "2m", "2m0s"},
		{"value", func(fs *flag.FlagSet, short string) func() string {
			v := listValue{"initial"}
			flagx.Var(fs, &v, "option", short, "an option")
			return v.String
		}, "initial", "changed", "initial,changed"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			for _, name := range []string{"-option", "--option", "-o", "--o"} {
				t.Run(name, func(t *testing.T) {
					fs := flag.NewFlagSet("test", flag.ContinueOnError)
					fs.SetOutput(io.Discard)
					value := tt.register(fs, "o")
					if got := value(); got != tt.initial {
						t.Fatalf("initial value = %q, want %q", got, tt.initial)
					}
					if err := fs.Parse([]string{name + "=" + tt.input}); err != nil {
						t.Fatal(err)
					}
					if got := value(); got != tt.want {
						t.Errorf("value = %q, want %q", got, tt.want)
					}
					for _, alias := range []string{"option", "o"} {
						if got := fs.Lookup(alias).DefValue; got != tt.initial {
							t.Errorf("%s default = %q, want %q", alias, got, tt.initial)
						}
					}
				})
			}
			t.Run("no short name", func(t *testing.T) {
				fs := flag.NewFlagSet("test", flag.ContinueOnError)
				tt.register(fs, "")
				var names []string
				fs.VisitAll(func(f *flag.Flag) { names = append(names, f.Name) })
				if len(names) != 1 || names[0] != "option" {
					t.Errorf("registered names = %v, want [option]", names)
				}
			})
		})
	}
}

func TestBoolWithoutArgument(t *testing.T) {
	fs := flag.NewFlagSet("test", flag.ContinueOnError)
	var verbose bool
	flagx.BoolVar(fs, &verbose, "verbose", "v", false, "verbose output")
	if err := fs.Parse([]string{"-v", "--verbose=false", "--verbose"}); err != nil {
		t.Fatal(err)
	}
	if !verbose {
		t.Fatal("boolean flag did not become true")
	}
}
