// Package flagx adds flag aliases and formatted help to the standard
// library's flag.FlagSet without changing its parsing rules.
package flagx

import (
	"encoding"
	"flag"
	"time"
)

// Var registers value under each name in the |-separated names string.
// Names omit leading hyphens and are displayed in the specified order.
// Empty names panic; duplicate and invalid names follow flag.FlagSet behavior.
func Var(fs *flag.FlagSet, value flag.Value, names, usage string) {
	registerNames(fs, names, func(name string) {
		fs.Var(value, name, usage)
	})
}

// StringVar registers a string flag under the |-separated names, as in Var.
// All names write to p, which is initialized to value.
func StringVar(fs *flag.FlagSet, p *string, names, value, usage string) {
	registerNames(fs, names, func(name string) {
		fs.StringVar(p, name, value, usage)
	})
}

// String registers a string flag under the |-separated names, as in Var.
// It returns a pointer to the shared value, initialized to value.
func String(fs *flag.FlagSet, names, value, usage string) *string {
	p := new(string)
	StringVar(fs, p, names, value, usage)
	return p
}

// BoolVar registers a boolean flag under the |-separated names, as in Var.
// All names write to p, which is initialized to value.
func BoolVar(fs *flag.FlagSet, p *bool, names string, value bool, usage string) {
	registerNames(fs, names, func(name string) {
		fs.BoolVar(p, name, value, usage)
	})
}

// Bool registers a boolean flag under the |-separated names, as in Var.
// It returns a pointer to the shared value, initialized to value.
func Bool(fs *flag.FlagSet, names string, value bool, usage string) *bool {
	p := new(bool)
	BoolVar(fs, p, names, value, usage)
	return p
}

// IntVar registers an integer flag under the |-separated names, as in Var.
// All names write to p, which is initialized to value.
func IntVar(fs *flag.FlagSet, p *int, names string, value int, usage string) {
	registerNames(fs, names, func(name string) {
		fs.IntVar(p, name, value, usage)
	})
}

// Int registers an integer flag under the |-separated names, as in Var.
// It returns a pointer to the shared value, initialized to value.
func Int(fs *flag.FlagSet, names string, value int, usage string) *int {
	p := new(int)
	IntVar(fs, p, names, value, usage)
	return p
}

// Int64Var registers an int64 flag under the |-separated names, as in Var.
// All names write to p, which is initialized to value.
func Int64Var(fs *flag.FlagSet, p *int64, names string, value int64, usage string) {
	registerNames(fs, names, func(name string) {
		fs.Int64Var(p, name, value, usage)
	})
}

// Int64 registers an int64 flag under the |-separated names, as in Var.
// It returns a pointer to the shared value, initialized to value.
func Int64(fs *flag.FlagSet, names string, value int64, usage string) *int64 {
	p := new(int64)
	Int64Var(fs, p, names, value, usage)
	return p
}

// UintVar registers an unsigned integer flag under the |-separated names, as in Var.
// All names write to p, which is initialized to value.
func UintVar(fs *flag.FlagSet, p *uint, names string, value uint, usage string) {
	registerNames(fs, names, func(name string) {
		fs.UintVar(p, name, value, usage)
	})
}

// Uint registers an unsigned integer flag under the |-separated names, as in Var.
// It returns a pointer to the shared value, initialized to value.
func Uint(fs *flag.FlagSet, names string, value uint, usage string) *uint {
	p := new(uint)
	UintVar(fs, p, names, value, usage)
	return p
}

// Uint64Var registers a uint64 flag under the |-separated names, as in Var.
// All names write to p, which is initialized to value.
func Uint64Var(fs *flag.FlagSet, p *uint64, names string, value uint64, usage string) {
	registerNames(fs, names, func(name string) {
		fs.Uint64Var(p, name, value, usage)
	})
}

// Uint64 registers a uint64 flag under the |-separated names, as in Var.
// It returns a pointer to the shared value, initialized to value.
func Uint64(fs *flag.FlagSet, names string, value uint64, usage string) *uint64 {
	p := new(uint64)
	Uint64Var(fs, p, names, value, usage)
	return p
}

// Float64Var registers a float64 flag under the |-separated names, as in Var.
// All names write to p, which is initialized to value.
func Float64Var(fs *flag.FlagSet, p *float64, names string, value float64, usage string) {
	registerNames(fs, names, func(name string) {
		fs.Float64Var(p, name, value, usage)
	})
}

// Float64 registers a float64 flag under the |-separated names, as in Var.
// It returns a pointer to the shared value, initialized to value.
func Float64(fs *flag.FlagSet, names string, value float64, usage string) *float64 {
	p := new(float64)
	Float64Var(fs, p, names, value, usage)
	return p
}

// DurationVar registers a duration flag under the |-separated names, as in Var.
// All names write to p, which is initialized to value.
func DurationVar(fs *flag.FlagSet, p *time.Duration, names string, value time.Duration, usage string) {
	registerNames(fs, names, func(name string) {
		fs.DurationVar(p, name, value, usage)
	})
}

// Duration registers a duration flag under the |-separated names, as in Var.
// It returns a pointer to the shared value, initialized to value.
func Duration(fs *flag.FlagSet, names string, value time.Duration, usage string) *time.Duration {
	p := new(time.Duration)
	DurationVar(fs, p, names, value, usage)
	return p
}

// TextVar registers a text flag under the |-separated names, as in Var.
// All names write to p, which is initialized to value.
// The default value must have the same type as p, as required by flag.FlagSet.TextVar.
func TextVar(fs *flag.FlagSet, p encoding.TextUnmarshaler, names string, value encoding.TextMarshaler, usage string) {
	registerNames(fs, names, func(name string) {
		fs.TextVar(p, name, value, usage)
	})
}

// Func registers a callback flag under the |-separated names, as in Var.
// Each occurrence calls fn with its value; a returned error stops parsing.
func Func(fs *flag.FlagSet, names, usage string, fn func(string) error) {
	registerNames(fs, names, func(name string) {
		fs.Func(name, usage, fn)
	})
}

// BoolFunc registers a callback flag under the |-separated names, as in Var.
// Values are optional and default to "true" when omitted.
// Each occurrence calls fn with its value; a returned error stops parsing.
func BoolFunc(fs *flag.FlagSet, names, usage string, fn func(string) error) {
	registerNames(fs, names, func(name string) {
		fs.BoolFunc(name, usage, fn)
	})
}
