package cfitsio

// #include "fitsio.h"
// #include <string.h>
// #include <stdlib.h>
import "C"

import (
	"unsafe"
)

type Mode int

const (
	ReadOnly Mode = C.READONLY
	ReadWrite Mode = C.READWRITE
)

type File struct {
	c *C.fitsfile
}

func OpenFile(fname string, mode Mode) (f File, err error) {
	c_status := C.int(0)
	c_fname := C.CString(fname)
	defer C.free(unsafe.Pointer(c_fname))
	
	C.ffopen(&f.c, c_fname, C.int(mode), &c_status)
	if c_status > 0 {
		err = to_err(c_status)
		return 
	}

	return
}

func OpenDiskFile(fname string, mode Mode) (f File, err error) {
	c_status := C.int(0)
	c_fname := C.CString(fname)
	defer C.free(unsafe.Pointer(c_fname))
	
	C.ffdkopn(&f.c, c_fname, C.int(mode), &c_status)
	if c_status > 0 {
		err = to_err(c_status)
		return 
	}

	return
}

// eof
