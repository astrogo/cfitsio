package cfitsio

import (
	"fmt"
	"io/ioutil"
	"os"
	"reflect"
	"testing"
)

func TestTable(t *testing.T) {
	for _, table := range g_tables {
		fname := table.fname
		f, err := Open(fname, ReadOnly)
		if err != nil {
			t.Fatalf("error opening file [%v]: %v", fname, err)
		}

		for i := range f.HDUs() {
			hdu, ok := f.HDU(i).(*Table)
			if !ok {
				continue
			}
			for irow := int64(0); irow < hdu.NumRows(); irow++ {
				err := hdu.readRow(irow)
				if err != nil {
					t.Fatalf(
						"error reading row [%v] (fname=%v, table=%v): %v",
						irow, fname, hdu.Name(), err,
					)
				}
			}
		}
	}
}

func TestTableNext(t *testing.T) {
	for _, table := range g_tables {
		fname := table.fname
		f, err := Open(fname, ReadOnly)
		if err != nil {
			t.Fatalf("error opening file [%v]: %v", fname, err)
		}

		for i := range f.HDUs() {
			hdu, ok := f.HDU(i).(*Table)
			if !ok {
				continue
			}

			nrows := hdu.NumRows()
			// iter over all rows
			rows, err := hdu.Read(0, nrows)
			if err != nil {
				t.Fatalf("table.Read: %v", err)
			}
			count := int64(0)
			for rows.Next() {
				count++
			}
			err = rows.Err()
			if err != nil {
				t.Fatalf("rows.Err: %v", err)
			}
			if count != nrows {
				t.Fatalf("rows.Next: expected [%d] rows. got %d.", nrows, count)
			}

			// iter over no row
			rows, err = hdu.Read(0, 0)
			if err != nil {
				t.Fatalf("table.Read: %v", err)
			}
			count = int64(0)
			for rows.Next() {
				count++
			}
			err = rows.Err()
			if err != nil {
				t.Fatalf("rows.Err: %v", err)
			}
			if count != 0 {
				t.Fatalf("rows.Next: expected [%d] rows. got %d.", 0, count)
			}

			// iter over 1 row
			rows, err = hdu.Read(0, 1)
			if err != nil {
				t.Fatalf("table.Read: %v", err)
			}
			count = int64(0)
			for rows.Next() {
				count++
			}
			err = rows.Err()
			if err != nil {
				t.Fatalf("rows.Err: %v", err)
			}
			if count != 1 {
				t.Fatalf("rows.Next: expected [%d] rows. got %d.", 1, count)
			}

			// iter over all rows
			rows, err = hdu.Read(0, nrows)
			if err != nil {
				t.Fatalf("table.Read: %v", err)
			}
			count = int64(0)
			for rows.Next() {
				count++
			}
			err = rows.Err()
			if err != nil {
				t.Fatalf("rows.Err: %v", err)
			}
			if count != nrows {
				t.Fatalf("rows.Next: expected [%d] rows. got %d.", nrows, count)
			}

			// iter over all rows +1
			rows, err = hdu.Read(0, nrows+1)
			if err != nil {
				t.Fatalf("table.Read: %v", err)
			}
			count = int64(0)
			for rows.Next() {
				count++
			}
			err = rows.Err()
			if err != nil {
				t.Fatalf("rows.Err: %v", err)
			}
			if count != nrows {
				t.Fatalf("rows.Next: expected [%d] rows. got %d.", nrows, count)
			}

			// iter over all rows -1
			rows, err = hdu.Read(1, nrows)
			if err != nil {
				t.Fatalf("table.Read: %v", err)
			}
			count = int64(0)
			for rows.Next() {
				count++
			}
			err = rows.Err()
			if err != nil {
				t.Fatalf("rows.Err: %v", err)
			}
			if count != nrows-1 {
				t.Fatalf("rows.Next: expected [%d] rows. got %d.", nrows-1, count)
			}

			// iter over [1,1+maxrows -1)
			rows, err = hdu.Read(1, nrows-1)
			if err != nil {
				t.Fatalf("table.Read: %v", err)
			}
			count = int64(0)
			for rows.Next() {
				count++
			}
			err = rows.Err()
			if err != nil {
				t.Fatalf("rows.Err: %v", err)
			}
			exp := nrows - 2
			if exp <= 0 {
				exp = 0
			}
			if count != exp {
				t.Fatalf("rows.Next: expected [%d] rows. got %d.", exp, count)
			}

			// iter over last row
			rows, err = hdu.Read(nrows-1, nrows)
			if err != nil {
				t.Fatalf("table.Read: %v", err)
			}
			count = int64(0)
			for rows.Next() {
				count++
			}
			err = rows.Err()
			if err != nil {
				t.Fatalf("rows.Err: %v", err)
			}
			if count != 1 {
				t.Fatalf("rows.Next: expected [%d] rows. got %d.", 1, count)
			}

		}
	}
}

