package cfitsio

// #include "go-cfitsio.h"
// #include "go-cfitsio-utils.h"
import "C"
import (
	"fmt"
	"strconv"
	"strings"
	"unsafe"
)

type Keyword struct {
	Name    string
	Value   interface{}
	Comment string
	hdu     *Hdu
}

func (key *Keyword) Delete() {
	key.hdu = nil
}

// KeywordByName returns the Keyword with name `name` held by this HDU
func (hdu *Hdu) KeywordByName(name string) (*Keyword, error) {

	if err := hdu.makeCurrent(); err != nil {
		return nil, err
	}

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

// Keyword returns the i-th Keyword held by this HDU
func (hdu *Hdu) Keyword(i int) (*Keyword, error) {
	if err := hdu.makeCurrent(); err != nil {
		return nil, err
	}

	c_status := C.int(0)
	c_key := C.char_buf_array(C.FLEN_KEYWORD)
	defer C.free(unsafe.Pointer(c_key))
	c_value := C.char_buf_array(C.FLEN_VALUE)
	defer C.free(unsafe.Pointer(c_value))
	c_com := C.char_buf_array(C.FLEN_COMMENT)
	defer C.free(unsafe.Pointer(c_com))

	c_keyn := C.int(i)

	C.fits_read_keyn(hdu.File.c, c_keyn, c_key, c_value, c_com, &c_status)
	if c_status > 0 {
		return nil, to_err(c_status)
	}

	name := C.GoString(c_key)
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

	keyclass := KeyType(C.fits_get_keyclass(c_key))
	switch keyclass {
	case TYP_COMM_KEY, TYP_CONT_KEY:
		return nil, nil
	}

	return hdu.parseRecord(name, value, comment)
}

func (hdu *Hdu) parseRecord(name, value, comment string) (*Keyword, error) {
	var key *Keyword
	var err error

	c_status := C.int(0)
	c_value := C.CString(value)
	defer C.free(unsafe.Pointer(c_value))
	var c_type C.char
	C.fits_get_keytype(c_value, &c_type, &c_status)
	if c_status > 0 {
		return nil, to_err(c_status)
	}

	if value[0] == '\'' {
		value = value[1 : len(value)-1]
	}

	dtype := string(c_type)[0]
	switch dtype {
	case 'L':
		vv := value == "T"
		key = &Keyword{
			Name:    name,
			Value:   vv,
			Comment: comment,
		}

	case 'F':
		var vv float64
		vv, err = strconv.ParseFloat(value, 64)
		if err != nil {
			return nil, err
		}
		key = &Keyword{
			Name:    name,
			Value:   vv,
			Comment: comment,
		}

	case 'I', 'T':
		var vv int64
		vv, err = strconv.ParseInt(value, 10, 64)
		if err != nil {
			return nil, err
		}
		key = &Keyword{
			Name:    name,
			Value:   vv,
			Comment: comment,
		}

	case 'X':
		var vv complex128
		_, err = fmt.Scanf(value, &vv)
		if err != nil {
			return nil, err
		}
		key = &Keyword{
			Name:    name,
			Value:   vv,
			Comment: comment,
		}

	case 'C':
		vv := strings.TrimRight(value, " ")
		key = &Keyword{
			Name:    name,
			Value:   vv,
			Comment: comment,
		}

	default:
		return nil, fmt.Errorf("invalid key-type value (%v)", string(dtype))
	}
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
