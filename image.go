package cfitsio

// #include "go-cfitsio.h"
// #include "go-cfitsio-utils.h"
import "C"
import (
	"fmt"
	"reflect"
	"unsafe"
)

// ImageHDU is a Header-Data-Unit extension holding an image as data payload.
type ImageHDU struct {
	f      *File
	header Header
}

// Close closes this HDU, cleaning up cycles for the proper garbage collection.
func (hdu *ImageHDU) Close() error {
	hdu.f = nil
	return nil
}

// Header returns the Header part of this "Header Data-Unit" block.
func (hdu *ImageHDU) Header() Header {
	return hdu.header
}

// Type returns the HDUType for this HDU.
func (hdu *ImageHDU) Type() HDUType {
	return hdu.header.htype
}

// Name returns the value of the 'EXTNAME' Card (or "" if none)
func (hdu *ImageHDU) Name() string {
	card := hdu.header.Get("EXTNAME")
	if card == nil {
		return ""
	}
	return card.Value.(string)
}

// Version returns the value of the 'EXTVER' Card (or 1 if none)
func (hdu *ImageHDU) Version() int {
	card := hdu.header.Get("EXTVER")
	if card == nil {
		return 1
	}
	return card.Value.(int)
}

// Data loads the image data associated with this HDU into data, which should
// be a pointer to a slice []T.
// cfitsio will return an error if the image payload can not be converted into Ts.
// It panics if data isn't addressable.
func (hdu *ImageHDU) Data(data interface{}) error {

	rv := reflect.ValueOf(data).Elem()
	if !rv.CanAddr() {
		return fmt.Errorf("%T is not addressable", data)
	}

	err := hdu.load(rv)
	return err
}

// load loads the image data associated with this HDU into v.
func (hdu *ImageHDU) load(v reflect.Value) error {
	var err error
	hdr := hdu.Header()
	naxes := len(hdr.Axes())
	if naxes == 0 {
		return nil
	}
	nelmts := 1
	for _, dim := range hdr.Axes() {
		nelmts *= int(dim)
	}
	rv := reflect.MakeSlice(v.Type(), nelmts, nelmts)

	c_start := C.LONGLONG(0)
	c_nelmts := C.LONGLONG(nelmts)
	c_anynull := C.int(0)
	c_status := C.int(0)
	c_imgtype := C.int(0)
	var c_ptr unsafe.Pointer
	switch rv.Interface().(type) {
	case []byte:
		c_imgtype = C.TBYTE
		data := rv.Interface().([]byte)
		c_ptr = unsafe.Pointer(&data[0])

	case []int8:
		c_imgtype = C.TBYTE
		data := rv.Interface().([]int8)
		c_ptr = unsafe.Pointer(&data[0])

	case []int16:
		c_imgtype = C.TSHORT
		data := rv.Interface().([]int16)
		c_ptr = unsafe.Pointer(&data[0])

	case []uint16:
		c_imgtype = C.TUSHORT
		data := rv.Interface().([]uint16)
		c_ptr = unsafe.Pointer(&data[0])

	case []int32:
		c_imgtype = C.TINT
		data := rv.Interface().([]int32)
		c_ptr = unsafe.Pointer(&data[0])

	case []int64:
		c_imgtype = C.TLONGLONG
		data := rv.Interface().([]int64)
		c_ptr = unsafe.Pointer(&data[0])

	case []float32:
		c_imgtype = C.TFLOAT
		data := rv.Interface().([]float32)
		c_ptr = unsafe.Pointer(&data[0])

	case []float64:
		c_imgtype = C.TDOUBLE
		data := rv.Interface().([]float64)
		c_ptr = unsafe.Pointer(&data[0])

	default:
		panic(fmt.Errorf("invalid image type [%T]", rv.Interface()))
	}
	C.fits_read_img(hdu.f.c, c_imgtype, c_start+1, c_nelmts, c_ptr, c_ptr, &c_anynull, &c_status)
	if c_status > 0 {
		return to_err(c_status)
	}

	n := reflect.Copy(v, rv)
	if n != nelmts {
		err = fmt.Errorf("cfitsio: copied [%v] elements. expected [%v]", n, nelmts)
	}
	return err
}

// Write writes the image to disk
// data should be a pointer to a slice []T.
func (hdu *ImageHDU) Write(data interface{}) error {
	var err error
	rv := reflect.ValueOf(data).Elem()
	if !rv.CanAddr() {
		return fmt.Errorf("%T is not addressable", data)
	}

	hdr := hdu.Header()
	naxes := len(hdr.Axes())
	if naxes == 0 {
		return nil
	}
	nelmts := 1
	for _, dim := range hdr.Axes() {
		nelmts *= int(dim)
	}

	c_start := C.LONGLONG(0)
	c_nelmts := C.LONGLONG(nelmts)
	c_status := C.int(0)
	c_imgtype := C.int(0)
	var c_ptr unsafe.Pointer

	switch rv.Interface().(type) {
	case []byte:
		c_imgtype = C.TBYTE
		data := rv.Interface().([]byte)
		c_ptr = unsafe.Pointer(&data[0])

	case []int8:
		c_imgtype = C.TBYTE
		data := rv.Interface().([]int8)
		c_ptr = unsafe.Pointer(&data[0])

	case []int16:
		c_imgtype = C.TSHORT
		data := rv.Interface().([]int16)
		c_ptr = unsafe.Pointer(&data[0])

	case []uint16:
		c_imgtype = C.TUSHORT
		data := rv.Interface().([]uint16)
		c_ptr = unsafe.Pointer(&data[0])

	case []int32:
		c_imgtype = C.TINT
		data := rv.Interface().([]int32)
		c_ptr = unsafe.Pointer(&data[0])

	case []int64:
		c_imgtype = C.TLONGLONG
		data := rv.Interface().([]int64)
		c_ptr = unsafe.Pointer(&data[0])

	case []float32:
		c_imgtype = C.TFLOAT
		data := rv.Interface().([]float32)
		c_ptr = unsafe.Pointer(&data[0])

	case []float64:
		c_imgtype = C.TDOUBLE
		data := rv.Interface().([]float64)
		c_ptr = unsafe.Pointer(&data[0])

	default:
		panic(fmt.Errorf("invalid image type [%T]", rv.Interface()))
	}

	C.fits_write_img(hdu.f.c, c_imgtype, c_start+1, c_nelmts, c_ptr, &c_status)
	if c_status > 0 {
		return to_err(c_status)
	}

	return err
}

// newImageHDU returns the i-th HDU from file f.
// if i==0, the returned ImageHDU is actually the primary HDU.
func newImageHDU(f *File, hdr Header, i int) (hdu HDU, err error) {
	switch i {
	case 0:
		hdu, err = newPrimaryHDU(f, hdr)
	default:
		hdu = &ImageHDU{
			f:      f,
			header: hdr,
		}
	}
	return hdu, err
}

func init() {
	g_hdus[IMAGE_HDU] = newImageHDU
}

// EOF
