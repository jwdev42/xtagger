package main

import (
	"fmt"
	"github.com/jwdev42/xtagger/internal/program"
	"os"
)

func main() {
	if err := program.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "[FATAL] %s\n", err)
		os.Exit(1)
	}
}
