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

// EOF
