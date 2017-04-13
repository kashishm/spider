package main

import (
	"fmt"

	"github.com/getgauge/spider/debug"
	"github.com/getgauge/spider/version"
)

func main() {
	fmt.Printf("Version: %s\n", version.Version)
	debug.Start()
}
