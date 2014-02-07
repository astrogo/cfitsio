package cfitsio

import (
	"database/sql"
	"database/sql/driver"
	"fmt"
)

// fitsdriver adapts a FITS table to the database/sql/driver interface
type fitsdriver struct {
}

// Open returns a new connection to the FITS table
func (drv *fitsdriver) Open(name string) (driver.Conn, error) {
	f, err := Open(name, ReadOnly)
	if err != nil {
		return nil, err
	}
	hdu := f.CHDU()
	tbl, ok := hdu.(*Table)
	if !ok {
		return nil, fmt.Errorf("cfitsio: current HDU isn't a Table")
	}
	conn := &fitsconn{
		f: f,
		t: tbl,
	}
	return conn, err
}

// fitsconn adapts a FITS table to the database/sql/driver Conn interface
type fitsconn struct {
	f File
	t *Table
}

// Prepare returns a prepared statement, bound to this connection
func (conn *fitsconn) Prepare(query string) (driver.Stmt, error) {
	var stmt driver.Stmt
	var err error

	return stmt, err
}

// Close invalidates and potentially stops any current prepared statements
// and transactions, marking this connection as no longer in use.
func (conn *fitsconn) Close() error {
	err := conn.t.Close()
	if err != nil {
		return err
	}
	err = conn.f.Close()
	return err
}

// Begin starts and returns a new transaction
func (conn *fitsconn) Begin() (driver.Tx, error) {
	var tx driver.Tx = &fitstx{conn}
	var err error
	return tx, err
}

// fitstx is a transaction on a FITS table
type fitstx struct {
	conn *fitsconn
}

func (tx *fitstx) Commit() error {
	if tx.conn == nil {
		return fmt.Errorf("cfitsio: invalid FITS connection")
	}
	panic("not implemented")
}

func (tx *fitstx) Rollback() error {
	if tx.conn == nil {
		return fmt.Errorf("cfitsio: invalid FITS connection")
	}
	panic("not implemented")
}

func init() {
	sql.Register("fits", &fitsdriver{})
}

var _ driver.Driver = (*fitsdriver)(nil)

// EOF
