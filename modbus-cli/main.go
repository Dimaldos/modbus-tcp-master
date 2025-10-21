package main

import (
	"fmt"
	"modbus-cli/cmd"
	"os"
)

func main() {
	if err := cmd.Execute(); err != nil {
		fmt.Print(err.Error())
		os.Exit(1)
	}
}
