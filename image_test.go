package cfitsio

import (
	"fmt"
	"io/ioutil"
	"math"
	"os"
	"reflect"
	"testing"
)

func TestImageRW(t *testing.T) {
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

	for ii, table := range []struct {
		name     string
		version  int
		cards    []Card
		bitpix   int64
		bzero    uint64
		unsigned bool
		axes     []int64
		image    interface{}
	}{
		{
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
			},
			bitpix: 8,
			axes:   []int64{3, 4},
			image: []int8{
				0, 1, 2, 3,
				4, 5, 6, 7,
				8, 9, 0, 1,
			},
		},
		{
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
			},
			bitpix: 16,
			axes:   []int64{3, 4},
			image: []int16{
				0, 1, 2, 3,
				4, 5, 6, 7,
				8, 9, 0, 1,
			},
		},
		{
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
			},
			bitpix:   16,
			bzero:    -math.MinInt16,
			unsigned: true,
			axes:     []int64{3, 4},
			image: []uint16{
				0, 1, 2, 3,
				4, 5, 6, 7,
				8, 9, 0, math.MaxUint16,
			},
		},
		{
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
			},
			bitpix: 32,
			axes:   []int64{3, 4},
			image: []int32{
				0, 1, 2, 3,
				4, 5, 6, 7,
				8, 9, 0, 1,
			},
		},
		{
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
			},
			bitpix:   32,
			bzero:    -math.MinInt32,
			unsigned: true,
			axes:     []int64{3, 4},
			image: []uint32{
				0, 1, 2, 3,
				4, 5, 6, 7,
				8, 9, 0, math.MaxUint32,
			},
		},
		{
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
			},
			bitpix: 64,
			axes:   []int64{3, 4},
			image: []int64{
				0, 1, 2, 3,
				4, 5, 6, 7,
				8, 9, 0, 1,
			},
		},
		{
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
			},
			bitpix:   64,
			bzero:    -math.MinInt64,
			unsigned: true,
			axes:     []int64{3, 4},
			image: []uint64{
				0, 1, 2, 3,
				4, 5, 6, 7,
				8, 9, 0, math.MaxUint64,
			},
		},
		{
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
			},
			bitpix: -32,
			axes:   []int64{3, 4},
			image: []float32{
				0, 1, 2, 3,
				4, 5, 6, 7,
				8, 9, 0, 1,
			},
		},
		{
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
			},
			bitpix: -64,
			axes:   []int64{3, 4},
			image: []float64{
				0, 1, 2, 3,
				4, 5, 6, 7,
				8, 9, 0, 1,
			},
		},
	} {
		fname := fmt.Sprintf("%03d_%s", ii, table.name)
		for i := 0; i < 2; i++ {
			func(i int) {
				var f File
				var err error
				var hdu HDU

				switch i {

				case 0: // create
					f, err = Create(fname)
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
					if table.bzero != 0 {
						phdr.Set("BZERO", table.bzero, "offset data range")
					}
					phdu, err := NewPrimaryHDU(&f, phdr)
					if err != nil {
						t.Fatalf("error creating PHDU: %v", err)
					}
					defer phdu.Close()
					hdu = phdu

					err = phdu.(*PrimaryHDU).Write(&table.image)
					if err != nil {
						t.Fatalf("error writing image: %v", err)
					}

				case 1: // read
					f, err = Open(fname, ReadOnly)
					if err != nil {
						t.Fatalf("error opening file [%v]: %v", fname, err)
					}
					defer f.Close()

					hdu = f.HDU(0)
					hdr := hdu.Header()
					nelmts := 1
					for _, axe := range hdr.Axes() {
						nelmts *= int(axe)
					}

					var data interface{}
					switch hdr.Bitpix() {
					case 8:
						v := make([]int8, nelmts)
						data = v
						err = hdu.Data(&v)

					case 16:
						if table.unsigned {
							v := make([]uint16, nelmts)
							data = v
							err = hdu.Data(&v)
						} else {
							v := make([]int16, nelmts)
							data = v
							err = hdu.Data(&v)
						}

					case 32:
						if table.unsigned {
							v := make([]uint32, nelmts)
							data = v
							err = hdu.Data(&v)
						} else {
							v := make([]int32, nelmts)
							data = v
							err = hdu.Data(&v)
						}

					case 64:
						if table.unsigned {
							v := make([]uint64, nelmts)
							data = v
							err = hdu.Data(&v)
						} else {
							v := make([]int64, nelmts)
							data = v
							err = hdu.Data(&v)
						}

					case -32:
						v := make([]float32, nelmts)
						data = v
						err = hdu.Data(&v)

					case -64:
						v := make([]float64, nelmts)
						data = v
						err = hdu.Data(&v)
					}

					if err != nil {
						t.Fatalf("error reading image: %v", err)
					}

					if !reflect.DeepEqual(data, table.image) {
						t.Fatalf("expected image:\nref=%v\ngot=%v", table.image, data)
					}
				}

				hdr := hdu.Header()
				if hdr.bitpix != table.bitpix {
					t.Fatalf("expected BITPIX=%v. got %v", table.bitpix, hdr.bitpix)
				}

				if !reflect.DeepEqual(hdr.Axes(), table.axes) {
					t.Fatalf("expected AXES==%v. got %v (i=%v)", table.axes, hdr.Axes(), i)
				}

				name := hdu.Name()
				if name != "primary hdu" {
					t.Fatalf("expected EXTNAME==%q. got %q", "primary hdu", name)
				}

				vers := hdu.Version()
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
			}(i)
		}
	}
}

// EOF
