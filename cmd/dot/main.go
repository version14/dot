package main

import (
	"fmt"
	"os"
)

func main() {
	if err := run(os.Args[1:]); err != nil {
		fmt.Fprintf(os.Stderr, "dot: %v\n", err)
		os.Exit(1)
	}
}

func run(args []string) error {
	if len(args) == 0 {
		printUsage()
		return nil
	}

	switch args[0] {
	case "init":
		return cmdInit()
	case "help":
		return cmdHelp()
	case "self-update", "update":
		return cmdSelfUpdate()
	case "version", "--version", "-v":
		fmt.Printf("dot %s\n", buildVersion)
		return nil
	default:
		return fmt.Errorf("unknown command %q — run 'dot help' for usage", args[0])
	}
}

func printUsage() {
	fmt.Print(`dot — universal project companion

Usage:
  dot init          scaffold a new project
  dot help          list available commands
  dot version       print version
  dot self-update   update dot to the latest release
`)
}
