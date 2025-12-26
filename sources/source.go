package sources

import (
	"runtime"

	ffms "github.com/GreatValueCreamSoda/goffms2"
	"github.com/GreatValueCreamSoda/gometrics/comparator"
	"github.com/GreatValueCreamSoda/gopixfmts"
	vship "github.com/GreatValueCreamSoda/govship"
)

type ffmsSource struct {
	currentIndex int
	video        *ffms.VideoSource
	numFrame     int
	colorspace   vship.Colorspace
	planeSizes   [3]int
	planeStrides [3]int
}

func NewFFms2Reader(path string) (comparator.Source, error) {
	var err error

	var indexer *ffms.Indexer
	if indexer, _, err = ffms.CreateIndexer(path); err != nil {
		return nil, err
	}

	var index *ffms.Index
	if index, _, err = indexer.DoIndexing(ffms.IEHAbort); err != nil {
		return nil, err
	}

	track, _, err := index.GetFirstTrackOfType(ffms.TypeVideo)
	if err != nil {
		return nil, err
	}

	var decThreads int = runtime.NumCPU() / 2
	video, _, err := ffms.CreateVideoSource(path, index, track, decThreads,
		ffms.SeekNormal)
	if err != nil {
		return nil, err
	}

	props, err := video.GetVideoProperties()
	if err != nil {
		return nil, err
	}

	ff, _, err := video.GetFrame(0)
	if err != nil {
		return nil, err
	}

	video.SetOutputFormatV2([]int{ff.EncodedPixelFormat}, ff.EncodedWidth,
		ff.EncodedHeight, ffms.ResizerBicubic)

	ff, _, err = video.GetFrame(0)
	if err != nil {
		return nil, err
	}

	var planeSizes, planeStrides [3]int

	for i := range 3 {
		planeSizes[i] = len(ff.Data[i])
		planeStrides[i] = ff.Linesize[i]
	}

	colorspace, err := convertFfmsFrameToVshipColorspace(&ff)

	return &ffmsSource{0, video, props.NumFrames, colorspace, planeSizes,
		planeStrides}, nil
}

func (s *ffmsSource) GetFrame(frame *comparator.Frame) error {
	ffmsFrame, _, err := s.video.GetFrame(s.currentIndex)
	if err != nil {
		return err
	}

	frame.Write(
		[3][]byte{ffmsFrame.Data[0], ffmsFrame.Data[1], ffmsFrame.Data[2]},
		[3]int64{int64(ffmsFrame.Linesize[0]), int64(ffmsFrame.Linesize[1]),
			int64(ffmsFrame.Linesize[2])})

	s.currentIndex++
	return nil
}

func (s *ffmsSource) GetColorspace() *vship.Colorspace { return &s.colorspace }
func (s *ffmsSource) GetNumFrames() int                { return s.numFrame }

func (c *ffmsSource) GetPlaneSizes() ([3]int, [3]int) {
	return c.planeSizes, c.planeStrides
}

// convertFfmsFrameToVshipColorspace extracts colorspace information from an
// FFMS2 frame and converts it into the equivalent vship.Colorspace
// representation used by this package heavily.
//
// This function inspects the frame's pixel format (after any conversion
// applied via SetOutputFormat), color range, matrix coefficients, transfer
// characteristics, primaries, chroma subsampling, and chroma location to
// populate a vship.Colorspace struct accurately.
//
// Default assumptions (common in video encoding):
//   - If ColorSpace is unset (0), assumes BT.709 for YUV and unspecified (0)
//     for RGB.
//   - If TransferCharacteristics or ColorPrimaries are unset (0), defaults to
//     BT.709.
//   - If ChromaLocation is unset (0), defaults to Left (value 1, MPEG2/4
//     style).
//
// The function panics on unsupported bit depths (outside 8–16 bits). Currently
// supported depths: 8, 9, 10, 12, 14, 16 bits per component.
//
// Cropping fields (CropTop, CropBottom, CropLeft, CropRight) are always set to
// 0.
//
// Returns the populated Colorspace and any error encountered while querying
// the pixel format descriptor.
func convertFfmsFrameToVshipColorspace(frame *ffms.Frame) (vship.Colorspace,
	error) {

	var colorspace vship.Colorspace

	colorspace.Width = int64(frame.ScaledWidth)
	colorspace.TargetWidth = colorspace.Width
	colorspace.Height = int64(frame.ScaledHeight)
	colorspace.TargetHeight = colorspace.Height

	videoPixelFormat, err := gopixfmts.PixFmtDescGet(gopixfmts.PixelFormat(
		frame.ConvertedPixelFormat))
	if err != nil {
		return colorspace, err
	}

	comp, err := videoPixelFormat.Component(0)
	if err != nil {
		return colorspace, err
	}

	var videoDepth vship.SamplingFormat

	switch comp.Depth {
	case 8:
		videoDepth = vship.SamplingFormatUInt8
	case 9:
		videoDepth = vship.SamplingFormatUInt9
	case 10:
		videoDepth = vship.SamplingFormatUInt10
	case 12:
		videoDepth = vship.SamplingFormatUInt12
	case 14:
		videoDepth = vship.SamplingFormatUInt14
	case 16:
		videoDepth = vship.SamplingFormatUInt16
	default:
		panic("UNKNOWN PIXEL FORMAT")
	}

	colorspace.SamplingFormat = videoDepth

	if frame.ColorRange == int(gopixfmts.ColorRangeMPEG) ||
		frame.ColorRange == 0 {
		colorspace.ColorRange = vship.ColorRangeLimited
	} else {
		colorspace.ColorRange = vship.ColorRangeFull
	}

	colorspace.ChromaSubsamplingHeight = videoPixelFormat.Log2ChromaH()
	colorspace.ChromaSubsamplingWidth = videoPixelFormat.Log2ChromaW()

	if frame.ChromaLocation > 0 {
		colorspace.ChromaLocation = vship.ChromaLocation(frame.ChromaLocation)
	} else {
		colorspace.ChromaLocation = 1 // Default assumption
	}

	if videoPixelFormat.Flags()&uint64(gopixfmts.PixFmtFlagRGB) == 0 {
		colorspace.ColorFamily = vship.ColorFamilyYUV
	} else {
		colorspace.ColorFamily = vship.ColorFamilyRGB
	}

	if frame.ColorSpace > 0 {
		colorspace.ColorMatrix = vship.ColorMatrix(frame.ColorSpace)
	} else {
		if colorspace.ColorFamily == vship.ColorFamilyYUV {
			colorspace.ColorMatrix = 1 // BT.709 assumed
		} else {
			colorspace.ColorMatrix = 0
		}
	}

	if frame.TransferCharateristics > 0 {
		colorspace.ColorTransfer = vship.ColorTransfer(
			frame.TransferCharateristics)
	} else {
		colorspace.ColorTransfer = 1 // BT.709
	}

	if frame.ColorPrimaries > 0 {
		colorspace.ColorPrimaries = vship.ColorPrimaries(
			frame.ColorPrimaries)
	} else {
		colorspace.ColorPrimaries = 1 // BT.709
	}

	colorspace.CropTop, colorspace.CropBottom, colorspace.CropLeft = 0, 0, 0
	colorspace.CropRight = 0

	return colorspace, nil
}
