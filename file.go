package cfitsio

// #include <string.h>
// #include <stdlib.h>
// #include "go-cfitsio.h"
// #include "go-cfitsio-utils.h"
import "C"

import (
	"unsafe"
)

type Mode int

const (
	ReadOnly  Mode = C.READONLY
	ReadWrite Mode = C.READWRITE
)

type File struct {
	c    *C.fitsfile
	hdus []HDU
}

func (f *File) HDUs() []HDU {
	return f.hdus
}

// Open an existing FITS file
func Open(fname string, mode Mode) (File, error) {
	var f File
	var err error

	c_status := C.int(0)
	c_fname := C.CString(fname)
	defer C.free(unsafe.Pointer(c_fname))

	C.ffopen(&f.c, c_fname, C.int(mode), &c_status)
	if c_status > 0 {
		return f, to_err(c_status)
	}

	// ffopen might have moved to specific HDU (via fname specifications)
	// remember it and go back to that one after we've dealt with "our" HDUs
	ihdu := f.HDUNum()
	defer f.SeekHDU(ihdu, 0)

	nhdus, err := f.NumHDUs()
	if err != nil {
		return f, err
	}

	f.hdus = make([]HDU, 0, nhdus)
	for i := 0; i < nhdus; i++ {
		hdu, err := f.readHDU(i)
		if err != nil {
			return f, err
		}
		f.hdus = append(f.hdus, hdu)
	}
	if err != nil {
		return f, err
	}

	return f, err
}

// Create creates and opens a new empty output FITS file.
func Create(fname string) (f File, err error) {
	c_status := C.int(0)
	c_fname := C.CString(fname)
	defer C.free(unsafe.Pointer(c_fname))

	C.fits_create_file(&f.c, c_fname, &c_status)
	if c_status > 0 {
		err = to_err(c_status)
		return
	}

	return
}

// Close closes a previously opened FITS file.
func (f *File) Close() error {
	c_status := C.int(0)
	C.fits_close_file(f.c, &c_status)
	err := to_err(c_status)
	if err != nil {
		return err
	}

	for _, hdu := range f.hdus {
		err2 := hdu.Close()
		if err2 != nil {
			err = err2
		}
	}
	return err
}

// Delete closes a previously opened FITS file and also DELETES the file.
func (f *File) Delete() (err error) {
	c_status := C.int(0)
	C.fits_delete_file(f.c, &c_status)
	err = to_err(c_status)
	return
}

// Name returns the name of a FITS file
func (f *File) Name() (string, error) {
	c_name := C.char_buf_array(C.FLEN_FILENAME)
	defer C.free(unsafe.Pointer(c_name))
	c_status := C.int(0)
	C.fits_file_name(f.c, c_name, &c_status)
	if c_status > 0 {
		return "", to_err(c_status)
	}
	return C.GoString(c_name), nil
}

// Mode returns the mode of a FITS file (ReadOnly or ReadWrite)
func (f *File) Mode() (Mode, error) {
	c_mode := C.int(0)
	c_status := C.int(0)
	C.fits_file_mode(f.c, &c_mode, &c_status)
	if c_status > 0 {
		return Mode(0), to_err(c_status)
	}
	return Mode(c_mode), nil
}

// UrlType returns the type of a FITS file (e.g. ftp:// or file://)
func (f *File) UrlType() (string, error) {
	c_url := C.char_buf_array(C.FLEN_KEYWORD) //FIXME: correct length ?
	defer C.free(unsafe.Pointer(c_url))
	c_status := C.int(0)
	C.fits_url_type(f.c, c_url, &c_status)
	if c_status > 0 {
		return "", to_err(c_status)
	}
	return C.GoString(c_url), nil
}

// eof
