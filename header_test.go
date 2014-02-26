package cfitsio

import (
	"io/ioutil"
	"os"
	"reflect"
	"testing"
)

func TestHeaderRW(t *testing.T) {
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

	table := struct {
		name    string
		version int
		cards   []Card
		bitpix  int64
		axes    []int64
		image   interface{}
	}{
		name:    "new.fits",
		version: 2,
		cards: []Card{
			{
				"EXTNAME",
				"primary hdu",
				"the primary HDU",
			},
			{
				"EXTVER",
				2,
				"the primary hdu version",
			},
			{
				"card_uint8",
				byte(42),
				"an uint8",
			},
			{
				"card_uint16",
				uint16(42),
				"an uint16",
			},
			{
				"card_uint32",
				uint32(42),
				"an uint32",
			},
			{
				"card_uint64",
				uint64(42),
				"an uint64",
			},
			{
				"card_int8",
				int8(42),
				"an int8",
			},
			{
				"card_int16",
				int16(42),
				"an int16",
			},
			{
				"card_int32",
				int32(42),
				"an int32",
			},
			{
				"card_int64",
				int64(42),
				"an int64",
			},
			{
				"card_int3264",
				int(42),
				"an int",
			},
			{
				"card_uintxx",
				uint(42),
				"an uint",
			},
			{
				"card_float32",
				float32(666),
				"a float32",
			},
			{
				"card_float64",
				float64(666),
				"a float64",
			},
			{
				"card_complex64",
				complex(float32(42), float32(66)),
				"a complex64",
			},
			{
				"card_complex128",
				complex(float64(42), float64(66)),
				"a complex128",
			},
		},
		bitpix: 8,
		axes:   []int64{3, 4},
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

			phdr := NewHeader(
				table.cards,
				IMAGE_HDU,
				table.bitpix,
				table.axes,
			)
			phdu, err := NewPrimaryHDU(&f, phdr)
			if err != nil {
				t.Fatalf("error creating PHDU: %v", err)
			}
			defer phdu.Close()

			hdr := phdu.Header()
			if hdr.bitpix != table.bitpix {
				t.Fatalf("expected BITPIX=%v. got %v", table.bitpix, hdr.bitpix)
			}

			name := phdu.Name()
			if name != "primary hdu" {
				t.Fatalf("expected EXTNAME==%q. got %q", "primary hdu", name)
			}

			vers := phdu.Version()
			if vers != table.version {
				t.Fatalf("expected EXTVER==%v. got %v", table.version, vers)
			}

			card := hdr.Get("EXTNAME")
			if card == nil {
				t.Fatalf("error retrieving card [EXTNAME]")
			}
			if card.Comment != "the primary HDU" {
				t.Fatalf("expected EXTNAME.Comment==%q. got %q", "the primary HDU", card.Comment)
			}

			card = hdr.Get("EXTVER")
			if card == nil {
				t.Fatalf("error retrieving card [EXTVER]")
			}
			if card.Comment != "the primary hdu version" {
				t.Fatalf("expected EXTVER.Comment==%q. got %q", "the primary hdu version", card.Comment)

			}

			for _, ref := range table.cards {
				card := hdr.Get(ref.Name)
				if card == nil {
					t.Fatalf("error retrieving card [%v]", ref.Name)
				}
				rv := reflect.ValueOf(ref.Value)
				var val interface{}
				switch rv.Type().Kind() {
				case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
					val = rv.Int()
				case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
					val = int64(rv.Uint())
				case reflect.Float32, reflect.Float64:
					val = rv.Float()
				case reflect.Complex64, reflect.Complex128:
					val = rv.Complex()
				case reflect.String:
					val = ref.Value.(string)
				}
				if !reflect.DeepEqual(card.Value, val) {
					t.Fatalf(
						"card %q. expected [%v](%T). got [%v](%T)",
						ref.Name,
						val, val,
						card.Value, card.Value,
					)
				}
				if card.Comment != ref.Comment {
					t.Fatalf("card %q. comment differ. expected %q. got %q", ref.Name, ref.Comment, card.Comment)
				}
			}

			card = hdr.Get("NOT THERE")
			if card != nil {
				t.Fatalf("expected no card. got [%v]", card)
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
			if hdr.bitpix != table.bitpix {
				t.Fatalf("expected BITPIX=%v. got %v", 8, hdr.bitpix)
			}

			name := hdu.Name()
			if name != "primary hdu" {
				t.Fatalf("expected EXTNAME==%q. got %q", "primary hdu", name)
			}

			vers := hdu.Version()
			if vers != table.version {
				t.Fatalf("expected EXTVER==%v. got %v", 2, vers)
			}

			card := hdr.Get("EXTNAME")
			if card == nil {
				t.Fatalf("error retrieving card [EXTNAME]")
			}
			if card.Comment != "the primary HDU" {
				t.Fatalf("expected EXTNAME.Comment==%q. got %q", "the primary HDU", card.Comment)
			}

			card = hdr.Get("EXTVER")
			if card == nil {
				t.Fatalf("error retrieving card [EXTVER]")
			}
			if card.Comment != "the primary hdu version" {
				t.Fatalf("expected EXTVER.Comment==%q. got %q", "the primary hdu version", card.Comment)

			}

			for _, ref := range table.cards {
				card := hdr.Get(ref.Name)
				if card == nil {
					t.Fatalf("error retrieving card [%v]", ref.Name)
				}

				rv := reflect.ValueOf(ref.Value)
				var val interface{}
				switch rv.Type().Kind() {
				case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
					val = rv.Int()
				case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
					val = int64(rv.Uint())
				case reflect.Float32, reflect.Float64:
					val = rv.Float()
				case reflect.Complex64, reflect.Complex128:
					val = rv.Complex()
				case reflect.String:
					val = ref.Value.(string)
				}
				if !reflect.DeepEqual(card.Value, val) {
					t.Fatalf(
						"card %q. expected [%v](%T). got [%v](%T)",
						ref.Name,
						val, val,
						card.Value, card.Value,
					)
				}

				if card.Comment != ref.Comment {
					t.Fatalf("card %q. comment differ. expected %q. got %q", ref.Name, ref.Comment, card.Comment)
				}
			}

			card = hdr.Get("NOT THERE")
			if card != nil {
				t.Fatalf("expected no card. got [%v]", card)
			}
		},
	} {
		fct()
	}
}

// EOF
