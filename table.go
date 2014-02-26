package cfitsio

// #include "go-cfitsio.h"
// #include "go-cfitsio-utils.h"
import "C"
import (
	"fmt"
	"reflect"
	"strconv"
	"strings"
	"unsafe"
)

type TypeCode int

const (
	TBIT        TypeCode = C.TBIT /* codes for FITS table data types */
	TBYTE       TypeCode = C.TBYTE
	TSBYTE      TypeCode = C.TSBYTE
	TLOGICAL    TypeCode = C.TLOGICAL
	TSTRING     TypeCode = C.TSTRING
	TUSHORT     TypeCode = C.TUSHORT
	TSHORT      TypeCode = C.TSHORT
	TUINT       TypeCode = C.TUINT
	TINT        TypeCode = C.TINT
	TULONG      TypeCode = C.TULONG
	TLONG       TypeCode = C.TLONG
	TINT32BIT   TypeCode = C.TINT32BIT /* used when returning datatype of a column */
	TFLOAT      TypeCode = C.TFLOAT
	TLONGLONG   TypeCode = C.TLONGLONG
	TDOUBLE     TypeCode = C.TDOUBLE
	TCOMPLEX    TypeCode = C.TCOMPLEX
	TDBLCOMPLEX TypeCode = C.TDBLCOMPLEX

	// variable length arrays
	TVLABIT        TypeCode = -C.TBIT /* codes for FITS table data types */
	TVLABYTE       TypeCode = -C.TBYTE
	TVLASBYTE      TypeCode = -C.TSBYTE
	TVLALOGICAL    TypeCode = -C.TLOGICAL
	TVLASTRING     TypeCode = -C.TSTRING
	TVLAUSHORT     TypeCode = -C.TUSHORT
	TVLASHORT      TypeCode = -C.TSHORT
	TVLAUINT       TypeCode = -C.TUINT
	TVLAINT        TypeCode = -C.TINT
	TVLAULONG      TypeCode = -C.TULONG
	TVLALONG       TypeCode = -C.TLONG
	TVLAINT32BIT   TypeCode = -C.TINT32BIT /* used when returning datatype of a column */
	TVLAFLOAT      TypeCode = -C.TFLOAT
	TVLALONGLONG   TypeCode = -C.TLONGLONG
	TVLADOUBLE     TypeCode = -C.TDOUBLE
	TVLACOMPLEX    TypeCode = -C.TCOMPLEX
	TVLADBLCOMPLEX TypeCode = -C.TDBLCOMPLEX
)

func govalue_from_typecode(t TypeCode) Value {
	var v Value
	switch t {
	case TBIT, TBYTE:
		var vv byte
		v = vv

	case TSBYTE:
		var vv int8
		v = vv

	case TLOGICAL:
		var vv bool
		v = vv

	case TSTRING:
		var vv string
		v = vv

	case TUSHORT:
		var vv uint16
		v = vv

	case TSHORT:
		var vv int16
		v = vv

	case TUINT:
		var vv uint32
		v = vv

	case TINT:
		var vv int32
		v = vv

	case TULONG:
		var vv uint64
		v = vv

	case TLONG:
		var vv int64
		v = vv

	case TFLOAT:
		var vv float32
		v = vv

	case TLONGLONG:
		var vv int64
		v = vv

	case TDOUBLE:
		var vv float64
		v = vv

	case TCOMPLEX:
		var vv complex64
		v = vv

	case TDBLCOMPLEX:
		var vv complex128
		v = vv

	case TVLABIT, TVLABYTE:
		var vv = make([]byte, 0)
		v = vv

	case TVLASBYTE:
		var vv = make([]int8, 0)
		v = vv

	case TVLALOGICAL:
		var vv = make([]bool, 0)
		v = vv

	case TVLASTRING:
		var vv = make([]string, 0)
		v = vv

	case TVLAUSHORT:
		var vv = make([]uint16, 0)
		v = vv

	case TVLASHORT:
		var vv = make([]int16, 0)
		v = vv

	case TVLAUINT:
		var vv = make([]uint32, 0)
		v = vv

	case TVLAINT:
		var vv = make([]int32, 0)
		v = vv

	case TVLAULONG:
		var vv = make([]uint64, 0)
		v = vv

	case TVLALONG:
		var vv = make([]int64, 0)
		v = vv

	case TVLAFLOAT:
		var vv = make([]float32, 0)
		v = vv

	case TVLALONGLONG:
		var vv = make([]int64, 0)
		v = vv

	case TVLADOUBLE:
		var vv = make([]float64, 0)
		v = vv

	case TVLACOMPLEX:
		var vv = make([]complex64, 0)
		v = vv

	case TVLADBLCOMPLEX:
		var vv = make([]complex128, 0)
		v = vv

	default:
		panic(fmt.Errorf("cfitsio: invalid TypeCode value [%v]", int(t)))
	}
	return v
}

