flagx
=======

[![Test Status](https://github.com/Songmu/flagx/actions/workflows/test.yaml/badge.svg?branch=main)][actions]
[![Coverage Status](https://codecov.io/gh/Songmu/flagx/branch/main/graph/badge.svg)][codecov]
[![MIT License](https://img.shields.io/github/license/Songmu/flagx)][license]
[![PkgGoDev](https://pkg.go.dev/badge/github.com/Songmu/flagx)][PkgGoDev]

[actions]: https://github.com/Songmu/flagx/actions?workflow=test
[codecov]: https://codecov.io/gh/Songmu/flagx
[license]: https://github.com/Songmu/flagx/blob/main/LICENSE
[PkgGoDev]: https://pkg.go.dev/github.com/Songmu/flagx

Flag aliases and formatted help, built on Go's
standard `flag.FlagSet`. No dependencies or changes to flag parsing.

## Synopsis

```go
fs := flagx.NewFlagSet("myapp", flag.ExitOnError)

var output string
var verbose bool
var jobs int
var timeout time.Duration

fs.StringVar(&output, "output|o", "", "write to `file`")
fs.BoolVar(&verbose, "verbose|v", false, "verbose output")
fs.IntVar(&jobs, "jobs|j", 4, "number of `workers`")
fs.DurationVar(&timeout, "timeout|t", 30*time.Second, "timeout after `duration`")

fs.Usage = func() {
    fs.PrintUsage(fs.Output(), "myapp [options] <input>", "Process input files.")
}
fs.Parse(os.Args[1:])
```

Import `github.com/Songmu/flagx` alongside `flag`, `os`, and `time`.
Running `myapp --help` prints:

```text
Usage: myapp [options] <input>

Process input files.

Options:
  -output, -o file       write to file
  -verbose, -v           verbose output
  -jobs, -j workers      number of workers (default 4)
  -timeout, -t duration  timeout after duration (default 30s)
  -help, -h              Show this help message
```

## Description

`NewFlagSet(name, errorHandling)` returns a `*flagx.FlagSet` embedding
`flag.FlagSet`. Its registration methods accept one or more `|`-separated names.
`String`, `Bool`, `Int`, `Int64`, `Uint`,
`Uint64`, `Float64`, and `Duration` return a pointer shared by all names;
their `Var` forms bind an existing variable. `TextVar` binds a value implementing
`encoding.TextUnmarshaler`, and `Var` accepts a custom `flag.Value`.
`Func` and `BoolFunc` invoke a callback for each occurrence; `BoolFunc` allows
omitting the value. Defaults and parsing follow the corresponding standard
`flag.FlagSet` methods.

Default help uses the flag set's current name as the synopsis and prints the
formatted option list. Assign `fs.Usage` to customize it, as above.
The zero value is usable with an empty name and `flag.ContinueOnError`;
`Init` sets the name and error policy without clearing flags or custom help.
Do not copy a `FlagSet` after first use.

The package-level registration and help functions remain available for existing
`*flag.FlagSet` values, taking the flag set as an argument:

```go
fs := flag.NewFlagSet("myapp", flag.ExitOnError)
flagx.String(fs, "output|o", "", "write to `file`")
fs.Usage = flagx.Usage(fs, "myapp [options]", "")
```

Use `&fs.FlagSet` when passing a `*flagx.FlagSet` to these functions or other
APIs accepting a `*flag.FlagSet`.
If another API parses `&fs.FlagSet` while `Usage` is nil (including a zero-value
set), call `fs.Init(fs.Name(), fs.ErrorHandling())` before handing it off to
install flagx's default help.

Names omit leading hyphens and are split literally, without
trimming whitespace. Empty names panic; duplicate and invalid names otherwise
follow the standard library's behavior.

Names are aliases, not a required long/short pair: `"output"`, `"output|o"`,
`"o|output"`, `"aa|bb"`, and `"output|out|o"` are all accepted. Help preserves the
specified name order, so `"o|output"` displays as `-o, -output string`.

Both `-output` and `--output` work, as do `-o` and `--o`. Standard `flag`
semantics still apply: parsing stops at the first positional argument, short
flags cannot be bundled, and bool flags do not require a value.

Help combines names registered together and also values with the same pointer
and concrete type, including flags registered directly through `flag.FlagSet`.
Within each option, flagx-registered names retain their registration order;
directly registered names follow, sorted long before short and lexicographically
within each kind because their registration order is unavailable. All names
have one hyphen.

Options registered through flagx come first, in registration order. Options
registered only through `flag.FlagSet` follow, sorted by their lexicographically
smallest long name, or their smallest short name if there is no long name.
An option combining both kinds of registration uses its first flagx registration
position. The first displayed name supplies the usage and default value when
aliases have different metadata.

Pointer-based grouping of separately registered flags is a heuristic.
Non-pointer values and pointers to zero-sized values are grouped only when their
names are registered together.
Backquoted placeholders use `flag.UnquoteUsage`; bool placeholders are omitted.
Nonzero registration-time defaults are shown, with string defaults quoted.
Custom values are compared with their type's zero value; if its `String` method
panics, the default is shown without interrupting help output.

`help` and `h` are not registered: standard `flag` handles undefined help names.
The option list displays missing help names and preserves explicitly registered
ones. An explicitly registered help flag has its normal parsing behavior, so
the application is responsible for handling it.

### Custom help

Use `PrintDefaults` to embed just the option list in your own help:

```go
fs.Usage = func() {
    fmt.Fprintln(fs.Output(), "Usage: myapp [options] <input>")
    fmt.Fprintln(fs.Output(), "\nOptions:")
    fs.PrintDefaults()
}
```

Default help and `fs.PrintDefaults()` use the current `fs.Output()`, including
later `SetOutput` changes. `fs.PrintUsage(w, synopsis, description)` writes a
complete help message to an explicit writer without changing the flag set.
For a standard `*flag.FlagSet`, use `flagx.PrintDefaults(fs)` and
`flagx.PrintUsage(w, fs, synopsis, description)` instead.

flagx does not manage subcommands, groups, required/hidden flags, or Markdown,
man page, and shell completion generation.

## Installation

```console
% go get github.com/Songmu/flagx
```

## Development

After changing GitHub Actions dependencies, run `gh actions-lock` and commit
`.github/workflows/actions.lock` with the workflow changes. Install the tool
with `gh extension install github/gh-actions-lock`.

## Author

[Songmu](https://github.com/Songmu)
