package cfitsio

import (
	"fmt"
)

// Rows is the result of a query on a FITS Table.
// Its cursors starts before the first row of the result set.
// Use Next to advance through the rows:
//
//  rows, err := table.Read(0, -1)
//  ...
//  for rows.Next() {
//      var id int
//      var x float64
//      err = rows.Scan(&id, &x)
//      ...
//  }
//  err = rows.Err() // get any error encountered during iteration
//  ...
//
type Rows struct {
	table  *Table
	cols   []int // list of (active) column indices
	nrows  int64 // number of rows this iterator iters over
	irow   int64 // current row index
	closed bool
	err    error // last error
}

// Err returns the error, if any, that was encountered during iteration.
// Err may be called after an explicit or implicit Close.
func (rows *Rows) Err() error {
	return rows.err
}

// Close closes the Rows, preventing further enumeration.
// Close is idempotent and does not affect the result of Err.
func (rows *Rows) Close() error {
	if rows.closed {
		return nil
	}
	rows.closed = true
	rows.table = nil
	return nil
}

// Scan copies the columns in the current row into the values pointed at by
// dest.
func (rows *Rows) Scan(args ...interface{}) error {
	var err error
	defer func() {
		rows.err = err
	}()

	if len(args) != len(rows.cols) {
		return fmt.Errorf(
			"cfitsio.Rows.Scan: invalid number of arguments (got %d. expected %d)",
			len(args),
			len(rows.cols),
		)
	}
	for i, icol := range rows.cols {
		err = rows.table.cols[icol].read(rows.table.f, icol, rows.irow)
		if err != nil {
			return err
		}
		args[i] = rows.table.cols[icol].Value
	}
	return err
}

// Next prepares the next result row for reading with the Scan method.
// It returns true on success, false if there is no next result row.
// Every call to Scan, even the first one, must be preceded by a call to Next.
func (rows *Rows) Next() bool {
	if rows.closed {
		return false
	}
	rows.irow += 1
	next := rows.irow < rows.nrows
	if !next {
		rows.err = rows.Close()
	}
	return next
}

// EOF
