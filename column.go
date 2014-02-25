package cfitsio

// #include <complex.h>
// #include "go-cfitsio.h"
// #include "go-cfitsio-utils.h"
import "C"
import "unsafe"

// Value is a value in a FITS table
type Value interface{}

// Column represents a column in a FITS table
type Column struct {
	Name    string  // column name, corresponding to ``TTYPE`` keyword
	Format  string  // column format, corresponding to ``TFORM`` keyword
	Unit    string  // column unit, corresponding to ``TUNIT`` keyword
	Null    string  // null value, corresponding to ``TNULL`` keyword
	Bscale  float64 // bscale value, corresponding to ``TSCAL`` keyword
	Bzero   float64 // bzero value, corresponding to ``TZERO`` keyword
	Display string  // display format, corresponding to ``TDISP`` keyword
	Dim     []int64 // column dimension corresponding to ``TDIM`` keyword
	Start   int64   // column starting position, corresponding to ``TBCOL`` keyword
	IsVLA   bool    // whether this is a variable length array
	Value   Value   // value at current row
}

func (col *Column) read(f *File, icol int, irow int64) error {
	var err error

	c_type := C.int(0)
	c_icol := C.int(icol + 1)      // 0-based to 1-based index
	c_irow := C.LONGLONG(irow + 1) // 0-based to 1-based index
	c_anynul := C.int(0)
	c_status := C.int(0)

	switch col.Value.(type) {
	case byte:
		c_type = C.TBYTE
		var c_value C.char
		c_ptr := unsafe.Pointer(&c_value)
		C.fits_read_col(f.c, c_type, c_icol, c_irow, 1, 1, c_ptr, c_ptr, &c_anynul, &c_status)
		col.Value = byte(c_value)

	case uint16:
		c_type = C.TUSHORT
		var c_value C.ushort
		c_ptr := unsafe.Pointer(&c_value)
		C.fits_read_col(f.c, c_type, c_icol, c_irow, 1, 1, c_ptr, c_ptr, &c_anynul, &c_status)
		col.Value = uint16(c_value)

	case uint32:
		c_type = C.TUINT
		var c_value C.uint
		c_ptr := unsafe.Pointer(&c_value)
		C.fits_read_col(f.c, c_type, c_icol, c_irow, 1, 1, c_ptr, c_ptr, &c_anynul, &c_status)
		col.Value = uint32(c_value)

	case uint64:
		c_type = C.TULONG
		var c_value C.ulong
		c_ptr := unsafe.Pointer(&c_value)
		C.fits_read_col(f.c, c_type, c_icol, c_irow, 1, 1, c_ptr, c_ptr, &c_anynul, &c_status)
		col.Value = uint64(c_value)

	case int8:
		c_type = C.TSBYTE
		var c_value C.char
		c_ptr := unsafe.Pointer(&c_value)
		C.fits_read_col(f.c, c_type, c_icol, c_irow, 1, 1, c_ptr, c_ptr, &c_anynul, &c_status)
		col.Value = int8(c_value)

	case int16:
		c_type = C.TSHORT
		var c_value C.short
		c_ptr := unsafe.Pointer(&c_value)
		C.fits_read_col(f.c, c_type, c_icol, c_irow, 1, 1, c_ptr, c_ptr, &c_anynul, &c_status)
		col.Value = int16(c_value)

	case int32:
		c_type = C.TINT
		var c_value C.int
		c_ptr := unsafe.Pointer(&c_value)
		C.fits_read_col(f.c, c_type, c_icol, c_irow, 1, 1, c_ptr, c_ptr, &c_anynul, &c_status)
		col.Value = int32(c_value)

	case int64:
		c_type = C.TLONG
		var c_value C.long
		c_ptr := unsafe.Pointer(&c_value)
		C.fits_read_col(f.c, c_type, c_icol, c_irow, 1, 1, c_ptr, c_ptr, &c_anynul, &c_status)
		col.Value = int64(c_value)

	case float32:
		c_type = C.TFLOAT
		var c_value C.float
		c_ptr := unsafe.Pointer(&c_value)
		C.fits_read_col(f.c, c_type, c_icol, c_irow, 1, 1, c_ptr, c_ptr, &c_anynul, &c_status)
		col.Value = float32(c_value)

	case float64:
		c_type = C.TDOUBLE
		var c_value C.double
		c_ptr := unsafe.Pointer(&c_value)
		C.fits_read_col(f.c, c_type, c_icol, c_irow, 1, 1, c_ptr, c_ptr, &c_anynul, &c_status)
		col.Value = float64(c_value)

	case complex64:
		c_type = C.TCOMPLEX
		var c_value C.complexfloat
		c_ptr := unsafe.Pointer(&c_value)
		C.fits_read_col(f.c, c_type, c_icol, c_irow, 1, 1, c_ptr, c_ptr, &c_anynul, &c_status)
		col.Value = complex(
			float32(C.crealf(c_value)),
			float32(C.cimagf(c_value)),
		)

	case complex128:
		c_type = C.TDBLCOMPLEX
		var c_value C.complexdouble
		c_ptr := unsafe.Pointer(&c_value)
		C.fits_read_col(f.c, c_type, c_icol, c_irow, 1, 1, c_ptr, c_ptr, &c_anynul, &c_status)
		col.Value = complex(
			float64(C.creal(c_value)),
			float64(C.cimag(c_value)),
		)

	case string:
		c_type = C.TSTRING
		// FIXME: get correct maximum size from card
		c_value := C.CStringN(C.FLEN_FILENAME)
		defer C.free(unsafe.Pointer(c_value))
		c_ptr := unsafe.Pointer(c_value)
		C.fits_read_col(f.c, c_type, c_icol, c_irow, 1, 1, c_ptr, c_ptr, &c_anynul, &c_status)
		col.Value = C.GoString(c_value)

	}

	if c_status > 0 {
		err = to_err(c_status)
	}

	return err
}

// ColDefs is a list of Column definitions
type ColDefs []Column

// EOF
