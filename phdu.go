package cfitsio

// #include "go-cfitsio.h"
// #include "go-cfitsio-utils.h"
import "C"

import (
	"fmt"
	"reflect"
	"unsafe"
)

// PrimaryHDU is the primary HDU
type PrimaryHDU struct {
	ImageHDU
}

func (hdu *PrimaryHDU) Name() string {
	card := hdu.header.Get("EXTNAME")
	if card == nil {
		return "PRIMARY"
	}
	return card.Value.(string)
}

func (hdu *PrimaryHDU) Version() int {
	card := hdu.header.Get("EXTVER")
	if card == nil {
		return 1
	}
	rv := reflect.ValueOf(card.Value)
	return int(rv.Int())
}

func newPrimaryHDU(f *File, hdr Header) (HDU, error) {
	var err error
	hdu := &PrimaryHDU{
		ImageHDU{
			f:      f,
			header: hdr,
		},
	}
	return hdu, err
}

func NewPrimaryHDU(f *File, hdr Header) (HDU, error) {
	var err error

	naxes := len(hdr.axes)
	c_naxes := C.int(naxes)
	c_axes := C.long_array_new(c_naxes)
	defer C.free(unsafe.Pointer(c_axes))
	c_status := C.int(0)

	C.fits_create_img(f.c, C.int(hdr.bitpix), c_naxes, c_axes, &c_status)
	if c_status > 0 {
		return nil, to_err(c_status)
	}

	for icard := range hdr.slice {
		card := &hdr.slice[icard]
		c_name := C.CString(card.Name)
		defer C.free(unsafe.Pointer(c_name))
		c_type := C.int(0)
		c_status := C.int(0)
		c_comm := C.CString(card.Comment)
		defer C.free(unsafe.Pointer(c_comm))
		var c_ptr unsafe.Pointer

		switch v := card.Value.(type) {
		case byte:
			c_type = C.TBYTE
			c_ptr = unsafe.Pointer(&v)

		case uint16:
			c_type = C.TUSHORT
			c_ptr = unsafe.Pointer(&v)

		case uint32:
			c_type = C.TUINT
			c_ptr = unsafe.Pointer(&v)

		case uint64:
			c_type = C.TULONG
			c_ptr = unsafe.Pointer(&v)

		case int8:
			c_type = C.TSBYTE
			c_ptr = unsafe.Pointer(&v)

		case int16:
			c_type = C.TSHORT
			c_ptr = unsafe.Pointer(&v)

		case int32:
			c_type = C.TINT
			c_ptr = unsafe.Pointer(&v)

		case int64:
			c_type = C.TLONG
			c_ptr = unsafe.Pointer(&v)

		case float32:
			c_type = C.TFLOAT
			c_ptr = unsafe.Pointer(&v)

		case float64:
			c_type = C.TDOUBLE
			c_ptr = unsafe.Pointer(&v)

		case complex64:
			c_type = C.TCOMPLEX
			c_ptr = unsafe.Pointer(&v) // FIXME: assumes same memory layout than C

		case complex128:
			c_type = C.TDBLCOMPLEX
			c_ptr = unsafe.Pointer(&v) // FIXME: assumes same memory layout than C

		case string:
			c_type = C.TSTRING
			c_value := C.CString(v)
			defer C.free(unsafe.Pointer(c_value))
			c_ptr = unsafe.Pointer(&v)

		default:
			panic(fmt.Errorf("cfitsio: invalid card type (%T)", v))
		}

		C.fits_update_key(f.c, c_type, c_name, c_ptr, c_comm, &c_status)

		if c_status > 0 {
			return nil, to_err(c_status)
		}
	}

	if len(f.hdus) > 0 {
		return nil, fmt.Errorf("cfitsio: File has already a Primary HDU")
	}

	hdu, err := f.readHDU(0)
	if err != nil {
		return nil, err
	}
	f.hdus = append(f.hdus, hdu)

	return hdu, err
}

// EOF
