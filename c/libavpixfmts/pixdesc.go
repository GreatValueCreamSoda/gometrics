package libavpixfmts

/*
#cgo pkg-config: libavutil
#include <libavutil/pixdesc.h>
#include <stdlib.h>
*/
import "C"
import (
	"errors"
	"fmt"
	"unsafe"
)

type PixFmtFlag int

const (
	PixFmtFlagBigEndian PixFmtFlag = C.AV_PIX_FMT_FLAG_BE        // Pixel format is big-endian.
	PixFmtFlagPAL       PixFmtFlag = C.AV_PIX_FMT_FLAG_PAL       // Pixel format has a palette in data[1], values are indexes in this palette.
	PixFmtFlagBitstream PixFmtFlag = C.AV_PIX_FMT_FLAG_BITSTREAM // All values of a component are bit-wise packed end to end.
	PixFmtFlagHWAccel   PixFmtFlag = C.AV_PIX_FMT_FLAG_HWACCEL   // Pixel format is an HW accelerated format.
	PixFmtFlagPlanar    PixFmtFlag = C.AV_PIX_FMT_FLAG_PLANAR    // At least one pixel component is not in the first data plane.
	PixFmtFlagRGB       PixFmtFlag = C.AV_PIX_FMT_FLAG_RGB       // The pixel format contains RGB-like data (as opposed to YUV/grayscale).
	PixFmtFlagAlpha     PixFmtFlag = C.AV_PIX_FMT_FLAG_ALPHA     //
	PixFmtFlagBayer     PixFmtFlag = C.AV_PIX_FMT_FLAG_BAYER     // The pixel format is following a Bayer pattern
	PixFmtFlagFloat     PixFmtFlag = C.AV_PIX_FMT_FLAG_FLOAT     // The pixel format contains IEEE-754 floating point values. Precision (double, single, or half) should be determined by the pixel size (64, 32, or 16 bits).
	PixFmtFlagXYZ       PixFmtFlag = C.AV_PIX_FMT_FLAG_XYZ       // The pixel format contains XYZ-like data (as opposed to YUV/RGB/grayscale).
)

// PixFmtDescRef is a Descriptor that unambiguously describes how the bits of a
// pixel are stored in the up to 4 data planes of an image. It also stores the
// subsampling factors and number of components.
type PixFmtDescRef struct {
	cptr *C.AVPixFmtDescriptor
}

// ComponentDescriptor is a Descriptor of one of the 4 planes within a
// PixFmtDescRef.
type ComponentDescriptor struct {
	Plane  int
	Step   int
	Offset int
	Shift  int
	Depth  int
}

var (
	ErrUnknownPixelFormat = errors.New("unknown pixel format")
	ErrInvalidArgument    = errors.New("invalid argument")
)

// ---------------- Descriptor acquisition ----------------

// PixFmtDescGet returns a description of the given pixel format.
//
// The returned PixFmtDescRef provides information about the format’s
// structure, including chroma subsampling, component count, and pixel layout.
//
// If the pixel format is unknown, the function returns ErrUnknownPixelFormat.
// The returned PixFmtDescRef is read-only and managed internally; it does not
// require manual cleanup.
func PixFmtDescGet(pf PixelFormat) (*PixFmtDescRef, error) {
	cdesc := C.av_pix_fmt_desc_get(cPixelFormat(pf))
	if cdesc == nil {
		return nil, ErrUnknownPixelFormat
	}
	return &PixFmtDescRef{cptr: cdesc}, nil
}

// PixFmtDescNext returns the next pixel format descriptor.
//
// Calling this on a nil receiver returns the first descriptor. Calling it on
// an existing descriptor returns the one that follows it.
//
// If there are no more descriptors, it returns nil. The returned descriptor
// is read-only and does not require manual cleanup.
func (prev *PixFmtDescRef) PixFmtDescNext() *PixFmtDescRef {
	var cprev *C.AVPixFmtDescriptor
	if prev != nil {
		cprev = prev.cptr
	}
	next := C.av_pix_fmt_desc_next(cprev)
	if next == nil {
		return nil
	}
	return &PixFmtDescRef{cptr: next}
}

