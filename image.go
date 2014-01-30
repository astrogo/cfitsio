package cfitsio

//
type ImageHDU struct {
	header Header
	data   interface{}
}

func (hdu *ImageHDU) Header() Header {
	return hdu.header
}

func (hdu *ImageHDU) Type() HduType {
	return hdu.header.htype
}

func (hdu *ImageHDU) Name() string {
	card := hdu.header.Get("EXTNAME")
	if card == nil {
		return ""
	}
	return card.Value.(string)
}

func (hdu *ImageHDU) Version() int {
	card := hdu.header.Get("EXTVER")
	if card == nil {
		return 1
	}
	return card.Value.(int)
}

func (hdu *ImageHDU) Data() (interface{}, error) {
	var err error
	if hdu.data == nil {
		err = hdu.load()
	}
	return hdu.data, err
}

func (hdu *ImageHDU) load() error {
	return nil
}

func newImageHDU(f *File, hdr Header, i int) (hdu HDU, err error) {
	switch i {
	case 0:
		hdu = &PrimaryHDU{
			header: hdr,
			data:   nil,
		}
	default:
		hdu = &ImageHDU{
			header: hdr,
			data:   nil,
		}
	}
	return hdu, err
}

func init() {
	g_hdus[IMAGE_HDU] = newImageHDU
}

// EOF