type Table struct {
	f      *File
	id     C.int
	header Header
	nrows  int64
	cols   []Column
	data   interface{}
}

func (hdu *Table) Close() error {
	hdu.f = nil
	return nil
}

func (hdu *Table) Header() Header {
	return hdu.header
}

func (hdu *Table) Type() HDUType {
	return hdu.header.htype
}

func (hdu *Table) Name() string {
	card := hdu.header.Get("EXTNAME")
	if card == nil {
		return ""
	}
	return card.Value.(string)
}

func (hdu *Table) Version() int {
	card := hdu.header.Get("EXTVER")
	if card == nil {
		return 1
	}
	return card.Value.(int)
}

func (hdu *Table) Data(interface{}) error {
	var err error
	if hdu.data == nil {
		err = hdu.load()
	}
	return err
}

func (hdu *Table) load() error {
	return nil
}

func (hdu *Table) NumRows() int64 {
	return hdu.nrows
}

func (hdu *Table) NumCols() int {
	return len(hdu.cols)
}

func (hdu *Table) Cols() []Column {
	return hdu.cols
}

func (hdu *Table) Col(i int) *Column {
	return &hdu.cols[i]
}

// Index returns the index of the first column with name `n` or -1
func (hdu *Table) Index(n string) int {
	for i := range hdu.cols {
		col := &hdu.cols[i]
		if col.Name == n {
			return i
		}
	}
	return -1
}

func (hdu *Table) readRow(irow int64) error {
	err := hdu.seekHDU()
	if err != nil {
		return err
	}
	for icol := range hdu.cols {
		err = hdu.cols[icol].read(hdu.f, icol, irow)
		if err != nil {
			return err
		}
	}
	return err
}

// ReadRange reads rows over the range [beg, end) and returns the corresponding iterator.
// if end > maxrows, the iteration will stop at maxrows
// ReadRange has the same semantics than a `for i=0; i < max; i+=inc {...}` loop
func (hdu *Table) ReadRange(beg, end, inc int64) (*Rows, error) {
	var rows *Rows
	err := hdu.seekHDU()
	if err != nil {
		return rows, err
	}

	maxrows := hdu.NumRows()
	if end > maxrows {
		end = maxrows
	}

	if beg < 0 {
		beg = 0
	}

	cols := make([]int, len(hdu.cols))
	for i := range hdu.cols {
		cols[i] = i
	}

	rows = &Rows{
		table: hdu,
		cols:  cols,
		i:     beg,
		n:     end,
		inc:   inc,
		cur:   beg - inc,
		err:   nil,
	}
	return rows, err
}

// Read reads rows over the range [beg, end) and returns the corresponding iterator.
// if end > maxrows, the iteration will stop at maxrows
// ReadRange has the same semantics than a `for i=0; i < max; i++ {...}` loop
func (hdu *Table) Read(beg, end int64) (*Rows, error) {
	return hdu.ReadRange(beg, end, 1)
}

