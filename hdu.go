package cfitsio

// #include <stdio.h>
// #include "go-cfitsio.h"
import "C"
import (
	"io"
	"io/ioutil"
	"unsafe"
)

type HduType int

const (
	ImageHdu  HduType = C.IMAGE_HDU  // Primary Array or IMAGE HDU
	AsciiTbl  HduType = C.ASCII_TBL  // ASCII table HDU
	BinaryTbl HduType = C.BINARY_TBL // Binary table HDU
	AnyHdy    HduType = C.ANY_HDU    // matches any HDU type
)

// MovAbsHdu moves to a different HDU in the file
func (f *File) MovAbsHdu(hdu int) (HduType, error) {
	c_hdu := C.int(hdu)
	c_htype := C.int(0)
	c_status := C.int(0)

	C.fits_movabs_hdu(f.c, c_hdu, &c_htype, &c_status)
	if c_status > 0 {
		return HduType(c_htype), to_err(c_status)
	}

	return HduType(c_htype), nil
}

// MovRelHdu moves to a different HDU in the file
func (f *File) MovRelHdu(n int) (HduType, error) {
	c_n := C.int(n)
	c_htype := C.int(0)
	c_status := C.int(0)

	C.fits_movrel_hdu(f.c, c_n, &c_htype, &c_status)
	if c_status > 0 {
		return HduType(c_htype), to_err(c_status)
	}

	return HduType(c_htype), nil
}

// MovNamHdu moves to a different HDU in the file
func (f *File) MovNamHdu(hdu HduType, extname string, extvers int) error {
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

// Copy the current HDU from the FITS file associated with infptr and append it to the end of the FITS file associated with outfptr. Space may be reserved for MOREKEYS additional keywords in the output header.
func (f *File) CopyHdu(out *File, morekeys int) error {
	c_morekeys := C.int(morekeys)
	c_status := C.int(0)
	C.fits_copy_hdu(f.c, out.c, c_morekeys, &c_status)
	if c_status > 0 {
		return to_err(c_status)
	}
	return nil
}

// Write the current HDU in the input FITS file to the output FILE stream (e.g., to stdout).
func (f *File) WriteHdu(w io.Writer) error {
	tmp, err := ioutil.TempFile("", "go-cfitsio-")
	if err != nil {
		return err
	}
	defer tmp.Close()

	c_mode := C.CString("rw+")
	defer C.free(unsafe.Pointer(c_mode))
	fstream := C.fdopen(C.int(tmp.Fd()), c_mode)
	c_status := C.int(0)
	C.fits_write_hdu(f.c, fstream, &c_status)
	if c_status > 0 {
		return to_err(c_status)
	}
	C.fflush(fstream)

	_, err = tmp.Seek(0, 0)
	if err != nil {
		return err
	}

	_, err = io.Copy(w, tmp)
	return err
}

// Copy the header (and not the data) from the CHDU associated with infptr to the CHDU associated with outfptr. If the current output HDU is not completely empty, then the CHDU will be closed and a new HDU will be appended to the output file. An empty output data unit will be created with all values initially = 0).
func (f *File) CopyHeader(out *File) error {
	c_status := C.int(0)
	C.fits_copy_header(f.c, out.c, &c_status)
	if c_status > 0 {
		return to_err(c_status)
	}
	return nil
}

// Delete the CHDU in the FITS file. Any following HDUs will be shifted forward in the file, to fill in the gap created by the deleted HDU. In the case of deleting the primary array (the first HDU in the file) then the current primary array will be replace by a null primary array containing the minimum set of required keywords and no data. If there are more extensions in the file following the one that is deleted, then the the CHDU will be redefined to point to the following extension. If there are no following extensions then the CHDU will be redefined to point to the previous HDU. The output hdutype parameter returns the type of the new CHDU. A null pointer may be given for hdutype if the returned value is not needed.
func (f *File) DeleteHdu() (HduType, error) {
	c_hdu := C.int(0)
	c_status := C.int(0)
	C.fits_delete_hdu(f.c, &c_hdu, &c_status)
	if c_status > 0 {
		return 0, to_err(c_status)
	}
	return HduType(c_hdu), nil
}

// EOF
