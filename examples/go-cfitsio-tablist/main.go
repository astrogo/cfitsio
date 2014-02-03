package main

import (
	"flag"
	"fmt"
	"os"

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
	hdu := f.HDUs()[ihdu]

	if hdu.Type() == cfitsio.IMAGE_HDU {
		fmt.Printf("Error: this program only displays tables, not images\n")
		os.Exit(1)
	}

	table := hdu.(*cfitsio.Table)
	nrows := table.NumRows()
	//ncols := table.NumCols()

	w := os.Stdout
	for icol, col := range table.Cols() {
		fmt.Fprintf(w, "\n[row] -+- %s (%T)\n", col.Name, col.Value)

		for i := int64(0); i < nrows; i++ {
			fmt.Fprintf(w, "%4d ", i)
			err := table.ReadRow(i)
			if err != nil {
				panic(err)
			}
			fmt.Fprintf(w, "  |  %-10v\n", table.Cols()[icol].Value)
		}
	}
}
