package cfitsio

// #include <complex.h>
// #include "go-cfitsio.h"
// #include "go-cfitsio-utils.h"
import "C"
import (
	"fmt"
	"reflect"
	"unsafe"
)

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

func (col *Column) inferFormat(htype HDUType) error {
	var err error
	if col.Format != "" {
		return nil
	}

	str := gotype2FITS(col.Value, htype)
	if str == "" {
		return fmt.Errorf("cfitsio: %v can not handle [%T]", htype, col.Value)
	}
	col.Format = str
	return err
}

func (col *Column) read(f *File, icol int, irow int64, ptr interface{}) error {
	var err error

	c_type := C.int(0)
	c_icol := C.int(icol + 1)      // 0-based to 1-based index
	c_irow := C.LONGLONG(irow + 1) // 0-based to 1-based index
	c_anynul := C.int(0)
	c_status := C.int(0)

	rv := reflect.ValueOf(ptr).Elem()
	value := rv.Interface()
	switch value.(type) {
	case bool:
		c_type = C.TLOGICAL
		c_value := C.char(0) // 'F'
		c_ptr := unsafe.Pointer(&c_value)
		C.fits_read_col(f.c, c_type, c_icol, c_irow, 1, 1, c_ptr, c_ptr, &c_anynul, &c_status)
		value = c_value == 1

	case byte:
		c_type = C.TBYTE
		var c_value C.char
		c_ptr := unsafe.Pointer(&c_value)
		C.fits_read_col(f.c, c_type, c_icol, c_irow, 1, 1, c_ptr, c_ptr, &c_anynul, &c_status)
		value = byte(c_value)

	case uint16:
		c_type = C.TUSHORT
		var c_value C.ushort
		c_ptr := unsafe.Pointer(&c_value)
		C.fits_read_col(f.c, c_type, c_icol, c_irow, 1, 1, c_ptr, c_ptr, &c_anynul, &c_status)
		value = uint16(c_value)

	case uint32:
		c_type = C.TUINT
		var c_value C.uint
		c_ptr := unsafe.Pointer(&c_value)
		C.fits_read_col(f.c, c_type, c_icol, c_irow, 1, 1, c_ptr, c_ptr, &c_anynul, &c_status)
		value = uint32(c_value)

	case uint64:
		c_type = C.TULONG
		var c_value C.ulong
		c_ptr := unsafe.Pointer(&c_value)
		C.fits_read_col(f.c, c_type, c_icol, c_irow, 1, 1, c_ptr, c_ptr, &c_anynul, &c_status)
		value = uint64(c_value)

	case uint:
		c_type = C.TULONG
		var c_value C.ulong
		c_ptr := unsafe.Pointer(&c_value)
		C.fits_read_col(f.c, c_type, c_icol, c_irow, 1, 1, c_ptr, c_ptr, &c_anynul, &c_status)
		value = uint(c_value)

	case int8:
		c_type = C.TSBYTE
		var c_value C.char
		c_ptr := unsafe.Pointer(&c_value)
		C.fits_read_col(f.c, c_type, c_icol, c_irow, 1, 1, c_ptr, c_ptr, &c_anynul, &c_status)
		value = int8(c_value)

	case int16:
		c_type = C.TSHORT
		var c_value C.short
		c_ptr := unsafe.Pointer(&c_value)
		C.fits_read_col(f.c, c_type, c_icol, c_irow, 1, 1, c_ptr, c_ptr, &c_anynul, &c_status)
		value = int16(c_value)

	case int32:
		c_type = C.TINT
		var c_value C.int
		c_ptr := unsafe.Pointer(&c_value)
		C.fits_read_col(f.c, c_type, c_icol, c_irow, 1, 1, c_ptr, c_ptr, &c_anynul, &c_status)
		value = int32(c_value)

	case int64:
		c_type = C.TLONG
		var c_value C.long
		c_ptr := unsafe.Pointer(&c_value)
		C.fits_read_col(f.c, c_type, c_icol, c_irow, 1, 1, c_ptr, c_ptr, &c_anynul, &c_status)
		value = int64(c_value)

	case int:
		c_type = C.TLONG
		var c_value C.long
		c_ptr := unsafe.Pointer(&c_value)
		C.fits_read_col(f.c, c_type, c_icol, c_irow, 1, 1, c_ptr, c_ptr, &c_anynul, &c_status)
		value = int(c_value)

	case float32:
		c_type = C.TFLOAT
		var c_value C.float
		c_ptr := unsafe.Pointer(&c_value)
		C.fits_read_col(f.c, c_type, c_icol, c_irow, 1, 1, c_ptr, c_ptr, &c_anynul, &c_status)
		value = float32(c_value)

	case float64:
		c_type = C.TDOUBLE
		var c_value C.double
		c_ptr := unsafe.Pointer(&c_value)
		C.fits_read_col(f.c, c_type, c_icol, c_irow, 1, 1, c_ptr, c_ptr, &c_anynul, &c_status)
		value = float64(c_value)

	case complex64:
		c_type = C.TCOMPLEX
		var c_value C.complexfloat
		c_ptr := unsafe.Pointer(&c_value)
		C.fits_read_col(f.c, c_type, c_icol, c_irow, 1, 1, c_ptr, c_ptr, &c_anynul, &c_status)
		value = complex(
			float32(C.crealf(c_value)),
			float32(C.cimagf(c_value)),
		)

	case complex128:
		c_type = C.TDBLCOMPLEX
		var c_value C.complexdouble
		c_ptr := unsafe.Pointer(&c_value)
		C.fits_read_col(f.c, c_type, c_icol, c_irow, 1, 1, c_ptr, c_ptr, &c_anynul, &c_status)
		value = complex(
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
		value = C.GoString(c_value)

	case []bool:
		c_type = C.TLOGICAL
		c_len := C.long(0)
		c_off := C.long(0)
		C.fits_read_descript(f.c, c_icol, c_irow, &c_len, &c_off, &c_status)
		if c_status > 0 {
			err = to_err(c_status)
			return err
		}
		v := make([]bool, int(c_len), int(c_len))
		value = v
		slice := (*reflect.SliceHeader)((unsafe.Pointer(&v)))
		c_ptr := unsafe.Pointer(slice.Data)
		C.fits_read_col(f.c, c_type, c_icol, c_irow, 1, C.LONGLONG(c_len), nil, c_ptr, &c_anynul, &c_status)

	case []uint8:
		c_type = C.TBYTE
		c_len := C.long(0)
		c_off := C.long(0)
		C.fits_read_descript(f.c, c_icol, c_irow, &c_len, &c_off, &c_status)
		if c_status > 0 {
			err = to_err(c_status)
			return err
		}
		v := make([]uint8, int(c_len), int(c_len))
		value = v
		slice := (*reflect.SliceHeader)((unsafe.Pointer(&v)))
		c_ptr := unsafe.Pointer(slice.Data)
		C.fits_read_col(f.c, c_type, c_icol, c_irow, 1, C.LONGLONG(c_len), nil, c_ptr, &c_anynul, &c_status)

	case []uint16:
		c_type = C.TUSHORT
		c_len := C.long(0)
		c_off := C.long(0)
		C.fits_read_descript(f.c, c_icol, c_irow, &c_len, &c_off, &c_status)
		if c_status > 0 {
			err = to_err(c_status)
			return err
		}
		v := make([]uint16, int(c_len), int(c_len))
		value = v
		slice := (*reflect.SliceHeader)((unsafe.Pointer(&v)))
		c_ptr := unsafe.Pointer(slice.Data)
		C.fits_read_col(f.c, c_type, c_icol, c_irow, 1, C.LONGLONG(c_len), nil, c_ptr, &c_anynul, &c_status)

	case []uint32:
		c_type = C.TUINT
		c_len := C.long(0)
		c_off := C.long(0)
		C.fits_read_descript(f.c, c_icol, c_irow, &c_len, &c_off, &c_status)
		if c_status > 0 {
			err = to_err(c_status)
			return err
		}
		v := make([]uint32, int(c_len), int(c_len))
		value = v
		slice := (*reflect.SliceHeader)((unsafe.Pointer(&v)))
		c_ptr := unsafe.Pointer(slice.Data)
		C.fits_read_col(f.c, c_type, c_icol, c_irow, 1, C.LONGLONG(c_len), nil, c_ptr, &c_anynul, &c_status)

	case []uint64:
		c_type = C.TULONG
		c_len := C.long(0)
		c_off := C.long(0)
		C.fits_read_descript(f.c, c_icol, c_irow, &c_len, &c_off, &c_status)
		if c_status > 0 {
			err = to_err(c_status)
			return err
		}
		v := make([]uint64, int(c_len), int(c_len))
		value = v
		slice := (*reflect.SliceHeader)((unsafe.Pointer(&v)))
		c_ptr := unsafe.Pointer(slice.Data)
		C.fits_read_col(f.c, c_type, c_icol, c_irow, 1, C.LONGLONG(c_len), nil, c_ptr, &c_anynul, &c_status)

	case []uint:
		c_type = C.TULONG
		c_len := C.long(0)
		c_off := C.long(0)
		C.fits_read_descript(f.c, c_icol, c_irow, &c_len, &c_off, &c_status)
		if c_status > 0 {
			err = to_err(c_status)
			return err
		}
		v := make([]uint, int(c_len), int(c_len))
		value = v
		slice := (*reflect.SliceHeader)((unsafe.Pointer(&v)))
		c_ptr := unsafe.Pointer(slice.Data)
		C.fits_read_col(f.c, c_type, c_icol, c_irow, 1, C.LONGLONG(c_len), nil, c_ptr, &c_anynul, &c_status)

	case []int8:
		c_type = C.TSBYTE
		c_len := C.long(0)
		c_off := C.long(0)
		C.fits_read_descript(f.c, c_icol, c_irow, &c_len, &c_off, &c_status)
		if c_status > 0 {
			err = to_err(c_status)
			return err
		}
		v := make([]int8, int(c_len), int(c_len))
		value = v
		slice := (*reflect.SliceHeader)((unsafe.Pointer(&v)))
		c_ptr := unsafe.Pointer(slice.Data)
		C.fits_read_col(f.c, c_type, c_icol, c_irow, 1, C.LONGLONG(c_len), nil, c_ptr, &c_anynul, &c_status)

	case []int16:
		c_type = C.TSHORT
		c_len := C.long(0)
		c_off := C.long(0)
		C.fits_read_descript(f.c, c_icol, c_irow, &c_len, &c_off, &c_status)
		if c_status > 0 {
			err = to_err(c_status)
			return err
		}
		v := make([]int16, int(c_len), int(c_len))
		value = v
		slice := (*reflect.SliceHeader)((unsafe.Pointer(&v)))
		c_ptr := unsafe.Pointer(slice.Data)
		C.fits_read_col(f.c, c_type, c_icol, c_irow, 1, C.LONGLONG(c_len), nil, c_ptr, &c_anynul, &c_status)

	case []int32:
		c_type = C.TINT
		c_len := C.long(0)
		c_off := C.long(0)
		C.fits_read_descript(f.c, c_icol, c_irow, &c_len, &c_off, &c_status)
		if c_status > 0 {
			err = to_err(c_status)
			return err
		}
		v := make([]int32, int(c_len), int(c_len))
		value = v
		slice := (*reflect.SliceHeader)((unsafe.Pointer(&v)))
		c_ptr := unsafe.Pointer(slice.Data)
		C.fits_read_col(f.c, c_type, c_icol, c_irow, 1, C.LONGLONG(c_len), nil, c_ptr, &c_anynul, &c_status)

	case []int64:
		c_type = C.TLONG
		c_len := C.long(0)
		c_off := C.long(0)
		C.fits_read_descript(f.c, c_icol, c_irow, &c_len, &c_off, &c_status)
		if c_status > 0 {
			err = to_err(c_status)
			return err
		}
		v := make([]int64, int(c_len), int(c_len))
		value = v
		slice := (*reflect.SliceHeader)((unsafe.Pointer(&v)))
		c_ptr := unsafe.Pointer(slice.Data)
		C.fits_read_col(f.c, c_type, c_icol, c_irow, 1, C.LONGLONG(c_len), nil, c_ptr, &c_anynul, &c_status)

	case []int:
		c_type = C.TLONG
		c_len := C.long(0)
		c_off := C.long(0)
		C.fits_read_descript(f.c, c_icol, c_irow, &c_len, &c_off, &c_status)
		if c_status > 0 {
			err = to_err(c_status)
			return err
		}
		v := make([]int, int(c_len), int(c_len))
		value = v
		slice := (*reflect.SliceHeader)((unsafe.Pointer(&v)))
		c_ptr := unsafe.Pointer(slice.Data)
		C.fits_read_col(f.c, c_type, c_icol, c_irow, 1, C.LONGLONG(c_len), nil, c_ptr, &c_anynul, &c_status)

	case []float32:
		c_type = C.TFLOAT
		c_len := C.long(0)
		c_off := C.long(0)
		C.fits_read_descript(f.c, c_icol, c_irow, &c_len, &c_off, &c_status)
		if c_status > 0 {
			err = to_err(c_status)
			return err
		}
		v := make([]float32, int(c_len), int(c_len))
		value = v
		slice := (*reflect.SliceHeader)((unsafe.Pointer(&v)))
		c_ptr := unsafe.Pointer(slice.Data)
		C.fits_read_col(f.c, c_type, c_icol, c_irow, 1, C.LONGLONG(c_len), nil, c_ptr, &c_anynul, &c_status)

	case []float64:
		c_type = C.TDOUBLE
		c_len := C.long(0)
		c_off := C.long(0)
		C.fits_read_descript(f.c, c_icol, c_irow, &c_len, &c_off, &c_status)
		if c_status > 0 {
			err = to_err(c_status)
			return err
		}
		v := make([]float64, int(c_len), int(c_len))
		value = v
		slice := (*reflect.SliceHeader)((unsafe.Pointer(&v)))
		c_ptr := unsafe.Pointer(slice.Data)
		C.fits_read_col(f.c, c_type, c_icol, c_irow, 1, C.LONGLONG(c_len), nil, c_ptr, &c_anynul, &c_status)

	case []complex64:
		c_type = C.TCOMPLEX
		c_len := C.long(0)
		c_off := C.long(0)
		C.fits_read_descript(f.c, c_icol, c_irow, &c_len, &c_off, &c_status)
		if c_status > 0 {
			err = to_err(c_status)
			return err
		}
		v := make([]complex64, int(c_len), int(c_len)) // FIXME: assume same binary layout
		value = v
		slice := (*reflect.SliceHeader)((unsafe.Pointer(&v)))
		c_ptr := unsafe.Pointer(slice.Data)
		C.fits_read_col(f.c, c_type, c_icol, c_irow, 1, C.LONGLONG(c_len), nil, c_ptr, &c_anynul, &c_status)

	case []complex128:
		c_type = C.TDBLCOMPLEX
		c_len := C.long(0)
		c_off := C.long(0)
		C.fits_read_descript(f.c, c_icol, c_irow, &c_len, &c_off, &c_status)
		if c_status > 0 {
			err = to_err(c_status)
			return err
		}
		v := make([]complex128, int(c_len), int(c_len)) // FIXME: assume same binary layout
		value = v
		slice := (*reflect.SliceHeader)((unsafe.Pointer(&v)))
		c_ptr := unsafe.Pointer(slice.Data)
		C.fits_read_col(f.c, c_type, c_icol, c_irow, 1, C.LONGLONG(c_len), nil, c_ptr, &c_anynul, &c_status)

	default:
		panic(fmt.Errorf("invalid type [%T]", value))
	}

	if c_status > 0 {
		err = to_err(c_status)
	}

	rv.Set(reflect.ValueOf(value))
	col.Value = value
	return err
}

func (col *Column) write(f *File, icol int, irow int64) error {
	var err error

	c_type := C.int(0)
	c_icol := C.int(icol + 1)      // 0-based to 1-based index
	c_irow := C.LONGLONG(irow + 1) // 0-based to 1-based index
	c_status := C.int(0)

	switch value := col.Value.(type) {
	case bool:
		c_type = C.TLOGICAL
		c_value := C.char(0) // 'F'
		if value {
			c_value = 1
		}
		c_ptr := unsafe.Pointer(&c_value)
		C.fits_write_col(f.c, c_type, c_icol, c_irow, 1, 1, c_ptr, &c_status)

	case byte:
		c_type = C.TBYTE
		c_value := C.char(value)
		c_ptr := unsafe.Pointer(&c_value)
		C.fits_write_col(f.c, c_type, c_icol, c_irow, 1, 1, c_ptr, &c_status)

	case uint16:
		c_type = C.TUSHORT
		c_value := C.ushort(value)
		c_ptr := unsafe.Pointer(&c_value)
		C.fits_write_col(f.c, c_type, c_icol, c_irow, 1, 1, c_ptr, &c_status)

	case uint32:
		c_type = C.TUINT
		c_value := C.uint(value)
		c_ptr := unsafe.Pointer(&c_value)
		C.fits_write_col(f.c, c_type, c_icol, c_irow, 1, 1, c_ptr, &c_status)

	case uint64:
		c_type = C.TULONG
		c_value := C.ulong(value)
		c_ptr := unsafe.Pointer(&c_value)
		C.fits_write_col(f.c, c_type, c_icol, c_irow, 1, 1, c_ptr, &c_status)

	case uint:
		c_type = C.TULONG
		c_value := C.ulong(value)
		c_ptr := unsafe.Pointer(&c_value)
		C.fits_write_col(f.c, c_type, c_icol, c_irow, 1, 1, c_ptr, &c_status)

	case int8:
		c_type = C.TSBYTE
		c_value := C.char(value)
		c_ptr := unsafe.Pointer(&c_value)
		C.fits_write_col(f.c, c_type, c_icol, c_irow, 1, 1, c_ptr, &c_status)

	case int16:
		c_type = C.TSHORT
		c_value := C.short(value)
		c_ptr := unsafe.Pointer(&c_value)
		C.fits_write_col(f.c, c_type, c_icol, c_irow, 1, 1, c_ptr, &c_status)

	case int32:
		c_type = C.TINT
		c_value := C.int(value)
		c_ptr := unsafe.Pointer(&c_value)
		C.fits_write_col(f.c, c_type, c_icol, c_irow, 1, 1, c_ptr, &c_status)

	case int64:
		c_type = C.TLONG
		c_value := C.long(value)
		c_ptr := unsafe.Pointer(&c_value)
		C.fits_write_col(f.c, c_type, c_icol, c_irow, 1, 1, c_ptr, &c_status)

	case int:
		c_type = C.TLONG
		c_value := C.long(value)
		c_ptr := unsafe.Pointer(&c_value)
		C.fits_write_col(f.c, c_type, c_icol, c_irow, 1, 1, c_ptr, &c_status)

	case float32:
		c_type = C.TFLOAT
		c_value := C.float(value)
		c_ptr := unsafe.Pointer(&c_value)
		C.fits_write_col(f.c, c_type, c_icol, c_irow, 1, 1, c_ptr, &c_status)

	case float64:
		c_type = C.TDOUBLE
		c_value := C.double(value)
		c_ptr := unsafe.Pointer(&c_value)
		C.fits_write_col(f.c, c_type, c_icol, c_irow, 1, 1, c_ptr, &c_status)

	case complex64:
		c_type = C.TCOMPLEX
		c_ptr := unsafe.Pointer(&value) // FIXME: assumes same memory layout
		C.fits_write_col(f.c, c_type, c_icol, c_irow, 1, 1, c_ptr, &c_status)

	case complex128:
		c_type = C.TDBLCOMPLEX
		c_ptr := unsafe.Pointer(&value) // FIXME: assumes same memory layout
		C.fits_write_col(f.c, c_type, c_icol, c_irow, 1, 1, c_ptr, &c_status)

	case string:
		c_type = C.TSTRING
		c_value := C.CString(value)
		defer C.free(unsafe.Pointer(c_value))
		c_ptr := unsafe.Pointer(c_value)
		C.fits_write_col(f.c, c_type, c_icol, c_irow, 1, 1, c_ptr, &c_status)

	case []bool:
		c_type = C.TLOGICAL
		slice := (*reflect.SliceHeader)((unsafe.Pointer(&value)))
		c_ptr := unsafe.Pointer(slice.Data)
		C.fits_write_col(f.c, c_type, c_icol, c_irow, 1, C.LONGLONG(slice.Len), c_ptr, &c_status)

	case []uint8:
		c_type = C.TBYTE
		slice := (*reflect.SliceHeader)((unsafe.Pointer(&value)))
		c_ptr := unsafe.Pointer(slice.Data)
		C.fits_write_col(f.c, c_type, c_icol, c_irow, 1, C.LONGLONG(slice.Len), c_ptr, &c_status)

	case []uint16:
		c_type = C.TUSHORT
		slice := (*reflect.SliceHeader)((unsafe.Pointer(&value)))
		c_ptr := unsafe.Pointer(slice.Data)
		C.fits_write_col(f.c, c_type, c_icol, c_irow, 1, C.LONGLONG(slice.Len), c_ptr, &c_status)

	case []uint32:
		c_type = C.TUINT
		slice := (*reflect.SliceHeader)((unsafe.Pointer(&value)))
		c_ptr := unsafe.Pointer(slice.Data)
		C.fits_write_col(f.c, c_type, c_icol, c_irow, 1, C.LONGLONG(slice.Len), c_ptr, &c_status)

	case []uint64:
		c_type = C.TULONG
		slice := (*reflect.SliceHeader)((unsafe.Pointer(&value)))
		c_ptr := unsafe.Pointer(slice.Data)
		C.fits_write_col(f.c, c_type, c_icol, c_irow, 1, C.LONGLONG(slice.Len), c_ptr, &c_status)

	case []uint:
		c_type = C.TULONG
		slice := (*reflect.SliceHeader)((unsafe.Pointer(&value)))
		c_ptr := unsafe.Pointer(slice.Data)
		C.fits_write_col(f.c, c_type, c_icol, c_irow, 1, C.LONGLONG(slice.Len), c_ptr, &c_status)

	case []int8:
		c_type = C.TSBYTE
		slice := (*reflect.SliceHeader)((unsafe.Pointer(&value)))
		c_ptr := unsafe.Pointer(slice.Data)
		C.fits_write_col(f.c, c_type, c_icol, c_irow, 1, C.LONGLONG(slice.Len), c_ptr, &c_status)

	case []int16:
		c_type = C.TSHORT
		slice := (*reflect.SliceHeader)((unsafe.Pointer(&value)))
		c_ptr := unsafe.Pointer(slice.Data)
		C.fits_write_col(f.c, c_type, c_icol, c_irow, 1, C.LONGLONG(slice.Len), c_ptr, &c_status)

	case []int32:
		c_type = C.TINT
		slice := (*reflect.SliceHeader)((unsafe.Pointer(&value)))
		c_ptr := unsafe.Pointer(slice.Data)
		C.fits_write_col(f.c, c_type, c_icol, c_irow, 1, C.LONGLONG(slice.Len), c_ptr, &c_status)

	case []int64:
		c_type = C.TLONG
		slice := (*reflect.SliceHeader)((unsafe.Pointer(&value)))
		c_ptr := unsafe.Pointer(slice.Data)
		C.fits_write_col(f.c, c_type, c_icol, c_irow, 1, C.LONGLONG(slice.Len), c_ptr, &c_status)

	case []int:
		c_type = C.TLONG
		slice := (*reflect.SliceHeader)((unsafe.Pointer(&value)))
		c_ptr := unsafe.Pointer(slice.Data)
		C.fits_write_col(f.c, c_type, c_icol, c_irow, 1, C.LONGLONG(slice.Len), c_ptr, &c_status)

	case []float32:
		c_type = C.TFLOAT
		slice := (*reflect.SliceHeader)((unsafe.Pointer(&value)))
		c_ptr := unsafe.Pointer(slice.Data)
		C.fits_write_col(f.c, c_type, c_icol, c_irow, 1, C.LONGLONG(slice.Len), c_ptr, &c_status)

	case []float64:
		c_type = C.TDOUBLE
		slice := (*reflect.SliceHeader)((unsafe.Pointer(&value)))
		c_ptr := unsafe.Pointer(slice.Data)
		C.fits_write_col(f.c, c_type, c_icol, c_irow, 1, C.LONGLONG(slice.Len), c_ptr, &c_status)

	case []complex64:
		c_type = C.TCOMPLEX
		slice := (*reflect.SliceHeader)((unsafe.Pointer(&value))) // FIXME: assume same bin-layout
		c_ptr := unsafe.Pointer(slice.Data)
		C.fits_write_col(f.c, c_type, c_icol, c_irow, 1, C.LONGLONG(slice.Len), c_ptr, &c_status)

	case []complex128:
		c_type = C.TDBLCOMPLEX
		slice := (*reflect.SliceHeader)((unsafe.Pointer(&value))) // FIXME: assume same bin-layout
		c_ptr := unsafe.Pointer(slice.Data)
		C.fits_write_col(f.c, c_type, c_icol, c_irow, 1, C.LONGLONG(slice.Len), c_ptr, &c_status)

	default:
		panic(fmt.Errorf("unhandled type [%T]", col.Value))
	}

	if c_status > 0 {
		err = to_err(c_status)
	}

	return err
}

// ColDefs is a list of Column definitions
type ColDefs []Column

// EOF
