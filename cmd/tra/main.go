package main

import (
	"context"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"strings"

	"github.com/masp/awktree/eval"
	"github.com/masp/awktree/token"
)

var (
	flagHelp     = flag.Bool("h", false, "Show help")
	flagVersion  = flag.Bool("v", false, "Show version")
	flagVerbose  = flag.Bool("d", false, "Verbose mode")
	flagProgFile = flag.String("f", "", "Path to a tra file to execute instead of inline")
)

func usage() {
	fmt.Fprintf(os.Stderr, "usage: %s [options] [-f progfile | program] [file ...]\n", os.Args[0])
	flag.PrintDefaults()
	os.Exit(1)
}

func main() {
	flag.Usage = usage
	flag.Parse()

	if *flagHelp {
		fmt.Fprintf(os.Stderr, "AwkTree: a tree-sitter based awk-like tool\n")
		flag.Usage()
		fmt.Fprintf(os.Stderr, "GitHub: https://github.com/masp/AwkTree\n")
	}

	if *flagVersion {
		fmt.Fprintf(os.Stderr, "AwkTree v0.1\n")
		os.Exit(0)
	}

	if *flagVerbose {
		log.SetOutput(os.Stderr)
	} else {
		log.SetOutput(io.Discard)
	}

	var (
		err        error
		programSrc io.Reader
		inputs     []input
	)
	if *flagProgFile != "" {
		programSrc, err = os.Open(*flagProgFile)
		if err != nil {
			fatalf("tra: error: open program file (-f): %v", err)
		}
		if flag.NArg() == 0 {
			inputs = append(inputs, input{filename: "<stdin>", rd: os.Stdin})
		} else {
			inputs = readInputs(flag.Args())
		}
	} else {
		if len(flag.Args()) == 0 {
			usage()
		}
		programSrc = strings.NewReader(flag.Arg(0))
		if flag.NArg() > 1 {
			inputs = readInputs(flag.Args()[1:])
		} else {
			inputs = append(inputs, input{filename: "<stdin>", rd: os.Stdin})
		}
	}

	progFilename := *flagProgFile
	if progFilename == "" {
		progFilename = "<inline>"
	}
	programFullSrc, err := io.ReadAll(programSrc)
	if err != nil {
		fatalf("error: reading program: %v", err)
	}
	prog, err := eval.Compile(progFilename, programFullSrc)
	if err != nil {
		if list, ok := err.(token.ErrorList); ok {
			for _, err := range list {
				fmt.Fprintf(os.Stderr, "tra: syntax error: %v\n", err)
			}
		} else {
			fatalf("tra: error: %v\n", err)
		}
		os.Exit(1)
	}

	for _, input := range inputs {
		src, err := io.ReadAll(input.rd)
		if err != nil {
			fmt.Fprintf(os.Stderr, "tra: can't open file %s: %v\n", input.filename, err)
			continue
		}
		err = prog.Eval(context.Background(), src, &eval.Options{
			Filename: input.filename,
			Stdout:   os.Stdout,
		})
		if err != nil {
			fatalf("tra: error: %v\n", err)
		}
	}
}

func readInputs(inputs []string) []input {
	var readers []input
	for _, f := range inputs {
		r, err := os.Open(f)
		if err != nil {
			fmt.Fprintf(os.Stderr, "tra: can't open file: %v\n", err)
			continue
		}
		readers = append(readers, input{filename: f, rd: r})
	}
	return readers
}

type input struct {
	filename string
	rd       io.Reader
}

func fatalf(format string, args ...any) {
	fmt.Fprintf(os.Stderr, format, args...)
	os.Exit(1)
}
