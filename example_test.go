package flagx_test

import (
	"flag"
	"fmt"
	"os"

	"github.com/Songmu/flagx"
)

func Example() {
	fs := flag.NewFlagSet("myapp", flag.ContinueOnError)
	fs.SetOutput(os.Stdout)
	var output string
	var verbose bool
	flagx.StringVar(fs, &output, "output|o", "", "write to `file`")
	flagx.BoolVar(fs, &verbose, "verbose|v", false, "verbose output")
	fs.Usage = flagx.Usage(fs, "myapp [options] <input>", "Process input files.")
	if err := fs.Parse([]string{"-o", "result.txt", "--verbose", "input.txt"}); err != nil {
		return
	}

	fmt.Println(output, verbose, fs.Args())
	fs.Usage()
	// Output:
	// result.txt true [input.txt]
	// Usage: myapp [options] <input>
	//
	// Process input files.
	//
	// Options:
	//   -output, -o file  write to file
	//   -verbose, -v      verbose output
	//   -help, -h         Show this help message
}

func ExampleNewFlagSet() {
	fs := flagx.NewFlagSet("myapp", flag.ContinueOnError)
	fs.SetOutput(os.Stdout)
	output := fs.String("output|o", "", "write to `file`")
	verbose := fs.Bool("verbose|v", false, "verbose output")
	if err := fs.Parse([]string{"-o", "result.txt", "--verbose", "input.txt"}); err != nil {
		return
	}
	fmt.Println(*output, *verbose, fs.Args())
	fs.Usage()
	// Output:
	// result.txt true [input.txt]
	// Usage: myapp
	//
	// Options:
	//   -output, -o file  write to file
	//   -verbose, -v      verbose output
	//   -help, -h         Show this help message
}
