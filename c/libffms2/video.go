package libffms2

//#cgo LDFLAGS: -lffms2
//#cgo CFLAGS: -I/usr/include
//#include <ffms.h>
//#include <stdlib.h>
import "C"
import (
	"errors"
	"unsafe"

	pixfmts "github.com/GreatValueCreamSoda/gometrics/c/libavpixfmts"
)

var (
	ErrInvalidOrNilVideoSource error = errors.New("video source was consumed, failed to create, or was destroyed")
)

// A struct representing a Video source that can be read from and have it's
// properties listed.
type VideoSource struct {
	source *C.FFMS_VideoSource
}

func CreateVideoSource(sourceFile string, index *Index, track,
	threads int, seekMode SeekMode) (*VideoSource, *ErrorInfo, error) {

	if err := index.checkValidity(); err != nil {
		return nil, nil, err
	}

	var sourceFileC *C.char = C.CString(sourceFile)
	defer safeFree(sourceFileC)

	fn := func(c *C.FFMS_ErrorInfo) *C.FFMS_VideoSource {
		return C.FFMS_CreateVideoSource(sourceFileC, C.int(track), index.index,
			C.int(threads), C.int(seekMode), c)
	}

	res, info, err := withErrorInfo(fn)
	if err != nil {
		return nil, info, err
	}

	return &VideoSource{res}, info, nil
}

func (vs *VideoSource) GetVideoProperties() (VideoProperties, error) {
	if err := vs.checkValidity(); err != nil {
		return VideoProperties{}, err
	}

	cVideoProperties := C.FFMS_GetVideoProperties(vs.source)
	if cVideoProperties == nil {
		return VideoProperties{}, ErrFFmsNilPtrReturn
	}

	var videoProperties VideoProperties
	videoProperties.videoPropertiesFromC(cVideoProperties)

	return videoProperties, nil
}

func (vs *VideoSource) GetFrame(frameNumber int) (Frame, *ErrorInfo, error) {
	if err := vs.checkValidity(); err != nil {
		return Frame{}, nil, err
	}

	res, info, err := withErrorInfo(func(c *C.FFMS_ErrorInfo) *C.FFMS_Frame {
		return C.FFMS_GetFrame(vs.source, C.int(frameNumber), c)
	})

	var frame Frame
	frame.fromCFrame(res)

	return frame, info, err
}

func (vs *VideoSource) GetFrameByTime(timeStamp float64) (Frame, *ErrorInfo, error) {
	if err := vs.checkValidity(); err != nil {
		return Frame{}, nil, err
	}

	res, info, err := withErrorInfo(func(c *C.FFMS_ErrorInfo) *C.FFMS_Frame {
		return C.FFMS_GetFrameByTime(vs.source, C.double(timeStamp), c)
	})

	var frame Frame
	frame.fromCFrame(res)

	return frame, info, err
}

func (vs *VideoSource) SetOutputFormatV2(TargetFormats []int, width,
	height int, resizer Resizers) (int, *ErrorInfo, error) {
	if err := vs.checkValidity(); err != nil {
		return 0, nil, nil
	}

	cTargetFormats := (*C.int)(C.malloc(C.size_t(unsafe.Sizeof(C.int(0))) *
		C.size_t(len(TargetFormats))))
	defer safeFree(cTargetFormats)

	array := (*[1 << 30]C.int)(unsafe.Pointer(cTargetFormats))

	for i := range TargetFormats {
		array[i] = C.int(TargetFormats[i])
	}

	res, info, err := withErrorInfo(func(c *C.FFMS_ErrorInfo) C.int {
		return C.FFMS_SetOutputFormatV2(vs.source, cTargetFormats,
			C.int(width), C.int(height), C.int(resizer), c)
	})

	return int(res), info, err
}

func (vs *VideoSource) ResetOutpuFormat() error {
	if err := vs.checkValidity(); err != nil {
		return err
	}

	C.FFMS_ResetOutputFormatV(vs.source)
	return nil
}

func (vs *VideoSource) SetInputFormat(colorSpace int, colorRange ColorRange,
	format int) (int, *ErrorInfo, error) {
	if err := vs.checkValidity(); err != nil {
		return 0, nil, nil
	}

	res, info, err := withErrorInfo(func(c *C.FFMS_ErrorInfo) C.int {
		return C.FFMS_SetInputFormatV(vs.source, C.int(colorSpace),
			C.int(colorRange), C.int(format), c)
	})

	return int(res), info, err
}

