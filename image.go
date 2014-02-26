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

func (hdu *ImageHDU) Close() error {
	hdu.f = nil
	return nil
}

func (hdu *ImageHDU) Header() Header {
	return hdu.header
}

func (hdu *ImageHDU) Type() HDUType {
	return hdu.header.htype
}

func (hdu *ImageHDU) Name() string {
	card := hdu.header.Get("EXTNAME")
	if card == nil {
		return ""
	}
	return card.Value.(string)
}

func (hdu *ImageHDU) Version() int {
	card := hdu.header.Get("EXTVER")
	if card == nil {
		return 1
	}
	return card.Value.(int)
}

func (hdu *ImageHDU) Data(data interface{}) error {
	//rt := reflect.TypeOf(data)
	rv := reflect.ValueOf(data).Elem()
	if !rv.CanAddr() {
		return fmt.Errorf("%T is not addressable", data)
	}

	err := hdu.load(rv)
	return err
}

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

	v.Set(rv)
	return err
}

// Write writes the image to disk
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
