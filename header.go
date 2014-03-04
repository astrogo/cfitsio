package cfitsio

// #include "go-cfitsio.h"
// #include "go-cfitsio-utils.h"
import "C"
import (
	"fmt"
	"unsafe"
)

type Header struct {
	slice  []Card
	cards  map[string]int
	htype  HDUType // header type
	bitpix int64   // character information
	axes   []int64 // dimensions of image data array
}

// NewHeader creates a Header from a set of Cards, HDUType, bitpix and axes.
func NewHeader(cards []Card, htype HDUType, bitpix int64, axes []int64) Header {
	hdr := Header{
		slice:  make([]Card, 0, len(cards)),
		cards:  make(map[string]int, len(cards)),
		htype:  htype,
		bitpix: bitpix,
		axes:   make([]int64, len(axes)),
	}
	copy(hdr.axes, axes)
	hdr.Append(cards...)
	return hdr
}

// NewDefaultHeader creates a Header with CFITSIO default Cards, of type IMAGE_HDU and
// bitpix=8, no axes.
func NewDefaultHeader() Header {
	return NewHeader(
		[]Card{},
		IMAGE_HDU,
		8,
		[]int64{},
	)
}

func (h *Header) AddComment(v string) {
	panic("not implemented")
}

func (h *Header) AddHistory(v string) {
	panic("not implemented")
}

// Append appends a set of Cards to this Header
func (h *Header) Append(cards ...Card) *Header {
	h.slice = append(h.slice, cards...)
	for i := range h.slice {
		h.cards[h.slice[i].Name] = i
	}
	return h
}

func (h *Header) Clear() {
	h.slice = make([]Card, 0)
	h.cards = make(map[string]int)
	h.bitpix = 8
	h.axes = make([]int64, 0)
}

func (h *Header) Get(n string) *Card {
	idx, ok := h.cards[n]
	if ok {
		return &h.slice[idx]
	}
	return nil
}

func (h *Header) Comment() string {
	card := h.Get("COMMENT")
	if card != nil {
		return card.Value.(string)
	}
	return ""
}

func (h *Header) History() string {
	card := h.Get("HISTORY")
	if card != nil {
		return card.Value.(string)
	}
	return ""
}

func (h *Header) Bitpix() int64 {
	return h.bitpix
}

func (h *Header) Axes() []int64 {
	return h.axes
}

func (h *Header) Index(n string) int {
	idx, ok := h.cards[n]
	if ok {
		return idx
	}
	return -1
}

func (h *Header) Keys() []string {
	keys := make([]string, 0, len(h.slice))
	for i := range h.slice {
		keys = append(keys, h.slice[i].Name)
	}
	return keys
}

func (h *Header) Set(n string, v interface{}, comment string) {
	card := h.Get(n)
	if card == nil {
		panic(fmt.Errorf("Header.Set: no such card name %q", n))
	}
	card.Value = v
	card.Comment = comment
}

func readHeader(f *File, i int) (Header, error) {
	var err error
	var hdr Header

	htype, err := f.seekHDU(i, 0)
	if err != nil {
		return hdr, err
	}

	var bitpix int64
	var axes []int64

	switch htype {
	case IMAGE_HDU:
		// do not try to read BITPIX or NAXIS directly.
		// if this is a compressed image, the relevant information will need
		// to come from ZBITPX and ZNAXIS.
		c_status := C.int(0)
		c_val := C.int(0)
		C.fits_get_img_type(f.c, &c_val, &c_status)
		if c_status > 0 {
			return hdr, to_err(c_status)
		}
		bitpix = int64(c_val)

		c_val = 0
		C.fits_get_img_dim(f.c, &c_val, &c_status)
		if c_status > 0 {
			return hdr, to_err(c_status)
		}
		naxes := int(c_val)
		if naxes > 0 {
			axes = make([]int64, naxes)
			c_status := C.int(0)
			c_axes := (*C.long)(unsafe.Pointer(&axes[0]))
			C.fits_get_img_size(f.c, C.int(naxes), c_axes, &c_status)
			if c_status > 0 {
				return hdr, to_err(c_status)
			}
		}
	default:
		bitpix = 8
		// read NAXIS keyword
		c_status := C.int(0)
		c_val := C.long(0)
		c_key := C.CString("NAXIS")
		defer C.free(unsafe.Pointer(c_key))
		var c_comment *C.char = nil
		C.fits_read_key_lng(f.c, c_key, &c_val, c_comment, &c_status)
		if c_status > 0 {
			return hdr, to_err(c_status)
		}
		naxes := int(c_val)
		if naxes > 0 {
			axes = make([]int64, naxes)
			for i := range axes {
				c_status := C.int(0)
				c_key := C.CString(fmt.Sprintf("NAXIS%d", i+1))
				defer C.free(unsafe.Pointer(c_key))
				c_axis := (*C.long)(unsafe.Pointer(&axes[i]))
				C.fits_read_key_lng(f.c, c_key, c_axis, c_comment, &c_status)
				if c_status > 0 {
					return hdr, to_err(c_status)
				}
			}
		}
	}

	c_status := C.int(0)
	c_n := C.int(0)
	c_dummy := C.int(0)

	C.fits_get_hdrpos(f.c, &c_n, &c_dummy, &c_status)
	if c_status > 0 {
		return hdr, to_err(c_status)
	}

	hdr = Header{
		slice:  make([]Card, 0, int(c_n)),
		cards:  make(map[string]int),
		htype:  htype,
		bitpix: bitpix,
		axes:   axes,
	}
	for i := 0; i <= int(c_n); i++ {
		// if the reading of a particular Card fails (most likely
		// due to an undefined value) simple skip and continue to next Card
		card, e := newCard(f, i)
		if e != nil {
			continue
		}
		hdr.Append(card)
	}
	return hdr, err
}

// EOf
