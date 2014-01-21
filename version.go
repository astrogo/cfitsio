package cfitsio

// #include "go-cfitsio.h"
import "C"

// Version returns the revision number of the CFITSIO library. The revision number will be incremented with each new release of CFITSIO.
func Version() float32 {
	v := C.float(0)
	C.ffvers(&v)
	return float32(v)
}

// EOF
