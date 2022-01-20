package main

import (
	"fmt"
	"os"
)

func main() {
	err := AppRun()

	if err != nil {
		fmt.Errorf(err.Error())
		os.Exit(-1)
	}
}
