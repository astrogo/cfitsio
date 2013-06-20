package cfitsio

// #include "fitsio.h"
// #include <string.h>
// #include <stdlib.h>
// #include "go-cfitsio-utils.h"
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

func OpenData(fname string, mode Mode) (f File, err error) {
	c_status := C.int(0)
	c_fname := C.CString(fname)
	defer C.free(unsafe.Pointer(c_fname))
	
	C.ffdopn(&f.c, c_fname, C.int(mode), &c_status)
	if c_status > 0 {
		err = to_err(c_status)
		return 
	}

	return
}

func OpenTable(fname string, mode Mode) (f File, err error) {
	c_status := C.int(0)
	c_fname := C.CString(fname)
	defer C.free(unsafe.Pointer(c_fname))
	
	C.fftopn(&f.c, c_fname, C.int(mode), &c_status)
	if c_status > 0 {
		err = to_err(c_status)
		return 
	}

	return
}

func OpenImage(fname string, mode Mode) (f File, err error) {
	c_status := C.int(0)
	c_fname := C.CString(fname)
	defer C.free(unsafe.Pointer(c_fname))
	
	C.ffiopn(&f.c, c_fname, C.int(mode), &c_status)
	if c_status > 0 {
		err = to_err(c_status)
		return 
	}

	return
}

func NewFile(fname string) (f File, err error) {
	c_status := C.int(0)
	c_fname := C.CString(fname)
	defer C.free(unsafe.Pointer(c_fname))
	
	C.ffinit(&f.c, c_fname, &c_status)
	if c_status > 0 {
		err = to_err(c_status)
		return 
	}

	return
}

func NewDiskFile(fname string) (f File, err error) {
	c_status := C.int(0)
	c_fname := C.CString(fname)
	defer C.free(unsafe.Pointer(c_fname))
	
	C.ffdkinit(&f.c, c_fname, &c_status)
	if c_status > 0 {
		err = to_err(c_status)
		return 
	}

	return
}

func (f *File) Close() (err error) {
	c_status := C.int(0)
	C.ffclos(f.c, &c_status)
	err = to_err(c_status)
	return
}

func (f *File) Delete() (err error) {
	c_status := C.int(0)
	C.ffdelt(f.c, &c_status)
	err = to_err(c_status)
	return
}

func (f *File) Name() (string, error) {
	c_name := C.char_buf_array(C.FLEN_FILENAME)
	defer C.free(unsafe.Pointer(c_name))
	c_status := C.int(0)
	C.ffflnm(f.c, c_name, &c_status)
	if c_status > 0 {
		return "", to_err(c_status)
	}
	return C.GoString(c_name), nil
}

func (f *File) Mode() (Mode, error) {
	c_mode := C.int(0)
	c_status := C.int(0)
	C.ffflmd(f.c, &c_mode, &c_status)
	if c_status > 0 {
		return Mode(0), to_err(c_status)
	}
	return Mode(c_mode), nil
}

func (f *File) UrlType() (string, error) {
	c_url := C.char_buf_array(C.FLEN_KEYWORD) //FIXME: correct length ?
	defer C.free(unsafe.Pointer(c_url))
	c_status := C.int(0)
	C.ffurlt(f.c, c_url, &c_status)
	if c_status > 0 {
		return "", to_err(c_status)
	}
	return C.GoString(c_url), nil
}

// eof
