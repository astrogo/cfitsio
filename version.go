package cfitsio

// #include "fitsio.h"
import "C"

func Version() float32 {
	v := C.float(0)
	C.ffvers(&v)
	return float32(v)
}

// EOF
