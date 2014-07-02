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
	Type    TypeCode
	Len     int   // repeat. if <= 1: scalar
	Value   Value // value at current row
}

// inferFormat infers the FITS format associated with a Column, according to its HDUType and Go type.
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

// read reads the value at column number icol and row irow, into ptr.
// icol and irow are 0-based indices.
func (col *Column) read(f *File, icol int, irow int64, ptr interface{}) error {
	var err error

	c_type := C.int(0)
	c_icol := C.int(icol + 1)      // 0-based to 1-based index
	c_irow := C.LONGLONG(irow + 1) // 0-based to 1-based index

	c_len := C.LONGLONG(col.Len)
	if col.Len <= 1 {
		c_len = 1
	}

	var value interface{}
	rv := reflect.Indirect(reflect.ValueOf(ptr))
	rt := reflect.TypeOf(rv.Interface())

	decode := func(c_type C.int, c_len C.LONGLONG, rt reflect.Type, scalar bool) (interface{}, error) {
		var value interface{}
		c_status := C.int(0)
		c_anynul := C.int(0)

		switch rt.Kind() {
		case reflect.Bool:
			c_type = C.TLOGICAL
			v := make([]C.char, int(c_len), int(c_len))
			slice := (*reflect.SliceHeader)((unsafe.Pointer(&v)))
			c_ptr := unsafe.Pointer(slice.Data)
			C.fits_read_col(f.c, c_type, c_icol, c_irow, 1, c_len, nil, c_ptr, &c_anynul, &c_status)
			if scalar {
				value = v[0] == 1
			} else {
				vv := make([]bool, len(v), len(v))
				for ci, cv := range v {
					vv[ci] = cv == 1
				}
				value = vv
			}

		case reflect.Uint8:
			c_type = C.TBYTE
			v := make([]uint8, int(c_len), int(c_len))
			value = v
			slice := (*reflect.SliceHeader)((unsafe.Pointer(&v)))
			c_ptr := unsafe.Pointer(slice.Data)
			C.fits_read_col(f.c, c_type, c_icol, c_irow, 1, c_len, nil, c_ptr, &c_anynul, &c_status)
			if scalar {
				value = v[0]
			}

		case reflect.Uint16:
			c_type = C.TUSHORT
			v := make([]uint16, int(c_len), int(c_len))
			value = v
			slice := (*reflect.SliceHeader)((unsafe.Pointer(&v)))
			c_ptr := unsafe.Pointer(slice.Data)
			C.fits_read_col(f.c, c_type, c_icol, c_irow, 1, c_len, nil, c_ptr, &c_anynul, &c_status)
			if scalar {
				value = v[0]
			}

		case reflect.Uint32:
			c_type = C.TUINT
			v := make([]uint32, int(c_len), int(c_len))
			value = v
			slice := (*reflect.SliceHeader)((unsafe.Pointer(&v)))
			c_ptr := unsafe.Pointer(slice.Data)
			C.fits_read_col(f.c, c_type, c_icol, c_irow, 1, c_len, nil, c_ptr, &c_anynul, &c_status)
			if scalar {
				value = v[0]
			}

		case reflect.Uint64:
			c_type = C.TULONG
			v := make([]uint64, int(c_len), int(c_len))
			value = v
			slice := (*reflect.SliceHeader)((unsafe.Pointer(&v)))
			c_ptr := unsafe.Pointer(slice.Data)
			C.fits_read_col(f.c, c_type, c_icol, c_irow, 1, c_len, nil, c_ptr, &c_anynul, &c_status)
			if scalar {
				value = v[0]
			}

		case reflect.Uint:
			c_type = C.TULONG
			v := make([]uint, int(c_len), int(c_len))
			value = v
			slice := (*reflect.SliceHeader)((unsafe.Pointer(&v)))
			c_ptr := unsafe.Pointer(slice.Data)
			C.fits_read_col(f.c, c_type, c_icol, c_irow, 1, c_len, nil, c_ptr, &c_anynul, &c_status)
			if scalar {
				value = v[0]
			}

		case reflect.Int8:
			c_type = C.TSBYTE
			v := make([]int8, int(c_len), int(c_len))
			value = v
			slice := (*reflect.SliceHeader)((unsafe.Pointer(&v)))
			c_ptr := unsafe.Pointer(slice.Data)
			C.fits_read_col(f.c, c_type, c_icol, c_irow, 1, c_len, nil, c_ptr, &c_anynul, &c_status)
			if scalar {
				value = v[0]
			}

		case reflect.Int16:
			c_type = C.TSHORT
			v := make([]int16, int(c_len), int(c_len))
			value = v
			slice := (*reflect.SliceHeader)((unsafe.Pointer(&v)))
			c_ptr := unsafe.Pointer(slice.Data)
			C.fits_read_col(f.c, c_type, c_icol, c_irow, 1, c_len, nil, c_ptr, &c_anynul, &c_status)
			if scalar {
				value = v[0]
			}

		case reflect.Int32:
			c_type = C.TINT
			v := make([]int32, int(c_len), int(c_len))
			value = v
			slice := (*reflect.SliceHeader)((unsafe.Pointer(&v)))
			c_ptr := unsafe.Pointer(slice.Data)
			C.fits_read_col(f.c, c_type, c_icol, c_irow, 1, c_len, nil, c_ptr, &c_anynul, &c_status)
			if scalar {
				value = v[0]
			}

		case reflect.Int64:
			c_type = C.TLONG
			v := make([]int64, int(c_len), int(c_len))
			value = v
			slice := (*reflect.SliceHeader)((unsafe.Pointer(&v)))
			c_ptr := unsafe.Pointer(slice.Data)
			C.fits_read_col(f.c, c_type, c_icol, c_irow, 1, c_len, nil, c_ptr, &c_anynul, &c_status)
			if scalar {
				value = v[0]
			}

		case reflect.Int:
			c_type = C.TLONG
			v := make([]int, int(c_len), int(c_len))
			value = v
			slice := (*reflect.SliceHeader)((unsafe.Pointer(&v)))
			c_ptr := unsafe.Pointer(slice.Data)
			C.fits_read_col(f.c, c_type, c_icol, c_irow, 1, c_len, nil, c_ptr, &c_anynul, &c_status)
			if scalar {
				value = v[0]
			}

		case reflect.Float32:
			c_type = C.TFLOAT
			v := make([]float32, int(c_len), int(c_len))
			value = v
			slice := (*reflect.SliceHeader)((unsafe.Pointer(&v)))
			c_ptr := unsafe.Pointer(slice.Data)
			C.fits_read_col(f.c, c_type, c_icol, c_irow, 1, c_len, nil, c_ptr, &c_anynul, &c_status)
			if scalar {
				value = v[0]
			}

		case reflect.Float64:
			c_type = C.TDOUBLE
			v := make([]float64, int(c_len), int(c_len))
			value = v
			slice := (*reflect.SliceHeader)((unsafe.Pointer(&v)))
			c_ptr := unsafe.Pointer(slice.Data)
			C.fits_read_col(f.c, c_type, c_icol, c_irow, 1, c_len, nil, c_ptr, &c_anynul, &c_status)
			if scalar {
				value = v[0]
			}

		case reflect.Complex64:
			c_type = C.TCOMPLEX
			v := make([]complex64, int(c_len), int(c_len)) // FIXME: assume same binary layout
			value = v
			slice := (*reflect.SliceHeader)((unsafe.Pointer(&v)))
			c_ptr := unsafe.Pointer(slice.Data)
			C.fits_read_col(f.c, c_type, c_icol, c_irow, 1, c_len, nil, c_ptr, &c_anynul, &c_status)
			if scalar {
				value = v[0]
			}

		case reflect.Complex128:
			c_type = C.TDBLCOMPLEX
			v := make([]complex128, int(c_len), int(c_len)) // FIXME: assume same binary layout
			value = v
			slice := (*reflect.SliceHeader)((unsafe.Pointer(&v)))
			c_ptr := unsafe.Pointer(slice.Data)
			C.fits_read_col(f.c, c_type, c_icol, c_irow, 1, c_len, nil, c_ptr, &c_anynul, &c_status)
			if scalar {
				value = v[0]
			}

		case reflect.String:
			c_type = C.TSTRING
			// FIXME: get correct maximum size from card
			c_value := C.CStringN(C.FLEN_FILENAME)
			defer C.free(unsafe.Pointer(c_value))
			c_ptr := unsafe.Pointer(&c_value)
			C.fits_read_col(f.c, c_type, c_icol, c_irow, 1, 1, c_ptr, c_ptr, &c_anynul, &c_status)
			value = C.GoString(c_value)

		default:
			return value, fmt.Errorf("cfitsio: invalid type [%v]", rt)
		}

		if c_status > 0 {
			err = to_err(c_status)
			return value, err
		}
		return value, nil
	}

	switch rt.Kind() {

	case reflect.Slice:
		c_len := C.long(0)
		if col.Type < 0 {
			c_off := C.long(0)
			c_status := C.int(0)
			C.fits_read_descript(f.c, c_icol, c_irow, &c_len, &c_off, &c_status)
			if c_status > 0 {
				err = to_err(c_status)
				return err
			}
		} else {
			c_len = C.long(col.Len)
		}
		scalar := false
		value, err = decode(c_type, C.LONGLONG(c_len), rt.Elem(), scalar)

		rv.Set(reflect.ValueOf(value))
		col.Value = rv.Interface()

	case reflect.Array:
		c_len := C.long(rt.Len())
		scalar := false
		value, err = decode(c_type, C.LONGLONG(c_len), rt.Elem(), scalar)

		// FIXME: unnecessary copy
		array := reflect.New(rt).Elem()
		reflect.Copy(array, reflect.ValueOf(value))
		rv.Set(array)
		col.Value = rv.Interface()

	case reflect.Bool,
		reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
		reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64,
		reflect.Float32, reflect.Float64,
		reflect.Complex64, reflect.Complex128,
		reflect.String:

		scalar := true
		value, err = decode(c_type, C.LONGLONG(c_len), rt, scalar)

		rv.Set(reflect.ValueOf(value))
		col.Value = rv.Interface()

	default:
		return fmt.Errorf("cfitsio: invalid type [%v]", rt)
	}

	return err
}

