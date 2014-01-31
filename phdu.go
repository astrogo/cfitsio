package cfitsio

// PrimaryHDU is the primary HDU
type PrimaryHDU struct {
	f      *File
	header Header
	data   interface{}
}

func (hdu *PrimaryHDU) Close() error {
	hdu.f = nil
	return nil
}

func (hdu *PrimaryHDU) Header() Header {
	return hdu.header
}

func (hdu *PrimaryHDU) Type() HduType {
	return hdu.header.htype
}

func (hdu *PrimaryHDU) Name() string {
	card := hdu.header.Get("EXTNAME")
	if card == nil {
		return "PRIMARY"
	}
	return card.Value.(string)
}

func (hdu *PrimaryHDU) Version() int {
	card := hdu.header.Get("EXTVER")
	if card == nil {
		return 1
	}
	return card.Value.(int)
}

func (hdu *PrimaryHDU) Data() (interface{}, error) {
	var err error
	if hdu.data == nil {
		err = hdu.load()
	}
	return hdu.data, err
}

func (hdu *PrimaryHDU) load() error {
	return nil
}

// EOF
