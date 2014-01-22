package cfitsio

// #include "go-cfitsio.h"
import "C"

func init() {
	c_status := C.fits_init_cfitsio()
	if c_status > 0 {
		err := to_err(c_status)
		panic(err)
	}
}
