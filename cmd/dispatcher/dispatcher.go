package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/Arheon/markus-backend/internal"
)

func main() {
	env := flag.String("e", "dev", "")
	flag.Usage = func() {
		fmt.Println("Usage: server -e {mode}")
		os.Exit(1)
	}
	flag.Parse()

	err := internal.InitDispatcher(*env)
	if err != nil {
		log.Fatal(err)
	}
}
