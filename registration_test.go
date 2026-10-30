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

type registrationValues struct {
	initial string
	input   string
	want    string
}

var registrationData = struct {
	stringValue   registrationValues
	boolValue     registrationValues
	intValue      registrationValues
	int64Value    registrationValues
	uintValue     registrationValues
	uint64Value   registrationValues
	float64Value  registrationValues
	durationValue registrationValues
	textValue     registrationValues
}{
	stringValue:   registrationValues{"initial", "changed", "changed"},
	boolValue:     registrationValues{"false", "true", "true"},
	intValue:      registrationValues{"4", "12", "12"},
	int64Value:    registrationValues{"-4", "-1099511627776", "-1099511627776"},
	uintValue:     registrationValues{"4", "0xff", "255"},
	uint64Value:   registrationValues{"4", "18446744073709551615", "18446744073709551615"},
	float64Value:  registrationValues{"1.5", "-2.5e2", "-250"},
	durationValue: registrationValues{"1s", "2m", "2m0s"},
	textValue:     registrationValues{"127.0.0.1", "::1", "::1"},
}

func testRegistration(t *testing.T, values registrationValues, register func(string) (registrationFlagSet, func() string)) {
	t.Helper()
	for _, names := range []string{"option", "o|other|option"} {
		t.Run(names, func(t *testing.T) {
			aliases := strings.Split(names, "|")
			wantNames := slices.Clone(aliases)
			slices.Sort(wantNames)
			for _, alias := range aliases {
				t.Run(alias, func(t *testing.T) {
					fs, value := register(names)
					if got := value(); got != values.initial {
						t.Fatalf("initial value = %q, want %q", got, values.initial)
					}
					var gotNames []string
					fs.VisitAll(func(f *flag.Flag) { gotNames = append(gotNames, f.Name) })
					if !slices.Equal(gotNames, wantNames) {
						t.Fatalf("registered names = %v, want %v", gotNames, wantNames)
					}
					if err := fs.Parse([]string{"--" + alias + "=" + values.input}); err != nil {
						t.Fatal(err)
					}
					if got := value(); got != values.want {
						t.Errorf("value = %q, want %q", got, values.want)
					}
					for _, name := range aliases {
						f := fs.Lookup(name)
						if f.DefValue != values.initial {
							t.Errorf("%s default = %q, want %q", name, f.DefValue, values.initial)
						}
						if got := f.Value.String(); got != values.want {
							t.Errorf("%s shared value = %q, want %q", name, got, values.want)
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
