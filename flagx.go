// Package flagx adds flag aliases and formatted help to the standard
// library's flag.FlagSet without changing its parsing rules.
package flagx

import (
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

// BoolVar registers a boolean flag under the |-separated names, as in Var.
// All names write to p, which is initialized to value.
func BoolVar(fs *flag.FlagSet, p *bool, names string, value bool, usage string) {
	registerNames(fs, names, func(name string) {
		fs.BoolVar(p, name, value, usage)
	})
}

// IntVar registers an integer flag under the |-separated names, as in Var.
// All names write to p, which is initialized to value.
func IntVar(fs *flag.FlagSet, p *int, names string, value int, usage string) {
	registerNames(fs, names, func(name string) {
		fs.IntVar(p, name, value, usage)
	})
}

// DurationVar registers a duration flag under the |-separated names, as in Var.
// All names write to p, which is initialized to value.
func DurationVar(fs *flag.FlagSet, p *time.Duration, names string, value time.Duration, usage string) {
	registerNames(fs, names, func(name string) {
		fs.DurationVar(p, name, value, usage)
	})
}
