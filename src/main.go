package main

import (
	"fmt"
	"os"
)

func main() {
	err := AppRun()

	if err != nil {
		fmt.Printf(err.Error())
		os.Exit(-1)
	}
}
