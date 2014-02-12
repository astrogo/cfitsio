package cfitsio

import "testing"

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

// EOF
