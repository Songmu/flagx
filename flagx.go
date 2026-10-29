// Package flagx adds long/short flag aliases and formatted help to the standard
// library's flag.FlagSet without changing its parsing rules.
package flagx

import (
	"flag"
	"time"
)

// Var registers value under long and, if nonempty, short.
// Names omit leading hyphens. Both names share value.
func Var(fs *flag.FlagSet, value flag.Value, long, short, usage string) {
	fs.Var(value, long, usage)
	if short != "" {
		fs.Var(value, short, usage)
	}
}

// StringVar registers a string flag under long and, if nonempty, short.
// Both names write to p, which is initialized to value.
func StringVar(fs *flag.FlagSet, p *string, long, short, value, usage string) {
	fs.StringVar(p, long, value, usage)
	if short != "" {
		fs.StringVar(p, short, value, usage)
	}
}

// BoolVar registers a boolean flag under long and, if nonempty, short.
// Both names write to p, which is initialized to value.
func BoolVar(fs *flag.FlagSet, p *bool, long, short string, value bool, usage string) {
	fs.BoolVar(p, long, value, usage)
	if short != "" {
		fs.BoolVar(p, short, value, usage)
	}
}

// IntVar registers an integer flag under long and, if nonempty, short.
// Both names write to p, which is initialized to value.
func IntVar(fs *flag.FlagSet, p *int, long, short string, value int, usage string) {
	fs.IntVar(p, long, value, usage)
	if short != "" {
		fs.IntVar(p, short, value, usage)
	}
}

// DurationVar registers a duration flag under long and, if nonempty, short.
// Both names write to p, which is initialized to value.
func DurationVar(fs *flag.FlagSet, p *time.Duration, long, short string, value time.Duration, usage string) {
	fs.DurationVar(p, long, value, usage)
	if short != "" {
		fs.DurationVar(p, short, value, usage)
	}
}
