package cfitsio

import (
	"reflect"
	"testing"
)

type any interface{}

var g_tables = []struct {
	fname string
	hdus  []HDU
	tuple [][][]any
}{
	{
		fname: "testdata/swp06542llg.fits",
		hdus: []HDU{
			&PrimaryHDU{ImageHDU{
				header: NewHeader(
					[]Card{
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
							Name:    "FILENAME",
							Value:   "swp06542llg",
							Comment: "original name of input file",
						},
						{
							Name:    "TELESCOP",
							Value:   "IUE",
							Comment: "International Ultraviolet Explorer",
						},
						{
							Name:    "ORIGIN",
							Value:   "GODDARD",
							Comment: "Tape writing location",
						},
						{
							Name:    "CAMERA",
							Value:   int64(3),
							Comment: "IUE camera number",
						},
						{
							Name:    "IMAGE",
							Value:   int64(6542),
							Comment: "IUE image sequence number",
						},
						{
							Name:    "APERTURE",
							Value:   "",
							Comment: "Aperture",
						},
						{
							Name:    "DISPERSN",
							Value:   "LOW",
							Comment: "IUE spectrograph dispersion",
						},
						{
							Name:    "DATE-OBS",
							Value:   "nn/nn/nn",
							Comment: "Observation date (dd/mm/yy)",
						},
						{
							Name:    "DATE-PRO",
							Value:   "nn/nn/nn",
							Comment: "Processing date (dd/mm/yy)",
						},
						{
							Name:    "DATE",
							Value:   "18-Feb-1993",
							Comment: "Date file was written (dd/mm/yy)",
						},
						{
							Name:    "RA",
							Value:   0.0,
							Comment: "Right Ascension in degrees",
						},
						{
							Name:    "DEC",
							Value:   0.0,
							Comment: "Declination in degrees",
						},
						{
							Name:    "EQUINOX",
							Value:   1950.0,
							Comment: "Epoch for coordinates (years)",
						},
						{
							Name:    "THDA-RES",
							Value:   0.0,
							Comment: "THDA at time of read",
						},
						{
							Name:    "THDA-SPE",
							Value:   0.0,
							Comment: "THDA at end of exposure",
						},
					},
					IMAGE_HDU,
					8,
					[]int64{},
				),
			}},
			&Table{
				header: NewHeader(
					[]Card{
						{
							Name:    "XTENSION",
							Value:   "BINTABLE",
							Comment: "Extension type",
						},
						{
							Name:    "BITPIX",
							Value:   int64(8),
							Comment: "binary data",
						},
						{
							Name:    "NAXIS",
							Value:   int64(2),
							Comment: "Number of Axes",
						},
						{
							Name:    "NAXIS1",
							Value:   int64(7532),
							Comment: "width of table in bytes",
						},
						{
							Name:    "NAXIS2",
							Value:   int64(1),
							Comment: "Number of entries in table",
						},
						{
							Name:    "PCOUNT",
							Value:   int64(0),
							Comment: "Number of parameters/group",
						},
						{
							Name:    "GCOUNT",
							Value:   int64(1),
							Comment: "Number of groups",
						},
						{
							Name:    "TFIELDS",
							Value:   int64(9),
							Comment: "Number of fields in each row",
						},
						{
							Name:    "EXTNAME",
							Value:   "IUE MELO",
							Comment: "name of table (?)",
						},
						{
							Name:    "TFORM1",
							Value:   "1I",
							Comment: "Count and data type of field 1",
						},
						{
							Name:    "TTYPE1",
							Value:   "ORDER",
							Comment: "spectral order (low dispersion = 1)",
						},
						{
							Name:    "TUNIT1",
							Value:   "",
							Comment: "unitless",
						},
						{
							Name:    "TFORM2",
							Value:   "1I",
							Comment: "field 2 has one 2-byte integer",
						},
						{
							Name:    "TTYPE2",
							Value:   "NPTS",
							Comment: "number of non-zero points in each vector",
						},
						{
							Name:    "TUNIT2",
							Value:   "",
							Comment: "unitless",
						},
						{
							Name:    "TFORM3",
							Value:   "1E",
							Comment: "Count and data type of field 3",
						},
						{
							Name:    "TTYPE3",
							Value:   "LAMBDA",
							Comment: "3rd field is starting wavelength",
						},
						{
							Name:    "TUNIT3",
							Value:   "ANGSTROM",
							Comment: "unit is angstrom",
						},
						{
							Name:    "TFORM4",
							Value:   "1E",
							Comment: "Count and Type of 4th field",
						},
						{
							Name:    "TTYPE4",
							Value:   "DELTAW",
							Comment: "4th field is wavelength increment",
						},
						{
							Name:    "TUNIT4",
							Value:   "ANGSTROM",
							Comment: "unit is angstrom",
						},
						{
							Name:    "TFORM5",
							Value:   "376E",
							Comment: "Count and Type of 5th field",
						},
						{
							Name:    "TTYPE5",
							Value:   "GROSS",
							Comment: "5th field is gross flux array",
						},
						{
							Name:    "TUNIT5",
							Value:   "FN",
							Comment: "unit is IUE FN",
						},
						{
							Name:    "TFORM6",
							Value:   "376E",
							Comment: "Count and Type of 6th field",
						},
						{
							Name:    "TTYPE6",
							Value:   "BACK",
							Comment: "6th field is background flux array",
						},
						{
							Name:    "TUNIT6",
							Value:   "FN",
							Comment: "unit is IUE FN",
						},
						{
							Name:    "TFORM7",
							Value:   "376E",
							Comment: "Count and Type of 7th field",
						},
						{
							Name:    "TTYPE7",
							Value:   "NET",
							Comment: "7th field is net flux array",
						},
						{
							Name:    "TUNIT7",
							Value:   "ERGS",
							Comment: "unit is IUE FN",
						},
						{
							Name:    "TFORM8",
							Value:   "376E",
							Comment: "Count and Type of 8th field",
						},
						{
							Name:    "TTYPE8",
							Value:   "ABNET",
							Comment: "absolutely calibrated net flux array",
						},
						{
							Name:    "TUNIT8",
							Value:   "ERGS",
							Comment: "unit is ergs/cm2/sec/angstrom",
						},
						{
							Name:    "TFORM9",
							Value:   "376E",
							Comment: "Count and Type of 9th field",
						},
						{
							Name:    "TTYPE9",
							Value:   "EPSILONS",
							Comment: "9th field is epsilons",
						},
						{
							Name:    "TUNIT9",
							Value:   "",
							Comment: "unitless",
						},
					},
					BINARY_TBL,
					8,
					[]int64{},
				),
			},
		},
		tuple: [][][]any{
			nil,
			[][]any{
				// row-0
				[]any{
					int16(1),
					int16(376), float32(1000.8), float32(2.6515958), float32(19286.426),
					float32(22999.92), float32(1001.04297), float32(1001.04297),
					float32(88),
				},
			},
		},
	},
	{
		fname: "testdata/file001.fits",
		hdus: []HDU{
			&PrimaryHDU{ImageHDU{
				header: NewHeader(
					[]Card{
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
							Name:    "EXTEND",
							Value:   true,
							Comment: "THERE IS AN EXTENSION",
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
					IMAGE_HDU,
					8,
					[]int64{},
				),
			}},
			&Table{
				header: NewHeader(
					[]Card{
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
							Name:    "PCOUNT",
							Value:   int64(0),
							Comment: "RANDOM PARAMETER COUNT",
						},
						{
							Name:    "GCOUNT",
							Value:   int64(1),
							Comment: "GROUP COUNT",
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
					ASCII_TBL,
					8,
					[]int64{98, 10},
				),
			},
		},
		tuple: [][][]any{
			nil,
			[][]any{
				[]any{
					-1116.59, 11.28, 59.5667, float64(3),
					77.89999, float64(12), float64(0),
				},
				[]any{
					-1109.54, 11.15667, 54.300000000000004, float64(3),
					float64(40), float64(12), float64(0),
				},
				[]any{
					-340.285, 3.6683, -28.0167, float64(3),
					float64(70), float64(30), 4060.0000000000005,
				},
				[]any{
					536.0, 1.3933300000000002, 34.449999999999996,
					3.5, 222.89999999999998, 59.9, float64(5160),
				},
				[]any{
					float64(3177), 10.23, 21.366699999999998, float64(3),
					99.60001000000001, float64(36), float64(1220),
				},
				[]any{
					float64(3627), 11.29333, 13.2667, 3.2, 522.5999999999999, 59.9, float64(697),
				},
				[]any{
					float64(3756), 11.56667, 54.566700000000004, 4.2,
					261.90000000000003, 56.699999999999996, float64(1071),
				},
				[]any{
					float64(5457), 14.025000000000002, 54.583000000000006, 6.2, 1614.8999999999999,
					12.2, 266.0,
				},
				[]any{
					7292.0, 22.435, 30.049999999999997, 9.5, 128.3, 35.6, 934.0,
				},
				[]any{
					14237.0, 13.3633, -20.867, 3.5, 117.0, 42.199999999999996, 0.0,
				},
			},
		},
	},
}

func TestHDU(t *testing.T) {
	for _, table := range g_tables {
		f, err := Open(table.fname, ReadOnly)
		if err != nil {
			t.Fatalf("could not open FITS file [%s]: %v", table.fname, err)
		}
		defer f.Close()

		nhdus, err := f.NumHDUs()
		if err != nil {
			t.Fatalf("error num-hdus [%v]: %v", table.fname, err)
		}
		if nhdus != len(table.hdus) {
			t.Fatalf("file [%v]: expected %v hdus. got %v", table.fname, nhdus, len(table.hdus))
		}

		for i := 0; i < len(table.hdus); i++ {
			ref := table.hdus[i]
			hdu := f.HDUs()[i]

			if hdu.Type() != ref.Type() {
				t.Fatalf("expected HDU-type [%v]. got [%v] (fname=%v)", ref.Type(), hdu.Type(), table.fname)
			}

			if hdu.Name() != ref.Name() {
				t.Fatalf("expected HDU-name [%v]. got [%v] (fname=%v)", ref.Name(), hdu.Name(), table.fname)
			}

			if hdu.Version() != ref.Version() {
				t.Fatalf("expected HDU-version [%v]. got [%v] (fname=%v)", ref.Version(), hdu.Version(), table.fname)
			}

			xhdr := hdu.Header()
			rhdr := ref.Header()
			if len(xhdr.slice) != len(rhdr.slice) {
				t.Fatalf("#cards differ: ref=%v chk=%v (fname=%v)", len(rhdr.slice), len(xhdr.slice), table.fname)
			}
			for ii := 0; ii < len(rhdr.slice); ii++ {
				if !reflect.DeepEqual(xhdr.slice[ii], rhdr.slice[ii]) {
					t.Fatalf("cards differ (fname=%v).\nexpected:\n%v\ngot:\n%v", table.fname, rhdr.slice[ii], xhdr.slice[ii])
				}

			}
		}

	}
}

// EOF