func (vs *VideoSource) ResetInputFormat() error {
	if err := vs.checkValidity(); err != nil {
		return err
	}

	C.FFMS_ResetInputFormatV(vs.source)
	return nil
}

// checkValidity simply checks if the c ptr to the wrapped *C.FFMS_VideoSource
// is nil or not. Any other checks that need to be preformed before the type
// can be used should be added here.
func (vs VideoSource) checkValidity() error {
	if vs.source == nil {
		return ErrInvalidOrNilVideoSource
	}

	return nil
}

// Destroys the VideoSource object if it still exists. Invalidates any further
// usage of the VideoSource.
//
// Note: This must be called to avoid memory leaks as the VideoSource exists
// within C allocated memory. Therefore it will not be automatically cleaned up
// by GO! once the object leaves scope. (Nor does GO! ever guarentee any
// finalizer will ever be called).
func (vs *VideoSource) Close() error {
	if err := vs.checkValidity(); err != nil {
		return err
	}

	C.FFMS_DestroyVideoSource(vs.source)
	vs.source = nil

	return nil

}

type VideoProperties struct {
	// The nominal framerate of the track, as a rational number. For Matroska
	// files, this number is based on the average frame duration of all frames,
	// while for everything else it's based on the duration of the first frame.
	// While it might seem tempting to use these values to extrapolate
	// wallclock timestamps for each frame, you really shouldn't do that since
	// it makes your code unable to handle variable framerate properly. The
	// ugly reality is that these values are pretty much only useful for
	// informational purposes; they are only somewhat reliable for antiquated
	// containers like AVI. Normally they should never be used for practical
	// purposes; generate individual frame timestamps from FrameInfo->PTS
	// instead.
	FPSDenominator int
	// The nominal framerate of the track, as a rational number. For Matroska
	// files, this number is based on the average frame duration of all frames,
	// while for everything else it's based on the duration of the first frame.
	// While it might seem tempting to use these values to extrapolate
	// wallclock timestamps for each frame, you really shouldn't do that since
	// it makes your code unable to handle variable framerate properly. The
	// ugly reality is that these values are pretty much only useful for
	// informational purposes; they are only somewhat reliable for antiquated
	// containers like AVI. Normally they should never be used for practical
	// purposes; generate individual frame timestamps from FrameInfo->PTS
	// instead.
	FPSNumerator int
	// The special RFF timebase, as a rational number. See RepeatPict in the
	// Frame documentation for more information.
	RFFDenominator int
	// The special RFF timebase, as a rational number. See RepeatPict in the
	// Frame documentation for more information.
	RFFNumerator int
	// The number of frames in the video track.
	NumFrames int
	// The sample aspect ratio of the video frames, as a rational number where
	// SARNum is the numerator and SARDen is the denominator. Note that this is
	// a metadata setting that you are free to ignore, but if you want the
	// proper display aspect ratio with anamorphic material, you should honor
	// it. On the other hand, there are situations (like when encoding) where
	// you should probably ignore it because the user expects it to be ignored.
	SARNum int
	// The sample aspect ratio of the video frames, as a rational number where
	// SARNum is the numerator and SARDen is the denominator. Note that this is
	// a metadata setting that you are free to ignore, but if you want the
	// proper display aspect ratio with anamorphic material, you should honor
	// it. On the other hand, there are situations (like when encoding) where
	// you should probably ignore it because the user expects it to be ignored.
	SARDen int
	// The number of pixels in each direction you should crop the frame before
	// displaying it. Note that like the SAR, this is a metadata setting and
	// you are free to ignore it, but if you want things to display 100%
	// correctly you should honor it.
	CropTop int
	// The number of pixels in each direction you should crop the frame before
	// displaying it. Note that like the SAR, this is a metadata setting and
	// you are free to ignore it, but if you want things to display 100%
	// correctly you should honor it.
	CropBottom int
	// The number of pixels in each direction you should crop the frame before
	// displaying it. Note that like the SAR, this is a metadata setting and
	// you are free to ignore it, but if you want things to display 100%
	// correctly you should honor it.
	CropLeft int
	// The number of pixels in each direction you should crop the frame before
	// displaying it. Note that like the SAR, this is a metadata setting and
	// you are free to ignore it, but if you want things to display 100%
	// correctly you should honor it.
	CropRight int
	// Nonzero if the stream has the top field first, zero if it has the bottom
	// field first.
	TopFieldFirst int
	// Identifies the YUV color coefficients used in the stream. Same as in the
	// MPEG-2 specs; see the ColorSpaces enum. The ColorSpace property in Frame
	// should be instead of this, as this can vary between frames unless you've
	// asked for everything to be converted to a single value with
	// SetOutputFormatV2.
	//
	// This value is Deprecated and should never be used. Exists for backwards
	// compatibility.
	ColorSpace int
	// Identifies the luma range of the stream. See the ColorRanges enum. The
	// ColorRange property in Frame should be instead of this, as this can vary
	// between frames unless you've asked for everything to be converted to a
	// single value with SetOutputFormatV2.
	//
	// This value is Deprecated and should never be used. Exists for backwards
	// compatibility.
	ColorRange int
	// The first and last timestamp of the stream respectively, in seconds.
	// Useful if you want to know if the stream has a delay, or for quickly
	// determining its length in seconds.
	FirstTime float64
	// The first and last timestamp of the stream respectively, in seconds.
	// Useful if you want to know if the stream has a delay, or for quickly
	// determining its length in seconds.
	LastTime float64
	// Rotation specifies how many degrees clockwise the decoded video frame
	// should be rotated to match its intended display orientation. Many
	// cameras and mobile devices store orientation separately from the raw
	// pixel data, so decoders rely on this field to present the image upright.
	// Typical values are multiples of ninety degrees, although arbitrary
	// values are possible depending on the container or encoder.
	Rotation int
	// The type of stereo 3D the video is. Corresponts to entries in
	// Stereo3DType. This information allows downstream software to apply
	// correct 3D reconstruction or to decide whether the content is compatible
	// with a given display pipeline.
	Stereo3DType int
	// Stereo 3D flags. Corresponds to entries in Stereo3DFlags.
	Stereo3DFlags int
	// LastEndTime records the timestamp, in seconds, of the end of the final
	// packet of this stream. This is effectively the stream’s duration as
	// detected by the demuxer. The value is useful for seeking, progress
	// reporting, and for tools that need an authoritative end-of-media marker
	// when other duration metadata is unreliable or missing.
	LastEndTime float64
	// HasMasteringDisplayPrimaries is non-zero when the mastering display’s
	// chromaticity information is present. When set, the following fields
	// (MasteringDisplayPrimariesX, MasteringDisplayPrimariesY,
	// MasteringDisplayWhitePointX, and MasteringDisplayWhitePointY) contain
	// valid HDR metadata describing the color characteristics of the display
	// used during the video’s grading process.
	HasMasteringDisplayPrimaries int
	// MasteringDisplayPrimariesX provides the x-coordinate chromaticity values
	// for the red, green, and blue primaries of the mastering display. These
	// are normalized CIE 1931 coordinates used in HDR metadata. They define
	// the color gamut the content was authored for, which is essential for
	// accurate tone and gamut mapping.
	MasteringDisplayPrimariesX [3]float64
	// MasteringDisplayPrimariesY provides the y-coordinate chromaticity values
	// for the red, green, and blue primaries of the mastering display. As with
	// the X coordinates, these values describe the display’s color gamut in
	// the CIE 1931 chromaticity space and help renderers determine how to
	// interpret HDR color information.
	MasteringDisplayPrimariesY [3]float64
	// MasteringDisplayWhitePointX is the x-coordinate of the display’s white
	// point in CIE 1931 space. This defines the “reference white” the colorist
	// saw when authoring the content and is critical for accurate reproduction
	// of intent, particularly in HDR workflows.
	MasteringDisplayWhitePointX float64
	// MasteringDisplayWhitePointY is the y-coordinate of the mastering
	// display’s white point. Used in conjunction with
	// MasteringDisplayWhitePointX, it describes the chromaticity of the white
	// reference used during HDR grading.
	MasteringDisplayWhitePointY float64
	// HasMasteringDisplayLuminance is non-zero when the mastering display’s
	// luminance range metadata is present. When set, the values in
	// MasteringDisplayMinLuminance and MasteringDisplayMaxLuminance describe
	// the minimum and maximum luminance capabilities of the original grading
	// display, expressed in cd/m² (nits).
	HasMasteringDisplayLuminance int
	// MasteringDisplayMinLuminance specifies the darkest luminance level, in
	// cd/m², that the mastering display can reproduce. HDR content often
	// relies on this value when performing tone mapping so that shadow detail
	// is preserved relative to the colorist’s intent.
	MasteringDisplayMinLuminance float64
	// MasteringDisplayMaxLuminance specifies the peak luminance level, in
	// cd/m², of the mastering display. This value indicates how bright
	// highlights were allowed to be during grading. Tone-mapping and HDR
	// renderers use this field to determine how to map very bright regions
	// into the capabilities of a target display.
	MasteringDisplayMaxLuminance float64
	// HasContentLightLevel is non-zero when the content light level metadata
	// is present. When set, ContentLightLevelMax and ContentLightLevelAverage
	// provide additional HDR display hints describing the brightness
	// characteristics of the video frames themselves rather than the mastering
	// display.
	HasContentLightLevel int
	// ContentLightLevelMax is the maximum content light level (MaxCLL),
	// expressed in cd/m², observed across the entire video. It indicates the
	// brightest pixel recorded in any frame and is often used by HDR
	// tone-mapping systems to avoid clipping excessively intense highlights.
	ContentLightLevelMax uint32
	// ContentLightLevelAverage is the maximum frame-average light level
	// (MaxFALL), again expressed in cd/m². It represents the brightest average
	// luminance of any single frame. This helps HDR systems decide how
	// aggressively to adjust brightness without creating sudden jumps or
	// violating display limits.
	ContentLightLevelAverage uint32
	// Flip direction to be applied before rotation: 0 for no operation, >0 for
	// horizontal flip, <0 for vertical flip.
	Flip int
}

