package flagx

import (
	"flag"
	"fmt"
	"io"
	"reflect"
	"sort"
	"strings"
	"unicode/utf8"
)

// Usage returns a help function suitable for fs.Usage.
// It uses fs.Output at call time, including subsequent SetOutput changes.
func Usage(fs *flag.FlagSet, synopsis, description string) func() {
	return func() {
		PrintUsage(fs.Output(), fs, synopsis, description)
	}
}

// PrintUsage writes a synopsis, an optional description, and an Options section
// to w. It does not change fs or its output writer.
func PrintUsage(w io.Writer, fs *flag.FlagSet, synopsis, description string) {
	fmt.Fprintf(w, "Usage: %s\n\n", synopsis)
	if description != "" {
		fmt.Fprintln(w, description)
		fmt.Fprintln(w)
	}
	fmt.Fprintln(w, "Options:")
	printDefaults(w, fs)
}

// PrintDefaults writes only the formatted option list to fs.Output.
// Names registered together are displayed in their specified order.
// Options registered through flagx come first, in registration order, followed
// by directly registered options in lexical order.
// Pointers to non-zero-sized values with the same type and address are also
// displayed as aliases.
// Missing help and h names are included for flag's built-in help handling,
// not registered.
func PrintDefaults(fs *flag.FlagSet) {
	printDefaults(fs.Output(), fs)
}

type option struct {
	flags []*flag.Flag
	order int
}

func (o option) sortName() string {
	var name string
	for _, f := range o.flags {
		if name == "" || lexicalNameLess(f.Name, name) {
			name = f.Name
		}
	}
	return name
}

func (o option) label() string {
	names := make([]string, len(o.flags))
	for i, f := range o.flags {
		names[i] = "-" + f.Name
	}
	label := strings.Join(names, ", ")
	name, _ := flag.UnquoteUsage(o.flags[0])
	if o.flags[0].Value == nil {
		name = ""
	}
	if b, ok := o.flags[0].Value.(interface{ IsBoolFlag() bool }); ok && b.IsBoolFlag() {
		name = ""
	}
	if name != "" {
		label += " " + name
	}
	return label
}

func (o option) description() string {
	f := o.flags[0]
	_, usage := flag.UnquoteUsage(f)
	if f.Value == nil || isZeroDefault(f) {
		return usage
	}
	value := f.DefValue
	typ := reflect.TypeOf(f.Value)
	if typ.Kind() == reflect.Pointer && typ.Elem().Kind() == reflect.String {
		value = fmt.Sprintf("%q", value)
	}
	return usage + " (default " + value + ")"
}

func (o option) less(other option) bool {
	if o.order == other.order {
		return o.sortName() < other.sortName()
	}
	if o.order < 0 {
		return false
	}
	return other.order < 0 || o.order < other.order
}

func nameLess(a, b string, declared map[string]nameOrder) bool {
	aOrder, aKnown := declared[a]
	bOrder, bKnown := declared[b]
	if aKnown != bKnown {
		return aKnown
	}
	if aKnown {
		return aOrder.rank < bOrder.rank
	}
	return lexicalNameLess(a, b)
}

func lexicalNameLess(a, b string) bool {
	aShort, bShort := utf8.RuneCountInString(a) == 1, utf8.RuneCountInString(b) == 1
	if aShort != bShort {
		return !aShort
	}
	return a < b
}

type valueIdentity struct {
	typ reflect.Type
	ptr uintptr
}

func pointerIdentity(value flag.Value) (valueIdentity, bool) {
	v := reflect.ValueOf(value)
	if !v.IsValid() || v.Kind() != reflect.Pointer || v.IsNil() || v.Type().Elem().Size() == 0 {
		return valueIdentity{}, false
	}
	return valueIdentity{v.Type(), v.Pointer()}, true
}

func groupOptions(fs *flag.FlagSet, declared map[string]nameOrder) []option {
	var result []option
	byPointer := make(map[valueIdentity]int)
	byGroup := make(map[int]int)
	fs.VisitAll(func(f *flag.Flag) {
		index := -1
		order, known := declared[f.Name]
		if known {
			if i, ok := byGroup[order.group]; ok {
				index = i
			}
		}
		key, pointer := pointerIdentity(f.Value)
		if pointer {
			if i, ok := byPointer[key]; ok {
				index = i
			}
		}
		if index < 0 {
			index = len(result)
			result = append(result, option{order: -1})
		}
		result[index].flags = append(result[index].flags, f)
		if known {
			byGroup[order.group] = index
			if result[index].order < 0 || order.group < result[index].order {
				result[index].order = order.group
			}
		}
		if pointer {
			byPointer[key] = index
		}
	})
	return result
}

func options(fs *flag.FlagSet) []option {
	declared := registeredNames(fs)
	result := groupOptions(fs, declared)
	help := option{order: -1}
	for _, name := range []string{"help", "h"} {
		if fs.Lookup(name) == nil {
			help.flags = append(help.flags, &flag.Flag{Name: name, Usage: "Show this help message"})
		}
	}
	if len(help.flags) > 0 {
		result = append(result, help)
	}
	for _, o := range result {
		sort.Slice(o.flags, func(i, j int) bool {
			return nameLess(o.flags[i].Name, o.flags[j].Name, declared)
		})
	}
	sort.Slice(result, func(i, j int) bool {
		return result[i].less(result[j])
	})
	return result
}

func printDefaults(w io.Writer, fs *flag.FlagSet) {
	opts := options(fs)
	labels := make([]string, len(opts))
	width := 0
	for i, o := range opts {
		labels[i] = o.label()
		if n := utf8.RuneCountInString(labels[i]); n > width {
			width = n
		}
	}
	const maxWidth = 28
	if width > maxWidth {
		width = maxWidth
	}
	indent := strings.Repeat(" ", width+4)
	for i, o := range opts {
		usage := o.description()
		fmt.Fprintf(w, "  %s", labels[i])
		if n := utf8.RuneCountInString(labels[i]); n > width {
			fmt.Fprint(w, "\n", indent)
		} else {
			fmt.Fprint(w, strings.Repeat(" ", width-n+2))
		}
		fmt.Fprintln(w, strings.ReplaceAll(usage, "\n", "\n"+indent))
	}
}

func isZeroDefault(f *flag.Flag) (zero bool) {
	// A custom Value may not support String on its zero value. Keep its default
	// visible rather than letting it prevent the rest of the help from printing.
	defer func() { _ = recover() }()
	typ := reflect.TypeOf(f.Value)
	var value reflect.Value
	if typ.Kind() == reflect.Pointer {
		value = reflect.New(typ.Elem())
	} else {
		value = reflect.Zero(typ)
	}
	return f.DefValue == value.Interface().(flag.Value).String()
}
