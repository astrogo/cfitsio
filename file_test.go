package cfitsio_test

import (
	"bytes"
	"io/ioutil"
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

func TestAsciiTbl(t *testing.T) {
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
	hdu_ref, err := ioutil.ReadFile("testdata/ref.file001.hdu")
	if err != nil {
		t.Fatalf("error reading ref file: %v", err)
	}
	if string(hdu_ref) != hdu_tmp {
		t.Fatalf("expected hdu-write:\n%v\ngot:\n%v", string(hdu_ref), hdu_tmp)
	}

	// move to next header
	_, err = f.MovRelHdu(1)
	if err != nil {
		t.Fatalf("error next-hdu: %v", err)
	}

	ihdu = f.HduNum()
	if ihdu != 2 {
		t.Fatalf("expected hdu number [%v]. got [%v]", 2, ihdu)
	}

	hdutype, err = f.HduType()
	if err != nil {
		t.Fatalf("error hdu-type: %v", err)
	}
	if hdutype != cfitsio.ASCII_TBL {
		t.Fatalf("expected hdu type [%v]. got [%v]", cfitsio.ASCII_TBL, hdutype)
	}

	// move to first header
	_, err = f.MovAbsHdu(1)
	if err != nil {
		t.Fatalf("error abs-hdu: %v", err)
	}

	ihdu = f.HduNum()
	if ihdu != 1 {
		t.Fatalf("expected hdu number [%v]. got [%v]", 1, ihdu)
	}

	hdutype, err = f.HduType()
	if err != nil {
		t.Fatalf("error hdu-type: %v", err)
	}
	if hdutype != cfitsio.IMAGE_HDU {
		t.Fatalf("expected hdu type [%v]. got [%v]", cfitsio.IMAGE_HDU, hdutype)
	}

	// move to second header
	_, err = f.MovAbsHdu(2)
	if err != nil {
		t.Fatalf("error abs-hdu: %v", err)
	}

	ihdu = f.HduNum()
	if ihdu != 2 {
		t.Fatalf("expected hdu number [%v]. got [%v]", 2, ihdu)
	}

	hdutype, err = f.HduType()
	if err != nil {
		t.Fatalf("error hdu-type: %v", err)
	}
	if hdutype != cfitsio.ASCII_TBL {
		t.Fatalf("expected hdu type [%v]. got [%v]", cfitsio.ASCII_TBL, hdutype)
	}

	// move to non-existent header
	_, err = f.MovAbsHdu(3)
	if err == nil {
		t.Fatalf("expected error")
	}
	if err != cfitsio.END_OF_FILE {
		t.Fatalf("expected error [%v]. got [%v]", cfitsio.END_OF_FILE, err)
	}
	ihdu = f.HduNum()
	if ihdu != 2 {
		t.Fatalf("expected hdu number [%v]. got [%v]", 2, ihdu)
	}

	// move to non-existent header
	_, err = f.MovAbsHdu(0)
	if err == nil {
		t.Fatalf("expected error")
	}
	if err != cfitsio.BAD_HDU_NUM {
		t.Fatalf("expected error [%v]. got [%v]", cfitsio.BAD_HDU_NUM, err)
	}
	ihdu = f.HduNum()
	if ihdu != 2 {
		t.Fatalf("expected hdu number [%v]. got [%v]", 2, ihdu)
	}

	// move to first header
	_, err = f.MovAbsHdu(1)
	if err != nil {
		t.Fatalf("error abs-hdu: %v", err)
	}

	ihdu = f.HduNum()
	if ihdu != 1 {
		t.Fatalf("expected hdu number [%v]. got [%v]", 1, ihdu)
	}

	// move to hdu by name
	err = f.MovNamHdu(cfitsio.ASCII_TBL, "TABLE", 1)
	if err != cfitsio.BAD_HDU_NUM {
		t.Fatalf("expected error [%v]. got [%v]", cfitsio.BAD_HDU_NUM, err)
	}
}

func TestBinTable(t *testing.T) {
	const fname = "testdata/swp06542llg.fits"
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
	hdu_ref, err := ioutil.ReadFile("testdata/ref.swp06542llg.hdu")
	if err != nil {
		t.Fatalf("error reading ref file: %v", err)
	}
	if string(hdu_ref) != hdu_tmp {
		t.Fatalf("expected hdu-write:\n%v\ngot:\n%v", string(hdu_ref), hdu_tmp)
	}

	// move to next header
	_, err = f.MovRelHdu(1)
	if err != nil {
		t.Fatalf("error next-hdu: %v", err)
	}

	ihdu = f.HduNum()
	if ihdu != 2 {
		t.Fatalf("expected hdu number [%v]. got [%v]", 2, ihdu)
	}

	hdutype, err = f.HduType()
	if err != nil {
		t.Fatalf("error hdu-type: %v", err)
	}
	if hdutype != cfitsio.BINARY_TBL {
		t.Fatalf("expected hdu type [%v]. got [%v]", cfitsio.BINARY_TBL, hdutype)
	}

	// move to first header
	_, err = f.MovAbsHdu(1)
	if err != nil {
		t.Fatalf("error abs-hdu: %v", err)
	}

	ihdu = f.HduNum()
	if ihdu != 1 {
		t.Fatalf("expected hdu number [%v]. got [%v]", 1, ihdu)
	}

	hdutype, err = f.HduType()
	if err != nil {
		t.Fatalf("error hdu-type: %v", err)
	}
	if hdutype != cfitsio.IMAGE_HDU {
		t.Fatalf("expected hdu type [%v]. got [%v]", cfitsio.IMAGE_HDU, hdutype)
	}

	// move to second header
	_, err = f.MovAbsHdu(2)
	if err != nil {
		t.Fatalf("error abs-hdu: %v", err)
	}

	ihdu = f.HduNum()
	if ihdu != 2 {
		t.Fatalf("expected hdu number [%v]. got [%v]", 2, ihdu)
	}

	hdutype, err = f.HduType()
	if err != nil {
		t.Fatalf("error hdu-type: %v", err)
	}
	if hdutype != cfitsio.BINARY_TBL {
		t.Fatalf("expected hdu type [%v]. got [%v]", cfitsio.BINARY_TBL, hdutype)
	}

	// move to non-existent header
	_, err = f.MovAbsHdu(3)
	if err == nil {
		t.Fatalf("expected error")
	}
	if err != cfitsio.END_OF_FILE {
		t.Fatalf("expected error [%v]. got [%v]", cfitsio.END_OF_FILE, err)
	}
	ihdu = f.HduNum()
	if ihdu != 2 {
		t.Fatalf("expected hdu number [%v]. got [%v]", 2, ihdu)
	}

	// move to non-existent header
	_, err = f.MovAbsHdu(0)
	if err == nil {
		t.Fatalf("expected error")
	}
	if err != cfitsio.BAD_HDU_NUM {
		t.Fatalf("expected error [%v]. got [%v]", cfitsio.BAD_HDU_NUM, err)
	}
	ihdu = f.HduNum()
	if ihdu != 2 {
		t.Fatalf("expected hdu number [%v]. got [%v]", 2, ihdu)
	}

	// move to first header
	_, err = f.MovAbsHdu(1)
	if err != nil {
		t.Fatalf("error abs-hdu: %v", err)
	}

	ihdu = f.HduNum()
	if ihdu != 1 {
		t.Fatalf("expected hdu number [%v]. got [%v]", 1, ihdu)
	}

	// move to hdu by name
	err = f.MovNamHdu(cfitsio.BINARY_TBL, "IUE MELO", 1)
	if err != nil {
		t.Fatalf("error hdu-nam: %v", err)
	}

	ihdu = f.HduNum()
	if ihdu != 2 {
		t.Fatalf("expected hdu number [%v]. got [%v]", 2, ihdu)
	}

}

// EOF