func (props *VideoProperties) videoPropertiesFromC(cProps *C.FFMS_VideoProperties) {
	props.FPSDenominator = int(cProps.FPSDenominator)
	props.FPSNumerator = int(cProps.FPSNumerator)
	props.RFFDenominator = int(cProps.RFFDenominator)
	props.RFFNumerator = int(cProps.RFFNumerator)
	props.NumFrames = int(cProps.NumFrames)
	props.SARNum = int(cProps.SARNum)
	props.SARDen = int(cProps.SARDen)
	props.CropTop = int(cProps.CropTop)
	props.CropBottom = int(cProps.CropBottom)
	props.CropLeft = int(cProps.CropLeft)
	props.CropRight = int(cProps.CropRight)
	props.TopFieldFirst = int(cProps.TopFieldFirst)
	props.ColorSpace = int(cProps.ColorSpace)
	props.ColorRange = int(cProps.ColorRange)
	props.FirstTime = float64(cProps.FirstTime)
	props.LastTime = float64(cProps.LastTime)
	props.Rotation = int(cProps.Rotation)
	props.Stereo3DType = int(cProps.Stereo3DType)
	props.Stereo3DFlags = int(cProps.Stereo3DFlags)
	props.LastEndTime = float64(cProps.LastEndTime)
	props.HasMasteringDisplayPrimaries = int(cProps.HasMasteringDisplayPrimaries)
	for i := 0; i < 3; i++ {
		props.MasteringDisplayPrimariesX[i] = float64(cProps.MasteringDisplayPrimariesX[i])
		props.MasteringDisplayPrimariesY[i] = float64(cProps.MasteringDisplayPrimariesY[i])
	}
	props.MasteringDisplayWhitePointX = float64(cProps.MasteringDisplayWhitePointX)
	props.MasteringDisplayWhitePointY = float64(cProps.MasteringDisplayWhitePointY)
	props.HasMasteringDisplayLuminance = int(cProps.HasMasteringDisplayLuminance)
	props.MasteringDisplayMinLuminance = float64(cProps.MasteringDisplayMinLuminance)
	props.MasteringDisplayMaxLuminance = float64(cProps.MasteringDisplayMaxLuminance)
	props.HasContentLightLevel = int(cProps.HasContentLightLevel)
	props.ContentLightLevelMax = uint32(cProps.ContentLightLevelMax)
	props.ContentLightLevelAverage = uint32(cProps.ContentLightLevelAverage)
	props.Flip = int(cProps.Flip)
}

