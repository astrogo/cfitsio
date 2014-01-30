package cfitsio

// #include "go-cfitsio.h"
// #include "go-cfitsio-utils.h"
import "C"
import (
	"strings"
	"unsafe"
)

// strIsContinued checks whether the last non-whitespace char in the string
// is an ampersand, which indicates the value string is to be continued on the
// next line
func strIsContinued(v string) bool {
	vv := strings.Trim(v, " \n\t'")
	if len(vv) == 0 {
		return false
	}
	return vv[len(vv)-1] == '&'
}

func getLongValueString(f *File, key string) (string, error) {
	c_key := C.CString(key)
	defer C.free(unsafe.Pointer(c_key))
	c_status := C.int(0)
	var c_value *C.char = nil
	defer C.free(unsafe.Pointer(c_value))
	C.fits_read_key_longstr(f.c, c_key, &c_value, nil, &c_status)
	if c_status > 0 {
		return "", to_err(c_status)
	}
	return C.GoString(c_value), nil
}

// EOF
