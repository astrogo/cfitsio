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

	reflect.Bool: map[HDUType]string{
		ASCII_TBL:  "",
		BINARY_TBL: "L",
	},

	reflect.Int: map[HDUType]string{
		ASCII_TBL:  "I64",
		BINARY_TBL: "K",
	},

	reflect.Int8: map[HDUType]string{
		ASCII_TBL:  "I8",
		BINARY_TBL: "S",
	},

	reflect.Int16: map[HDUType]string{
		ASCII_TBL:  "I16",
		BINARY_TBL: "I",
	},

	reflect.Int32: map[HDUType]string{
		ASCII_TBL:  "I32",
		BINARY_TBL: "J",
	},

	reflect.Int64: map[HDUType]string{
		ASCII_TBL:  "I64",
		BINARY_TBL: "K",
	},

	reflect.Uint: map[HDUType]string{
		ASCII_TBL:  "I64",
		BINARY_TBL: "V",
	},

	reflect.Uint8: map[HDUType]string{
		ASCII_TBL:  "I8",
		BINARY_TBL: "B",
	},

	reflect.Uint16: map[HDUType]string{
		ASCII_TBL:  "I16",
		BINARY_TBL: "U",
	},

	reflect.Uint32: map[HDUType]string{
		ASCII_TBL:  "I32",
		BINARY_TBL: "V",
	},

	reflect.Uint64: map[HDUType]string{
		ASCII_TBL:  "I64",
		BINARY_TBL: "V",
	},

	reflect.Uintptr: map[HDUType]string{
		ASCII_TBL:  "",
		BINARY_TBL: "",
	},

	reflect.Float32: map[HDUType]string{
		ASCII_TBL:  "F32",
		BINARY_TBL: "E",
	},

	reflect.Float64: map[HDUType]string{
		ASCII_TBL:  "F64",
		BINARY_TBL: "D",
	},

	reflect.Complex64: map[HDUType]string{
		ASCII_TBL:  "",
		BINARY_TBL: "C",
	},

	reflect.Complex128: map[HDUType]string{
		ASCII_TBL:  "",
		BINARY_TBL: "M",
	},

	reflect.Array: map[HDUType]string{
		ASCII_TBL:  "",
		BINARY_TBL: "",
	},

	reflect.Chan: map[HDUType]string{
		ASCII_TBL:  "",
		BINARY_TBL: "",
	},

	reflect.Func: map[HDUType]string{
		ASCII_TBL:  "",
		BINARY_TBL: "",
	},

	reflect.Interface: map[HDUType]string{
		ASCII_TBL:  "",
		BINARY_TBL: "",
	},

	reflect.Map: map[HDUType]string{
		ASCII_TBL:  "",
		BINARY_TBL: "",
	},

	reflect.Ptr: map[HDUType]string{
		ASCII_TBL:  "",
		BINARY_TBL: "",
	},

	reflect.Slice: map[HDUType]string{
		ASCII_TBL:  "",
		BINARY_TBL: "",
	},

	reflect.String: map[HDUType]string{
		ASCII_TBL:  "A80",
		BINARY_TBL: "80A",
	},

	reflect.Struct: map[HDUType]string{
		ASCII_TBL:  "",
		BINARY_TBL: "",
	},
}

// EOF
