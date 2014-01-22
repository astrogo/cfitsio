package cfitsio

// #include "go-cfitsio.h"
import "C"
import "unsafe"

type HduType int

const (
	ImageHdu  HduType = C.IMAGE_HDU  // Primary Array or IMAGE HDU
	AsciiTbl  HduType = C.ASCII_TBL  // ASCII table HDU
	BinaryTbl HduType = C.BINARY_TBL // Binary table HDU
	AnyHdy    HduType = C.ANY_HDU    // matches any HDU type
)

// MoveAbsHdu moves to a different HDU in the file
func (f *File) MoveAbsHdu(hdu int) (HduType, error) {
	c_hdu := C.int(hdu)
	c_htype := C.int(0)
	c_status := C.int(0)

	C.fits_movabs_hdu(f.c, c_hdu, &c_htype, &c_status)
	if c_status > 0 {
		return HduType(c_htype), to_err(c_status)
	}

	return HduType(c_htype), nil
}

// MoveRelHdu moves to a different HDU in the file
func (f *File) MoveRelHdu(n int) (HduType, error) {
	c_n := C.int(n)
	c_htype := C.int(0)
	c_status := C.int(0)

	C.fits_movrel_hdu(f.c, c_n, &c_htype, &c_status)
	if c_status > 0 {
		return HduType(c_htype), to_err(c_status)
	}

	return HduType(c_htype), nil
}

// MoveNamHdu moves to a different HDU in the file
func (f *File) MoveNamHdu(hdu HduType, extname string, extvers int) error {
	c_hdu := C.int(hdu)
	c_name := C.CString(extname)
	defer C.free(unsafe.Pointer(c_name))
	c_vers := C.int(extvers)
	c_status := C.int(0)

	C.fits_movnam_hdu(f.c, c_hdu, c_name, c_vers, &c_status)
	if c_status > 0 {
		return to_err(c_status)
	}

	return nil
}

// NumHdus returns the total number of HDUs in the FITS file.
// This returns the number of completely defined HDUs in the file. If a new HDU has just been added to the FITS file, then that last HDU will only be counted if it has been closed, or if data has been written to the HDU. The current HDU remains unchanged by this routine.
func (f *File) NumHdus() (int, error) {
	c_n := C.int(0)
	c_status := C.int(0)
	C.fits_get_num_hdus(f.c, &c_n, &c_status)
	if c_status > 0 {
		return 0, to_err(c_status)
	}

	return int(c_n), nil
}

// HduNum returns the number of the current HDU (CHDU) in the FITS file (where the primary array = 1). This function returns the HDU number rather than a status value.
func (f *File) HduNum() int {

	c_n := C.int(0)
	C.fits_get_hdu_num(f.c, &c_n)
	return int(c_n)
}

// HduType returns the type of the current HDU in the FITS file. The possible values for hdutype are: IMAGE_HDU, ASCII_TBL, or BINARY_TBL.
func (f *File) HduType() (HduType, error) {
	c_hdu := C.int(0)
	c_status := C.int(0)
	C.fits_get_hdu_type(f.c, &c_hdu, &c_status)
	if c_status > 0 {
		return 0, to_err(c_status)
	}

	return HduType(c_hdu), nil
}

// Copy all or part of the HDUs in the FITS file associated with infptr and append them to the end of the FITS file associated with outfptr. If 'previous' is true, then any HDUs preceding the current HDU in the input file will be copied to the output file. Similarly, 'current' and 'following' determine whether the current HDU, and/or any following HDUs in the input file will be copied to the output file. Thus, if all 3 parameters are true, then the entire input file will be copied. On exit, the current HDU in the input file will be unchanged, and the last HDU in the output file will be the current HDU.
func (f *File) Copy(out *File, previous, current, following bool) error {
	c_previous := C.int(0)
	if previous {
		c_previous = C.int(1)
	}
	c_current := C.int(0)
	if current {
		c_current = C.int(1)
	}
	c_following := C.int(0)
	if following {
		c_following = C.int(1)
	}
	c_status := C.int(0)
	C.fits_copy_file(f.c, out.c, c_previous, c_current, c_following, &c_status)
	if c_status > 0 {
		return to_err(c_status)
	}
	return nil
}

// EOF
