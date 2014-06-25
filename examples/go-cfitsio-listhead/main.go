package main

import (
	"flag"
	"fmt"
	"os"
	"strings"

	cfitsio "github.com/astrogo/cfitsio"
)

func main() {
	var single bool

	flag.Usage = func() {
		const msg = `Usage: go-cfitsio-listhead filename[ext]


List the FITS header keywords in a single extension, or, if 
ext is not given, list the keywords in all the extensions. 

Examples:
   
   go-cfitsio-listhead file.fits      - list every header in the file 
   go-cfitsio-listhead file.fits[0]   - list primary array header 
   go-cfitsio-listhead file.fits[2]   - list header of 2nd extension 
   go-cfitsio-listhead file.fits+2    - same as above 
   go-cfitsio-listhead file.fits[GTI] - list header of GTI extension

Note that it may be necessary to enclose the input file
name in single quote characters on the Unix command line.
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

	// list only a single header if a specific extension was given
	if ihdu != 0 || strings.Contains(fname, "[") {
		single = true
	}

	for i := ihdu; i < len(f.HDUs()); i++ {
		hdu := f.CHDU()
		hdr := hdu.Header()
		fmt.Printf("Header listing for HDU #%d:\n", i)

		for _, n := range hdr.Keys() {
			card := hdr.Get(n)
			if card == nil {
				panic(fmt.Errorf("could not retrieve card [%v]", n))
			}
			fmt.Printf(
				"%-8s= %-29s / %s\n",
				card.Name,
				fmt.Sprintf("%v", card.Value),
				card.Comment)
		}
		fmt.Printf("END\n\n")

		// quit if only listing a single header
		if single {
			break
		}
		err = f.SeekHDU(1, 0)
	}
	if err != nil && err != cfitsio.END_OF_FILE {
		panic(err)
	}
}
