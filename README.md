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

A lightweight, dependency-free flag library with an API matching Go's standard
`flag.FlagSet` and the same parsing rules. It adds two conveniences:

- Register aliases together with `|`, such as `"output|o"`.
- Display help in registration order rather than alphabetical order, with aliases
  grouped on one line.

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

fs.Parse(os.Args[1:])
```

Import `github.com/Songmu/flagx` alongside `flag`, `os`, and `time`.
Running `myapp --help` prints:

```text
Usage: myapp

Options:
  -output, -o file       write to file
  -verbose, -v           verbose output
  -jobs, -j workers      number of workers (default 4)
  -timeout, -t duration  timeout after duration (default 30s)
  -help, -h              Show this help message
```

## Usage

`flagx.FlagSet` embeds `flag.FlagSet` and provides the same registration methods,
including custom values and callbacks. Use them as usual, with `|`-separated
names to register aliases sharing one value. Names omit leading hyphens; a single
name works too. Aliases need not be a long/short pair, and their specified order
is preserved in help: `"o|output"` displays as `-o, -output string`.

Options registered through flagx appear in registration order, as in the example
above: output, verbose, jobs, timeout. Flags registered directly through the
embedded `flag.FlagSet` follow in alphabetical order.

Parsing is unchanged: both `-output` and `--output` work, parsing stops at the
first positional argument, short flags cannot be bundled, and bool flags do not
require a value.

### Custom help

Default help prints the flag set's name and option list. Assign `fs.Usage` and
call `fs.PrintUsage` to add a synopsis and description, or use
`PrintDefaults` to embed just the option list:

```go
fs.Usage = func() {
    fmt.Fprintln(fs.Output(), "Usage: myapp [options] <input>")
    fmt.Fprintln(fs.Output(), "\nOptions:")
    fs.PrintDefaults()
}
```

Default help and `PrintDefaults` write to `fs.Output()`; change it with
`fs.SetOutput(w)`. Help includes `-help` and `-h` without registering them,
leaving their handling to standard `flag`. If you register either explicitly,
your application must handle it.

### Existing flag sets

Package-level functions add the same features to an existing `*flag.FlagSet`:

```go
fs := flag.NewFlagSet("myapp", flag.ExitOnError)
flagx.String(fs, "output|o", "", "write to `file`")
fs.Usage = flagx.Usage(fs, "myapp [options]", "")
```

Use `&fs.FlagSet` to pass a `*flagx.FlagSet` to APIs accepting a `*flag.FlagSet`.
See the [package reference][PkgGoDev] for the full API.

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