func TestTableErrScan(t *testing.T) {
	for _, table := range g_tables {
		fname := table.fname
		f, err := Open(fname, ReadOnly)
		if err != nil {
			t.Fatalf("error opening file [%v]: %v", fname, err)
		}

		for i := range f.HDUs() {
			hdu, ok := f.HDU(i).(*Table)
			if !ok {
				continue
			}
			nrows := hdu.NumRows()
			rows, err := hdu.Read(0, nrows)
			if err != nil {
				t.Fatalf("table.Read: %v", err)
			}
			count := int64(0)
			for rows.Next() {
				count++
				err = rows.Scan()
				if err != nil {
					t.Fatalf("rows.Scan: error: %v", err)
				}

				dummy := 0
				err = rows.Scan(&dummy) // none of the tables has only 1 field
				if err == nil {
					t.Fatalf("rows.Scan: expected a failure")
				}
			}
			err = rows.Err()
			if err != nil {
				t.Fatalf("rows.Err: %v", err)
			}
			if count != nrows {
				t.Fatalf("rows.Next: expected [%d] rows. got %d.", nrows, count)
			}
		}
	}
}

func TestTableScan(t *testing.T) {
	for _, table := range g_tables {
		fname := table.fname
		f, err := Open(fname, ReadOnly)
		if err != nil {
			t.Fatalf("error opening file [%v]: %v", fname, err)
		}

		for i := range f.HDUs() {
			hdu, ok := f.HDU(i).(*Table)
			if !ok {
				continue
			}
			nrows := hdu.NumRows()
			rows, err := hdu.Read(0, nrows)
			if err != nil {
				t.Fatalf("table.Read: %v", err)
			}
			count := int64(0)
			for rows.Next() {
				ref := make([]interface{}, len(table.tuple[i][count]))
				data := make([]interface{}, len(ref))
				for ii, vv := range table.tuple[i][count] {
					rt := reflect.TypeOf(vv)
					rv := reflect.New(rt)
					xx := rv.Interface()
					data[ii] = xx
					ref[ii] = vv
				}
				err = rows.Scan(data...)
				if err != nil {
					t.Fatalf("rows.Scan: %v", err)
				}
				// check data just read in is ok
				// check columns data is ok
				for ii, vv := range data {
					rv := reflect.ValueOf(vv).Elem().Interface()
					if !reflect.DeepEqual(rv, hdu.Col(ii).Value) {
						t.Fatalf("rows.Scan:\nexpected=%v\ngot=%v", hdu.Col(ii).Value, rv)
					}
					if !reflect.DeepEqual(rv, ref[ii]) {
						t.Fatalf("rows.Scan:\nexpected=%v\ngot=%v", ref[ii], rv)
					}
				}
				// modify value of first column
				switch vv := hdu.Col(0).Value.(type) {
				case float64:
					hdu.Col(0).Value = 1 + vv
				case int16:
					hdu.Col(0).Value = 1 + vv
				}
				// check data just read in is ok
				// check columns data is ok
				for ii, vv := range data {
					if ii == 0 {
						continue
					}
					rv := reflect.ValueOf(vv).Elem().Interface()
					if !reflect.DeepEqual(rv, hdu.Col(ii).Value) {
						t.Fatalf("rows.Scan:\nexpected=%v\ngot=%v", hdu.Col(ii).Value, rv)
					}
					if !reflect.DeepEqual(rv, ref[ii]) {
						t.Fatalf("rows.Scan:\nexpected=%v\ngot=%v", ref[ii], rv)
					}
				}
				// but column data has changed
				if reflect.DeepEqual(reflect.ValueOf(data[0]).Elem().Interface(), hdu.Col(0).Value) {
					t.Fatalf("expected different values!")
				}
				count++
			}
			err = rows.Err()
			if err != nil {
				t.Fatalf("rows.Err: %v", err)
			}
			if count != nrows {
				t.Fatalf("rows.Next: expected [%d] rows. got %d.", nrows, count)
			}
		}
	}
}

