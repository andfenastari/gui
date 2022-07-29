package x

type Window uint32
type Pixmap uint32
type Cursor uint32
type Font uint32
type GContext uint32
type Colormap uint32
type Drawable uint32
type Fontable uint32
type Atom uint32
type VisualId uint32
type Value uint32
type Byte uint8
type Int8 int8
type Int16 int16
type Int32 int32
type Card8 uint8
type Card16 uint16
type Card32 uint32
type Timestamp uint32
type String8 []byte

type BitGravity uint8

const (
	BGForget BitGravity = iota
	BGStatic
	BGNorthWest
	BGNorth
	BGNorthEast
	BGWest
	BGCenter
	BGEast
	BGSouthWest
	BGSouth
	BGSouthEast
)

type WinGravity uint8

const (
	WGUnmap WinGravity = iota
	WGStatic
	WGNorthWest
	WGNorth
	WGNorthEast
	WGWest
	WGCenter
	WGEast
	WGSouthWest
	WGSouth
	WGSouthEast
)

type Bool uint8

const (
	True Bool = iota
	False
)

type Event uint8

const (
	EVKeyPress Event = iota
	EVKeyRelease
	EVOwnerGrabButton
	EVButtonPress
	EVButtonRelease
	EVEnterWindow
	EVLeaveWindow
	EVPointerMotion
	EVPointerMotionHint
	EVButton1Motion
	EVButton2Motion
	EVButton3Motion
	EVButton4Motion
	EVButton5Motion
	EVButtonMotion
	EVExposure
	EVVisibilityChange
	EVStructureNotify
	EVResizeRedirect
	EVSubstructureNotify
	EVSubstructureRedirect
	EVFocusChange
	EVPropertyChange
	EVColormapChange
	EVKeymapState
)

type ByteOrder Card8

const (
	LSBFirst ByteOrder = iota
	MSBFirst
)

type KeyCode Card8

type InitResponse struct {
	ProtocolMajorVersion Card16
	ProtocolMinorVersion Card16
	ReleaseNumber        Card32
	ResourceIdBase       Card32
	ResourceIdMask       Card32
	MotionBufferSize     Card32
	MaximumRequestLength Card16
	ImageByteOrder       ByteOrder
	BitmapBitOrder       ByteOrder
	BitmapScanlineUnit   Card8
	BitmapScanlinePad    Card8
	MinKeyCode           KeyCode
	MaxKeyCode           KeyCode
	Vendor               string
	PixmapFormats        []Format
	Roots                []Screen
}

type Format struct {
	Depth        Card8
	BitsPerPixel Card8
	ScanlinePad  Card8
}

type Screen struct {
	Root                Window
	DefaultColormap     Colormap
	WhitePixel          Card32
	BlackPixel          Card32
	CurrentInputMasks   []Event
	WidthInPixels       Card16
	HeightInPixels      Card16
	WidthInMillimiters  Card16
	HeightInMillimiters Card16
	MinInstalledMaps    Card16
	MaxInstalledMaps    Card16
	RootVisual          VisualId
	BackingStores       Card8
	SaveUnders          Bool
	RootDepth           Card8
	AllowedDepths       []Depth
}

type Depth struct {
	Depth   Card8
	Visuals []VisualType
}

type VisualType struct {
	VisualId        VisualId
	Class           VisualClass
	BitsPerRgbValue Card8
	ColormapEntries Card16
	RedMask         Card32
	GreenMask       Card32
	BlueMask        Card32
}

type VisualClass Card8

const (
	StaticGray VisualClass = iota
	GrayScale
	StaticColor
	PseudoColor
	TrueColor
	DirectColor
)
