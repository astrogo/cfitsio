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

	nhdus, err := f.NumHdus()
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

	// go back at beginning of file
	_, err = f.MovAbsHdu(1)
	return f, err
}

// Open an existing data file.
func OpenDiskFile(fname string, mode Mode) (f File, err error) {
	c_status := C.int(0)
	c_fname := C.CString(fname)
	defer C.free(unsafe.Pointer(c_fname))

	C.fits_open_diskfile(&f.c, c_fname, C.int(mode), &c_status)
	if c_status > 0 {
		err = to_err(c_status)
		return
	}

	return
}

// Open an existing data file.
func OpenData(fname string, mode Mode) (f File, err error) {
	c_status := C.int(0)
	c_fname := C.CString(fname)
	defer C.free(unsafe.Pointer(c_fname))

	C.fits_open_data(&f.c, c_fname, C.int(mode), &c_status)
	if c_status > 0 {
		err = to_err(c_status)
		return
	}

	return
}

// Open an existing data file.
func OpenTable(fname string, mode Mode) (f File, err error) {
	c_status := C.int(0)
	c_fname := C.CString(fname)
	defer C.free(unsafe.Pointer(c_fname))

	C.fits_open_table(&f.c, c_fname, C.int(mode), &c_status)
	if c_status > 0 {
		err = to_err(c_status)
		return
	}

	return
}

// Open an existing data file.
func OpenImage(fname string, mode Mode) (f File, err error) {
	c_status := C.int(0)
	c_fname := C.CString(fname)
	defer C.free(unsafe.Pointer(c_fname))

	C.fits_open_image(&f.c, c_fname, C.int(mode), &c_status)
	if c_status > 0 {
		err = to_err(c_status)
		return
	}

	return
}

// Create and open a new empty output FITS file.
func NewFile(fname string) (f File, err error) {
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

// NewDiskFile creates and opens a new empty output FITS file.
func NewDiskFile(fname string) (f File, err error) {
	c_status := C.int(0)
	c_fname := C.CString(fname)
	defer C.free(unsafe.Pointer(c_fname))

	C.fits_create_diskfile(&f.c, c_fname, &c_status)
	if c_status > 0 {
		err = to_err(c_status)
		return
	}

	return
}

// Close closes a previously opened FITS file.
func (f *File) Close() (err error) {
	c_status := C.int(0)
	C.fits_close_file(f.c, &c_status)
	err = to_err(c_status)
	return
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
