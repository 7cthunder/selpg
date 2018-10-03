package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
)

func selpgByLine(scanner *bufio.Scanner, args *SelpgArgs, data *string) {

	lineCounter := 0

	for scanner.Scan() {
		text := scanner.Text()
		text += "\n"
		pageCounter := lineCounter/args.pageLen + 1
		if pageCounter >= args.startPage && pageCounter <= args.endPage {
			*data += text
		}
		lineCounter++
	}
}

func selpgByF(reader *bufio.Reader, args *SelpgArgs, data *string, progname *string) {

	pageCounter := 1

	for {
		text, err := reader.ReadString('\f')
		if err == io.EOF {
			*data += text[:len(text)-1]
			break
		}
		if err != nil {
			fmt.Fprintf(os.Stderr, "%s: %s", *progname, err)
			os.Exit(4)
		}
		if pageCounter >= args.startPage && pageCounter <= args.endPage {
			*data += text[:len(text)-1]
		}
		if pageCounter > args.endPage {
			break
		}
		pageCounter++
	}
}
