package sources

import (
	"runtime"

	pixfmts "github.com/GreatValueCreamSoda/gometrics/c/libavpixfmts"
	ffms "github.com/GreatValueCreamSoda/gometrics/c/libffms2"
	"github.com/GreatValueCreamSoda/gometrics/video"
)

type ffmsSource struct {
	currentIndex int
	video        *ffms.VideoSource
	numFrame     int
	colorspace   video.ColorProperties
	planeSizes   [3]int
	planeStrides [3]int
	frameRate    float32
}

func NewFFms2Reader(path string) (video.Source, error) {
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

	var decThreads int = runtime.NumCPU()
	source, _, err := ffms.CreateVideoSource(path, index, track, decThreads,
		ffms.SeekNormal)
	if err != nil {
		return nil, err
	}

	props, err := source.GetVideoProperties()
	if err != nil {
		return nil, err
	}

	ff, _, err := source.GetFrame(0)
	if err != nil {
		return nil, err
	}

	// Causes ffms2 to randomly segfault. Need to figure out why.

	// video.SetOutputFormatV2([]int{ff.EncodedPixelFormat}, ff.EncodedWidth,
	//	ff.EncodedHeight, ffms.ResizerBicubic)

	// ff, _, err = video.GetFrame(0)
	// if err != nil {
	// 	return nil, err
	// }

	var planeSizes, planeStrides [3]int

	for i := range 3 {
		planeSizes[i] = len(ff.Data[i])
		planeStrides[i] = ff.Linesize[i]
	}

	colorProps := video.ColorProperties{
		Width:          ff.EncodedWidth,
		Height:         ff.EncodedHeight,
		PixelFormat:    pixfmts.PixelFormat(ff.EncodedPixelFormat),
		ColorRange:     pixfmts.ColorRange(ff.ColorRange),
		ColorSpace:     pixfmts.ColorSpace(ff.ColorSpace),
		ColorTransfer:  pixfmts.ColorTransferCharacteristic(ff.TransferCharateristics),
		ColorPrimaries: pixfmts.ColorPrimaries(ff.ColorPrimaries),
		ChromaLocation: pixfmts.ChromaLocation(ff.ChromaLocation),
	}

	return &ffmsSource{0, source, props.NumFrames, colorProps, planeSizes,
		planeStrides, float32(props.FPSNumerator) / float32(props.FPSDenominator)}, nil
}

func (s *ffmsSource) GetFrame(frame *video.Frame) error {
	ffmsFrame, _, err := s.video.GetFrame(s.currentIndex)
	if err != nil {
		return err
	}

	frame.Data = [3][]byte{
		ffmsFrame.Data[0], ffmsFrame.Data[1], ffmsFrame.Data[2]}
	frame.LineSize = [3]int64{int64(ffmsFrame.Linesize[0]), int64(ffmsFrame.Linesize[1]),
		int64(ffmsFrame.Linesize[2])}

	s.currentIndex++
	return nil
}

func (s *ffmsSource) GetColorProps() *video.ColorProperties { return &s.colorspace }
func (s *ffmsSource) GetNumFrames() int                     { return s.numFrame }
func (s *ffmsSource) GetFrameRate() float32                 { return s.frameRate }

func (c *ffmsSource) GetPlaneSizes() ([3]int, [3]int) {
	return c.planeSizes, c.planeStrides
}