// write writes the current value of this Column into file f at column icol and row irow.
// icol and irow are 0-based indices.
func (col *Column) write(f *File, icol int, irow int64, value interface{}) error {
	var err error

	c_type := C.int(0)
	c_icol := C.int(icol + 1)      // 0-based to 1-based index
	c_irow := C.LONGLONG(irow + 1) // 0-based to 1-based index
	c_status := C.int(0)

	rv := reflect.ValueOf(value)
	rt := reflect.TypeOf(value)

	switch rt.Kind() {
	case reflect.Bool:
		c_type = C.TLOGICAL
		c_value := C.char(0) // 'F'
		if value.(bool) {
			c_value = 1
		}
		c_ptr := unsafe.Pointer(&c_value)
		C.fits_write_col(f.c, c_type, c_icol, c_irow, 1, 1, c_ptr, &c_status)

	case reflect.Uint8:
		c_type = C.TBYTE
		c_value := C.char(value.(byte))
		c_ptr := unsafe.Pointer(&c_value)
		C.fits_write_col(f.c, c_type, c_icol, c_irow, 1, 1, c_ptr, &c_status)

	case reflect.Uint16:
		c_type = C.TUSHORT
		c_value := C.ushort(value.(uint16))
		c_ptr := unsafe.Pointer(&c_value)
		C.fits_write_col(f.c, c_type, c_icol, c_irow, 1, 1, c_ptr, &c_status)

	case reflect.Uint32:
		c_type = C.TUINT
		c_value := C.uint(value.(uint32))
		c_ptr := unsafe.Pointer(&c_value)
		C.fits_write_col(f.c, c_type, c_icol, c_irow, 1, 1, c_ptr, &c_status)

	case reflect.Uint64:
		c_type = C.TULONG
		c_value := C.ulong(value.(uint64))
		c_ptr := unsafe.Pointer(&c_value)
		C.fits_write_col(f.c, c_type, c_icol, c_irow, 1, 1, c_ptr, &c_status)

	case reflect.Uint:
		c_type = C.TULONG
		c_value := C.ulong(value.(uint))
		c_ptr := unsafe.Pointer(&c_value)
		C.fits_write_col(f.c, c_type, c_icol, c_irow, 1, 1, c_ptr, &c_status)

	case reflect.Int8:
		c_type = C.TSBYTE
		c_value := C.char(value.(int8))
		c_ptr := unsafe.Pointer(&c_value)
		C.fits_write_col(f.c, c_type, c_icol, c_irow, 1, 1, c_ptr, &c_status)

	case reflect.Int16:
		c_type = C.TSHORT
		c_value := C.short(value.(int16))
		c_ptr := unsafe.Pointer(&c_value)
		C.fits_write_col(f.c, c_type, c_icol, c_irow, 1, 1, c_ptr, &c_status)

	case reflect.Int32:
		c_type = C.TINT
		c_value := C.int(value.(int32))
		c_ptr := unsafe.Pointer(&c_value)
		C.fits_write_col(f.c, c_type, c_icol, c_irow, 1, 1, c_ptr, &c_status)

	case reflect.Int64:
		c_type = C.TLONG
		c_value := C.long(value.(int64))
		c_ptr := unsafe.Pointer(&c_value)
		C.fits_write_col(f.c, c_type, c_icol, c_irow, 1, 1, c_ptr, &c_status)

	case reflect.Int:
		c_type = C.TLONG
		c_value := C.long(value.(int))
		c_ptr := unsafe.Pointer(&c_value)
		C.fits_write_col(f.c, c_type, c_icol, c_irow, 1, 1, c_ptr, &c_status)

	case reflect.Float32:
		c_type = C.TFLOAT
		c_value := C.float(value.(float32))
		c_ptr := unsafe.Pointer(&c_value)
		C.fits_write_col(f.c, c_type, c_icol, c_irow, 1, 1, c_ptr, &c_status)

	case reflect.Float64:
		c_type = C.TDOUBLE
		c_value := C.double(value.(float64))
		c_ptr := unsafe.Pointer(&c_value)
		C.fits_write_col(f.c, c_type, c_icol, c_irow, 1, 1, c_ptr, &c_status)

	case reflect.Complex64:
		c_type = C.TCOMPLEX
		value := value.(complex64)
		c_ptr := unsafe.Pointer(&value) // FIXME: assumes same memory layout
		C.fits_write_col(f.c, c_type, c_icol, c_irow, 1, 1, c_ptr, &c_status)

	case reflect.Complex128:
		c_type = C.TDBLCOMPLEX
		value := value.(complex128)
		c_ptr := unsafe.Pointer(&value) // FIXME: assumes same memory layout
		C.fits_write_col(f.c, c_type, c_icol, c_irow, 1, 1, c_ptr, &c_status)

	case reflect.String:
		c_type = C.TSTRING
		c_value := C.CString(value.(string))
		defer C.free(unsafe.Pointer(c_value))
		c_ptr := unsafe.Pointer(&c_value)
		C.fits_write_col(f.c, c_type, c_icol, c_irow, 1, 1, c_ptr, &c_status)

	case reflect.Slice:
		switch rt.Elem().Kind() {
		case reflect.Bool:
			c_type = C.TLOGICAL
			value := value.([]bool)
			slice := (*reflect.SliceHeader)((unsafe.Pointer(&value)))
			c_ptr := unsafe.Pointer(slice.Data)
			C.fits_write_col(f.c, c_type, c_icol, c_irow, 1, C.LONGLONG(slice.Len), c_ptr, &c_status)

		case reflect.Uint8:
			c_type = C.TBYTE
			value := value.([]uint8)
			slice := (*reflect.SliceHeader)((unsafe.Pointer(&value)))
			c_ptr := unsafe.Pointer(slice.Data)
			C.fits_write_col(f.c, c_type, c_icol, c_irow, 1, C.LONGLONG(slice.Len), c_ptr, &c_status)

		case reflect.Uint16:
			c_type = C.TUSHORT
			value := value.([]uint16)
			slice := (*reflect.SliceHeader)((unsafe.Pointer(&value)))
			c_ptr := unsafe.Pointer(slice.Data)
			C.fits_write_col(f.c, c_type, c_icol, c_irow, 1, C.LONGLONG(slice.Len), c_ptr, &c_status)

		case reflect.Uint32:
			c_type = C.TUINT
			value := value.([]uint32)
			slice := (*reflect.SliceHeader)((unsafe.Pointer(&value)))
			c_ptr := unsafe.Pointer(slice.Data)
			C.fits_write_col(f.c, c_type, c_icol, c_irow, 1, C.LONGLONG(slice.Len), c_ptr, &c_status)

		case reflect.Uint64:
			c_type = C.TULONG
			value := value.([]uint64)
			slice := (*reflect.SliceHeader)((unsafe.Pointer(&value)))
			c_ptr := unsafe.Pointer(slice.Data)
			C.fits_write_col(f.c, c_type, c_icol, c_irow, 1, C.LONGLONG(slice.Len), c_ptr, &c_status)

		case reflect.Uint:
			c_type = C.TULONG
			value := value.([]uint)
			slice := (*reflect.SliceHeader)((unsafe.Pointer(&value)))
			c_ptr := unsafe.Pointer(slice.Data)
			C.fits_write_col(f.c, c_type, c_icol, c_irow, 1, C.LONGLONG(slice.Len), c_ptr, &c_status)

		case reflect.Int8:
			c_type = C.TSBYTE
			value := value.([]int8)
			slice := (*reflect.SliceHeader)((unsafe.Pointer(&value)))
			c_ptr := unsafe.Pointer(slice.Data)
			C.fits_write_col(f.c, c_type, c_icol, c_irow, 1, C.LONGLONG(slice.Len), c_ptr, &c_status)

		case reflect.Int16:
			c_type = C.TSHORT
			value := value.([]int16)
			slice := (*reflect.SliceHeader)((unsafe.Pointer(&value)))
			c_ptr := unsafe.Pointer(slice.Data)
			C.fits_write_col(f.c, c_type, c_icol, c_irow, 1, C.LONGLONG(slice.Len), c_ptr, &c_status)

		case reflect.Int32:
			c_type = C.TINT
			value := value.([]int32)
			slice := (*reflect.SliceHeader)((unsafe.Pointer(&value)))
			c_ptr := unsafe.Pointer(slice.Data)
			C.fits_write_col(f.c, c_type, c_icol, c_irow, 1, C.LONGLONG(slice.Len), c_ptr, &c_status)

		case reflect.Int64:
			c_type = C.TLONG
			value := value.([]int64)
			slice := (*reflect.SliceHeader)((unsafe.Pointer(&value)))
			c_ptr := unsafe.Pointer(slice.Data)
			C.fits_write_col(f.c, c_type, c_icol, c_irow, 1, C.LONGLONG(slice.Len), c_ptr, &c_status)

		case reflect.Int:
			c_type = C.TLONG
			value := value.([]int)
			slice := (*reflect.SliceHeader)((unsafe.Pointer(&value)))
			c_ptr := unsafe.Pointer(slice.Data)
			C.fits_write_col(f.c, c_type, c_icol, c_irow, 1, C.LONGLONG(slice.Len), c_ptr, &c_status)

		case reflect.Float32:
			c_type = C.TFLOAT
			value := value.([]float32)
			slice := (*reflect.SliceHeader)((unsafe.Pointer(&value)))
			c_ptr := unsafe.Pointer(slice.Data)
			C.fits_write_col(f.c, c_type, c_icol, c_irow, 1, C.LONGLONG(slice.Len), c_ptr, &c_status)

		case reflect.Float64:
			c_type = C.TDOUBLE
			value := value.([]float64)
			slice := (*reflect.SliceHeader)((unsafe.Pointer(&value)))
			c_ptr := unsafe.Pointer(slice.Data)
			C.fits_write_col(f.c, c_type, c_icol, c_irow, 1, C.LONGLONG(slice.Len), c_ptr, &c_status)

		case reflect.Complex64:
			c_type = C.TCOMPLEX
			value := value.([]complex64)
			slice := (*reflect.SliceHeader)((unsafe.Pointer(&value))) // FIXME: assume same bin-layout
			c_ptr := unsafe.Pointer(slice.Data)
			C.fits_write_col(f.c, c_type, c_icol, c_irow, 1, C.LONGLONG(slice.Len), c_ptr, &c_status)

		case reflect.Complex128:
			c_type = C.TDBLCOMPLEX
			value := value.([]complex128)
			slice := (*reflect.SliceHeader)((unsafe.Pointer(&value))) // FIXME: assume same bin-layout
			c_ptr := unsafe.Pointer(slice.Data)
			C.fits_write_col(f.c, c_type, c_icol, c_irow, 1, C.LONGLONG(slice.Len), c_ptr, &c_status)
		default:
			panic(fmt.Errorf("unhandled type '%T'", value))
		}

	case reflect.Array:

		// FIXME: unnecessary copy
		arr := reflect.New(rv.Type()).Elem()
		reflect.Copy(arr, rv)
		ptr := arr.Addr()
		c_ptr := unsafe.Pointer(ptr.Pointer())
		c_len := C.LONGLONG(rt.Len())

		switch rt.Elem().Kind() {
		case reflect.Bool:
			c_type = C.TLOGICAL
			C.fits_write_col(f.c, c_type, c_icol, c_irow, 1, c_len, c_ptr, &c_status)

		case reflect.Uint8:
			c_type = C.TBYTE
			C.fits_write_col(f.c, c_type, c_icol, c_irow, 1, c_len, c_ptr, &c_status)

		case reflect.Uint16:
			c_type = C.TUSHORT
			C.fits_write_col(f.c, c_type, c_icol, c_irow, 1, c_len, c_ptr, &c_status)

		case reflect.Uint32:
			c_type = C.TUINT
			C.fits_write_col(f.c, c_type, c_icol, c_irow, 1, c_len, c_ptr, &c_status)

		case reflect.Uint64:
			c_type = C.TULONG
			C.fits_write_col(f.c, c_type, c_icol, c_irow, 1, c_len, c_ptr, &c_status)

		case reflect.Uint:
			c_type = C.TULONG
			C.fits_write_col(f.c, c_type, c_icol, c_irow, 1, c_len, c_ptr, &c_status)

		case reflect.Int8:
			c_type = C.TSBYTE
			C.fits_write_col(f.c, c_type, c_icol, c_irow, 1, c_len, c_ptr, &c_status)

		case reflect.Int16:
			c_type = C.TSHORT
			C.fits_write_col(f.c, c_type, c_icol, c_irow, 1, c_len, c_ptr, &c_status)

		case reflect.Int32:
			c_type = C.TINT
			C.fits_write_col(f.c, c_type, c_icol, c_irow, 1, c_len, c_ptr, &c_status)

		case reflect.Int64:
			c_type = C.TLONG
			C.fits_write_col(f.c, c_type, c_icol, c_irow, 1, c_len, c_ptr, &c_status)

		case reflect.Int:
			c_type = C.TLONG
			C.fits_write_col(f.c, c_type, c_icol, c_irow, 1, c_len, c_ptr, &c_status)

		case reflect.Float32:
			c_type = C.TFLOAT
			C.fits_write_col(f.c, c_type, c_icol, c_irow, 1, c_len, c_ptr, &c_status)

		case reflect.Float64:
			c_type = C.TDOUBLE
			C.fits_write_col(f.c, c_type, c_icol, c_irow, 1, c_len, c_ptr, &c_status)

		case reflect.Complex64:
			c_type = C.TCOMPLEX
			C.fits_write_col(f.c, c_type, c_icol, c_irow, 1, c_len, c_ptr, &c_status)

		case reflect.Complex128:
			c_type = C.TDBLCOMPLEX
			C.fits_write_col(f.c, c_type, c_icol, c_irow, 1, c_len, c_ptr, &c_status)
		default:
			panic(fmt.Errorf("unhandled type '%T'", value))
		}

	default:
		panic(fmt.Errorf("unhandled type '%T' (%v)", value, rt.Kind()))
	}

	if c_status > 0 {
		err = to_err(c_status)
	}

	return err
}

// EOF
