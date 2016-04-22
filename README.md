cfitsio
=======

[![Build Status](https://drone.io/github.com/astrogo/cfitsio/status.png)](https://drone.io/github.com/astrogo/cfitsio/latest)

Naive CGo bindings for ``FITSIO``.

## Installation

```sh
$ go get github.com/astrogo/cfitsio
```

You, of course, need the ``C`` library ``CFITSIO`` installed and available through ``pkg-config``.

## Documentation

http://godoc.org/github.com/astrogo/cfitsio

## Example

```go
import fits "github.com/astrogo/cfitsio"

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
    
    // using a struct
    xx := struct{
        Id int     `fits:"ID"`
        X  float64 `fits:"x"`
        Y  float64 `fits:"y"`
    }{}
    // using a map
    yy := make(map[string]interface{})
    
    rows, err = table.Read(0, nrows)
    if err != nil {
        panic(err)
    }
    defer rows.Close()
	for rows.Next() {
        err = rows.Scan(&xx)
        if err != nil {
            panic(err)
        }
        fmt.Printf(">>> %v\n", xx)

        err = rows.Scan(&yy)
        if err != nil {
            panic(err)
        }
        fmt.Printf(">>> %v\n", yy)
	}
    err = rows.Err()
    if err != nil { panic(err) }
    
}

```

## TODO

- ``[DONE]`` add support for writing tables from structs
- ``[DONE]`` add support for writing tables from maps
- ``[DONE]`` add support for variable length array
- provide benchmarks _wrt_ ``CFITSIO``

## Contribute

`astrogo/cfitsio` is released under the `BSD-3` license.
Please send a pull request to [astrogo/license](https://github.com/astrogo/license), adding
yourself to the `AUTHORS` and/or `CONTRIBUTORS` files.
