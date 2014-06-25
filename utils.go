package cfitsio

// #include "go-cfitsio.h"
// #include "go-cfitsio-utils.h"
import "C"
import (
	"reflect"
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

// gotype2FITS returns the FITS format corresponding to a Go type
// it returns "" if there is no corresponding FITS format.
func gotype2FITS(v interface{}, hdu HDUType) string {
	rt := reflect.TypeOf(v)
	hdr := ""
	switch rt.Kind() {
	case reflect.Slice:
		hdr = "Q"
		rt = rt.Elem()
	case reflect.Array:
		hdr = "Q"
		rt = rt.Elem()
	default:
		// no-op
	}

	t, ok := g_gotype2FITS[rt.Kind()]
	if !ok {
		return ""
	}
	fitsfmt, ok := t[hdu]
	if !ok {
		return ""
	}
	return hdr + fitsfmt
}

var g_gotype2FITS = map[reflect.Kind]map[HDUType]string{

	reflect.Bool: {
		ASCII_TBL:  "",
		BINARY_TBL: "L",
	},

	reflect.Int: {
		ASCII_TBL:  "I21",
		BINARY_TBL: "K",
	},

	reflect.Int8: {
		ASCII_TBL:  "I7",
		BINARY_TBL: "S",
	},

	reflect.Int16: {
		ASCII_TBL:  "I7",
		BINARY_TBL: "I",
	},

	reflect.Int32: {
		ASCII_TBL:  "I12",
		BINARY_TBL: "J",
	},

	reflect.Int64: {
		ASCII_TBL:  "I21",
		BINARY_TBL: "K",
	},

	reflect.Uint: {
		ASCII_TBL:  "I21",
		BINARY_TBL: "V",
	},

	reflect.Uint8: {
		ASCII_TBL:  "I7",
		BINARY_TBL: "B",
	},

	reflect.Uint16: {
		ASCII_TBL:  "I7",
		BINARY_TBL: "U",
	},

	reflect.Uint32: {
		ASCII_TBL:  "I12",
		BINARY_TBL: "V",
	},

	reflect.Uint64: {
		ASCII_TBL:  "I21",
		BINARY_TBL: "V",
	},

	reflect.Uintptr: {
		ASCII_TBL:  "",
		BINARY_TBL: "",
	},

	reflect.Float32: {
		ASCII_TBL:  "E26.17", // must write as float64 since we can only read as such
		BINARY_TBL: "E",
	},

	reflect.Float64: {
		ASCII_TBL:  "E26.17",
		BINARY_TBL: "D",
	},

	reflect.Complex64: {
		ASCII_TBL:  "",
		BINARY_TBL: "C",
	},

	reflect.Complex128: {
		ASCII_TBL:  "",
		BINARY_TBL: "M",
	},

	reflect.Array: {
		ASCII_TBL:  "",
		BINARY_TBL: "",
	},

	reflect.Chan: {
		ASCII_TBL:  "",
		BINARY_TBL: "",
	},

	reflect.Func: {
		ASCII_TBL:  "",
		BINARY_TBL: "",
	},

	reflect.Interface: {
		ASCII_TBL:  "",
		BINARY_TBL: "",
	},

	reflect.Map: {
		ASCII_TBL:  "",
		BINARY_TBL: "",
	},

	reflect.Ptr: {
		ASCII_TBL:  "",
		BINARY_TBL: "",
	},

	reflect.Slice: {
		ASCII_TBL:  "",
		BINARY_TBL: "",
	},

	reflect.String: {
		ASCII_TBL:  "A80",
		BINARY_TBL: "80A",
	},

	reflect.Struct: {
		ASCII_TBL:  "",
		BINARY_TBL: "",
	},
}

// EOF
