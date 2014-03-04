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

// Card is a record block (meta data) in a Header.
type Card struct {
	Name    string
	Value   interface{}
	Comment string
}

// newCard returns the i-th Card from the current HDU of file f.
func newCard(f *File, i int) (Card, error) {
	var card Card
	var err error

	c_status := C.int(0)
	c_key := C.CStringN(C.FLEN_KEYWORD)
	defer C.free(unsafe.Pointer(c_key))
	c_value := C.CStringN(C.FLEN_VALUE)
	defer C.free(unsafe.Pointer(c_value))
	c_com := C.CStringN(C.FLEN_COMMENT)
	defer C.free(unsafe.Pointer(c_com))

	c_keyn := C.int(i)

	C.fits_read_keyn(f.c, c_keyn, c_key, c_value, c_com, &c_status)
	if c_status > 0 {
		return card, to_err(c_status)
	}

	name := C.GoString(c_key)
	value := C.GoString(c_value)
	if strIsContinued(value) {
		restoreQuote := value[0] == '\''
		if restoreQuote {
			const q = string('\'')
			v, err := getLongValueString(f, name)
			if err != nil {
				return card, err
			}
			value = q + v + q
		}
	}
	comment := C.GoString(c_com)

	keyclass := C.fits_get_keyclass(c_key)
	switch keyclass {
	case C.TYP_COMM_KEY, C.TYP_CONT_KEY:
		return card, fmt.Errorf("comm key | continue key")
	}

	err = parseRecord(name, value, comment, &card)
	return card, err
}

func parseRecord(name, value, comment string, card *Card) error {
	var err error

	card.Name = name
	card.Comment = comment

	c_status := C.int(0)
	c_value := C.CString(value)
	defer C.free(unsafe.Pointer(c_value))
	var c_type C.char
	C.fits_get_keytype(c_value, &c_type, &c_status)
	if c_status > 0 {
		return to_err(c_status)
	}

	if value[0] == '\'' {
		value = value[1 : len(value)-1]
	}

	dtype := string(c_type)[0]
	switch dtype {
	case 'L':
		card.Value = value == "T"

	case 'F':
		var vv float64
		vv, err = strconv.ParseFloat(value, 64)
		if err != nil {
			return err
		}
		card.Value = vv

	case 'I', 'T':
		var vv int64
		vv, err = strconv.ParseInt(value, 10, 64)
		if err != nil {
			return err
		}
		card.Value = vv

	case 'X':
		toks := strings.Split(value[1:len(value)-1], ",")
		var vv0 float64
		vv0, err = strconv.ParseFloat(strings.Trim(toks[0], " \t\n"), 64)
		if err != nil {
			return err
		}
		var vv1 float64
		vv1, err = strconv.ParseFloat(strings.Trim(toks[1], " \t\n"), 64)
		if err != nil {
			return err
		}

		card.Value = complex(vv0, vv1)

	case 'C':
		card.Value = strings.TrimRight(value, " ")

	default:
		return fmt.Errorf("invalid keyword-type value (%v)", string(dtype))
	}
	return err
}

// EOF
