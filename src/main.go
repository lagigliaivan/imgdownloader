package main

import (
	"flag"
	"fmt"
	"os"
	"strconv"
)

func main() {
	var amount string

	flag.StringVar(&amount, "amount", "10", "Default quantity 10")
	flag.Parse()

	q, err := strconv.Atoi(amount)
	if err != nil {
		goto exit
	}

	err = AppRun(q)

exit:
	if err != nil {
		fmt.Printf(err.Error())
		os.Exit(-1)
	}
}
