go-cfitsio
==========

[![Build Status](https://drone.io/github.com/sbinet/go-cfitsio/status.png)](https://drone.io/github.com/sbinet/go-cfitsio/latest)

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
	table := f.HDU(1).(*fits.Table)
	nrows := table.NumRows()
    rows, err := table.Read(0, nrows)
    if err != nil {
        panic(err)
    }
    defer rows.Close()
	for rows.Next() {
        var x, y float64
        var id int64
        err = rows.Scan(&id, &x, &y)
        if err != nil {
            panic(err)
        }
        fmt.Printf(">>> %v %v %v\n", id, x, y)
	}
    err = rows.Err()
    if err != nil { panic(err) }
}

```