func TestTableScanMap(t *testing.T) {
	for _, table := range g_tables {
		fname := table.fname
		f, err := Open(fname, ReadOnly)
		if err != nil {
			t.Fatalf("error opening file [%v]: %v", fname, err)
		}

		for i := range f.HDUs() {
			hdu, ok := f.HDU(i).(*Table)
			if !ok {
				continue
			}
			refmap := table.maps[i]
			if len(refmap) <= 0 {
				continue
			}
			nrows := hdu.NumRows()
			rows, err := hdu.Read(0, nrows)
			if err != nil {
				t.Fatalf("table.Read: %v", err)
			}
			count := int64(0)
			for rows.Next() {
				ref := make(map[string]interface{}, len(refmap))
				data := make(map[string]interface{}, len(refmap))
				for kk, vv := range table.maps[i] {
					rt := reflect.TypeOf(vv)
					rv := reflect.New(rt)
					xx := rv.Interface()
					data[kk] = xx
					ii := hdu.Index(kk)
					if ii < 0 {
						t.Fatalf("could not find index of [%v]", kk)
					}
					ref[kk] = table.tuple[i][count][ii]
				}
				err = rows.Scan(&data)
				if err != nil {
					t.Fatalf("rows.Scan: %v", err)
				}
				// check data just read in is ok
				if !reflect.DeepEqual(data, ref) {
					t.Fatalf("rows.Scan:\nexpected=%v\ngot=%v", ref, data)
				}
				// check columns data is ok
				for kk, vv := range data {
					ii := hdu.Index(kk)
					if !reflect.DeepEqual(vv, hdu.Col(ii).Value) {
						t.Fatalf("rows.Scan:\nexpected=%v\ngot=%v", hdu.Col(ii).Value, vv)
					}
				}
				// modify value of first column
				kk := hdu.Col(0).Name
				switch vv := hdu.Col(0).Value.(type) {
				case float64:
					hdu.Col(0).Value = 1 + vv
				case int16:
					hdu.Col(0).Value = 1 + vv
				}
				// check data is still ok
				if !reflect.DeepEqual(data, ref) {
					t.Fatalf("rows.Scan:\nexpected=%v\ngot=%v", ref, data)
				}
				// but column data has changed
				if reflect.DeepEqual(data[kk], hdu.Col(0).Value) {
					t.Fatalf("expected different values!")
				}
				count++
			}
			err = rows.Err()
			if err != nil {
				t.Fatalf("rows.Err: %v", err)
			}
			if count != nrows {
				t.Fatalf("rows.Next: expected [%d] rows. got %d.", nrows, count)
			}
		}
	}
}

