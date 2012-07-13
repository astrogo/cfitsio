package cfitsio

// #include "fitsio.h"
import "C"

func Version() float32 {
	var v C.float(0)
	C.fits_get_version(&v)
	return float32(v)
}

// eof