// PixFmtDescID returns the pixel format identifier described by this
// descriptor.
//
// If the descriptor is nil or invalid, it returns ErrInvalidArgument. If the
// descriptor corresponds to an unknown format, it returns
// ErrUnknownPixelFormat.
//
// A successful call yields the PixelFormat value that this descriptor
// represents.
func (desc *PixFmtDescRef) PixFmtDescID() (PixelFormat, error) {
	if desc == nil || desc.cptr == nil {
		return PixelFormat(C.AV_PIX_FMT_NONE), ErrInvalidArgument
	}
	id := C.av_pix_fmt_desc_get_id(desc.cptr)
	if id == C.AV_PIX_FMT_NONE {
		return PixelFormat(id), ErrUnknownPixelFormat
	}
	return PixelFormat(id), nil
}

// Name returns the symbolic name of this pixel format.
//
// The name is a short identifier such as "yuv420p" or "rgba". If the
// descriptor is nil or does not provide a name, an empty string is returned.
func (r *PixFmtDescRef) Name() string {
	if r == nil || r.cptr == nil || r.cptr.name == nil {
		return ""
	}
	return C.GoString(r.cptr.name)
}

// NbComponents returns the number of color components in this pixel format.
//
// For example, RGB formats have three components, RGBA formats have four,
// and many YUV formats have three. If the descriptor is nil, zero is returned.
func (r *PixFmtDescRef) NbComponents() int {
	if r == nil || r.cptr == nil {
		return 0
	}
	return int(r.cptr.nb_components)
}

// Log2ChromaW returns the horizontal chroma subsampling factor in base-2.
//
// Many YUV formats reduce chroma resolution horizontally. For instance,
// 4:2:0 formats have a factor of 1 (corresponding to a halving of horizontal
// chroma resolution). Formats without horizontal subsampling report 0.
// A nil descriptor yields 0.
func (r *PixFmtDescRef) Log2ChromaW() int {
	if r == nil || r.cptr == nil {
		return 0
	}
	return int(r.cptr.log2_chroma_w)
}

// Log2ChromaH returns the vertical chroma subsampling factor in base-2.
//
// Similar to Log2ChromaW, many YUV formats reduce vertical chroma resolution.
// For example, a 4:2:0 format reports 1 (halved vertical chroma). A nil
// descriptor yields 0.
func (r *PixFmtDescRef) Log2ChromaH() int {
	if r == nil || r.cptr == nil {
		return 0
	}
	return int(r.cptr.log2_chroma_h)
}

// Flags returns a bitmask describing the properties of this pixel format.
//
// The flags indicate characteristics such as whether the format contains an
// alpha channel, is planar or packed, or other structural traits. A nil
// descriptor returns 0, meaning no flags are set.
func (r *PixFmtDescRef) Flags() uint64 {
	if r == nil || r.cptr == nil {
		return 0
	}
	return uint64(r.cptr.flags)
}

// Component returns the descriptor for the i-th color component of this pixel
// format.
//
// The descriptor provides information about the component’s plane, step size,
// bit offset, shift, and bit depth. The index i must be between 0 and 3.
// If the index is out of range, or if the descriptor is nil,
// ErrInvalidArgument is returned.
func (r *PixFmtDescRef) Component(i int) (ComponentDescriptor, error) {
	var zero ComponentDescriptor
	if r == nil || r.cptr == nil {
		return zero, ErrInvalidArgument
	}
	if i < 0 || i >= 4 {
		return zero, fmt.Errorf("component index out of range: %d", i)
	}
	ccomp := r.cptr.comp[i]
	cd := ComponentDescriptor{
		Plane:  int(ccomp.plane),
		Step:   int(ccomp.step),
		Offset: int(ccomp.offset),
		Shift:  int(ccomp.shift),
		Depth:  int(ccomp.depth),
	}
	return cd, nil
}

// Alias returns an alternate symbolic name for this pixel format, if one
// exists.
//
// Some pixel formats have historical or shorthand names. If an alternate name
// exists, it is returned as a Go string. If no alias exists or the descriptor
// is nil, an empty string is returned.
func (r *PixFmtDescRef) Alias() string {
	if r == nil || r.cptr == nil || r.cptr.alias == nil {
		return ""
	}
	return C.GoString(r.cptr.alias)
}

// GetBitsPerPixel returns the number of meaningful bits per pixel for this
// format.
//
// This counts all components in the format. It may differ from the actual
// storage size due to padding or alignment. If the descriptor is nil,
// ErrInvalidArgument is returned.
func GetBitsPerPixel(r *PixFmtDescRef) (int, error) {
	if r == nil || r.cptr == nil {
		return 0, ErrInvalidArgument
	}
	ret := C.av_get_bits_per_pixel(r.cptr)
	return int(ret), nil
}

