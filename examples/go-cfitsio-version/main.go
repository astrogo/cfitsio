package main

import (
	"fmt"

	"github.com/astrogo/cfitsio"
)

func main() {
	fmt.Printf("cfistio version=%v\n", cfitsio.Version())
}
