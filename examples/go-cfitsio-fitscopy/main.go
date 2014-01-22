package main

import (
	"flag"
	"fmt"
	"os"

	cfitsio "github.com/sbinet/go-cfitsio"
)

func main() {
	flag.Parse()
	if flag.NArg() != 2 {
		flag.Usage = func() {
			const msg = `Usage:  go-cfitsio-fitscopy inputfile outputfile

Copy an input file to an output file, optionally filtering
the file in the process.  This seemingly simple program can
apply powerful filters which transform the input file as
it is being copied.  Filters may be used to extract a
subimage from a larger image, select rows from a table,
filter a table with a GTI time extension or a SAO region file,
create or delete columns in a table, create an image by
binning (histogramming) 2 table columns, and convert IRAF
format *.imh or raw binary data files into FITS images.
See the CFITSIO User's Guide for a complete description of
the Extended File Name filtering syntax.

Examples:

go-cfitsio-fitscopy in.fit out.fit                   (simple file copy)
go-cfitsio-fitscopy - -                              (stdin to stdout)
go-cfitsio-fitscopy in.fit[11:50,21:60] out.fit      (copy a subimage)
go-cfitsio-fitscopy iniraf.imh out.fit               (IRAF image to FITS)
go-cfitsio-fitscopy in.dat[i512,512] out.fit         (raw array to FITS)
go-cfitsio-fitscopy in.fit[events][pi>35] out.fit    (copy rows with pi>35)
go-cfitsio-fitscopy in.fit[events][bin X,Y] out.fit  (bin an image) 
go-cfitsio-fitscopy in.fit[events][col x=.9*y] out.fit        (new x column)
go-cfitsio-fitscopy in.fit[events][gtifilter()] out.fit       (time filter)
go-cfitsio-fitscopy in.fit[2][regfilter(\"pow.reg\")] out.fit (spatial filter)

Note that it may be necessary to enclose the input file name
in single quote characters on the Unix command line.
`
			fmt.Fprintf(os.Stderr, "%v\n", msg)
		}
		flag.Usage()
		os.Exit(1)
	}

	// open input file
	in, err := cfitsio.OpenFile(flag.Arg(0), cfitsio.ReadOnly)
	if err != nil {
		panic(err)
	}
	defer in.Close()

	// create output file
	out, err := cfitsio.NewFile(flag.Arg(1))
	if err != nil {
		panic(err)
	}
	defer out.Close()

	// copy every HDU until we get an error
	for i := 1; err == nil; i++ {
		_, err = in.MovAbsHdu(i)
		if err != nil {
			break
		}
		err = in.CopyHdu(&out, 0)
	}
	if err != nil && err != cfitsio.END_OF_FILE {
		panic(err)
	}
}
