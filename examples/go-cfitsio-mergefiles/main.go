package main

import (
	"flag"
	"fmt"
	"os"
	"runtime/pprof"
	"time"

	fits "github.com/sbinet/go-cfitsio"
)

func main() {

	flag.Usage = func() {
		const msg = `Usage: go-cfitsio-mergefiles -o outfname file1 file2 [file3 ...]

Merge FITS tables into a single file/table.

`
		fmt.Fprintf(os.Stderr, "%v\n", msg)
		flag.PrintDefaults()
	}

	outfname := flag.String("o", "out.fits", "path to merged FITS file")
	doprofile := flag.Bool("profile", false, "enable CPU profiling")

	flag.Parse()
	if flag.NArg() < 2 {
		flag.Usage()
		os.Exit(1)
	}

	if *doprofile {
		f, err := os.Create("cpu.prof")
		if err != nil {
			panic(err)
		}
		defer f.Close()
		err = pprof.StartCPUProfile(f)
		if err != nil {
			panic(err)
		}
		defer pprof.StopCPUProfile()
	}
	_, err := os.Stat(*outfname)
	if err == nil {
		err = os.Remove(*outfname)
		if err != nil {
			panic(err)
		}
	}

	start := time.Now()
	defer func() {
		delta := time.Since(start)
		fmt.Printf("::: timing: %v\n", delta)
	}()

	fmt.Printf("::: creating merged file [%s]...\n", *outfname)
	out, err := fits.Create(*outfname)
	if err != nil {
		panic(err)
	}
	defer out.Close()

	infiles := make([]string, 0, flag.NArg())
	for i := 0; i < flag.NArg(); i++ {
		fname := flag.Arg(i)
		infiles = append(infiles, fname)
	}

	var table *fits.Table
	fmt.Printf("::: merging [%d] FITS files...\n", len(infiles))
	for i, fname := range infiles {
		f, err := fits.Open(fname, fits.ReadOnly)
		if err != nil {
			panic(err)
		}
		defer f.Close()
		hdu := f.HDU(1).(*fits.Table)
		nrows := hdu.NumRows()
		fmt.Printf("::: reading [%s] -> nrows=%d\n", fname, nrows)
		if i == 0 {
			// get header from first input file
			phdu, err := fits.NewPrimaryHDU(&out, f.HDU(0).Header())
			if err != nil {
				panic(err)
			}
			defer phdu.Close()

			// get schema from first input file
			cols := hdu.Cols()
			table, err = fits.NewTable(&out, hdu.Name(), cols, hdu.Type())
			if err != nil {
				panic(err)
			}
			defer table.Close()
		}

		err = fits.CopyTable(table, hdu)
		if err != nil {
			panic(err)
		}
	}
	fmt.Printf("::: merging [%d] FITS files... [done]\n", len(infiles))
	fmt.Printf("::: nrows: %d\n", table.NumRows())
}
