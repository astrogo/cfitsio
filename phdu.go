package cfitsio

// PrimaryHDU is the primary HDU
type PrimaryHDU struct {
	f      *File
	header Header
	data   interface{}
	read   bool // whether the image has been loaded from FITS
}

func (hdu *PrimaryHDU) Close() error {
	hdu.f = nil
	return nil
}

func (hdu *PrimaryHDU) Header() Header {
	return hdu.header
}

func (hdu *PrimaryHDU) Type() HDUType {
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
	if !hdu.read {
		err = hdu.load()
	}
	return hdu.data, err
}

func (hdu *PrimaryHDU) load() error {
	data, err := loadImageData(hdu.f, &hdu.header)
	if err != nil {
		return err
	}
	hdu.read = true
	hdu.data = data
	return err
}

// EOF