// GetPaddedBitsPerPixel returns the number of bits per pixel including
// padding.
//
// This accounts for alignment and packing rules in the format. The value is
// generally larger than the raw bits-per-pixel count. If the descriptor is
// nil, ErrInvalidArgument is returned.
func GetPaddedBitsPerPixel(r *PixFmtDescRef) (int, error) {
	if r == nil || r.cptr == nil {
		return 0, ErrInvalidArgument
	}
	ret := C.av_get_padded_bits_per_pixel(r.cptr)
	return int(ret), nil
}

// PixFmtGetChromaSubSample returns the horizontal and vertical chroma shift
// factors for the given pixel format.
//
// The returned values indicate how much chroma is downsampled relative to
// luma, expressed as base-2 logarithms. For example, many 4:2:0 formats return
// hShift == 1 and vShift == 1, meaning chroma has half the resolution
// horizontally and vertically.
//
// If the pixel format is invalid, an error is returned.
func PixFmtGetChromaSubSample(pf PixelFormat) (hShift int, vShift int, err error) {
	var ch C.int
	var cv C.int
	ret := C.av_pix_fmt_get_chroma_sub_sample(cPixelFormat(pf), &ch, &cv)
	if ret < 0 {
		return 0, 0, fmt.Errorf("av_pix_fmt_get_chroma_sub_sample error: %d", int(ret))
	}
	return int(ch), int(cv), nil
}

// PixFmtCountPlanes returns the number of image planes used by the given pixel
// format.
//
// For example, packed RGB formats use a single plane, planar YUV formats often
// use three planes, and certain specialized formats may have more.
//
// If the pixel format is invalid, an error is returned.
func PixFmtCountPlanes(pf PixelFormat) (int, error) {
	ret := C.av_pix_fmt_count_planes(cPixelFormat(pf))
	if ret < 0 {
		return 0, fmt.Errorf("av_pix_fmt_count_planes error: %d", int(ret))
	}
	return int(ret), nil
}

// ColorRangeName returns the canonical name of a color range value.
//
// Color range describes how pixel values map to light intensities. For
// example, "limited" or "tv" ranges reserve headroom for broadcast safety,
// while "full" or "pc" ranges use the entire value range.
//
// If the value is unrecognized, an empty string is returned.
func ColorRangeName(r int) string {
	c := C.av_color_range_name(C.enum_AVColorRange(r))
	if c == nil {
		return ""
	}
	return C.GoString(c)
}

// ColorRangeFromName converts a color range name into its corresponding value.
//
// Recognized names include "limited", "full", "tv", "pc", and similar standard
// ranges. The lookup is case-insensitive.
//
// If the name is empty or unrecognized,an error is returned.
func ColorRangeFromName(name string) (int, error) {
	if name == "" {
		return -1, ErrInvalidArgument
	}
	cs := C.CString(name)
	defer C.free(unsafe.Pointer(cs))
	ret := C.av_color_range_from_name(cs)
	if ret < 0 {
		return int(ret), fmt.Errorf("unknown color range: %s (err %d)", name, int(ret))
	}
	return int(ret), nil
}

// ColorPrimariesName returns the descriptive name for a given color primaries
// value.
//
// Color primaries define the RGB chromaticity coordinates used by a format,
// such as BT.709, BT.2020, or SMPTE 240M. They affect color reproduction and
// gamut conversion.
//
// If the value is unrecognized, an empty string is returned.
func ColorPrimariesName(p int) string {
	c := C.av_color_primaries_name(C.enum_AVColorPrimaries(p))
	if c == nil {
		return ""
	}
	return C.GoString(c)
}

// ColorPrimariesFromName converts a color primaries name into its
// corresponding value.
//
// Recognized names include standard sets like "BT.709", "BT.2020", and "SMPTE
// 170M". The lookup is case-insensitive.
//
// If the input string is empty or unrecognized, an error is returned. A
// successful call returns the integer identifier for the primaries.
func ColorPrimariesFromName(name string) (int, error) {
	if name == "" {
		return -1, ErrInvalidArgument
	}
	cs := C.CString(name)
	defer C.free(unsafe.Pointer(cs))
	ret := C.av_color_primaries_from_name(cs)
	if ret < 0 {
		return int(ret), fmt.Errorf("unknown color primaries: %s (err %d)", name, int(ret))
	}
	return int(ret), nil
}

