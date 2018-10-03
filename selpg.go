package main

import (
	"bufio"
	"fmt"
	flag "github.com/spf13/pflag"
	"os"
	"os/exec"
)

type SelpgArgs struct {
	startPage int
	endPage   int
	pageLen   int
	pageType  bool
	printDest string
	inputFile string
}

/**
 * init args:
 *
 * "page length": 72(default)
 * "page type": false(default) -> lNumber/page, true -> \f as next page
 * "print dest": ""(default) -> no print dest
 */

func initArgs(args *SelpgArgs) {
	flag.IntVarP(&args.startPage, "sNumber", "s", -1, "start page")
	flag.IntVarP(&args.endPage, "eNumber", "e", -1, "end page")
	flag.IntVarP(&args.pageLen, "lNumber", "l", 72, "lines/page")
	flag.BoolVarP(&args.pageType, "formFeed", "f", false, "form-feed-delimited")
	flag.StringVarP(&args.printDest, "dest", "d", "", "print dest")
	flag.Parse()
}

/**
 * check the CLI arguments is valid or not
 *
 * startPage and endPage are  mandatory_opts
 * Minimum cmd is "selpg -s start_page -e end_page"
 *
 * start_page shouldn't be greater than end_page
 */
func processArgs(args *SelpgArgs, progname *string) {

	if args.startPage == -1 || args.endPage == -1 {

		fmt.Fprintf(os.Stderr, "%s: not enough arguments\n", *progname)
		if args.startPage == -1 && args.endPage != -1 {
			fmt.Fprintf(os.Stderr, "%s: 1st arg should be -s start_page\n", *progname)
			os.Exit(1)
		} else if args.startPage != -1 && args.endPage == -1 {
			fmt.Fprintf(os.Stderr, "%s: 2st arg should be -e end_page\n", *progname)
			os.Exit(1)
		}
	} else {
		if args.startPage > args.endPage {
			fmt.Fprintf(os.Stderr, "%s: start_page(%d) shouldn't be greater than end_page(%d)\n", *progname, args.startPage, args.endPage)
			os.Exit(2)
		}

		if args.startPage < 1 {
			fmt.Fprintf(os.Stderr, "%s: invalid start page %d\n", *progname, args.startPage)
			os.Exit(2)
		}

		if args.endPage < 1 {
			fmt.Fprintf(os.Stderr, "%s: invalid end page %d\n", *progname, args.endPage)
			os.Exit(2)
		}
	}

}

func processInput(args *SelpgArgs, progname *string) {

	result := ""

	if flag.NArg() == 0 {
		// pageType: true -> by '\n', false -> by n lines/page
		if args.pageType {
			reader := bufio.NewReader(os.Stdin)
			selpgByF(reader, args, &result, progname)
		} else {
			scanner := bufio.NewScanner(os.Stdin)
			selpgByLine(scanner, args, &result)
		}
	} else {
		args.inputFile = flag.Arg(0)
		file, err := os.Open(args.inputFile)
		if err != nil {
			fmt.Fprintf(os.Stderr, "%s: %s\n", *progname, err)
			os.Exit(3)
		}

		// pageType: true -> by '\n', false -> by n lines/page
		if args.pageType {
			reader := bufio.NewReader(file)
			selpgByF(reader, args, &result, progname)
		} else {
			scanner := bufio.NewScanner(file)
			selpgByLine(scanner, args, &result)
		}

		// fmt.Print(result)

		file.Close()
	}

	/* here we have not print device, so we use "cat -n" cmd instead */
	if args.printDest == "" {
		fmt.Fprint(os.Stdout, result)
	} else {
		cmd := exec.Command("cat", "-n")
		stdin, err := cmd.StdinPipe()
		if err != nil {
			panic(err)
		}
		_, err = stdin.Write([]byte(result))
		if err != nil {
			fmt.Fprintf(os.Stderr, "%s: %s", err)
		}
		stdin.Close()
		cmd.Stdout = os.Stdout
		cmd.Run()
	}
}

func main() {

	progname := os.Args[0]

	var args SelpgArgs

	initArgs(&args)

	processArgs(&args, &progname)

	processInput(&args, &progname)
}
