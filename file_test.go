package cfitsio_test

import (
	"bytes"
	"strings"
	"testing"

	cfitsio "github.com/sbinet/go-cfitsio"
)

func TestOpenFile(t *testing.T) {
	const fname = "testdata/file001.fits"
	f, err := cfitsio.OpenFile(fname, cfitsio.ReadOnly)
	if err != nil {
		t.Fatalf("could not open FITS file [%s]: %v", fname, err)
	}
	defer f.Close()

	fmode, err := f.Mode()
	if err != nil {
		t.Fatalf("error mode: %v", err)
	}
	if fmode != cfitsio.ReadOnly {
		t.Fatalf("expected file-mode [%v]. got [%v]", cfitsio.ReadOnly, fmode)
	}

	name, err := f.Name()
	if err != nil {
		t.Fatalf("error name: %v", err)
	}
	if name != fname {
		t.Fatalf("expected file-name [%v]. got [%v]", fname, name)
	}

	furl, err := f.UrlType()
	if err != nil {
		t.Fatalf("error url: %v", err)
	}
	if furl != "file://" {
		t.Fatalf("expected file-url [%v]. got [%v]", "file://", furl)
	}

}

func TestHdu(t *testing.T) {
	const fname = "testdata/file001.fits"
	f, err := cfitsio.OpenFile(fname, cfitsio.ReadOnly)
	if err != nil {
		t.Fatalf("could not open FITS file [%s]: %v", fname, err)
	}
	defer f.Close()

	nhdus, err := f.NumHdus()
	if err != nil {
		t.Fatalf("error hdu: %v", err)
	}

	if nhdus != 2 {
		t.Fatalf("expected #hdus [%v]. got [%v]", 2, nhdus)
	}

	ihdu := f.HduNum()
	if ihdu != 1 {
		t.Fatalf("expected hdu number [%v]. got [%v]", 1, ihdu)
	}

	hdutype, err := f.HduType()
	if err != nil {
		t.Fatalf("error hdu-type: %v", err)
	}
	if hdutype != cfitsio.IMAGE_HDU {
		t.Fatalf("expected hdu type [%v]. got [%v]", cfitsio.IMAGE_HDU, hdutype)
	}

	buf := bytes.NewBuffer(nil)
	err = f.WriteHdu(buf)
	if err != nil {
		t.Fatalf("error write-hdu: %v", err)
	}
	hdu_tmp := strings.Trim(string(buf.Bytes()), " ")
	const hdu_ref = `SIMPLE  =                    T          / STANDARD FITS FORMAT (REV OCT 1981)   BITPIX  =                    8          / CHARACTER INFORMATION                 NAXIS   =                    0          / NO IMAGE DATA ARRAY PRESENT           EXTEND  =                    T          / THERE IS AN EXTENSION                 ORIGIN  = 'ESO     '                    / EUROPEAN SOUTHERN OBSERVATORY         OBJECT  = 'SNG - CAT.      '            / THE IDENTIFIER                        DATE    = '27/ 5/84'                    / DATE THIS TAPE WRITTEN DD/MM/YY       END`
	if hdu_ref != hdu_tmp {
		t.Fatalf("expected hdu-write:\n%v\ngot:\n%v", hdu_ref, hdu_tmp)
	}
}

// EOF
