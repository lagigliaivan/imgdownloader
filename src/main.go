package main

import (
	"flag"
	"fmt"
	"os"
)

func main() {
	var amount int

	flag.IntVar(&amount, "amount", 10, "Default quantity 10")
	flag.Parse()

	err := AppRun(amount)
	if err != nil {
		fmt.Printf(err.Error())
		os.Exit(-1)
	}
}
