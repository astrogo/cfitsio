go-cfitsio
==========

Naive CGo bindings for ``CFITSIO``.

## Installation

```sh
$ go get github.com/sbinet/go-cfitsio
```

You, of course, need the ``C`` library ``CFITSIO`` installed and available through ``pkg-config``.

## Documentation

http://godoc.org/github.com/sbinet/go-cfitsio

## Example

```go
import fits "github.com/sbinet/go-cfitsio"

func dumpFitsTable(fname string) {
	f, err := fits.Open(fname, fits.ReadOnly)
	if err != nil {
		panic(err)
	}
	defer f.Close()

	// get the second HDU
	table := f.HDUs()[1].(*fits.Table)
	nrows := table.NumRows()
	for i := int64(0); i < nrows; i++ {
		err := table.ReadRow(i)
		if err != nil {
			panic(err)
		}
		for icol := range table.Cols() {
			col := table.Col(icol)
			fmt.Printf("%s[%d]=%v\n", col.Name, i, col.Value)
		}
	}
}

```