// A struct representing a video frame.
type Frame struct {
	// An array of slices to the picture planes (fields containing actual
	// pixel data). Planar formats use more than one plane, for example YV12
	// uses one plane each for the Y, U and V data. Packed formats (such as the
	// various RGB32 flavors) use only the first plane. If you want to
	// determine if plane i contains data or not, check for Frame.Linesize[i]
	// != 0 or check the length of the current data slice.
	Data [4][]uint8
	// An array of integers representing the length of each scan line in each
	// of the four picture planes, in bytes. In alternative terminology, this
	// is the "pitch" of the plane. Usually, the total size in bytes of picture
	// plane i is Frame.Linesize[i] * Frame.EncodedHeight, but do note that
	// some pixel formats (most notably YV12) has vertical chroma subsampling,
	// and then the U/V planes may be of a different height than the primary
	// plane. This may be negative; if so the image is stored inverted in
	// memory and Data actually points of the last row of the data. You usually
	// do not need to worry about this, as it mostly works correctly by default
	// if you're processing the image correctly.
	Linesize [4]int
	// The original resolution of the frame (in pixels), as encoded in the
	// compressed file, before any scaling was applied. Note that must not
	// necessarily be the same for all frames in a stream.
	EncodedWidth int
	// The original resolution of the frame (in pixels), as encoded in the
	// compressed file, before any scaling was applied. Note that must not
	// necessarily be the same for all frames in a stream.
	EncodedHeight int
	// The original pixel format of the frame, as encoded in the compressed
	// file.
	EncodedPixelFormat int
	// The output resolution of the frame (in pixels), i.e. the resolution of
	// what is actually stored in the Data field. Set to -1 if no scaling is
	// done.
	ScaledWidth int
	// The output resolution of the frame (in pixels), i.e. the resolution of
	// what is actually stored in the Data field. Set to -1 if no scaling is
	// done.
	ScaledHeight int
	// The output pixel format of the frame, i.e. the pixel format of what is
	// actually stored in the Data field.
	ConvertedPixelFormat int
	// Nonzero if the frame is a keyframe, 0 otherwise.
	KeyFrame int
	// An integer repesenting the RFF flag for this frame; i.e. the frame shall
	// be displayed for 1+RepeatPict time units, where the time units are
	// expressed in the special RFF timebase available in
	// VideoProperties.RFFDenominator and VideoProperties.RFFNumerator.
	//
	// Note: if you actually end up using this, you need to ignore the usual
	// timestamps (calculated via the TrackTimeBase and the frame PTS) since
	// they are fundamentally incompatible with RFF flags.
	RepeatPict int
	// Nonzero if the frame was coded as interlaced, zero otherwise
	InterlacedFrame int
	// Nonzero if the frame has the top field first, zero if it has the bottom
	// field first. Only relevant if InterlacedFrame is nonzero.
	TopFieldFirst int
	// A single character denoting coding type (I/B/P etc) of the compressed
	// frame. See the Constants and Preprocessor Definitions section for more
	// information about what the different letters mean.
	PictType byte
	// Identifies the YUV color coefficients used in the stream. Same as in the
	// MPEG-2 specs; see the ColorSpaces enum. The ColorSpace property in Frame
	// should be instead of this, as this can vary between frames unless you've
	// asked for everything to be converted to a single value with
	// SetOutputFormatV2.
	ColorSpace int
	// Identifies the luma range of the stream. See the ColorRanges enum. The
	// ColorRange property in Frame should be instead of this, as this can vary
	// between frames unless you've asked for everything to be converted to a
	// single value with SetOutputFormatV2.
	ColorRange int
	// Identifies the color primaries of the frame.
	ColorPrimaries int
	// Identifies the transfer characteristics of the frame.
	TransferCharateristics int
	// Identifies the chroma location for the frame. Corresponds to
	// ChromaLocations.
	ChromaLocation int
	// HasMasteringDisplayPrimaries is non-zero when the mastering display’s
	// chromaticity information is present. When set, the following fields
	// (MasteringDisplayPrimariesX, MasteringDisplayPrimariesY,
	// MasteringDisplayWhitePointX, and MasteringDisplayWhitePointY) contain
	// valid HDR metadata describing the color characteristics of the display
	// used during the video’s grading process.
	HasMasteringDisplayPrimaries int
	// MasteringDisplayPrimariesX provides the x-coordinate chromaticity values
	// for the red, green, and blue primaries of the mastering display. These
	// are normalized CIE 1931 coordinates used in HDR metadata. They define
	// the color gamut the content was authored for, which is essential for
	// accurate tone and gamut mapping.
	MasteringDisplayPrimariesX [3]float64
	// MasteringDisplayPrimariesY provides the y-coordinate chromaticity values
	// for the red, green, and blue primaries of the mastering display. As with
	// the X coordinates, these values describe the display’s color gamut in
	// the CIE 1931 chromaticity space and help renderers determine how to
	// interpret HDR color information.
	MasteringDisplayPrimariesY [3]float64
	// MasteringDisplayWhitePointX is the x-coordinate of the display’s white
	// point in CIE 1931 space. This defines the “reference white” the colorist
	// saw when authoring the content and is critical for accurate reproduction
	// of intent, particularly in HDR workflows.
	MasteringDisplayWhitePointX float64
	// MasteringDisplayWhitePointY is the y-coordinate of the mastering
	// display’s white point. Used in conjunction with
	// MasteringDisplayWhitePointX, it describes the chromaticity of the white
	// reference used during HDR grading.
	MasteringDisplayWhitePointY float64
	// HasMasteringDisplayLuminance is non-zero when the mastering display’s
	// luminance range metadata is present. When set, the values in
	// MasteringDisplayMinLuminance and MasteringDisplayMaxLuminance describe
	// the minimum and maximum luminance capabilities of the original grading
	// display, expressed in cd/m² (nits).
	HasMasteringDisplayLuminance int
	// MasteringDisplayMinLuminance specifies the darkest luminance level, in
	// cd/m², that the mastering display can reproduce. HDR content often
	// relies on this value when performing tone mapping so that shadow detail
	// is preserved relative to the colorist’s intent.
	MasteringDisplayMinLuminance float64
	// MasteringDisplayMaxLuminance specifies the peak luminance level, in
	// cd/m², of the mastering display. This value indicates how bright
	// highlights were allowed to be during grading. Tone-mapping and HDR
	// renderers use this field to determine how to map very bright regions
	// into the capabilities of a target display.
	MasteringDisplayMaxLuminance float64
	// HasContentLightLevel is non-zero when the content light level metadata
	// is present. When set, ContentLightLevelMax and ContentLightLevelAverage
	// provide additional HDR display hints describing the brightness
	// characteristics of the video frames themselves rather than the mastering
	// display.
	HasContentLightLevel int
	// ContentLightLevelMax is the maximum content light level (MaxCLL),
	// expressed in cd/m², observed across the entire video. It indicates the
	// brightest pixel recorded in any frame and is often used by HDR
	// tone-mapping systems to avoid clipping excessively intense highlights.
	ContentLightLevelMax uint32
	// ContentLightLevelAverage is the maximum frame-average light level
	// (MaxFALL), again expressed in cd/m². It represents the brightest average
	// luminance of any single frame. This helps HDR systems decide how
	// aggressively to adjust brightness without creating sudden jumps or
	// violating display limits.
	ContentLightLevelAverage uint32
	DolbyVisionRPU           []byte
	HDR10Plus                []byte
}

