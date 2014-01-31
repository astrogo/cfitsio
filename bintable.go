package cfitsio

//
type BinTableHDU struct {
	f      *File
	header Header
	data   interface{}
}

func (hdu *BinTableHDU) Close() error {
	hdu.f = nil
	return nil
}

func (hdu *BinTableHDU) Header() Header {
	return hdu.header
}

func (hdu *BinTableHDU) Type() HduType {
	return hdu.header.htype
}

func (hdu *BinTableHDU) Name() string {
	card := hdu.header.Get("EXTNAME")
	if card == nil {
		return ""
	}
	return card.Value.(string)
}

func (hdu *BinTableHDU) Version() int {
	card := hdu.header.Get("EXTVER")
	if card == nil {
		return 1
	}
	return card.Value.(int)
}

func (hdu *BinTableHDU) Data() (interface{}, error) {
	var err error
	if hdu.data == nil {
		err = hdu.load()
	}
	return hdu.data, err
}

func (hdu *BinTableHDU) load() error {
	return nil
}

func newBinTableHDU(f *File, hdr Header, i int) (hdu HDU, err error) {
	hdu = &BinTableHDU{
		f:      f,
		header: hdr,
		data:   nil,
	}
	return hdu, err
}

func init() {
	//g_hdus[BINARY_TBL] = newBinTableHDU
}

// EOF
