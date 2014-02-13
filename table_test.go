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
				err := hdu.ReadRow(irow)
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

			// iter over all rows
			rows, err := hdu.Read(0, -1)
			if err != nil {
				t.Fatalf("table.Read: %v", err)
			}
			nrows := hdu.NumRows()
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
			rows, err := hdu.Read(0, -1)
			if err != nil {
				t.Fatalf("table.Read: %v", err)
			}
			nrows := hdu.NumRows()
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
			rows, err := hdu.Read(0, -1)
			if err != nil {
				t.Fatalf("table.Read: %v", err)
			}
			nrows := hdu.NumRows()
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

// EOF
