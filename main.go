package main

import (
	"fmt"
	"os"
	"sando/cmd/root"
)

func main() {
	rootCmd := root.NewCmdRoot()
	if _, err := rootCmd.ExecuteC(); err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err)
		os.Exit(1)
	}
}
