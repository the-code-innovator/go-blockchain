package main

import (
	"fmt"
	"os"

	"github.com/the-code-innovator/go-blockchain/line"
)

func main() {
	defer os.Exit(0)
	inter := line.Interface{}
	inter.Run()
	fmt.Println("exitting ....")
}