// ColorTransferName returns the name associated with a given color transfer
// characteristic.
//
// Color transfer characteristics define how digital values map to physical
// light, including gamma curves, HDR transfer functions, HLG, and legacy
// power-law encodings.
//
// If the value is unrecognized, an empty string is returned.
func ColorTransferName(t int) string {
	c := C.av_color_transfer_name(C.enum_AVColorTransferCharacteristic(t))
	if c == nil {
		return ""
	}
	return C.GoString(c)
}

// ColorTransferFromName converts a color transfer characteristic name into
// its corresponding value.
//
// Recognized names include standard transfer curves like "bt709", "smpte240m",
// "arib-std-b67" (HLG), and "smpte2084" (PQ). The lookup is case-insensitive.
//
// If the input string is empty or unrecognized, an error is returned. A
// successful call returns the integer identifier for the transfer
// characteristic.
func ColorTransferFromName(name string) (int, error) {
	if name == "" {
		return -1, ErrInvalidArgument
	}
	cs := C.CString(name)
	defer C.free(unsafe.Pointer(cs))
	ret := C.av_color_transfer_from_name(cs)
	if ret < 0 {
		return int(ret), fmt.Errorf("unknown color transfer: %s (err %d)", name, int(ret))
	}
	return int(ret), nil
}

// ColorSpaceName returns the name associated with a given color space value.
//
// A color space defines how RGB channels are converted to or from luma and
// chroma components. Standard examples include BT.601, BT.709, and BT.2020.
//
// If the value is unrecognized, an empty string is returned.
func ColorSpaceName(s int) string {
	c := C.av_color_space_name(C.enum_AVColorSpace(s))
	if c == nil {
		return ""
	}
	return C.GoString(c)
}

// ColorSpaceFromName converts a color-space name into its corresponding value.
//
// Recognized names include standard color spaces like "BT.601", "BT.709", and
// "BT.2020". The lookup is case-insensitive.
//
// If the input string is empty or unrecognized, an error is returned. A
// successful call returns the integer identifier for the color space.
func ColorSpaceFromName(name string) (int, error) {
	if name == "" {
		return -1, ErrInvalidArgument
	}
	cs := C.CString(name)
	defer C.free(unsafe.Pointer(cs))
	ret := C.av_color_space_from_name(cs)
	if ret < 0 {
		return int(ret), fmt.Errorf("unknown color space: %s (err %d)", name, int(ret))
	}
	return int(ret), nil
}

// ChromaLocationName returns the name associated with a given chroma location
// value.
//
// Chroma location defines the spatial alignment of chroma samples relative to
// luma samples in a subsampled image format (e.g., 4:2:0). Common names
// include "left", "center", and "top-left".
//
// If the value is unrecognized, an empty string is returned.
func ChromaLocationName(loc int) string {
	c := C.av_chroma_location_name(C.enum_AVChromaLocation(loc))
	if c == nil {
		return ""
	}
	return C.GoString(c)
}

// ChromaLocationFromName converts a chroma location name into its
// corresponding value.
//
// Recognized names include standard locations such as "left", "center", and
// "top-left". The lookup is case-insensitive.
//
// An empty input or unknown name produces an error. A successful call returns
// the integer identifier for the chroma location.
func ChromaLocationFromName(name string) (int, error) {
	if name == "" {
		return -1, ErrInvalidArgument
	}
	cs := C.CString(name)
	defer C.free(unsafe.Pointer(cs))
	ret := C.av_chroma_location_from_name(cs)
	if ret < 0 {
		return int(ret), fmt.Errorf("unknown chroma location: %s (err %d)", name, int(ret))
	}
	return int(ret), nil
}

// ChromaLocationEnumToPos converts a chroma-location enumeration value into
// its (x, y) sample position offsets.
//
// The returned xpos and ypos values indicate the horizontal and vertical
// placement of chroma samples relative to the top-left corner of a luma
// sample. If the enumeration is unrecognized, the function returns an error.
func ChromaLocationEnumToPos(pos int) (xpos, ypos int, err error) {
	var cx C.int
	var cy C.int
	ret := C.av_chroma_location_enum_to_pos(&cx, &cy, C.enum_AVChromaLocation(pos))
	if ret < 0 {
		return 0, 0, fmt.Errorf("av_chroma_location_enum_to_pos error: %d", int(ret))
	}
	return int(cx), int(cy), nil
}

