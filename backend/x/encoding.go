package x

import (
	// "fmt"
)

func (b *Backend) readInitResponse(ir *InitResponse) {
	b.readUnused(1)
	b.read(&ir.ProtocolMajorVersion)
	b.read(&ir.ProtocolMinorVersion)
	b.readUnused(2)
	b.read(&ir.ReleaseNumber)
	b.read(&ir.ResourceIdBase)
	b.read(&ir.ResourceIdMask)
	b.read(&ir.MotionBufferSize)
	var vendorLen Card16
	b.read(&vendorLen)
	b.read(&ir.MaximumRequestLength)
	var rootsLen, formatsLen Card8
	b.read(&rootsLen)
	b.read(&formatsLen)
	b.read(&ir.ImageByteOrder)
	b.read(&ir.BitmapBitOrder)
	b.read(&ir.BitmapScanlineUnit)
	b.read(&ir.BitmapScanlinePad)
	b.read(&ir.MinKeyCode)
	b.read(&ir.MaxKeyCode)
	b.readUnused(4)
	vendor := make([]byte, int(vendorLen))
	b.read(&vendor)
	ir.Vendor = string(vendor)
	b.readPadding()
	for i := Card8(0); i < formatsLen; i++ {
		var f Format
		b.readFormat(&f)
		ir.PixmapFormats = append(ir.PixmapFormats, f)
	}
	for i := Card8(0); i < rootsLen; i++ {
		var r Screen
		b.readScreen(&r)
		ir.Roots = append(ir.Roots, r)
	}
}

func (b *Backend) readFormat(f *Format) {
	b.read(&f.Depth)
	b.read(&f.BitsPerPixel)
	b.read(&f.ScanlinePad)
	b.readUnused(5)
}

func (b *Backend) readScreen(s *Screen) {
	b.read(&s.Root)
	b.read(&s.DefaultColormap)
	b.read(&s.WhitePixel)
	b.read(&s.BlackPixel)
	b.readUnused(4) // SetOfEvent not implemented
	b.read(&s.WidthInPixels)
	b.read(&s.HeightInPixels)
	b.read(&s.HeightInMillimiters)
	b.read(&s.WidthInMillimiters)
	b.read(&s.MinInstalledMaps)
	b.read(&s.MaxInstalledMaps)
	b.read(&s.RootVisual)
	b.read(&s.BackingStores)
	b.read(&s.SaveUnders)
	b.read(&s.RootDepth)
	var depthsLen Card8
	b.read(&depthsLen)
	for i := Card8(0); i < depthsLen; i++ {
		var d Depth
		b.readDepth(&d)
		s.AllowedDepths = append(s.AllowedDepths, d)
	}
}

func (b *Backend) readDepth(d *Depth) {
	b.read(&d.Depth)
	b.readUnused(1)
	var visualsLen Card16
	b.read(&visualsLen)
	b.readUnused(4)
	for k := Card16(0); k < visualsLen; k++ {
		var v VisualType
		b.readVisualType(&v)
		d.Visuals = append(d.Visuals, v)
	}
}

func (b *Backend) readVisualType(v *VisualType) {
	b.read(&v.VisualId)
	b.read(&v.Class)
	b.read(&v.BitsPerRgbValue)
	b.read(&v.ColormapEntries)
	b.read(&v.RedMask)
	b.read(&v.GreenMask)
	b.read(&v.BlueMask)
	b.readUnused(4)
}
