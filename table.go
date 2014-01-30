package cfitsio

//
type TableHDU struct {
	header Header
	data   interface{}
}

func (hdu *TableHDU) Header() Header {
	return hdu.header
}

func (hdu *TableHDU) Type() HduType {
	return hdu.header.htype
}

func (hdu *TableHDU) Name() string {
	card := hdu.header.Get("EXTNAME")
	if card == nil {
		return ""
	}
	return card.Value.(string)
}

func (hdu *TableHDU) Version() int {
	card := hdu.header.Get("EXTVER")
	if card == nil {
		return 1
	}
	return card.Value.(int)
}

func (hdu *TableHDU) Data() (interface{}, error) {
	var err error
	if hdu.data == nil {
		err = hdu.load()
	}
	return hdu.data, err
}

func (hdu *TableHDU) load() error {
	return nil
}

func newTableHDU(f *File, hdr Header, i int) (hdu HDU, err error) {
	hdu = &TableHDU{
		header: hdr,
		data:   nil,
	}
	return hdu, err
}

func init() {
	g_hdus[ASCII_TBL] = newTableHDU
}

// EOF
