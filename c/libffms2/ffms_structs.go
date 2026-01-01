package libffms2

//#cgo LDFLAGS: -lffms2
//#cgo CFLAGS: -I/usr/include
//#include <ffms.h>
import "C"

// A struct representing a Audio source that can be read from and have it's
// properties listed.
type AudioSource struct {
	source *C.FFMS_AudioSource
}

// A struct representing a FFMS track from an Index, Audio, or Video source.
type Track struct {
	track *C.FFMS_Track
}

// A struct representing a FFMS ErrorInfo made safe for usage in GO!.
type ErrorInfo struct {
	ErrorType int
	SubType   int
	Message   string
}

type AudioProperties struct {
	// An integer that represents the audio sample format. See SampleFormat.
	SampleFormat int
	// The audio samplerate, in samples per second.
	SampleRate int
	// The number of bits per audio sample. Note that this signifies the number
	// of bits actually used to code each sample, not the number of bits used
	// to store each sample, and may hence be different from what the
	// SampleFormat would imply. Figuring out which bytes are significant and
	// which aren't is left as an exercise for the reader.
	BitsPerSample int
	// The number of audio channels.
	Channels int
	// The channel layout of the audio stream. Constructed by binary OR'ing the
	// relevant integers from AudioChannel together, which means that if the
	// audio has the channel CH_EXAMPLE, the operation (ChannelOrder &
	// CH_EXAMPLE) will evaluate to true. The samples are interleaved in the
	// order the channels are listed in the AudioChannel enum.
	ChannelLayout int64
	// The number of samples in the audio track.
	NumSamples int64
	// The first and last timestamp of the stream respectively, in
	// milliseconds. Useful if you want to know if the stream has a delay, or
	// for quickly determining its length in seconds.
	FirstTime float64
	// The first and last timestamp of the stream respectively, in
	// milliseconds. Useful if you want to know if the stream has a delay, or
	// for quickly determining its length in seconds.
	LastTime float64
	// The end time of the last packet of the stream, in milliseconds.
	LastEndTime float64
}

// A struct representing the basic time unit of a track, as a rational number
// where Num is the numerator and Den is the denominator. Note that while this
// rational number may occasionally turn out to be equal to 1/framerate for
// some CFR video tracks, it really has no relation whatsoever with the video
// framerate and you should definitely not assume anything framerate-related
// based on it.
type TrackTimeBase struct {
	Num int64
	Den int64
}

type FrameInfo struct {
	// The decoding timestamp of the frame. To convert this to a timestamp in
	// wallclock milliseconds, use the relation var timestamp int64 =
	// int64(float64(FrameInfo.PTS * TrackTimeBase.Num) /
	// float64(TrackTimeBase.Den))
	PTS int64
	// RFF flag for the frame; same as in FFMS_Frame, see that structure for an
	// explanation.
	RepeatPict int
	// Non-zero if the frame is a keyframe, zero otherwise.
	KeyFrame    int
	OriginalPTS int64
}

type ResampleOptions struct {
	ChannelLayout          int64
	SampleFormat           SampleFormat
	SampleRate             int
	MixingCoefficientType  MixingCoefficientType
	CenterMixLevel         float64
	SurroundMixLevel       float64
	LFEMixLevel            float64
	Normalize              int
	ForceResample          int
	ResampleFilterSize     int
	ResamplePhaseShift     int
	LinearInterpolation    int
	CutoffFrequencyRatio   float64
	MatrixedStereoEncoding MatrixEncoding
	FilterType             ResampleFilterType
	KaiserBeta             int
	DitherMethod           AudioDitherMethod
}

func (opts *ResampleOptions) toC() C.FFMS_ResampleOptions {
	return C.FFMS_ResampleOptions{
		ChannelLayout:          C.int64_t(opts.ChannelLayout),
		SampleFormat:           C.FFMS_SampleFormat(opts.SampleFormat),
		SampleRate:             C.int(opts.SampleRate),
		MixingCoefficientType:  C.FFMS_MixingCoefficientType(opts.MixingCoefficientType),
		CenterMixLevel:         C.double(opts.CenterMixLevel),
		SurroundMixLevel:       C.double(opts.SurroundMixLevel),
		LFEMixLevel:            C.double(opts.LFEMixLevel),
		Normalize:              C.int(opts.Normalize),
		ForceResample:          C.int(opts.ForceResample),
		ResampleFilterSize:     C.int(opts.ResampleFilterSize),
		ResamplePhaseShift:     C.int(opts.ResamplePhaseShift),
		LinearInterpolation:    C.int(opts.LinearInterpolation),
		CutoffFrequencyRatio:   C.double(opts.CutoffFrequencyRatio),
		MatrixedStereoEncoding: C.FFMS_MatrixEncoding(opts.MatrixedStereoEncoding),
		FilterType:             C.FFMS_ResampleFilterType(opts.FilterType),
		KaiserBeta:             C.int(opts.KaiserBeta),
		DitherMethod:           C.FFMS_AudioDitherMethod(opts.DitherMethod),
	}
}

// ffmsAudioPropertiesFromC converts a C.FFMS_AudioProperties to a Go FFMSAudioProperties
func ffmsAudioPropertiesFromC(cProps *C.FFMS_AudioProperties) AudioProperties {
	return AudioProperties{
		SampleFormat:  int(cProps.SampleFormat),
		SampleRate:    int(cProps.SampleRate),
		BitsPerSample: int(cProps.BitsPerSample),
		Channels:      int(cProps.Channels),
		ChannelLayout: int64(cProps.ChannelLayout),
		NumSamples:    int64(cProps.NumSamples),
		FirstTime:     float64(cProps.FirstTime),
		LastTime:      float64(cProps.LastTime),
		LastEndTime:   float64(cProps.LastEndTime),
	}
}

// ffmsTrackTimeBaseFromC converts a C.FFMS_TrackTimeBase to a Go FFMSTrackTimeBase
func ffmsTrackTimeBaseFromC(cBase *C.FFMS_TrackTimeBase) TrackTimeBase {
	return TrackTimeBase{
		Num: int64(cBase.Num),
		Den: int64(cBase.Den),
	}
}

// ffmsFrameInfoFromC converts a C.FFMS_FrameInfo to a Go FFMSFrameInfo
func ffmsFrameInfoFromC(cInfo *C.FFMS_FrameInfo) FrameInfo {
	return FrameInfo{
		PTS:         int64(cInfo.PTS),
		RepeatPict:  int(cInfo.RepeatPict),
		KeyFrame:    int(cInfo.KeyFrame),
		OriginalPTS: int64(cInfo.OriginalPTS),
	}
}
