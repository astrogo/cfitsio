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

func (tc TypeCode) String() string {
	switch tc {
	case TBIT:
		return "TBIT"
	case TBYTE:
		return "TBYTE"
	case TSBYTE:
		return "TSBYTE"
	case TLOGICAL:
		return "TLOGICAL"
	case TSTRING:
		return "TSTRING"
	case TUSHORT:
		return "TUSHORT"
	case TSHORT:
		return "TSHORT"
	case TUINT:
		return "TUINT"
	case TINT:
		return "TINT"
	case TULONG:
		return "TULONG"
	case TLONG:
		return "TLONG"
	case TFLOAT:
		return "TFLOAT"
	case TLONGLONG:
		return "TLONGLONG"
	case TDOUBLE:
		return "TDOUBLE"
	case TCOMPLEX:
		return "TCOMPLEX"
	case TDBLCOMPLEX:
		return "TDBLCOMPLEX"

	case TVLABIT:
		return "TVLABIT"
	case TVLABYTE:
		return "TVLABYTE"
	case TVLASBYTE:
		return "TVLASBYTE"
	case TVLALOGICAL:
		return "TVLALOGICAL"
	case TVLASTRING:
		return "TVLASTRING"
	case TVLAUSHORT:
		return "TVLAUSHORT"
	case TVLASHORT:
		return "TVLASHORT"
	case TVLAUINT:
		return "TVLAUINT"
	case TVLAINT:
		return "TVLAINT"
	case TVLAULONG:
		return "TVLAULONG"
	case TVLALONG:
		return "TVLALONG"
	case TVLAFLOAT:
		return "TVLAFLOAT"
	case TVLALONGLONG:
		return "TVLALONGLONG"
	case TVLADOUBLE:
		return "TVLADOUBLE"
	case TVLACOMPLEX:
		return "TVLACOMPLEX"
	case TVLADBLCOMPLEX:
		return "TVLADBLCOMPLEX"
	}
	panic("unreachable")
}

var g_cfits2go map[TypeCode]reflect.Type
var g_go2cfits map[reflect.Type]TypeCode

func govalue_from_repeat(rt reflect.Type, n int) Value {
	var v Value
	switch n {
	case 1:
		rv := reflect.New(rt)
		v = rv.Elem().Interface()
	default:
		// FIXME: distinguish b/w reflect.Slice and reflect.Array
		// FIXME: use reflect.MakeArray + reflect.ArrayOf when available
		rv := reflect.MakeSlice(reflect.SliceOf(rt), n, n)
		v = rv.Interface()
	}
	return v
}

func govalue_from_typecode(t TypeCode, n int) Value {
	var v Value
	rt, ok := g_cfits2go[t]
	if !ok {
		panic(fmt.Errorf("cfitsio: invalid TypeCode value [%v]", int(t)))

	}
	// FIXME: distinguish b/w reflect.Slice and reflect.Array
	v = govalue_from_repeat(rt, n)
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
		col := &hdu.cols[icol]
		err = col.read(hdu.f, icol, irow, &col.Value)
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
			switch vv := card.Value.(type) {
			case float64:
				col.Bscale = vv
			case int64:
				col.Bscale = float64(vv)
			default:
				panic(fmt.Errorf("unhandled type [%T]", vv))
			}
			//col.Bscale = card.Value.(float64)
		} else {
			col.Bscale = 1.0
		}

		card = get("TZERO", ii)
		if card != nil {
			switch vv := card.Value.(type) {
			case float64:
				col.Bzero = vv
			case int64:
				col.Bzero = float64(vv)
			default:
				panic(fmt.Errorf("unhandled type [%T]", vv))
			}
			//col.Bzero = card.Value.(float64)
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
			col.Type = TypeCode(c_type)
			col.Len = int(c_repeat)
			col.Value = govalue_from_typecode(col.Type, col.Len)
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
		col := &cols[i]
		c_name := C.CString(col.Name)
		defer C.free(unsafe.Pointer(c_name))
		C.char_array_set(c_types, c_idx, c_name)

		err = col.inferFormat(hdutype)
		if err != nil {
			return table, err
		}
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

	err := hdu.seekHDU()
	if err != nil {
		return err
	}

	irow := hdu.NumRows()

	defer func() {
		// update nrows
		c_nrows := C.long(0)
		c_status := C.int(0)
		C.fits_get_num_rows(hdu.f.c, &c_nrows, &c_status)
		if c_status > 0 {
			return
		}
		hdu.nrows = int64(c_nrows)
	}()

	switch len(args) {
	case 0:
		return fmt.Errorf("cfitsio: Table.Write needs at least one argument")
	case 1:
		// maybe special case: map? struct?
		rt := reflect.TypeOf(args[0]).Elem()
		switch rt.Kind() {
		case reflect.Map:
			return hdu.writeMap(irow, *args[0].(*map[string]interface{}))
		case reflect.Struct:
			return hdu.writeStruct(irow, args[0])
		}
	}

	return hdu.write(irow, args...)
}

func (hdu *Table) writeMap(irow int64, data map[string]interface{}) error {
	var err error

	for k, v := range data {
		icol := hdu.Index(k)
		if icol < 0 {
			continue
		}
		col := &hdu.cols[icol]
		col.Value = v
		err = col.write(hdu.f, icol, irow, col.Value)
		if err != nil {
			return err
		}
	}

	return err
}

func (hdu *Table) writeStruct(irow int64, data interface{}) error {
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

	for _, icol := range icols {
		vv := rv.Field(icol[0])
		col := &hdu.cols[icol[1]]
		col.Value = vv.Interface()
		err = col.write(hdu.f, icol[1], irow, col.Value)
		if err != nil {
			return err
		}
	}
	return err
}

