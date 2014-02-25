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

type CardType int

const (
	TYP_STRUC_KEY  CardType = C.TYP_STRUC_KEY
	TYP_CMPRS_KEY  CardType = C.TYP_CMPRS_KEY
	TYP_SCAL_KEY   CardType = C.TYP_SCAL_KEY
	TYP_NULL_KEY   CardType = C.TYP_NULL_KEY
	TYP_DIM_KEY    CardType = C.TYP_DIM_KEY
	TYP_RANG_KEY   CardType = C.TYP_RANG_KEY
	TYP_UNIT_KEY   CardType = C.TYP_UNIT_KEY
	TYP_DISP_KEY   CardType = C.TYP_DISP_KEY
	TYP_HDUID_KEY  CardType = C.TYP_HDUID_KEY
	TYP_CKSUM_KEY  CardType = C.TYP_CKSUM_KEY
	TYP_WCS_KEY    CardType = C.TYP_WCS_KEY
	TYP_REFSYS_KEY CardType = C.TYP_REFSYS_KEY
	TYP_COMM_KEY   CardType = C.TYP_COMM_KEY
	TYP_CONT_KEY   CardType = C.TYP_CONT_KEY
	TYP_USER_KEY   CardType = C.TYP_USER_KEY
)

func (k CardType) String() string {
	switch k {
	case TYP_STRUC_KEY:
		return "TYP_STRUC_KEY"

	case TYP_CMPRS_KEY:
		return "TYP_CMPRS_KEY"

	case TYP_SCAL_KEY:
		return "TYP_SCAL_KEY"

	case TYP_NULL_KEY:
		return "TYP_NULL_KEY"

	case TYP_DIM_KEY:
		return "TYP_DIM_KEY"

	case TYP_RANG_KEY:
		return "TYP_RANG_KEY"

	case TYP_UNIT_KEY:
		return "TYP_UNIT_KEY"

	case TYP_DISP_KEY:
		return "TYP_DISP_KEY"

	case TYP_HDUID_KEY:
		return "TYP_HDUID_KEY"

	case TYP_CKSUM_KEY:
		return "TYP_CKSUM_KEY"

	case TYP_WCS_KEY:
		return "TYP_WCS_KEY"

	case TYP_REFSYS_KEY:
		return "TYP_REFSYS_KEY"

	case TYP_COMM_KEY:
		return "TYP_COMM_KEY"

	case TYP_CONT_KEY:
		return "TYP_CONT_KEY"

	case TYP_USER_KEY:
		return "TYP_USER_KEY"

	default:
		panic(fmt.Errorf("invalid CardType value (%v)", int(k)))
	}
}

type Card struct {
	Name    string
	Value   interface{}
	Comment string
}

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

	keyclass := CardType(C.fits_get_keyclass(c_key))
	switch keyclass {
	case TYP_COMM_KEY, TYP_CONT_KEY:
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
		var vv complex128
		_, err = fmt.Scanf(value, &vv)
		if err != nil {
			return err
		}
		card.Value = vv

	case 'C':
		card.Value = strings.TrimRight(value, " ")

	default:
		return fmt.Errorf("invalid keyword-type value (%v)", string(dtype))
	}
	return err
}

// EOF
