package main

import (
	"flag"
	"fmt"
	"os"
	"strings"

	cfitsio "github.com/sbinet/go-cfitsio"
)

func main() {

	flag.Usage = func() {
		const msg = `Usage: go-cfitsio-tablist filename[ext][col filter][row filter]

List the contents of a FITS table

Examples:
  tablist tab.fits[GTI]           - list the GTI extension
  tablist tab.fits[1][#row < 101] - list first 100 rows
  tablist tab.fits[1][col X;Y]    - list X and Y cols only
  tablist tab.fits[1][col -PI]    - list all but the PI col
  tablist tab.fits[1][col -PI][#row < 101]  - combined case

Display formats can be modified with the TDISPn keywords.
`
		fmt.Fprintf(os.Stderr, "%v\n", msg)
		flag.PrintDefaults()
	}

	flag.Parse()
	if flag.NArg() != 1 {
		flag.Usage()
		os.Exit(1)
	}

	fname := flag.Arg(0)
	f, err := cfitsio.Open(fname, cfitsio.ReadOnly)
	if err != nil {
		panic(err)
	}
	defer f.Close()

	ihdu := f.HDUNum() // get the current HDU position
	if ihdu == 0 {
		// this is the primary array.
		// try to move to the first extension and see if it is a table
		ihdu = 1
	}
	if ihdu >= len(f.HDUs()) {
		fmt.Printf("Error: input file has no extension\n")
		os.Exit(1)
	}
	hdu := f.HDU(ihdu)

	if hdu.Type() == cfitsio.IMAGE_HDU {
		fmt.Printf("Error: this program only displays tables, not images\n")
		os.Exit(1)
	}

	table := hdu.(*cfitsio.Table)
	nrows := table.NumRows()
	rows, err := table.Read(0, nrows)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}
	w := os.Stdout
	hdrline := strings.Repeat("=", 80-15)
	maxname := 10
	for _, col := range table.Cols() {
		if len(col.Name) > maxname {
			maxname = len(col.Name)
		}
	}
	rowfmt := fmt.Sprintf("%%-%ds | %%v\n", maxname)
	for irow := 0; rows.Next(); irow++ {
		err = rows.Scan()
		if err != nil {
			fmt.Printf("Error: (row=%v) %v\n", irow, err)
		}
		fmt.Fprintf(w, "== %05d/%05d %s\n", irow, nrows, hdrline)
		for _, col := range table.Cols() {
			fmt.Fprintf(w, rowfmt, col.Name, col.Value)
		}
	}

	err = rows.Err()
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}
}
