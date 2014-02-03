package cfitsio

// #include "go-cfitsio.h"
// #include "go-cfitsio-utils.h"
import "C"
import (
	"fmt"
	"unsafe"
)

//
type ImageHDU struct {
	f      *File
	header Header
	data   interface{}
	read   bool // whether the image has been loaded from FITS
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

func (hdu *ImageHDU) Data() (interface{}, error) {
	var err error
	if !hdu.read {
		err = hdu.load()
	}
	return hdu.data, err
}

func (hdu *ImageHDU) load() error {
	var err error
	hdr := hdu.Header()
	naxes := len(hdr.Axes())
	switch naxes {
	case 0:
		hdu.read = true
		hdu.data = []int64{}
	default:
		nelmts := 1
		for _, dim := range hdr.Axes() {
			nelmts *= int(dim)
		}
		c_start := C.LONGLONG(0)
		c_nelmts := C.LONGLONG(nelmts)
		c_anynull := C.int(0)
		c_status := C.int(0)
		c_imgtype := C.int(0)
		var c_ptr unsafe.Pointer
		switch hdr.Bitpix() {
		case 8:
			data := make([]byte, nelmts)
			c_ptr = unsafe.Pointer(&data[0])
			c_imgtype = C.TBYTE
			hdu.data = data

		case 16:
			data := make([]int16, nelmts)
			c_ptr = unsafe.Pointer(&data[0])
			c_imgtype = C.TSHORT
			hdu.data = data

		case 32:
			data := make([]int32, nelmts)
			c_ptr = unsafe.Pointer(&data[0])
			c_imgtype = C.TINT
			hdu.data = data

		case 64:
			data := make([]int64, nelmts)
			c_ptr = unsafe.Pointer(&data[0])
			c_imgtype = C.TLONGLONG
			hdu.data = data

		case -32:
			data := make([]float32, nelmts)
			c_ptr = unsafe.Pointer(&data[0])
			c_imgtype = C.TFLOAT
			hdu.data = data

		case -64:
			data := make([]float64, nelmts)
			c_ptr = unsafe.Pointer(&data[0])
			c_imgtype = C.TDOUBLE
			hdu.data = data

		default:
			panic(fmt.Errorf("invalid image type [%v]", hdr.Bitpix()))
		}
		C.fits_read_img(hdu.f.c, c_imgtype, c_start+1, c_nelmts, c_ptr, c_ptr, &c_anynull, &c_status)
		if c_status > 0 {
			return to_err(c_status)
		}
		hdu.read = true
	}
	return err
}

func newImageHDU(f *File, hdr Header, i int) (hdu HDU, err error) {
	switch i {
	case 0:
		hdu = &PrimaryHDU{
			f:      f,
			header: hdr,
			data:   nil,
		}
	default:
		hdu = &ImageHDU{
			f:      f,
			header: hdr,
			data:   nil,
		}
	}
	return hdu, err
}

func init() {
	g_hdus[IMAGE_HDU] = newImageHDU
}

// EOF
