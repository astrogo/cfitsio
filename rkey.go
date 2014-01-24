package cfitsio

// #include "go-cfitsio.h"
// #include "go-cfitsio-utils.h"
import "C"
import (
	"fmt"
	"reflect"
	"unsafe"
)

type Keyword struct {
	Name    string
	Value   interface{}
	Comment string
}

// Return the number of existing keywords (not counting the END keyword) and the amount of space currently available for more keywords. It returns morekeys = -1 if the header has not yet been closed. Note that CFITSIO will dynamically add space if required when writing new keywords to a header so in practice there is no limit to the number of keywords that can be added to a header. A null pointer may be entered for the morekeys parameter if it's value is not needed.
func (f *File) HdrSpace() (keysexist, morekeys int, err error) {
	c_key := C.int(0)
	c_more := C.int(0)
	c_status := C.int(0)
	C.fits_get_hdrspace(f.c, &c_key, &c_more, &c_status)
	if c_status > 0 {
		return 0, 0, to_err(c_status)
	}
	return int(c_key), int(c_more), nil
}

/*
// Return the specified keyword. In the first routine, the datatype parameter specifies the desired returned data type of the keyword value and can have one of the following symbolic constant values: TSTRING, TLOGICAL (== int), TBYTE, TSHORT, TUSHORT, TINT, TUINT, TLONG, TULONG, TLONGLONG, TFLOAT, TDOUBLE, TCOMPLEX, and TDBLCOMPLEX. Within the context of this routine, TSTRING corresponds to a 'char*' data type, i.e., a pointer to a character array. Data type conversion will be performed for numeric values if the keyword value does not have the same data type. If the value of the keyword is undefined (i.e., the value field is blank) then an error status = VALUE_UNDEFINED will be returned.
The second routine returns the keyword value as a character string (a literal copy of what is in the value field) regardless of the intrinsic data type of the keyword. The third routine returns the entire 80-character header record of the keyword, with any trailing blank characters stripped off. The fourth routine returns the (next) header record that contains the literal string of characters specified by the 'string' argument.

If a NULL comment pointer is supplied then the comment string will not be returned.

  int fits_read_key / ffgky
      (fitsfile *fptr, int datatype, char *keyname, > DTYPE *value,
       char *comment, int *status)

  int fits_read_keyword / ffgkey
      (fitsfile *fptr, char *keyname, > char *value, char *comment,
       int *status)

  int fits_read_card / ffgcrd
      (fitsfile *fptr, char *keyname, > char *card, int *status)

  int fits_read_str / ffgstr
      (fitsfile *fptr, char *string, > char *card, int *status)
*/

func (f *File) ReadKeyword(k *Keyword) error {
	v := &k.Value
	ptr := reflect.ValueOf(v)
	if ptr.Kind() != reflect.Ptr {
		return fmt.Errorf("cfitsio.ReadKey: invalid type (%T). expected reflect.Ptr", v)
	}

	c_key := C.CString(k.Name)
	defer C.free(unsafe.Pointer(c_key))

	c_ptr := unsafe.Pointer(ptr.Pointer())
	c_typ := C.int(0)
	c_comment := C.char_buf_array(C.FLEN_COMMENT)
	defer C.free(unsafe.Pointer(c_comment))

	c_status := C.int(0)

	switch ptr.Elem().Kind() {
	case reflect.Bool:
		c_typ = C.TLOGICAL

	case reflect.Int8:
		c_typ = C.TSBYTE
	case reflect.Int16:
		c_typ = C.TSHORT
	case reflect.Int32:
		c_typ = C.TINT
	case reflect.Int64:
		c_typ = C.TLONG

	case reflect.Uint8:
		c_typ = C.TBYTE
	case reflect.Uint16:
		c_typ = C.TUSHORT
	case reflect.Uint32:
		c_typ = C.TUINT
	case reflect.Uint64:
		c_typ = C.TULONG

	case reflect.Float32:
		c_typ = C.TFLOAT
	case reflect.Float64:
		c_typ = C.TDOUBLE

	case reflect.Complex64:
		c_typ = C.TCOMPLEX
	case reflect.Complex128:
		c_typ = C.TDBLCOMPLEX

	default:
		panic(fmt.Errorf("invalid type: %T", v))
	}

	C.fits_read_key(f.c, c_typ, c_key, c_ptr, c_comment, &c_status)
	if c_status > 0 {
		return to_err(c_status)
	}

	k.Comment = C.GoString(c_comment)
	return nil
}

// Return the nth header record in the CHU. The first keyword in the header is at keynum = 1; if keynum = 0 then these routines simply reset the internal CFITSIO pointer to the beginning of the header so that subsequent keyword operations will start at the top of the header (e.g., prior to searching for keywords using wild cards in the keyword name). The first routine returns the entire 80-character header record (with trailing blanks truncated), while the second routine parses the record and returns the name, value, and comment fields as separate (blank truncated) character strings. If a NULL comment pointer is given on input, then the comment string will not be returned.
func (f *File) ReadRecord(n int) (card string, err error) {
	c_key := C.int(n)
	c_status := C.int(0)
	c_card := C.char_buf_array(C.FLEN_CARD)
	defer C.free(unsafe.Pointer(c_card))

	C.fits_read_record(f.c, c_key, c_card, &c_status)
	if c_status > 0 {
		return "", to_err(c_status)
	}
	return C.GoString(c_card), nil
}

// Return the nth header record in the CHU. The first keyword in the header is at keynum = 1; if keynum = 0 then these routines simply reset the internal CFITSIO pointer to the beginning of the header so that subsequent keyword operations will start at the top of the header (e.g., prior to searching for keywords using wild cards in the keyword name). The first routine returns the entire 80-character header record (with trailing blanks truncated), while the second routine parses the record and returns the name, value, and comment fields as separate (blank truncated) character strings. If a NULL comment pointer is given on input, then the comment string will not be returned.
// int fits_read_keyn / ffgkyn
//     (fitsfile *fptr, int keynum, > char *keyname, char *value,
//      char *comment, int *status)

// EOF
