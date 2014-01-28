package cfitsio

// #include "go-cfitsio.h"
// #include "go-cfitsio-utils.h"
import "C"
import (
	"unsafe"
)

type Keyword struct {
	Name    string
	Value   interface{}
	Comment string
}

func (hdu *Hdu) Keyword(name string) (*Keyword, error) {

	c_status := C.int(0)
	c_card := C.char_buf_array(C.FLEN_CARD)
	defer C.free(unsafe.Pointer(c_card))
	c_value := C.char_buf_array(C.FLEN_VALUE)
	defer C.free(unsafe.Pointer(c_value))
	c_com := C.char_buf_array(C.FLEN_COMMENT)
	defer C.free(unsafe.Pointer(c_com))

	c_name := C.CString(name)
	defer C.free(unsafe.Pointer(c_name))
	C.fits_read_card(hdu.File.c, c_name, c_card, &c_status)
	if c_status > 0 {
		return nil, to_err(c_status)
	}

	C.fits_parse_value(c_card, c_value, c_com, &c_status)
	if c_status > 0 {
		return nil, to_err(c_status)
	}

	value := C.GoString(c_value)
	if strIsContinued(value) {
		restoreQuote := value[0] == '\''
		if restoreQuote {
			const q = string('\'')
			v, err := hdu.getLongValueString(name)
			if err != nil {
				return nil, err
			}
			value = q + v + q
		}
	}
	comment := C.GoString(c_com)

	return hdu.parseRecord(name, value, comment)
}

func (hdu *Hdu) parseRecord(name, value, comment string) (*Keyword, error) {
	var key *Keyword
	var err error

	return key, err
}

func (hdu *Hdu) getLongValueString(key string) (string, error) {
	c_key := C.CString(key)
	defer C.free(unsafe.Pointer(c_key))
	c_status := C.int(0)
	var c_value *C.char = nil
	defer C.free(unsafe.Pointer(c_value))
	C.fits_read_key_longstr(hdu.File.c, c_key, &c_value, nil, &c_status)
	if c_status > 0 {
		return "", to_err(c_status)
	}
	return C.GoString(c_value), nil
}

// EOF
