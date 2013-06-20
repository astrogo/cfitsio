package cfitsio

// #include "fitsio.h"
// #include <stdlib.h>
// #include <string.h>
// #include "go-cfitsio-utils.h"
import "C"

import (
	"fmt"
	"unsafe"
)

func Version() float32 {
	v := C.float(0)
	C.ffvers(&v)
	return float32(v)
}

func to_err(sc C.int) error {
	c_err := C.char_buf_array(C.FLEN_ERRMSG)
	defer C.free(unsafe.Pointer(c_err))
	C.ffgerr(sc, c_err)
	err := fmt.Errorf("%s", C.GoString(c_err))
	return err
}
// eof
