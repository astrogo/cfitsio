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
						t.Fatalf("could not find index of [%v.%v]", reftype.Name(), ft.Name)
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
						t.Fatalf("could not find index of [%v.%v]", reftype.Name(), ft.Name)
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

func TestTableBuiltinsRW(t *testing.T) {

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
					Name:  "int8s",
					Value: int8(42),
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
					Name:  "int16s",
					Value: int16(42),
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
					Name:  "int32s",
					Value: int32(42),
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
					Name:  "int64s",
					Value: int64(42),
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
					Name:  "ints",
					Value: int(42),
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
					Name:  "uint8s",
					Value: uint8(42),
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
					Name:  "uint16s",
					Value: uint16(42),
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
					Name:  "uint32s",
					Value: uint32(42),
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
					Name:  "uint64s",
					Value: uint64(42),
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
					Name:  "uints",
					Value: uint(42),
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
					Name:  "float32s",
					Value: float32(42),
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
					Name:  "float64s",
					Value: float64(42),
				},
			},
			htype: ASCII_TBL,
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
					Name:  "strings",
					Value: "",
				},
			},
			htype: ASCII_TBL,
			table: []string{
				"10", "11", "12", "13",
				"14", "15", "16", "17",
				"18", "19", "10", "11",
			},
		},
		// binary table
		{
			name: "new.fits",
			cols: []Column{
				{
					Name:  "bools",
					Value: false,
				},
			},
			htype: BINARY_TBL,
			table: []bool{
				true, true, true, true,
				false, false, false, false,
				true, false, true, false,
				false, true, false, true,
			},
		},
		{
			name: "new.fits",
			cols: []Column{
				{
					Name:  "int8s",
					Value: int8(42),
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
					Name:  "int16s",
					Value: int16(42),
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
					Name:  "int32s",
					Value: int32(42),
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
					Name:  "int64s",
					Value: int64(42),
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
					Name:  "ints",
					Value: int(42),
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
					Name:  "uint8s",
					Value: uint8(42),
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
					Name:  "uint16s",
					Value: uint16(42),
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
					Name:  "uint32s",
					Value: uint32(42),
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
					Name:  "uint64s",
					Value: uint64(42),
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
					Name:  "uints",
					Value: uint(42),
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
					Name:  "float32s",
					Value: float32(42),
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
					Name:  "float64s",
					Value: float64(42),
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
					Name:  "cplx64s",
					Value: complex(float32(42), float32(42)),
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
					Name:  "cplx128s",
					Value: complex(float64(42), float64(42)),
				},
			},
			htype: BINARY_TBL,
			table: []complex128{
				complex(10, 10), complex(11, 11), complex(12, 12), complex(13, 13),
				complex(14, 14), complex(15, 15), complex(16, 16), complex(17, 17),
				complex(18, 18), complex(19, 19), complex(10, 10), complex(11, 11),
			},
		},
		{
			name: "new.fits",
			cols: []Column{
				{
					Name:  "strings",
					Value: "",
				},
			},
			htype: BINARY_TBL,
			table: []string{
				"10", "11", "12", "13",
				"14", "15", "16", "17",
				"18", "19", "10", "11",
			},
		},
		{
			name: "new.fits",
			cols: []Column{
				{
					Name:  "float64s",
					Value: [2]float64{},
				},
			},
			htype: BINARY_TBL,
			table: [][2]float64{
				{10, 11},
				{12, 13},
				{14, 15},
				{16, 17},
				{18, 19},
				{10, 11},
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
					t.Fatalf("error creating new table: %v (%v)", err, table.cols[0].Name)
				}
				defer tbl.Close()

				rslice := reflect.ValueOf(table.table)
				for i := 0; i < rslice.Len(); i++ {
					data := rslice.Index(i).Addr()
					err = tbl.Write(data.Interface())
					if err != nil {
						t.Fatalf("error writing row [%v]: %v (data=%v %T)", i, err, data.Interface(), data.Interface())
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
						t.Fatalf("rows.Scan: %[3]s\nexp=%[1]v (%[1]T)\ngot=%[2]v (%[2]T)", ref, data, table.cols[0].Name)
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

func TestTableSliceRW(t *testing.T) {

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
		// binary table
		{
			name: "new.fits",
			cols: []Column{
				{
					Name:  "bools",
					Value: []bool{},
				},
			},
			htype: BINARY_TBL,
			table: [][]bool{
				{true, true, true, true},
				{false, false, false, false},
				{true, false, true, false},
				{false, true, false, true},
			},
		},
		{
			name: "new.fits",
			cols: []Column{
				{
					Name:  "int8s",
					Value: []int8{},
				},
			},
			htype: BINARY_TBL,
			table: [][]int8{
				{10, 11, 12, 13},
				{14, 15, 16, 17},
				{18, 19, 10, 11},
			},
		},
		{
			name: "new.fits",
			cols: []Column{
				{
					Name:  "int16s",
					Value: []int16{},
				},
			},
			htype: BINARY_TBL,
			table: [][]int16{
				{10, 11, 12, 13},
				{14, 15, 16, 17},
				{18, 19, 10, 11},
			},
		},
		{
			name: "new.fits",
			cols: []Column{
				{
					Name:  "int32s",
					Value: []int32{},
				},
			},
			htype: BINARY_TBL,
			table: [][]int32{
				{10, 11, 12, 13},
				{14, 15, 16, 17},
				{18, 19, 10, 11},
			},
		},
		{
			name: "new.fits",
			cols: []Column{
				{
					Name:  "int64s",
					Value: []int64{},
				},
			},
			htype: BINARY_TBL,
			table: [][]int64{
				{10, 11, 12, 13},
				{14, 15, 16, 17},
				{18, 19, 10, 11},
			},
		},
		{
			name: "new.fits",
			cols: []Column{
				{
					Name:  "ints",
					Value: []int{},
				},
			},
			htype: BINARY_TBL,
			table: [][]int{
				{10, 11, 12, 13},
				{14, 15, 16, 17},
				{18, 19, 10, 11},
			},
		},
		{
			name: "new.fits",
			cols: []Column{
				{
					Name:  "uint8s",
					Value: []uint8{},
				},
			},
			htype: BINARY_TBL,
			table: [][]uint8{
				{10, 11, 12, 13},
				{14, 15, 16, 17},
				{18, 19, 10, 11},
			},
		},
		{
			name: "new.fits",
			cols: []Column{
				{
					Name:  "uint16s",
					Value: []uint16{},
				},
			},
			htype: BINARY_TBL,
			table: [][]uint16{
				{10, 11, 12, 13},
				{14, 15, 16, 17},
				{18, 19, 10, 11},
			},
		},
		{
			name: "new.fits",
			cols: []Column{
				{
					Name:  "uint32s",
					Value: []uint32{},
				},
			},
			htype: BINARY_TBL,
			table: [][]uint32{
				{10, 11, 12, 13},
				{14, 15, 16, 17},
				{18, 19, 10, 11},
			},
		},
		{
			name: "new.fits",
			cols: []Column{
				{
					Name:  "uint64s",
					Value: []uint64{},
				},
			},
			htype: BINARY_TBL,
			table: [][]uint64{
				{10, 11, 12, 13},
				{14, 15, 16, 17},
				{18, 19, 10, 11},
			},
		},
		{
			name: "new.fits",
			cols: []Column{
				{
					Name:  "uints",
					Value: []uint{},
				},
			},
			htype: BINARY_TBL,
			table: [][]uint{
				{10, 11, 12, 13},
				{14, 15, 16, 17},
				{18, 19, 10, 11},
			},
		},
		{
			name: "new.fits",
			cols: []Column{
				{
					Name:  "float32s",
					Value: []float32{},
				},
			},
			htype: BINARY_TBL,
			table: [][]float32{
				{10, 11, 12, 13},
				{14, 15, 16, 17},
				{18, 19, 10, 11},
			},
		},
		{
			name: "new.fits",
			cols: []Column{
				{
					Name:  "float64s",
					Value: []float64{},
				},
			},
			htype: BINARY_TBL,
			table: [][]float64{
				{10, 11, 12, 13},
				{14, 15, 16, 17},
				{18, 19, 10, 11},
			},
		},
		{
			name: "new.fits",
			cols: []Column{
				{
					Name:  "cplx64s",
					Value: []complex64{},
				},
			},
			htype: BINARY_TBL,
			table: [][]complex64{
				{complex(float32(10), float32(10)), complex(float32(11), float32(11)),
					complex(float32(12), float32(12)), complex(float32(13), float32(13))},
				{complex(float32(14), float32(14)), complex(float32(15), float32(15)),
					complex(float32(16), float32(16)), complex(float32(17), float32(17))},
				{complex(float32(18), float32(18)), complex(float32(19), float32(19)),
					complex(float32(10), float32(10)), complex(float32(11), float32(11))},
			},
		},
		{
			name: "new.fits",
			cols: []Column{
				{
					Name:  "cplx128s",
					Value: []complex128{},
				},
			},
			htype: BINARY_TBL,
			table: [][]complex128{
				{complex(10, 10), complex(11, 11), complex(12, 12), complex(13, 13)},
				{complex(14, 14), complex(15, 15), complex(16, 16), complex(17, 17)},
				{complex(18, 18), complex(19, 19), complex(10, 10), complex(11, 11)},
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
					t.Fatalf("error creating new table: %v (%v)", err, table.cols[0].Name)
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
						t.Fatalf("rows.Scan: [%[3]s|%[4]v]\nexp=%[1]v (%[1]T)\ngot=%[2]v (%[2]T)", ref, data, table.cols[0].Name, table.htype)
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

func TestTableArrayRW(t *testing.T) {

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
		// binary table
		{
			name: "new.fits",
			cols: []Column{
				{
					Name:  "bools",
					Value: [4]bool{},
				},
			},
			htype: BINARY_TBL,
			table: [][4]bool{
				{true, true, true, true},
				{false, false, false, false},
				{true, false, true, false},
				{false, true, false, true},
			},
		},
		{
			name: "new.fits",
			cols: []Column{
				{
					Name:  "int8s",
					Value: [4]int8{},
				},
			},
			htype: BINARY_TBL,
			table: [][4]int8{
				{10, 11, 12, 13},
				{14, 15, 16, 17},
				{18, 19, 10, 11},
			},
		},
		{
			name: "new.fits",
			cols: []Column{
				{
					Name:  "int16s",
					Value: [4]int16{},
				},
			},
			htype: BINARY_TBL,
			table: [][4]int16{
				{10, 11, 12, 13},
				{14, 15, 16, 17},
				{18, 19, 10, 11},
			},
		},
		{
			name: "new.fits",
			cols: []Column{
				{
					Name:  "int32s",
					Value: [4]int32{},
				},
			},
			htype: BINARY_TBL,
			table: [][4]int32{
				{10, 11, 12, 13},
				{14, 15, 16, 17},
				{18, 19, 10, 11},
			},
		},
		{
			name: "new.fits",
			cols: []Column{
				{
					Name:  "int64s",
					Value: [4]int64{},
				},
			},
			htype: BINARY_TBL,
			table: [][4]int64{
				{10, 11, 12, 13},
				{14, 15, 16, 17},
				{18, 19, 10, 11},
			},
		},
		{
			name: "new.fits",
			cols: []Column{
				{
					Name:  "ints",
					Value: [4]int{},
				},
			},
			htype: BINARY_TBL,
			table: [][4]int{
				{10, 11, 12, 13},
				{14, 15, 16, 17},
				{18, 19, 10, 11},
			},
		},
		{
			name: "new.fits",
			cols: []Column{
				{
					Name:  "uint8s",
					Value: [4]uint8{},
				},
			},
			htype: BINARY_TBL,
			table: [][4]uint8{
				{10, 11, 12, 13},
				{14, 15, 16, 17},
				{18, 19, 10, 11},
			},
		},
		{
			name: "new.fits",
			cols: []Column{
				{
					Name:  "uint16s",
					Value: [4]uint16{},
				},
			},
			htype: BINARY_TBL,
			table: [][4]uint16{
				{10, 11, 12, 13},
				{14, 15, 16, 17},
				{18, 19, 10, 11},
			},
		},
		{
			name: "new.fits",
			cols: []Column{
				{
					Name:  "uint32s",
					Value: [4]uint32{},
				},
			},
			htype: BINARY_TBL,
			table: [][4]uint32{
				{10, 11, 12, 13},
				{14, 15, 16, 17},
				{18, 19, 10, 11},
			},
		},
		{
			name: "new.fits",
			cols: []Column{
				{
					Name:  "uint64s",
					Value: [4]uint64{},
				},
			},
			htype: BINARY_TBL,
			table: [][4]uint64{
				{10, 11, 12, 13},
				{14, 15, 16, 17},
				{18, 19, 10, 11},
			},
		},
		{
			name: "new.fits",
			cols: []Column{
				{
					Name:  "uints",
					Value: [4]uint{},
				},
			},
			htype: BINARY_TBL,
			table: [][4]uint{
				{10, 11, 12, 13},
				{14, 15, 16, 17},
				{18, 19, 10, 11},
			},
		},
		{
			name: "new.fits",
			cols: []Column{
				{
					Name:  "float32s",
					Value: [4]float32{},
				},
			},
			htype: BINARY_TBL,
			table: [][4]float32{
				{10, 11, 12, 13},
				{14, 15, 16, 17},
				{18, 19, 10, 11},
			},
		},
		{
			name: "new.fits",
			cols: []Column{
				{
					Name:  "float64s",
					Value: [4]float64{},
				},
			},
			htype: BINARY_TBL,
			table: [][4]float64{
				{10, 11, 12, 13},
				{14, 15, 16, 17},
				{18, 19, 10, 11},
			},
		},
		{
			name: "new.fits",
			cols: []Column{
				{
					Name:  "cplx64s",
					Value: [4]complex64{},
				},
			},
			htype: BINARY_TBL,
			table: [][4]complex64{
				{complex(float32(10), float32(10)), complex(float32(11), float32(11)),
					complex(float32(12), float32(12)), complex(float32(13), float32(13))},
				{complex(float32(14), float32(14)), complex(float32(15), float32(15)),
					complex(float32(16), float32(16)), complex(float32(17), float32(17))},
				{complex(float32(18), float32(18)), complex(float32(19), float32(19)),
					complex(float32(10), float32(10)), complex(float32(11), float32(11))},
			},
		},
		{
			name: "new.fits",
			cols: []Column{
				{
					Name:  "cplx128s",
					Value: [4]complex128{},
				},
			},
			htype: BINARY_TBL,
			table: [][4]complex128{
				{complex(10, 10), complex(11, 11), complex(12, 12), complex(13, 13)},
				{complex(14, 14), complex(15, 15), complex(16, 16), complex(17, 17)},
				{complex(18, 18), complex(19, 19), complex(10, 10), complex(11, 11)},
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
					t.Fatalf("error creating new table: %v (%v)", err, table.cols[0].Name)
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
						t.Fatalf("rows.Scan: [%[3]s|%[4]v]\nexp=%[1]v (%[1]T)\ngot=%[2]v (%[2]T)", ref, data, table.cols[0].Name, table.htype)
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

func TestTableStructsRW(t *testing.T) {

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

	type DataStruct struct {
		A   int64      `fits:"int64"`
		XX0 int        // hole
		B   float64    `fits:"float64"`
		XX1 int        // hole
		C   []int64    `fits:"int64s"`
		XX2 int        // hole
		D   []float64  `fits:"float64s"`
		XX3 int        // hole
		E   [2]float64 `fits:"arrfloat64s"`
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
					Name:  "int64",
					Value: int64(0),
				},
				{
					Name:  "float64",
					Value: float64(0),
				},
				{
					Name:  "int64s",
					Value: []int64{},
				},
				{
					Name:  "float64s",
					Value: []float64{},
				},
				{
					Name:  "arrfloat64s",
					Value: [2]float64{},
				},
			},
			htype: BINARY_TBL,
			table: []DataStruct{
				{A: 10, B: 10, C: []int64{10, 10}, D: []float64{10, 10}, E: [2]float64{10, 10}},
				{A: 11, B: 11, C: []int64{11, 11}, D: []float64{11, 11}, E: [2]float64{11, 11}},
				{A: 12, B: 12, C: []int64{12, 12}, D: []float64{12, 12}, E: [2]float64{12, 12}},
				{A: 13, B: 13, C: []int64{13, 13}, D: []float64{13, 13}, E: [2]float64{13, 13}},
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
					t.Fatalf("error creating new table: %v (%v)", err, table.cols[0].Name)
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
					rv := reflect.New(rt).Elem()
					data := rv.Interface().(DataStruct)
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

func TestTableMapsRW(t *testing.T) {

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
					Name:  "A",
					Value: int64(0),
				},
				{
					Name:  "B",
					Value: float64(0),
				},
				{
					Name:  "C",
					Value: []int64{},
				},
				{
					Name:  "D",
					Value: []float64{},
				},
				// FIXME: re-add when reflect.ArrayOf is available
				// {
				// 	Name:  "E",
				// 	Value: [2]float64{},
				// },
			},
			htype: BINARY_TBL,
			table: []map[string]interface{}{
				{
					"A": int64(10),
					"B": float64(10),
					"C": []int64{10, 10},
					"D": []float64{10, 10},
					// FIXME: re-add when reflect.ArrayOf is available
					// "E": [2]float64{10, 10},
				},
				{
					"A": int64(11),
					"B": float64(11),
					"C": []int64{11, 11},
					"D": []float64{11, 11},
					// FIXME: re-add when reflect.ArrayOf is available
					// "E": [2]float64{11, 11},
				},
				{
					"A": int64(12),
					"B": float64(12),
					"C": []int64{12, 12},
					"D": []float64{12, 12},
					// FIXME: re-add when reflect.ArrayOf is available
					// "E": [2]float64{12, 12},
				},
				{
					"A": int64(13),
					"B": float64(13),
					"C": []int64{13, 13},
					"D": []float64{13, 13},
					// FIXME: re-add when reflect.ArrayOf is available
					// "E": [2]float64{13, 13},
				},
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
					t.Fatalf("error creating new table: %v (%v)", err, table.cols[0].Name)
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
					ref := rslice.Index(int(count)).Interface().(map[string]interface{})
					data := map[string]interface{}{}
					err = rows.Scan(&data)
					if err != nil {
						t.Fatalf("rows.Scan: %v", err)
					}
					// check data just read in is ok
					if !reflect.DeepEqual(data, ref) {
						t.Fatalf("rows.Scan:\nexp=%[1]v (%[1]T)\ngot=%[2]v (%[2]T)", ref, data)
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