// ChromaLocationPosToEnum converts chroma sample coordinates into the
// corresponding AVChromaLocation enumeration value.
//
// xpos and ypos specify the horizontal and vertical chroma sample offsets.
// The function returns the enumeration value corresponding to these offsets.
func ChromaLocationPosToEnum(xpos, ypos int) (int, error) {
	ret := C.av_chroma_location_pos_to_enum(C.int(xpos), C.int(ypos))
	// This function returns an enum value; it doesn't return an error code in header.
	return int(ret), nil
}

// GetPixFmt retrieves the FFmpeg pixel-format enumeration corresponding to
// a textual name.
//
// Pixel formats define how pixel data is stored in memory, including bit depth
// and component ordering (e.g., "yuv420p", "rgb24").
//
// An empty name returns ErrInvalidArgument. Unknown names produce a Go error.
// On success, the function returns the corresponding PixelFormat value.
func GetPixFmt(name string) (PixelFormat, error) {
	if name == "" {
		return PixelFormat(C.AV_PIX_FMT_NONE), ErrInvalidArgument
	}
	cs := C.CString(name)
	defer C.free(unsafe.Pointer(cs))
	ret := C.av_get_pix_fmt(cs)
	if ret == C.AV_PIX_FMT_NONE {
		return PixelFormat(ret), fmt.Errorf("pixel format not found: %s", name)
	}
	return PixelFormat(ret), nil
}

// GetPixFmtName returns the canonical FFmpeg name for a given PixelFormat
// enumeration value.
//
// If the format is unrecognized, an empty string is returned. The result is a
// safe Go string copy.
func GetPixFmtName(pf PixelFormat) string {
	c := C.av_get_pix_fmt_name(cPixelFormat(pf))
	if c == nil {
		return ""
	}
	return C.GoString(c)
}

// GetPixFmtString returns a detailed string describing the pixel format.
//
// The returned string includes the format name, bit depth, component layout,
// and other descriptive information. This is similar to FFmpeg’s av_get_pix_fmt_string.
func GetPixFmtString(pf PixelFormat) (string, error) {
	buffer := C.malloc(1024)
	defer C.free(buffer)
	cbuf := (*C.char)(unsafe.Pointer(buffer))
	C.av_get_pix_fmt_string(cbuf, 1024, cPixelFormat(pf))
	return C.GoString(cbuf), nil
}

// ---------------- Image line read/write ----------------

// ReadImageLine2 reads a horizontal line of pixel data from an image into a
// destination buffer, supporting planar and paletted formats.
//
// dst is the output byte slice where the pixel data will be stored. data and
// linesize provide the source planes and their corresponding line sizes
// (stride) for each component.
//
// desc is a pointer to the pixel format descriptor (PixFmtDescRef), describing
// the layout, number of components, and bit depth.
//
// x and y specify the starting pixel coordinates in the source image, c
// selects the component to read, and w is the number of pixels to copy.
//
// If readPalComponent is true, the function will read the palette component
// for paletted formats.
//
// dstElementSize indicates the number of bytes per element in the destination
// buffer.
//
// The function automatically handles planar and packed formats, converting
// component offsets according to the pixel format descriptor. If desc is nil
// or its internal C pointer is nil, ErrInvalidArgument is returned. Otherwise,
// the function fills dst with the requested pixel data.
func ReadImageLine2(dst []byte, data [4][]byte, linesize [4]int,
	desc *PixFmtDescRef, x, y, c, w int, readPalComponent bool,
	dstElementSize int) error {
	if desc == nil || desc.cptr == nil {
		return ErrInvalidArgument
	}
	var dstPtr unsafe.Pointer
	if len(dst) > 0 {
		dstPtr = unsafe.Pointer(&dst[0])
	} else {
		dstPtr = nil
	}
	cdata, clinesize := buildPlanePointers(data, linesize)
	var readPal C.int
	if readPalComponent {
		readPal = 1
	}

	var cX, cY, cC, cW C.int = C.int(x), C.int(y), C.int(c), C.int(w)

	C.av_read_image_line2(dstPtr, (**C.uint8_t)(unsafe.Pointer(&cdata[0])),
		(*C.int)(unsafe.Pointer(&clinesize[0])), desc.cptr, cX, cY, cC, cW,
		readPal, C.int(dstElementSize))
	return nil
}

