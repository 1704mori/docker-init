package main

import (
	"fmt"
	"os"

	"github.com/1704mori/docker-init/ui"
)

func main() {
	if err := ui.Run(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
