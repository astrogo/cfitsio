package cfitsio

// #include "go-cfitsio.h"
import "C"

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

// EOF
