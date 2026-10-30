package flagx

import (
	"encoding"
	"flag"
	"io"
	"time"
)

// FlagSet extends flag.FlagSet with alias-aware registration and formatted help.
// Registration methods accept |-separated names, as in Var.
// The zero value has no name and uses flag.ContinueOnError.
// A FlagSet must not be copied after first use.
type FlagSet struct {
	flag.FlagSet
}

// NewFlagSet returns an empty flag set with the specified name and error policy.
// Its default Usage prints the name and the formatted option list.
func NewFlagSet(name string, errorHandling flag.ErrorHandling) *FlagSet {
	fs := new(FlagSet)
	fs.Init(name, errorHandling)
	return fs
}

// Init sets the name and error policy, installing formatted help if Usage is nil.
// It preserves registered flags and any existing Usage function.
func (fs *FlagSet) Init(name string, errorHandling flag.ErrorHandling) {
	fs.FlagSet.Init(name, errorHandling)
	if fs.Usage == nil {
		fs.Usage = fs.defaultUsage
	}
}

// Parse parses arguments using flag.FlagSet's rules.
// If Usage is nil, it installs the default formatted help before parsing.
func (fs *FlagSet) Parse(arguments []string) error {
	if fs.Usage == nil {
		fs.Usage = fs.defaultUsage
	}
	return fs.FlagSet.Parse(arguments)
}

func (fs *FlagSet) defaultUsage() {
	fs.PrintUsage(fs.Output(), fs.Name(), "")
}

// PrintDefaults writes the formatted option list to fs.Output.
func (fs *FlagSet) PrintDefaults() {
	PrintDefaults(&fs.FlagSet)
}

// PrintUsage writes a synopsis, optional description, and formatted option list
// to w without changing fs.Output.
func (fs *FlagSet) PrintUsage(w io.Writer, synopsis, description string) {
	PrintUsage(w, &fs.FlagSet, synopsis, description)
}

// Var registers value under each of the |-separated names, as in Var.
func (fs *FlagSet) Var(value flag.Value, names, usage string) {
	Var(&fs.FlagSet, value, names, usage)
}

// StringVar binds p to the |-separated names, as in StringVar.
func (fs *FlagSet) StringVar(p *string, names, value, usage string) {
	StringVar(&fs.FlagSet, p, names, value, usage)
}

// String registers the |-separated names, as in String.
func (fs *FlagSet) String(names, value, usage string) *string {
	return String(&fs.FlagSet, names, value, usage)
}

// BoolVar binds p to the |-separated names, as in BoolVar.
func (fs *FlagSet) BoolVar(p *bool, names string, value bool, usage string) {
	BoolVar(&fs.FlagSet, p, names, value, usage)
}

// Bool registers the |-separated names, as in Bool.
func (fs *FlagSet) Bool(names string, value bool, usage string) *bool {
	return Bool(&fs.FlagSet, names, value, usage)
}

// IntVar binds p to the |-separated names, as in IntVar.
func (fs *FlagSet) IntVar(p *int, names string, value int, usage string) {
	IntVar(&fs.FlagSet, p, names, value, usage)
}

// Int registers the |-separated names, as in Int.
func (fs *FlagSet) Int(names string, value int, usage string) *int {
	return Int(&fs.FlagSet, names, value, usage)
}

// Int64Var binds p to the |-separated names, as in Int64Var.
func (fs *FlagSet) Int64Var(p *int64, names string, value int64, usage string) {
	Int64Var(&fs.FlagSet, p, names, value, usage)
}

// Int64 registers the |-separated names, as in Int64.
func (fs *FlagSet) Int64(names string, value int64, usage string) *int64 {
	return Int64(&fs.FlagSet, names, value, usage)
}

// UintVar binds p to the |-separated names, as in UintVar.
func (fs *FlagSet) UintVar(p *uint, names string, value uint, usage string) {
	UintVar(&fs.FlagSet, p, names, value, usage)
}

// Uint registers the |-separated names, as in Uint.
func (fs *FlagSet) Uint(names string, value uint, usage string) *uint {
	return Uint(&fs.FlagSet, names, value, usage)
}

// Uint64Var binds p to the |-separated names, as in Uint64Var.
func (fs *FlagSet) Uint64Var(p *uint64, names string, value uint64, usage string) {
	Uint64Var(&fs.FlagSet, p, names, value, usage)
}

// Uint64 registers the |-separated names, as in Uint64.
func (fs *FlagSet) Uint64(names string, value uint64, usage string) *uint64 {
	return Uint64(&fs.FlagSet, names, value, usage)
}

// Float64Var binds p to the |-separated names, as in Float64Var.
func (fs *FlagSet) Float64Var(p *float64, names string, value float64, usage string) {
	Float64Var(&fs.FlagSet, p, names, value, usage)
}

// Float64 registers the |-separated names, as in Float64.
func (fs *FlagSet) Float64(names string, value float64, usage string) *float64 {
	return Float64(&fs.FlagSet, names, value, usage)
}

// DurationVar binds p to the |-separated names, as in DurationVar.
func (fs *FlagSet) DurationVar(p *time.Duration, names string, value time.Duration, usage string) {
	DurationVar(&fs.FlagSet, p, names, value, usage)
}

// Duration registers the |-separated names, as in Duration.
func (fs *FlagSet) Duration(names string, value time.Duration, usage string) *time.Duration {
	return Duration(&fs.FlagSet, names, value, usage)
}

// TextVar binds p to the |-separated names, as in TextVar.
func (fs *FlagSet) TextVar(p encoding.TextUnmarshaler, names string, value encoding.TextMarshaler, usage string) {
	TextVar(&fs.FlagSet, p, names, value, usage)
}

// Func registers a callback under the |-separated names, as in Func.
func (fs *FlagSet) Func(names, usage string, fn func(string) error) {
	Func(&fs.FlagSet, names, usage, fn)
}

// BoolFunc registers a boolean callback under the |-separated names, as in BoolFunc.
func (fs *FlagSet) BoolFunc(names, usage string, fn func(string) error) {
	BoolFunc(&fs.FlagSet, names, usage, fn)
}