func (*Frame) getSizePerPlane(cFrame *C.FFMS_Frame) ([]uint, error) {
	desc, err := pixfmts.PixFmtDescGet(pixfmts.PixelFormat(
		cFrame.ConvertedPixelFormat))
	if err != nil {
		return nil, err
	}

	var bHor, bver uint
	if cFrame.ScaledHeight == -1 {
		bHor, bver = uint(cFrame.EncodedWidth), uint(cFrame.EncodedHeight)
	} else {
		bHor, bver = uint(cFrame.ScaledWidth), uint(cFrame.ScaledHeight)
	}

	var res []uint

	for i := range desc.NbComponents() {
		comp, err := desc.Component(i)
		if err != nil {
			return nil, err
		}

		var horSub, verSub int = 1, 1
		if i > 0 {
			horSub, verSub = 1<<desc.Log2ChromaW(), 1<<desc.Log2ChromaH()
		}

		res = append(res, (bHor/uint(horSub))*(bver/uint(verSub))*
			uint(comp.Step))

	}

	return res, nil
}

func (frame *Frame) fromCFrame(cFrame *C.FFMS_Frame) error {
	if cFrame == nil {
		return nil
	}

	sizes, err := frame.getSizePerPlane(cFrame)
	if err != nil {
		return err
	}

	for i, size := range sizes {
		frame.Data[i] = sliceFromCPtr[C.uint8_t, uint8](cFrame.Data[i], size)
	}

	for i := 0; i < 4; i++ {
		frame.Linesize[i] = int(cFrame.Linesize[i])
	}

	for i := 0; i < 3; i++ {
		frame.MasteringDisplayPrimariesX[i] = float64(
			cFrame.MasteringDisplayPrimariesX[i])
		frame.MasteringDisplayPrimariesY[i] = float64(
			cFrame.MasteringDisplayPrimariesY[i])
	}

	frame.DolbyVisionRPU = sliceFromCPtr[C.uint8_t, byte](
		cFrame.DolbyVisionRPU, uint(cFrame.DolbyVisionRPUSize))
	frame.HDR10Plus = sliceFromCPtr[C.uint8_t, byte](
		cFrame.HDR10Plus, uint(cFrame.HDR10PlusSize))

	var (
		encodedWidth             = int(cFrame.EncodedWidth)
		encodedHeight            = int(cFrame.EncodedHeight)
		encodedPixelFormat       = int(cFrame.EncodedPixelFormat)
		scaledWidth              = int(cFrame.ScaledWidth)
		scaledheight             = int(cFrame.ScaledHeight)
		convertedPixelFormat     = int(cFrame.ConvertedPixelFormat)
		keyframe                 = int(cFrame.KeyFrame)
		repeatPic                = int(cFrame.RepeatPict)
		interlacedFrame          = int(cFrame.InterlacedFrame)
		topFeildFirst            = int(cFrame.TopFieldFirst)
		pictType                 = byte(cFrame.PictType)
		colorSpace               = int(cFrame.ColorSpace)
		colorRange               = int(cFrame.ColorRange)
		colorPrimaries           = int(cFrame.ColorPrimaries)
		transferCharateristics   = int(cFrame.TransferCharateristics)
		chromaLocation           = int(cFrame.ChromaLocation)
		hasDisplayPrimaries      = int(cFrame.HasMasteringDisplayPrimaries)
		masteringDisplayWhiteX   = float64(cFrame.MasteringDisplayWhitePointX)
		masteringDisplayWhiteY   = float64(cFrame.MasteringDisplayWhitePointY)
		hasMasteringDisplayLuma  = int(cFrame.HasMasteringDisplayLuminance)
		masteringDisplayMinLuma  = float64(cFrame.MasteringDisplayMinLuminance)
		masteringDisplayMaxLuma  = float64(cFrame.MasteringDisplayMaxLuminance)
		hasContentLightLevel     = int(cFrame.HasContentLightLevel)
		contentLightLevelMax     = uint32(cFrame.ContentLightLevelMax)
		contentLightLevelAverage = uint32(cFrame.ContentLightLevelAverage)
	)

	frame.EncodedWidth, frame.EncodedHeight = encodedWidth, encodedHeight
	frame.EncodedPixelFormat = encodedPixelFormat
	frame.ScaledWidth, frame.ScaledHeight = scaledWidth, scaledheight
	frame.ConvertedPixelFormat = convertedPixelFormat
	frame.KeyFrame, frame.RepeatPict = keyframe, repeatPic
	frame.InterlacedFrame, frame.TopFieldFirst = interlacedFrame, topFeildFirst
	frame.PictType = pictType
	frame.ColorSpace, frame.ColorRange = colorSpace, colorRange
	frame.ColorPrimaries = colorPrimaries
	frame.TransferCharateristics = transferCharateristics
	frame.ChromaLocation = chromaLocation
	frame.HasMasteringDisplayPrimaries = hasDisplayPrimaries
	frame.MasteringDisplayWhitePointX = masteringDisplayWhiteX
	frame.MasteringDisplayWhitePointY = masteringDisplayWhiteY
	frame.HasMasteringDisplayLuminance = hasMasteringDisplayLuma
	frame.MasteringDisplayMinLuminance = masteringDisplayMinLuma
	frame.MasteringDisplayMaxLuminance = masteringDisplayMaxLuma
	frame.HasContentLightLevel = hasContentLightLevel
	frame.ContentLightLevelMax = contentLightLevelMax
	frame.ContentLightLevelAverage = contentLightLevelAverage

	return nil
}