func TestTableScanStruct(t *testing.T) {
	for _, table := range g_tables {
		fname := table.fname
		f, err := Open(fname, ReadOnly)
		if err != nil {
			t.Fatalf("error opening file [%v]: %v", fname, err)
		}

		for i := range f.HDUs() {
			hdu, ok := f.HDU(i).(*Table)
			if !ok {
				continue
			}
			reftypes := table.types[i]
			if reftypes == nil {
				continue
			}
			reftype := reflect.TypeOf(reftypes)
			nrows := hdu.NumRows()
			rows, err := hdu.Read(0, nrows)
			if err != nil {
				t.Fatalf("table.Read: %v", err)
			}
			count := int64(0)
			for rows.Next() {
				ref := reflect.New(reftype)
				data := reflect.New(reftype)
				for ii := 0; ii < reftype.NumField(); ii++ {
					ft := reftype.Field(ii)
					kk := ft.Tag.Get("fits")
					if kk == "" {
						kk = ft.Name
					}
					idx := hdu.Index(kk)
					if idx < 0 {
						t.Fatalf("could not find index of [%v.%v]", reftype.Name, ft.Name)
					}
					vv := reflect.ValueOf(table.tuple[i][count][idx])
					reflect.Indirect(ref).Field(ii).Set(vv)
				}

				err = rows.Scan(data.Interface())
				if err != nil {
					t.Fatalf("rows.Scan: %v", err)
				}
				// check data just read in is ok
				if !reflect.DeepEqual(data.Interface(), ref.Interface()) {
					t.Fatalf("rows.Scan:\nexpected=%v\ngot=%v", ref, data)
				}
				// check columns data is ok
				for ii := 0; ii < reftype.NumField(); ii++ {
					ft := reftype.Field(ii)
					kk := ft.Tag.Get("fits")
					if kk == "" {
						kk = ft.Name
					}
					idx := hdu.Index(kk)
					if idx < 0 {
						t.Fatalf("could not find index of [%v.%v]", reftype.Name, ft.Name)
					}
					vv := data.Elem().Field(ii).Interface()
					if !reflect.DeepEqual(vv, hdu.Col(idx).Value) {
						t.Fatalf("rows.Scan:\nexpected=%v\ngot=%v", hdu.Col(idx).Value, vv)
					}
				}
				// modify value of first column
				switch vv := hdu.Col(0).Value.(type) {
				case float64:
					hdu.Col(0).Value = 1 + vv
				case int16:
					hdu.Col(0).Value = 1 + vv
				}
				// check data is still ok
				if !reflect.DeepEqual(data.Interface(), ref.Interface()) {
					t.Fatalf("rows.Scan:\nexpected=%v\ngot=%v", ref, data)
				}
				// but column data has changed
				if reflect.DeepEqual(data.Elem().Field(0).Interface(), hdu.Col(0).Value) {
					t.Fatalf("expected different values!")
				}
				count++
			}
			err = rows.Err()
			if err != nil {
				t.Fatalf("rows.Err: %v", err)
			}
			if count != nrows {
				t.Fatalf("rows.Next: expected [%d] rows. got %d.", nrows, count)
			}
		}
	}
}

