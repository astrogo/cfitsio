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

// Return a descriptive text string (30 char max.) corresponding to a CFITSIO error status code.
func to_err(sc C.int) error {
	c_err := C.char_buf_array(C.FLEN_ERRMSG)
	defer C.free(unsafe.Pointer(c_err))
	C.ffgerr(sc, c_err)
	err := fmt.Errorf("%s", C.GoString(c_err))
	return err
}

// eof
