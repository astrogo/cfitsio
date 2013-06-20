package main

import (
	"fmt"

	"github.com/sbinet/go-cfitsio"
)

func main() {
	fmt.Printf("cfistio version=%v\n", cfitsio.Version())
}
