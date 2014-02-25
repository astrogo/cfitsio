package cfitsio

import (
	"io/ioutil"
	"os"
	"testing"
)

func TestCreatePHDU(t *testing.T) {
	curdir, err := os.Getwd()
	if err != nil {
		t.Fatalf(err.Error())
	}
	defer os.Chdir(curdir)

	workdir, err := ioutil.TempDir("", "go-cfitsio-test-")
	if err != nil {
		t.Fatalf(err.Error())
	}
	defer os.RemoveAll(workdir)

	err = os.Chdir(workdir)
	if err != nil {
		t.Fatalf(err.Error())
	}

	fname := "new.fits"
	for _, fct := range []func(){
		// create
		func() {
			f, err := Create(fname)
			if err != nil {
				t.Fatalf("error creating new file [%v]: %v", fname, err)
			}
			defer f.Close()

			phdu, err := NewPrimaryHDU(&f, NewDefaultHeader())
			if err != nil {
				t.Fatalf("error creating PHDU: %v", err)
			}
			defer phdu.Close()

			hdr := phdu.Header()
			if hdr.bitpix != 8 {
				t.Fatalf("expected BITPIX=%v. got %v", 8, hdr.bitpix)
			}
		},
		// read-back
		func() {
			f, err := Open(fname, ReadOnly)
			if err != nil {
				t.Fatalf("error opening file [%v]: %v", fname, err)
			}
			defer f.Close()

			hdu := f.HDU(0)
			hdr := hdu.Header()
			if hdr.bitpix != 8 {
				t.Fatalf("expected BITPIX=%v. got %v", 8, hdr.bitpix)
			}
		},
	} {
		fct()
	}

}
