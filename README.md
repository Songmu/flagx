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

Small helpers for long/short flag aliases and formatted help, built on Go's
standard `flag.FlagSet`. No dependencies or changes to flag parsing.

## Synopsis

```go
fs := flag.NewFlagSet("myapp", flag.ExitOnError)

var output string
var verbose bool
var jobs int
var timeout time.Duration

flagx.StringVar(fs, &output, "output", "o", "", "write to `file`")
flagx.BoolVar(fs, &verbose, "verbose", "v", false, "verbose output")
flagx.IntVar(fs, &jobs, "jobs", "j", 4, "number of `workers`")
flagx.DurationVar(fs, &timeout, "timeout", "t", 30*time.Second, "timeout after `duration`")

fs.Usage = flagx.Usage(fs, "myapp [options] <input>", "Process input files.")
fs.Parse(os.Args[1:])
```

Import `github.com/Songmu/flagx` alongside `flag`, `os`, and `time`.
Running `myapp --help` prints:

```text
Usage: myapp [options] <input>

Process input files.

Options:
  -help, -h              Show this help message
  -jobs, -j workers      number of workers (default 4)
  -output, -o file       write to file
  -timeout, -t duration  timeout after duration (default 30s)
  -verbose, -v           verbose output
```

## Description

`StringVar`, `BoolVar`, `IntVar`, and `DurationVar` register both names against
the same variable. `Var` does the same for a custom `flag.Value`. Pass names
without hyphens; an empty short name registers only the long name. Duplicate
and invalid names follow the standard library's behavior.

Both `-output` and `--output` work, as do `-o` and `--o`. Standard `flag`
semantics still apply: parsing stops at the first positional argument, short
flags cannot be bundled, and bool flags do not require a value.

Help combines values with the same pointer and concrete type, including flags
registered directly through `flag.FlagSet`. Single-character names are shown
after longer names, all with one hyphen. Options are sorted by their first long
name, or their short name if there is no long name. Multiple names of each kind
are sorted lexicographically. The first displayed name supplies the usage and
default value when aliases have different metadata.

This pointer-based grouping is a heuristic, not an explicit alias registry.
Non-pointer values are not grouped, even if they share backing storage.
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
    flagx.PrintDefaults(fs)
}
```

`Usage` and `PrintDefaults` use the current `fs.Output()`, including later
`SetOutput` changes. `PrintUsage(w, fs, synopsis, description)` writes a complete
help message to an explicit writer without changing the flag set.

flagx does not manage subcommands, registration order, groups, required/hidden
flags, or Markdown, man page, and shell completion generation.

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
