package cfitsio_test

import (
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

// EOF
