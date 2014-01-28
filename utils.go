package cfitsio

import (
	"fmt"
	"strings"
)

// strIsContinued checks whether the last non-whitespace char in the string
// is an ampersand, which indicates the value string is to be continued on the
// next line
func strIsContinued(v string) bool {
	vv := strings.Trim(v, " \n\t'")
	fmt.Printf("vv=%q\n", vv)
	if len(vv) == 0 {
		return false
	}
	return vv[len(vv)-1] == '&'
}

// EOF
