package flagx_test

import (
	"flag"
	"io"
	"slices"
	"strings"
	"testing"

	"github.com/Songmu/flagx"
)

type registrationFlagSet interface {
	Lookup(string) *flag.Flag
	Parse([]string) error
	VisitAll(func(*flag.Flag))
}

func testRegistration(t *testing.T, initial, input, want string, register func(string) (registrationFlagSet, func() string)) {
	t.Helper()
	for _, names := range []string{"option", "o|other|option"} {
		t.Run(names, func(t *testing.T) {
			aliases := strings.Split(names, "|")
			wantNames := slices.Clone(aliases)
			slices.Sort(wantNames)
			for _, alias := range aliases {
				t.Run(alias, func(t *testing.T) {
					fs, value := register(names)
					if got := value(); got != initial {
						t.Fatalf("initial value = %q, want %q", got, initial)
					}
					var gotNames []string
					fs.VisitAll(func(f *flag.Flag) { gotNames = append(gotNames, f.Name) })
					if !slices.Equal(gotNames, wantNames) {
						t.Fatalf("registered names = %v, want %v", gotNames, wantNames)
					}
					if err := fs.Parse([]string{"--" + alias + "=" + input}); err != nil {
						t.Fatal(err)
					}
					if got := value(); got != want {
						t.Errorf("value = %q, want %q", got, want)
					}
					for _, name := range aliases {
						f := fs.Lookup(name)
						if f.DefValue != initial {
							t.Errorf("%s default = %q, want %q", name, f.DefValue, initial)
						}
						if got := f.Value.String(); got != want {
							t.Errorf("%s shared value = %q, want %q", name, got, want)
						}
					}
				})
			}
		})
	}
}

func TestRegistrationSyntax(t *testing.T) {
	for _, name := range []string{"-option", "--option", "-o", "--o", "-other", "--other"} {
		t.Run(name, func(t *testing.T) {
			fs := flag.NewFlagSet("test", flag.ContinueOnError)
			fs.SetOutput(io.Discard)
			value := flagx.String(fs, "o|other|option", "initial", "an option")
			if err := fs.Parse([]string{name + "=changed"}); err != nil {
				t.Fatal(err)
			}
			if *value != "changed" {
				t.Errorf("value = %q, want changed", *value)
			}
		})
	}
}
