package cfitsio_test

import (
	"reflect"
	"testing"

	cfitsio "github.com/sbinet/go-cfitsio"
)

func TestHdu(t *testing.T) {
	for _, table := range []struct {
		fname string
		hdus  []cfitsio.Hdu
	}{
		{
			fname: "testdata/swp06542llg.fits",
			hdus: []cfitsio.Hdu{
				{
					Id:     0,
					Type:   cfitsio.IMAGE_HDU,
					File:   nil,
					Bitpix: 8,
					Axes:   nil,
				},
				{
					Id:     1,
					Type:   cfitsio.BINARY_TBL,
					File:   nil,
					Bitpix: 8,
					Axes:   []int64{7532, 1},
				},
			},
		},
		{
			fname: "testdata/file001.fits",
			hdus: []cfitsio.Hdu{
				{
					Id:     0,
					Type:   cfitsio.IMAGE_HDU,
					File:   nil,
					Bitpix: 8,
					Axes:   nil,
				},
				{
					Id:     1,
					Type:   cfitsio.ASCII_TBL,
					File:   nil,
					Bitpix: 8,
					Axes:   []int64{98, 10},
				},
			},
		},
	} {
		f, err := cfitsio.OpenFile(table.fname, cfitsio.ReadOnly)
		if err != nil {
			t.Fatalf("could not open FITS file [%s]: %v", table.fname, err)
		}
		defer f.Close()

		nhdus, err := f.NumHdus()
		if err != nil {
			t.Fatalf("error num-hdus [%v]: %v", table.fname, err)
		}
		if nhdus != len(table.hdus) {
			t.Fatalf("file [%v]: expected %v hdus. got %v", table.fname, nhdus, len(table.hdus))
		}

		for i := 0; i < len(table.hdus); i++ {
			ref := &table.hdus[i]
			hdu, err := f.Hdu(i)
			if err != nil {
				t.Fatalf("error getting %d-th HDU: %v (fname=%v)", i, err, table.fname)
			}
			defer hdu.Delete()

			if hdu.Id != ref.Id {
				t.Fatalf("expected HDU index [%v]. got [%v] (fname=%v)", ref.Id, hdu.Id, table.fname)
			}
			if hdu.Type != ref.Type {
				t.Fatalf("expected HDU-type [%v]. got [%v] (fname=%v)", ref.Type, hdu.Type, table.fname)
			}
			if hdu.Bitpix != ref.Bitpix {
				t.Fatalf("expected bitpix [%v]. got [%v] (fname=%v[%d])", ref.Bitpix, hdu.Bitpix, table.fname, i)
			}
			if !reflect.DeepEqual(hdu.Axes, ref.Axes) {
				t.Fatalf("expected axes=%v. got %v (fname=%v)", ref.Axes, hdu.Axes, table.fname)
			}

			key, err := hdu.KeywordByName("NAXIS")
			if err != nil {
				t.Fatalf("error getting keyword 'NAXIS': %v (fname=%v)", err, table.fname)
			}
			if !reflect.DeepEqual(key.Value, int64(len(ref.Axes))) {
				t.Fatalf("expected Keyword-NAXIS %v. got %v (fname=%v)", int64(len(ref.Axes)), key.Value, table.fname)
			}
		}

	}
}

// EOF