// ReadImageLine reads a horizontal line of pixel data from an image into a
// uint16 destination buffer, supporting planar and paletted formats.
//
// dst is the output slice of uint16 values that will receive the pixel data.
// data and linesize provide the source image planes and their line sizes
// (stride) for each component.
//
// desc is a pointer to the pixel format descriptor (PixFmtDescRef), which
// describes the layout, number of components, and bit depth.
//
// x and y specify the starting pixel coordinates in the source image, c
// selects the component to read, and w is the number of pixels to copy.
//
// If readPalComponent is true, the function will read the palette component
// for paletted formats.
//
// The function automatically handles planar and packed formats, converting
// component offsets according to the pixel format descriptor. If desc is nil
// or its internal C pointer is nil, ErrInvalidArgument is returned. On
// success, dst will be filled with the requested pixel data.
func ReadImageLine(dst []uint16, data [4][]byte, linesize [4]int,
	desc *PixFmtDescRef, x, y, c, w int, readPalComponent bool) error {
	if desc == nil || desc.cptr == nil {
		return ErrInvalidArgument
	}
	var dstPtr unsafe.Pointer
	if len(dst) > 0 {
		dstPtr = unsafe.Pointer(&dst[0])
	} else {
		dstPtr = nil
	}
	cdata, clinesize := buildPlanePointers(data, linesize)
	var readPal C.int
	if readPalComponent {
		readPal = 1
	}
	C.av_read_image_line((*C.uint16_t)(dstPtr),
		(**C.uint8_t)(unsafe.Pointer(&cdata[0])),
		(*C.int)(unsafe.Pointer(&clinesize[0])),
		desc.cptr,
		C.int(x), C.int(y), C.int(c), C.int(w),
		readPal)
	return nil
}

// WriteImageLine2 writes a horizontal line of pixel data from a byte slice
// into an image, supporting planar and packed formats.
//
// src is the source byte slice containing pixel values to write. data and
// linesize provide the target image planes and their corresponding line sizes
// (stride) for each component.
//
// desc is a pointer to the pixel format descriptor (PixFmtDescRef), which
// describes the layout, number of components, and bit depth.
//
// x and y specify the starting pixel coordinates in the destination image, c
// selects the component to write, and w is the number of pixels to copy.
// srcElementSize indicates the size in bytes of each element in src.
//
// The function automatically handles planar and packed formats, writing each
// component according to the pixel format descriptor. If desc is nil or its
// internal C pointer is nil, ErrInvalidArgument is returned.
func WriteImageLine2(src []byte, data [4][]byte, linesize [4]int,
	desc *PixFmtDescRef, x, y, c, w int, srcElementSize int) error {
	if desc == nil || desc.cptr == nil {
		return ErrInvalidArgument
	}
	var srcPtr unsafe.Pointer
	if len(src) > 0 {
		srcPtr = unsafe.Pointer(&src[0])
	} else {
		srcPtr = nil
	}
	cdata, clinesize := buildPlanePointers(data, linesize)
	C.av_write_image_line2(srcPtr,
		(**C.uint8_t)(unsafe.Pointer(&cdata[0])),
		(*C.int)(unsafe.Pointer(&clinesize[0])),
		desc.cptr,
		C.int(x), C.int(y), C.int(c), C.int(w),
		C.int(srcElementSize))
	return nil
}

// WriteImageLine writes a horizontal line of pixel data from a uint16 slice
// into an image, supporting planar and packed formats.
//
// src is the source slice of uint16 values to write. data and linesize provide
// the target image planes and their corresponding line sizes (stride) for each
// component. desc is a pointer to the pixel format descriptor (PixFmtDescRef),
// describing the layout, number of components, and bit depth.
//
// x and y specify the starting pixel coordinates in the destination image, c
// selects the component to write, and w is the number of pixels to copy.
//
// The function automatically handles planar and packed formats. If desc is
// nil or its internal C pointer is nil, ErrInvalidArgument is returned. On
// success, the specified pixel data from src is written into the target image.
func WriteImageLine(src []uint16, data [4][]byte, linesize [4]int,
	desc *PixFmtDescRef, x, y, c, w int) error {
	if desc == nil || desc.cptr == nil {
		return ErrInvalidArgument
	}
	var srcPtr unsafe.Pointer
	if len(src) > 0 {
		srcPtr = unsafe.Pointer(&src[0])
	} else {
		srcPtr = nil
	}
	cdata, clinesize := buildPlanePointers(data, linesize)
	C.av_write_image_line((*C.uint16_t)(srcPtr),
		(**C.uint8_t)(unsafe.Pointer(&cdata[0])),
		(*C.int)(unsafe.Pointer(&clinesize[0])),
		desc.cptr,
		C.int(x), C.int(y), C.int(c), C.int(w))
	return nil
}