func (hdu *Table) seekHDU() error {
	c_status := C.int(0)
	c_htype := C.int(0)
	C.fits_movabs_hdu(hdu.f.c, hdu.id, &c_htype, &c_status)
	if c_status > 0 {
		return to_err(c_status)
	}
	return nil
}

func newTable(f *File, hdr Header, i int) (hdu HDU, err error) {
	c_status := C.int(0)
	c_id := C.int(0)
	C.fits_get_hdu_num(f.c, &c_id)
	if c_status > 0 {
		return nil, to_err(c_status)
	}

	c_nrows := C.long(0)
	C.fits_get_num_rows(f.c, &c_nrows, &c_status)
	if c_status > 0 {
		return nil, to_err(c_status)
	}

	c_ncols := C.int(0)
	C.fits_get_num_cols(f.c, &c_ncols, &c_status)
	if c_status > 0 {
		return nil, to_err(c_status)
	}

	ncols := int(c_ncols)
	cols := make([]Column, ncols)

	get := func(str string, ii int) *Card {
		return hdr.Get(fmt.Sprintf(str+"%d", ii+1))
	}
	for ii := 0; ii < ncols; ii++ {
		col := &cols[ii]
		// column name
		{
			c_status := C.int(0)
			c_tmpl := C.CString(fmt.Sprintf("%d", ii+1))
			defer C.free(unsafe.Pointer(c_tmpl))
			c_name := C.CStringN(C.FLEN_CARD)
			defer C.free(unsafe.Pointer(c_name))
			c_colnum := C.int(0)
			C.fits_get_colname(f.c, C.CASESEN, c_tmpl, c_name, &c_colnum, &c_status)
			if c_status > 0 {
				return nil, to_err(c_status)
			}
			col.Name = C.GoString(c_name)
		}

		card := get("TFORM", ii)
		if card != nil {
			col.Format = card.Value.(string)
		}

		card = get("TUNIT", ii)
		if card != nil {
			col.Unit = card.Value.(string)
		}

		card = get("TNULL", ii)
		if card != nil {
			col.Null = card.Value.(string)
		}

		card = get("TSCAL", ii)
		if card != nil {
			col.Bscale = card.Value.(float64)
		} else {
			col.Bscale = 1.0
		}

		card = get("TZERO", ii)
		if card != nil {
			col.Bzero = card.Value.(float64)
		} else {
			col.Bzero = 0.0
		}

		card = get("TDISP", ii)
		if card != nil {
			col.Display = card.Value.(string)
		}

		{
			// int fits_read_tdimll / ffgtdmll
			//(fitsfile *fptr, int colnum, int maxdim, > int *naxis,
			//LONGLONG *naxes, int *status)

		}
		card = get("TDIM", ii)
		if card != nil {
			dims := card.Value.(string)
			dims = strings.Replace(dims, "(", "", -1)
			dims = strings.Replace(dims, ")", "", -1)
			toks := make([]string, 0)
			for _, tok := range strings.Split(dims, ",") {
				tok = strings.Trim(tok, " \t\n")
				if tok == "" {
					continue
				}
				toks = append(toks, tok)
			}
			col.Dim = make([]int64, 0, len(toks))
			for _, tok := range toks {
				dim, err := strconv.ParseInt(tok, 10, 64)
				if err != nil {
					return nil, err
				}
				col.Dim = append(col.Dim, dim)
			}
		}

		card = get("TBCOL", ii)
		if card != nil {
			col.Start = card.Value.(int64)
		}

		{
			c_type := C.int(0)
			c_repeat := C.long(0)
			c_width := C.long(0)
			c_status := C.int(0)
			c_col := C.int(ii + 1) // 1-based index
			C.fits_get_coltype(f.c, c_col, &c_type, &c_repeat, &c_width, &c_status)
			if c_status > 0 {
				return nil, to_err(c_status)
			}
			col.Value = govalue_from_typecode(TypeCode(c_type))
		}
	}

	hdu = &Table{
		f:      f,
		id:     c_id,
		header: hdr,
		nrows:  int64(c_nrows),
		cols:   cols,
		data:   nil,
	}
	return hdu, err
}

