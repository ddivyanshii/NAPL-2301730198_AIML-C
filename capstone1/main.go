package main

import (
	"flag"
	"fmt"
	"os"

	"studentmgmt/service"
	"studentmgmt/store"
	"studentmgmt/ui"
)

func main() {
	dataFile := flag.String("data", "students.json", "path of the JSON data file")
	flag.Parse()

	svc, err := service.New(store.NewFileStore(*dataFile))
	if err != nil {
		fmt.Fprintln(os.Stderr, "Could not start:", err)
		os.Exit(1)
	}
	if err := ui.New(svc, os.Stdin, os.Stdout).Run(); err != nil {
		fmt.Fprintln(os.Stderr, "Fatal:", err)
		os.Exit(1)
	}
}
