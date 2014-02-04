package cfitsio

// PrimaryHDU is the primary HDU
type PrimaryHDU struct {
	ImageHDU
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

func newPrimaryHDU(f *File, hdr Header) (hdu HDU, err error) {
	hdu = &PrimaryHDU{
		ImageHDU{
			f:      f,
			header: hdr,
		},
	}
	return
}

// EOF