func TestTableRW(t *testing.T) {

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
		name  string
		cols  []Column
		htype HDUType
		table interface{}
	}{
		{
			name: "new.fits",
			cols: []Column{
				{
					Name:   "int8s",
					Format: "i8",
					Value:  int8(42),
				},
			},
			htype: ASCII_TBL,
			table: []int8{
				10, 11, 12, 13,
				14, 15, 16, 17,
				18, 19, 10, 11,
			},
		},
		{
			name: "new.fits",
			cols: []Column{
				{
					Name:   "int16s",
					Format: "i16",
					Value:  int16(42),
				},
			},
			htype: ASCII_TBL,
			table: []int16{
				10, 11, 12, 13,
				14, 15, 16, 17,
				18, 19, 10, 11,
			},
		},
		{
			name: "new.fits",
			cols: []Column{
				{
					Name:   "int32s",
					Format: "i32",
					Value:  int32(42),
				},
			},
			htype: ASCII_TBL,
			table: []int32{
				10, 11, 12, 13,
				14, 15, 16, 17,
				18, 19, 10, 11,
			},
		},
		{
			name: "new.fits",
			cols: []Column{
				{
					Name:   "int64s",
					Format: "i64",
					Value:  int64(42),
				},
			},
			htype: ASCII_TBL,
			table: []int64{
				10, 11, 12, 13,
				14, 15, 16, 17,
				18, 19, 10, 11,
			},
		},
		{
			name: "new.fits",
			cols: []Column{
				{
					Name:   "ints",
					Format: "i64",
					Value:  int(42),
				},
			},
			htype: ASCII_TBL,
			table: []int{
				10, 11, 12, 13,
				14, 15, 16, 17,
				18, 19, 10, 11,
			},
		},
		{
			name: "new.fits",
			cols: []Column{
				{
					Name:   "uint8s",
					Format: "i8",
					Value:  uint8(42),
				},
			},
			htype: ASCII_TBL,
			table: []uint8{
				10, 11, 12, 13,
				14, 15, 16, 17,
				18, 19, 10, 11,
			},
		},
		{
			name: "new.fits",
			cols: []Column{
				{
					Name:   "uint16s",
					Format: "i16",
					Value:  uint16(42),
				},
			},
			htype: ASCII_TBL,
			table: []uint16{
				10, 11, 12, 13,
				14, 15, 16, 17,
				18, 19, 10, 11,
			},
		},
		{
			name: "new.fits",
			cols: []Column{
				{
					Name:   "uint32s",
					Format: "i32",
					Value:  uint32(42),
				},
			},
			htype: ASCII_TBL,
			table: []uint32{
				10, 11, 12, 13,
				14, 15, 16, 17,
				18, 19, 10, 11,
			},
		},
		{
			name: "new.fits",
			cols: []Column{
				{
					Name:   "uint64s",
					Format: "i64",
					Value:  uint64(42),
				},
			},
			htype: ASCII_TBL,
			table: []uint64{
				10, 11, 12, 13,
				14, 15, 16, 17,
				18, 19, 10, 11,
			},
		},
		{
			name: "new.fits",
			cols: []Column{
				{
					Name:   "uints",
					Format: "i64",
					Value:  uint(42),
				},
			},
			htype: ASCII_TBL,
			table: []uint{
				10, 11, 12, 13,
				14, 15, 16, 17,
				18, 19, 10, 11,
			},
		},
		{
			name: "new.fits",
			cols: []Column{
				{
					Name:   "float32s",
					Format: "f32",
					Value:  float32(42),
				},
			},
			htype: ASCII_TBL,
			table: []float32{
				10, 11, 12, 13,
				14, 15, 16, 17,
				18, 19, 10, 11,
			},
		},
		{
			name: "new.fits",
			cols: []Column{
				{
					Name:   "float64s",
					Format: "f64",
					Value:  float64(42),
				},
			},
			htype: ASCII_TBL,
			table: []float64{
				10, 11, 12, 13,
				14, 15, 16, 17,
				18, 19, 10, 11,
			},
		},
		// binary table
		{
			name: "new.fits",
			cols: []Column{
				{
					Name:   "int8s",
					Format: "i8",
					Value:  int8(42),
				},
			},
			htype: BINARY_TBL,
			table: []int8{
				10, 11, 12, 13,
				14, 15, 16, 17,
				18, 19, 10, 11,
			},
		},
		{
			name: "new.fits",
			cols: []Column{
				{
					Name:   "int16s",
					Format: "i16",
					Value:  int16(42),
				},
			},
			htype: BINARY_TBL,
			table: []int16{
				10, 11, 12, 13,
				14, 15, 16, 17,
				18, 19, 10, 11,
			},
		},
		{
			name: "new.fits",
			cols: []Column{
				{
					Name:   "int32s",
					Format: "i32",
					Value:  int32(42),
				},
			},
			htype: BINARY_TBL,
			table: []int32{
				10, 11, 12, 13,
				14, 15, 16, 17,
				18, 19, 10, 11,
			},
		},
		{
			name: "new.fits",
			cols: []Column{
				{
					Name:   "int64s",
					Format: "i64",
					Value:  int64(42),
				},
			},
			htype: BINARY_TBL,
			table: []int64{
				10, 11, 12, 13,
				14, 15, 16, 17,
				18, 19, 10, 11,
			},
		},
		{
			name: "new.fits",
			cols: []Column{
				{
					Name:   "ints",
					Format: "i64",
					Value:  int(42),
				},
			},
			htype: BINARY_TBL,
			table: []int{
				10, 11, 12, 13,
				14, 15, 16, 17,
				18, 19, 10, 11,
			},
		},
		{
			name: "new.fits",
			cols: []Column{
				{
					Name:   "uint8s",
					Format: "b",
					Value:  uint8(42),
				},
			},
			htype: BINARY_TBL,
			table: []uint8{
				10, 11, 12, 13,
				14, 15, 16, 17,
				18, 19, 10, 11,
			},
		},
		{
			name: "new.fits",
			cols: []Column{
				{
					Name:   "uint16s",
					Format: "i16",
					Value:  uint16(42),
				},
			},
			htype: BINARY_TBL,
			table: []uint16{
				10, 11, 12, 13,
				14, 15, 16, 17,
				18, 19, 10, 11,
			},
		},
		{
			name: "new.fits",
			cols: []Column{
				{
					Name:   "uint32s",
					Format: "i32",
					Value:  uint32(42),
				},
			},
			htype: BINARY_TBL,
			table: []uint32{
				10, 11, 12, 13,
				14, 15, 16, 17,
				18, 19, 10, 11,
			},
		},
		{
			name: "new.fits",
			cols: []Column{
				{
					Name:   "uint64s",
					Format: "i64",
					Value:  uint64(42),
				},
			},
			htype: BINARY_TBL,
			table: []uint64{
				10, 11, 12, 13,
				14, 15, 16, 17,
				18, 19, 10, 11,
			},
		},
		{
			name: "new.fits",
			cols: []Column{
				{
					Name:   "uints",
					Format: "i64",
					Value:  uint(42),
				},
			},
			htype: BINARY_TBL,
			table: []uint{
				10, 11, 12, 13,
				14, 15, 16, 17,
				18, 19, 10, 11,
			},
		},
		{
			name: "new.fits",
			cols: []Column{
				{
					Name:   "float32s",
					Format: "e",
					Value:  float32(42),
				},
			},
			htype: BINARY_TBL,
			table: []float32{
				10, 11, 12, 13,
				14, 15, 16, 17,
				18, 19, 10, 11,
			},
		},
		{
			name: "new.fits",
			cols: []Column{
				{
					Name:   "float64s",
					Format: "d",
					Value:  float64(42),
				},
			},
			htype: BINARY_TBL,
			table: []float64{
				10, 11, 12, 13,
				14, 15, 16, 17,
				18, 19, 10, 11,
			},
		},
		{
			name: "new.fits",
			cols: []Column{
				{
					Name:   "cplx64s",
					Format: "64M",
					Value:  complex(float32(42), float32(42)),
				},
			},
			htype: BINARY_TBL,
			table: []complex64{
				complex(float32(10), float32(10)), complex(float32(11), float32(11)),
				complex(float32(12), float32(12)), complex(float32(13), float32(13)),
				complex(float32(14), float32(14)), complex(float32(15), float32(15)),
				complex(float32(16), float32(16)), complex(float32(17), float32(17)),
				complex(float32(18), float32(18)), complex(float32(19), float32(19)),
				complex(float32(10), float32(10)), complex(float32(11), float32(11)),
			},
		},
		{
			name: "new.fits",
			cols: []Column{
				{
					Name:   "cplx128s",
					Format: "128M",
					Value:  complex(float64(42), float64(42)),
				},
			},
			htype: BINARY_TBL,
			table: []complex128{
				complex(10, 10), complex(11, 11), complex(12, 12), complex(13, 13),
				complex(14, 14), complex(15, 15), complex(16, 16), complex(17, 17),
				complex(18, 18), complex(19, 19), complex(10, 10), complex(11, 11),
			},
		},
	} {
		fname := fmt.Sprintf("%03d_%s", ii, table.name)
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

				tbl, err := NewTable(&f, "test", table.cols, table.htype)
				if err != nil {
					t.Fatalf("error creating new table: %v", err)
				}
				defer tbl.Close()

				rslice := reflect.ValueOf(table.table)
				for i := 0; i < rslice.Len(); i++ {
					data := rslice.Index(i).Addr()
					err = tbl.Write(data.Interface())
					if err != nil {
						t.Fatalf("error writing row [%v]: %v", i, err)
					}
				}

				nrows := tbl.NumRows()
				if nrows != int64(rslice.Len()) {
					t.Fatalf("expected num rows [%v]. got [%v] (%v)", rslice.Len(), nrows, table.cols[0].Name)
				}
			},
			// read
			func() {
				f, err := Open(fname, ReadOnly)
				if err != nil {
					t.Fatalf("error opening file [%v]: %v", fname, err)
				}
				defer f.Close()

				hdu := f.HDU(1)
				tbl := hdu.(*Table)
				if tbl.Name() != "test" {
					t.Fatalf("expected table name==%q. got %q", "test", tbl.Name())
				}

				rslice := reflect.ValueOf(table.table)
				nrows := tbl.NumRows()
				if nrows != int64(rslice.Len()) {
					t.Fatalf("expected num rows [%v]. got [%v]", rslice.Len(), nrows)
				}

				rows, err := tbl.Read(0, nrows)
				if err != nil {
					t.Fatalf("table.Read: %v", err)
				}
				count := int64(0)
				for rows.Next() {
					ref := rslice.Index(int(count)).Interface()
					rt := reflect.TypeOf(ref)
					data := reflect.New(rt).Elem().Interface()
					err = rows.Scan(&data)
					if err != nil {
						t.Fatalf("rows.Scan: %v", err)
					}
					// check data just read in is ok
					if !reflect.DeepEqual(data, ref) {
						t.Fatalf("rows.Scan:\nexpected=%v\ngot=%v (%T)", ref, data, data)
					}
					count++
				}
				if count != nrows {
					t.Fatalf("expected [%v] rows. got [%v]", nrows, count)
				}
			},
		} {
			fct()
		}
	}
}

// EOF
