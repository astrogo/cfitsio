package main

import (
	"flag"
	"fmt"
	"os"

	fits "github.com/astrogo/cfitsio"
)

func main() {
	flag.Usage = func() {
		const msg = `Usage: go-cfitsio-imlist filename[ext][section filter] 

List the the pixel values in a FITS image 

Example: 
  imlist image.fits                    - list the whole image
  imlist image.fits[100:110,400:410]   - list a section
  imlist table.fits[2][bin (x,y) = 32] - list the pixels in
         an image constructed from a 2D histogram of X and Y
         columns in a table with a binning factor = 32
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
	f, err := fits.Open(fname, fits.ReadOnly)
	if err != nil {
		panic(err)
	}
	defer f.Close()

	hdu := f.CHDU()
	hdr := hdu.Header()

	if len(hdr.Axes()) > 2 || len(hdr.Axes()) == 0 {
		fmt.Printf("Error: only 1D or 2D images are supported\n")
		os.Exit(1)
	}

	axes := hdr.Axes()

	var data []float64
	err = hdu.Data(&data)
	if err != nil {
		panic(err)
	}

	// default output format string
	hdformat := " %15d"
	format := " %15.5f"
	if hdr.Bitpix() > 0 {
		hdformat = " %7d"
		format = " %7.0f"
	}
	// column header
	fmt.Printf("\n      ")
	row := int(axes[0])
	for ii := 0; ii < row; ii++ {
		fmt.Printf(hdformat, ii)
	}
	fmt.Printf("\n")

	// loop over all rows
	for jj := 0; jj < int(axes[1]); jj++ {
		fmt.Printf(" %4d ", jj)
		for ii := 0; ii < row; ii++ {
			fmt.Printf(format, data[row*jj+ii])
		}

		fmt.Printf("\n")
	}
	if err != nil && err != fits.END_OF_FILE {
		panic(err)
	}
}