// ---------------- Endianness swap ----------------

// PixFmtSwapEndianness returns the pixel format with swapped endianness for
// the given PixelFormat.
//
// Some pixel formats have different memory layouts depending on the system’s
// endianness. This function returns the equivalent format with the opposite
// endianness. If no swapped-endianness version exists, an error is returned.
func PixFmtSwapEndianness(pf PixelFormat) (PixelFormat, error) {
	ret := C.av_pix_fmt_swap_endianness(cPixelFormat(pf))
	if ret == C.AV_PIX_FMT_NONE {
		return PixelFormat(ret), fmt.Errorf("no swapped-endianness pixel format for: %d", int(pf))
	}
	return PixelFormat(ret), nil
}

// ---------------- Loss & best-format selection ----------------

const (
	FF_LOSS_RESOLUTION        = 0x0001
	FF_LOSS_DEPTH             = 0x0002
	FF_LOSS_COLORSPACE        = 0x0004
	FF_LOSS_ALPHA             = 0x0008
	FF_LOSS_COLORQUANT        = 0x0010
	FF_LOSS_CHROMA            = 0x0020
	FF_LOSS_EXCESS_RESOLUTION = 0x0040
	FF_LOSS_EXCESS_DEPTH      = 0x0080
)

// GetPixFmtLoss returns an integer representing the information loss when
// converting from the source pixel format to the destination pixel format.
//
// dst and src are the destination and source PixelFormat values, respectively.
// hasAlpha indicates whether the alpha channel should be considered in the
// comparison.
//
// The returned value is a non-negative integer where 0 indicates no loss and
// higher values indicate greater loss of color precision or components.
func GetPixFmtLoss(dst, src PixelFormat, hasAlpha bool) (int, error) {
	var ha C.int
	if hasAlpha {
		ha = 1
	} else {
		ha = 0
	}
	ret := C.av_get_pix_fmt_loss(cPixelFormat(dst), cPixelFormat(src), ha)
	return int(ret), nil
}

// FindBestPixFmtOf2 returns the best pixel format from two candidates for
// converting a source pixel format, along with the associated loss.
//
// dst1 and dst2 are candidate destination PixelFormats. src is the source
// PixelFormat. hasAlpha indicates whether the alpha channel should be
// considered.
//
// The function returns the chosen PixelFormat, an integer loss value, and an
// error if no suitable format could be found. The loss value quantifies the
// difference between the source and selected destination format.
func FindBestPixFmtOf2(dst1, dst2, src PixelFormat, hasAlpha bool) (best PixelFormat, loss int, err error) {
	var lossC C.int
	var ha C.int
	if hasAlpha {
		ha = 1
	}
	ret := C.av_find_best_pix_fmt_of_2(cPixelFormat(dst1), cPixelFormat(dst2), cPixelFormat(src), ha, &lossC)
	if ret == C.AV_PIX_FMT_NONE {
		return PixelFormat(ret), int(lossC), fmt.Errorf("no suitable pixel format found")
	}
	return PixelFormat(ret), int(lossC), nil
}

// Helpers

// helper: convert Go PixelFormat to C enum
func cPixelFormat(p PixelFormat) C.enum_AVPixelFormat {
	return C.enum_AVPixelFormat(p)
}

// helper: prepare C arrays for data[4] and linesize[4]
func buildPlanePointers(data [4][]byte, linesize [4]int) (cdata [4]*C.uint8_t, clinesize [4]C.int) {
	for i := 0; i < 4; i++ {
		if len(data[i]) > 0 {
			cdata[i] = (*C.uint8_t)(unsafe.Pointer(&data[i][0]))
		} else {
			cdata[i] = nil
		}
		clinesize[i] = C.int(linesize[i])
	}
	return
}
