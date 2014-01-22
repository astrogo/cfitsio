package main

import (
	"flag"
	"fmt"
	"os"
	"strings"

	cfitsio "github.com/sbinet/go-cfitsio"
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
	f, err := cfitsio.OpenFile(fname, cfitsio.ReadOnly)
	if err != nil {
		panic(err)
	}
	defer f.Close()

	hdu := f.HduNum() // get the current HDU position

	// list only a single header if a specific extension was given
	if hdu != 1 || strings.Contains(fname, "[") {
		single = true
	}

	for ; err == nil; hdu++ {
		var nkeys int
		// get number of keywords
		nkeys, _, err = f.HdrSpace()
		if err != nil {
			panic(err)
		}
		fmt.Printf("Header listing for HDU #%d:\n", hdu)

		for i := 1; i <= nkeys; i++ {
			var card string
			card, err = f.ReadRecord(i)
			if err != nil {
				panic(err)
			}
			fmt.Printf("%v\n", card)
		}
		fmt.Printf("END\n\n")

		// quit if only listing a single header
		if single {
			break
		}
		_, err = f.MovRelHdu(1)
	}
	if err != nil && err != cfitsio.END_OF_FILE {
		panic(err)
	}
}
