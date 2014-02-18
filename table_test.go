package cfitsio

import (
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
					rrv := reflect.Indirect(rv)
					xx := rrv.Interface()
					data[ii] = xx
					ref[ii] = vv
				}
				err = rows.Scan(data...)
				if err != nil {
					t.Fatalf("rows.Scan: %v", err)
				}
				// check data just read in is ok
				if !reflect.DeepEqual(data, ref) {
					t.Fatalf("rows.Scan:\nexpected=%v\ngot=%v", ref, data)
				}
				// check columns data is ok
				for ii, vv := range data {
					if !reflect.DeepEqual(vv, hdu.Col(ii).Value) {
						t.Fatalf("rows.Scan:\nexpected=%v\ngot=%v", hdu.Col(ii).Value, vv)
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
				if !reflect.DeepEqual(data, ref) {
					t.Fatalf("rows.Scan:\nexpected=%v\ngot=%v", ref, data)
				}
				// but column data has changed
				if reflect.DeepEqual(data[0], hdu.Col(0).Value) {
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

// EOF
