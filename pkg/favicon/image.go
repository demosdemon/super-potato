package favicon

import (
	"bytes"
	"encoding/binary"
	"image"
	"image/png"
	"io"
)

type direntry struct {
	Width   uint8
	Height  uint8
	Palette uint8
	_       byte
	Plane   uint16
	Bits    uint16
	Size    uint32
	Offset  uint32
}

type dir struct {
	_      uint16
	Type   uint16
	Number uint16
}

func Encode(w io.Writer, imgs ...image.Image) error {
	n := uint16(len(imgs))
	head := dir{
		Type:   1,
		Number: n,
	}

	entries := make([]direntry, n)
	images := make([]*bytes.Buffer, n)

	offset := uint32(binary.Size(head))

	for idx, img := range imgs {
		offset += uint32(binary.Size(entries[idx]))
		b := img.Bounds()
		//m := image.NewRGBA(b)
		//draw.Draw(m, m.Bounds(), img, b.Min, draw.Src)

		bb := new(bytes.Buffer)
		if err := png.Encode(bb, img); err != nil {
			return err
		}

		size := uint32(bb.Len())

		entries[idx] = direntry{
			Width:  uint8(b.Dx()),
			Height: uint8(b.Dy()),
			Plane:  1,
			Bits:   32,
			Size:   size,
			Offset: offset,
		}

		images[idx] = bb

		offset += size
	}

	if err := binary.Write(w, binary.LittleEndian, head); err != nil {
		return err
	}

	if err := binary.Write(w, binary.LittleEndian, entries); err != nil {
		return err
	}

	for _, bb := range images {
		if _, err := bb.WriteTo(w); err != nil {
			return err
		}
	}

	return nil
}