// NewTable creates a new table in the given FITS file
func NewTable(f *File, name string, cols []Column, hdutype HDUType) (*Table, error) {
	var err error
	var table *Table
	mode, err := f.Mode()
	if err != nil {
		return table, err
	}
	if mode == ReadOnly {
		return table, READONLY_FILE
	}

	nhdus := len(f.hdus)

	if len(cols) <= 0 {
		return table, fmt.Errorf("cfitsio.NewTable: invalid number of columns (%v)", len(cols))
	}

	c_status := C.int(0)
	c_sz := C.int(len(cols))
	c_types := C.char_array_new(c_sz)
	defer C.free(unsafe.Pointer(c_types))
	c_forms := C.char_array_new(c_sz)
	defer C.free(unsafe.Pointer(c_forms))
	c_units := C.char_array_new(c_sz)
	defer C.free(unsafe.Pointer(c_units))
	c_hduname := C.CString(name)
	defer C.free(unsafe.Pointer(c_hduname))

	for i := 0; i < len(cols); i++ {
		c_idx := C.int(i)
		col := cols[i]
		c_name := C.CString(col.Name)
		defer C.free(unsafe.Pointer(c_name))
		C.char_array_set(c_types, c_idx, c_name)

		c_form := C.CString(col.Format)
		defer C.free(unsafe.Pointer(c_form))
		C.char_array_set(c_forms, c_idx, c_form)

		c_unit := C.CString(col.Unit)
		defer C.free(unsafe.Pointer(c_unit))
		C.char_array_set(c_units, c_idx, c_unit)
	}

	C.fits_create_tbl(f.c, C.int(hdutype), 0, C.int(len(cols)), c_types, c_forms, c_units, c_hduname, &c_status)
	if c_status > 0 {
		return table, to_err(c_status)
	}

	hdu, err := f.readHDU(nhdus)
	if err != nil {
		return table, err
	}
	f.hdus = append(f.hdus, hdu)
	table = hdu.(*Table)

	return table, err
}

// Write writes a row to the table
func (hdu *Table) Write(args ...interface{}) error {

	switch len(args) {
	case 0:
		return fmt.Errorf("cfitsio: Table.Write needs at least one argument")
	case 1:
		// maybe special case: map? struct?
		rt := reflect.TypeOf(args[0]).Elem()
		switch rt.Kind() {
		case reflect.Map:
			return hdu.writeMap(*args[0].(*map[string]interface{}))
		case reflect.Struct:
			return hdu.writeStruct(args[0])
		}
	}

	return hdu.write(args...)
}

func (hdu *Table) writeMap(data map[string]interface{}) error {
	var err error
	return err
}

func (hdu *Table) writeStruct(data interface{}) error {
	var err error
	rt := reflect.TypeOf(data).Elem()
	rv := reflect.ValueOf(data).Elem()
	icols := make([][2]int, 0, rt.NumField())
	for i := 0; i < rt.NumField(); i++ {
		f := rt.Field(i)
		n := f.Tag.Get("fits")
		if n == "" {
			n = f.Name
		}
		icol := hdu.Index(n)
		if icol >= 0 {
			icols = append(icols, [2]int{i, icol})
		}
	}

	irow := hdu.NumRows()

	for _, icol := range icols {
		vv := reflect.ValueOf(hdu.cols[icol[1]].Value)
		vv.Field(icol[0]).Set(rv)
		err = hdu.cols[icol[1]].write(hdu.f, icol[1], irow)
		if err != nil {
			return err
		}
	}
	return err
}

func (hdu *Table) write(args ...interface{}) error {
	var err error
	return err
}

func init() {
	g_hdus[ASCII_TBL] = newTable
	g_hdus[BINARY_TBL] = newTable
}

// EOF
