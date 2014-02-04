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
	rv := reflect.ValueOf(data)
	if !rv.CanAddr() {
		return fmt.Errorf("%T is not addressable", data)
	}
	err := hdu.load(rv)
	return err
}

func (hdu *ImageHDU) load(rv reflect.Value) error {
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
	rv.SetLen(nelmts)

	c_start := C.LONGLONG(0)
	c_nelmts := C.LONGLONG(nelmts)
	c_anynull := C.int(0)
	c_status := C.int(0)
	c_imgtype := C.int(0)
	var c_ptr unsafe.Pointer
	switch rv.Interface().(type) {
	case []byte:
		c_ptr = unsafe.Pointer(rv.Index(0).Pointer())
		c_imgtype = C.TBYTE

	case []int16:
		c_ptr = unsafe.Pointer(rv.Index(0).Pointer())
		c_imgtype = C.TSHORT

	case []int32:
		c_ptr = unsafe.Pointer(rv.Index(0).Pointer())
		c_imgtype = C.TINT

	case []int64:
		c_ptr = unsafe.Pointer(rv.Index(0).Pointer())
		c_imgtype = C.TLONGLONG

	case []float32:
		c_ptr = unsafe.Pointer(rv.Index(0).Pointer())
		c_imgtype = C.TFLOAT

	case []float64:
		c_ptr = unsafe.Pointer(rv.Index(0).Pointer())
		c_imgtype = C.TDOUBLE

	default:
		panic(fmt.Errorf("invalid image type [%T]", rv.Interface()))
	}
	C.fits_read_img(hdu.f.c, c_imgtype, c_start+1, c_nelmts, c_ptr, c_ptr, &c_anynull, &c_status)
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