func (hdu *Table) write(irow int64, args ...interface{}) error {
	var err error

	nargs := len(args)
	if nargs > len(hdu.cols) {
		nargs = len(hdu.cols)
	}

	for i := 0; i < nargs; i++ {
		col := &hdu.cols[i]
		rv := reflect.ValueOf(args[i]).Elem()
		vv := reflect.ValueOf(&col.Value).Elem()
		vv.Set(rv)
		err = col.write(hdu.f, i, irow, col.Value)
		if err != nil {
			return err
		}
	}
	return err
}

// CopyTable copies all the rows from src into dst.
func CopyTable(dst, src *Table) error {
	return CopyTableRange(dst, src, 0, src.NumRows())
}

// CopyTableRange copies the rows interval [beg,end) from src into dst
func CopyTableRange(dst, src *Table, beg, end int64) error {
	var err error
	if dst == nil {
		return fmt.Errorf("cfitsio: dst pointer is nil")
	}
	if src == nil {
		return fmt.Errorf("cfitsio: src pointer is nil")
	}

	defer func() {
		// update nrows
		c_nrows := C.long(0)
		c_status := C.int(0)
		C.fits_get_num_rows(dst.f.c, &c_nrows, &c_status)
		if c_status > 0 {
			return
		}
		dst.nrows = int64(c_nrows)
	}()

	err = dst.seekHDU()
	if err != nil {
		return err
	}
	err = src.seekHDU()
	if err != nil {
		return err
	}

	hdr := src.Header()
	buf := make([]byte, int(hdr.Axes()[0]))
	slice := (*reflect.SliceHeader)((unsafe.Pointer(&buf)))
	c_ptr := (*C.uchar)(unsafe.Pointer(slice.Data))
	c_len := C.LONGLONG(len(buf))
	c_orow := C.LONGLONG(dst.nrows)
	for irow := beg; irow < end; irow++ {
		c_status := C.int(0)
		c_row := C.LONGLONG(irow) + 1 // from 0-based to 1-based index
		C.fits_read_tblbytes(src.f.c, c_row, 1, c_len, c_ptr, &c_status)
		if c_status > 0 {
			return to_err(c_status)
		}

		C.fits_write_tblbytes(dst.f.c, c_orow+c_row, 1, c_len, c_ptr, &c_status)
		if c_status > 0 {
			return to_err(c_status)
		}
	}

	return err
}

func init() {
	g_hdus[ASCII_TBL] = newTable
	g_hdus[BINARY_TBL] = newTable
	g_cfits2go = map[TypeCode]reflect.Type{
		TBIT:        reflect.TypeOf((*byte)(nil)).Elem(),
		TBYTE:       reflect.TypeOf((*byte)(nil)).Elem(),
		TSBYTE:      reflect.TypeOf((*int8)(nil)).Elem(),
		TLOGICAL:    reflect.TypeOf((*bool)(nil)).Elem(),
		TSTRING:     reflect.TypeOf((*string)(nil)).Elem(),
		TUSHORT:     reflect.TypeOf((*uint16)(nil)).Elem(),
		TSHORT:      reflect.TypeOf((*int16)(nil)).Elem(),
		TUINT:       reflect.TypeOf((*uint32)(nil)).Elem(),
		TINT:        reflect.TypeOf((*int32)(nil)).Elem(),
		TULONG:      reflect.TypeOf((*uint64)(nil)).Elem(),
		TLONG:       reflect.TypeOf((*int64)(nil)).Elem(),
		TFLOAT:      reflect.TypeOf((*float32)(nil)).Elem(),
		TLONGLONG:   reflect.TypeOf((*int64)(nil)).Elem(),
		TDOUBLE:     reflect.TypeOf((*float64)(nil)).Elem(),
		TCOMPLEX:    reflect.TypeOf((*complex64)(nil)).Elem(),
		TDBLCOMPLEX: reflect.TypeOf((*complex128)(nil)).Elem(),

		TVLABIT:        reflect.TypeOf((*[]byte)(nil)).Elem(),
		TVLABYTE:       reflect.TypeOf((*[]byte)(nil)).Elem(),
		TVLASBYTE:      reflect.TypeOf((*[]int8)(nil)).Elem(),
		TVLALOGICAL:    reflect.TypeOf((*[]bool)(nil)).Elem(),
		TVLASTRING:     reflect.TypeOf((*[]string)(nil)).Elem(),
		TVLAUSHORT:     reflect.TypeOf((*[]uint16)(nil)).Elem(),
		TVLASHORT:      reflect.TypeOf((*[]int16)(nil)).Elem(),
		TVLAUINT:       reflect.TypeOf((*[]uint32)(nil)).Elem(),
		TVLAINT:        reflect.TypeOf((*[]int32)(nil)).Elem(),
		TVLAULONG:      reflect.TypeOf((*[]uint64)(nil)).Elem(),
		TVLALONG:       reflect.TypeOf((*[]int64)(nil)).Elem(),
		TVLAFLOAT:      reflect.TypeOf((*[]float32)(nil)).Elem(),
		TVLALONGLONG:   reflect.TypeOf((*[]int64)(nil)).Elem(),
		TVLADOUBLE:     reflect.TypeOf((*[]float64)(nil)).Elem(),
		TVLACOMPLEX:    reflect.TypeOf((*[]complex64)(nil)).Elem(),
		TVLADBLCOMPLEX: reflect.TypeOf((*[]complex128)(nil)).Elem(),
	}

	g_go2cfits = make(map[reflect.Type]TypeCode, len(g_cfits2go))
	for k, v := range g_cfits2go {
		g_go2cfits[v] = k
	}
}

// EOF
