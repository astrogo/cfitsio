package cfitsio_test

import (
	"reflect"
	"testing"

	cfitsio "github.com/sbinet/go-cfitsio"
)

func TestKeyword(t *testing.T) {
	for _, table := range []struct {
		fname string
		keys  [][]cfitsio.Keyword
	}{
		{
			fname: "testdata/swp06542llg.fits",
			keys: [][]cfitsio.Keyword{
				{
					{
						Name:    "SIMPLE",
						Value:   true,
						Comment: "Standard FITS format",
					},
					{
						Name:    "BITPIX",
						Value:   int64(8),
						Comment: "",
					},
					{
						Name:    "NAXIS",
						Value:   int64(0),
						Comment: "no data in main file",
					},
					{
						Name:    "EXTEND",
						Value:   true,
						Comment: "Extensions may exist",
					},
					{
						Name:    "RA",
						Value:   0.0,
						Comment: "Right Ascension in degrees",
					},
					{
						Name:    "EQUINOX",
						Value:   1950.0,
						Comment: "Epoch for coordinates (years)",
					},
				},
			},
		},
		{
			fname: "testdata/file001.fits",
			keys: [][]cfitsio.Keyword{
				{
					{
						Name:    "SIMPLE",
						Value:   true,
						Comment: "STANDARD FITS FORMAT (REV OCT 1981)",
					},
					{
						Name:    "BITPIX",
						Value:   int64(8),
						Comment: "CHARACTER INFORMATION",
					},
					{
						Name:    "NAXIS",
						Value:   int64(0),
						Comment: "NO IMAGE DATA ARRAY PRESENT",
					},
					{
						Name:    "ORIGIN",
						Value:   "ESO",
						Comment: "EUROPEAN SOUTHERN OBSERVATORY",
					},
					{
						Name:    "OBJECT",
						Value:   "SNG - CAT.",
						Comment: "THE IDENTIFIER",
					},
					{
						Name:    "DATE",
						Value:   "27/ 5/84",
						Comment: "DATE THIS TAPE WRITTEN DD/MM/YY",
					},
				},
				{
					{
						Name:    "XTENSION",
						Value:   "TABLE",
						Comment: "TABLE EXTENSION",
					},
					{
						Name:    "BITPIX",
						Value:   int64(8),
						Comment: "CHARACTER INFORMATION",
					},
					{
						Name:    "NAXIS",
						Value:   int64(2),
						Comment: "SIMPLE 2-D MATRIX",
					},
					{
						Name:    "NAXIS1",
						Value:   int64(98),
						Comment: "NO. OF CHARACTERS PER ROW",
					},
					{
						Name:    "NAXIS2",
						Value:   int64(10),
						Comment: "NO. OF ROWS",
					},
					{
						Name:    "TFIELDS",
						Value:   int64(7),
						Comment: "NO. OF FIELDS PER ROW",
					},
					{
						Name:    "TTYPE1",
						Value:   "IDEN.",
						Comment: "NAME OF ROW",
					},
					{
						Name:    "TBCOL1",
						Value:   int64(1),
						Comment: "BEGINNING COLUMN OF THE FIELD",
					},
					{
						Name:    "TFORM1",
						Value:   "E14.7",
						Comment: "FORMAT",
					},
					{
						Name:    "TNULL1",
						Value:   "",
						Comment: "NULL VALUE",
					},
					{
						Name:    "TTYPE2",
						Value:   "RA",
						Comment: "NAME OF ROW",
					},
					{
						Name:    "TBCOL2",
						Value:   int64(15),
						Comment: "BEGINNING COLUMN OF THE FIELD",
					},
					{
						Name:    "TFORM2",
						Value:   "E14.7",
						Comment: "FORMAT",
					},
					{
						Name:    "TNULL2",
						Value:   "",
						Comment: "NULL VALUE",
					},
					{
						Name:    "TTYPE3",
						Value:   "DEC",
						Comment: "NAME OF ROW",
					},
					{
						Name:    "TBCOL3",
						Value:   int64(29),
						Comment: "BEGINNING COLUMN OF THE FIELD",
					},
					{
						Name:    "TFORM3",
						Value:   "E14.7",
						Comment: "FORMAT",
					},
					{
						Name:    "TNULL3",
						Value:   "",
						Comment: "NULL VALUE",
					},
					{
						Name:    "TTYPE4",
						Value:   "TYPE",
						Comment: "NAME OF ROW",
					},
					{
						Name:    "TBCOL4",
						Value:   int64(43),
						Comment: "BEGINNING COLUMN OF THE FIELD",
					},
					{
						Name:    "TFORM4",
						Value:   "E14.7",
						Comment: "FORMAT",
					},
					{
						Name:    "TNULL4",
						Value:   "",
						Comment: "NULL VALUE",
					},
					{
						Name:    "TTYPE5",
						Value:   "D25",
						Comment: "NAME OF ROW",
					},
					{
						Name:    "TBCOL5",
						Value:   int64(57),
						Comment: "BEGINNING COLUMN OF THE FIELD",
					},
					{
						Name:    "TFORM5",
						Value:   "E14.7",
						Comment: "FORMAT",
					},
					{
						Name:    "TNULL5",
						Value:   "",
						Comment: "NULL VALUE",
					},
					{
						Name:    "TTYPE6",
						Value:   "INCL.",
						Comment: "NAME OF ROW",
					},
					{
						Name:    "TBCOL6",
						Value:   int64(71),
						Comment: "BEGINNING COLUMN OF THE FIELD",
					},
					{
						Name:    "TFORM6",
						Value:   "E14.7",
						Comment: "FORMAT",
					},
					{
						Name:    "TNULL6",
						Value:   "",
						Comment: "NULL VALUE",
					},
					{
						Name:    "TTYPE7",
						Value:   "RV",
						Comment: "NAME OF ROW",
					},
					{
						Name:    "TBCOL7",
						Value:   int64(85),
						Comment: "BEGINNING COLUMN OF THE FIELD",
					},
					{
						Name:    "TFORM7",
						Value:   "E14.7",
						Comment: "FORMAT",
					},
					{
						Name:    "TNULL7",
						Value:   "",
						Comment: "NULL VALUE",
					},
				},
			},
		},
	} {
		fname := table.fname
		f, err := cfitsio.OpenFile(fname, cfitsio.ReadOnly)
		if err != nil {
			t.Fatalf("could not open FITS file [%s]: %v", fname, err)
		}
		defer f.Close()
		for ihdu, keys := range table.keys {
			hdu, err := f.Hdu(ihdu)
			if err != nil {
				t.Fatalf("error getting hdu [%v]: %v (fname=%v)", ihdu, err, fname)
			}
			for _, ref := range keys {
				key, err := hdu.KeywordByName(ref.Name)
				if err != nil {
					t.Fatalf("error getting key [%v]: %v (fname=%v)", ref.Name, err, fname)
				}
				if !reflect.DeepEqual(key.Value, ref.Value) {
					t.Fatalf("expected Key.Value [%v]. got [%v] (name=%v, fname=%v)", ref.Value, key.Value, ref.Name, fname)
				}
				if !reflect.DeepEqual(key.Comment, ref.Comment) {
					t.Fatalf("expected Key.Comment [%v]. got [%v] (name=%v, fname=%v)", ref.Comment, key.Comment, ref.Name, fname)
				}
			}
		}
	}
}

// EOF
