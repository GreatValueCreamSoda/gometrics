package libffms2

//#cgo LDFLAGS: -lffms2
//#cgo CFLAGS: -I/usr/include
//#include <ffms.h>
import "C"

type Errors int

const (
	// No Error

	ErrorSuccess Errors = C.FFMS_ERROR_SUCCESS
)

const (
	// Main types - where the error occurred

	// index file handling
	ErrorIndex          Errors = (iota - 1) + C.FFMS_ERROR_INDEX
	ErrorIndexing              // indexing
	ErrorPostProcessing        // video postprocessing (libpostproc)
	ErrorScaling               // image scaling (libswscale)
	ErrorDecoding              // audio/video decoding
	ErrorSeeking               // seeking
	ErrorParser                // file parsing
	ErrorTrack                 // track handling
	ErrorWaveWriter            // WAVE64 file writer
	ErrorCancelled             // operation aborted
	ErrorResampling            // audio resampling (libavresample)
)

const (
	// Subtypes - what caused the error

	// unknown error
	ErrorUnknown          Errors = (iota) + C.FFMS_ERROR_UNKNOWN
	ErrorUnsupported             // format or operation is not supported with this binary
	ErrorFileRead                // cannot read from file
	ErrorFileWrite               // cannot write to file
	ErrorNoFile                  // no such file or directory
	ErrroVersion                 // wrong version
	ErrorAllocationFailed        // out of memory
	ErrorInvalidArgument         // invalid or nonsensical argument
	ErrorCodec                   // decoder error
	ErrorNotAvailable            // requested mode or operation unavailable in this binary
	ErrorFileMismatch            // provided index does not match the file
	ErrorUser                    // problem exists between keyboard and chair
)

type SeekMode int

const (
	SeekLinearNoRw SeekMode = (iota) + C.FFMS_SEEK_LINEAR_NO_RW
	SeekLinear
	SeekNormal
	SeekUnsafe
	SeekAggressive
)

type IndexErrorHandling int

const (
	IEHAbort IndexErrorHandling = (iota) + C.FFMS_IEH_ABORT
	IEHClearTrack
	IEHStopTrack
	IEHIgnore
)

type TrackType int

const (
	TypeUnknown TrackType = (iota) + C.FFMS_TYPE_UNKNOWN
	TypeVideo
	TypeAudio
	TypeData
	TypeSubtitle
	TypeAttachment
)

type SampleFormat int

const (
	FmtU8 SampleFormat = (iota) + C.FFMS_FMT_U8
	FmtS16
	FmtS32
	FmtFlt
	FmtDbl
)

type AudioChannel int

const (
	ChannelFrontLeft          AudioChannel = C.FFMS_CH_FRONT_LEFT
	ChannelFrontRight         AudioChannel = C.FFMS_CH_FRONT_RIGHT
	ChannelCenter             AudioChannel = C.FFMS_CH_FRONT_CENTER
	ChannelLowFrequency       AudioChannel = C.FFMS_CH_LOW_FREQUENCY
	ChannelBackLeft           AudioChannel = C.FFMS_CH_BACK_LEFT
	ChannelBackRight          AudioChannel = C.FFMS_CH_BACK_RIGHT
	ChannelFrontLeftOfCenter  AudioChannel = C.FFMS_CH_FRONT_LEFT_OF_CENTER
	ChannelFrontRightOfCenter AudioChannel = C.FFMS_CH_FRONT_RIGHT_OF_CENTER
	ChannelBackCenter         AudioChannel = C.FFMS_CH_BACK_CENTER
	ChannelSideLeft           AudioChannel = C.FFMS_CH_SIDE_LEFT
	ChannelSideRight          AudioChannel = C.FFMS_CH_SIDE_RIGHT
	ChannelTopCenter          AudioChannel = C.FFMS_CH_TOP_CENTER
	ChannelTopFrontLeft       AudioChannel = C.FFMS_CH_TOP_FRONT_LEFT
	ChannelTopFrontCenter     AudioChannel = C.FFMS_CH_TOP_FRONT_CENTER
	ChannelTopFrontRight      AudioChannel = C.FFMS_CH_TOP_FRONT_RIGHT
	ChannelTopBackLeft        AudioChannel = C.FFMS_CH_TOP_BACK_LEFT
	ChannelTopBackCenter      AudioChannel = C.FFMS_CH_TOP_BACK_CENTER
	ChannelTopBackRight       AudioChannel = C.FFMS_CH_TOP_BACK_RIGHT
	ChannelStereoLeft         AudioChannel = C.FFMS_CH_STEREO_LEFT
	ChannelStereoRight        AudioChannel = C.FFMS_CH_STEREO_RIGHT
)

type Resizers int

const (
	ResizerFastBilinear Resizers = C.FFMS_RESIZER_FAST_BILINEAR
	ResizerBilinear     Resizers = C.FFMS_RESIZER_BILINEAR
	ResizerBicubic      Resizers = C.FFMS_RESIZER_BICUBIC
	ResizerX            Resizers = C.FFMS_RESIZER_X
	ResizerPoint        Resizers = C.FFMS_RESIZER_POINT
	ResizerArea         Resizers = C.FFMS_RESIZER_AREA
	ResizerBicublin     Resizers = C.FFMS_RESIZER_BICUBLIN
	ResizerGauss        Resizers = C.FFMS_RESIZER_GAUSS
	ResizerSinc         Resizers = C.FFMS_RESIZER_SINC
	ResizerLanczos      Resizers = C.FFMS_RESIZER_LANCZOS
	ResizerSpline       Resizers = C.FFMS_RESIZER_SPLINE
)

type AudioDelayMode int

const (
	DelayNoShift AudioDelayMode = (iota) + C.FFMS_DELAY_NO_SHIFT
	DelayTimeZero
	DelayFirstVideoTrack
)

type AudioGapFillModes int

const (
	GapFillAuto AudioGapFillModes = (iota) + C.FFMS_GAP_FILL_AUTO
	GapFillDisable
	GapFillEnabled
)

type ChromaLocations int

const (
	LocUnspecified ChromaLocations = (iota) + C.FFMS_LOC_UNSPECIFIED
	LocLeft
	LocCenter
	LocTopLeft
	LocTop
	LocBottomLeft
	LocBottom
)

type ColorRange int

const (
	ColorRangeUnspecified ColorRange = (iota) + C.FFMS_CR_UNSPECIFIED
	ColorRangeMpeg                   // 219*2^(n-8), i.e. 16-235 with 8-bit samples
	ColorRangeJpeg                   // 2^n-1, or "fullrange"
)

type Stereo3DType int

const (
	Stereo3d2d Stereo3DType = (iota) + C.FFMS_S3D_TYPE_2D
	Stereo3dSideBySide
	Stereo3dTopBottom
	Stereo3dFrameSequence
	Stereo3dCheckerBoard
	Stereo3dSideBySideQuincunx
	Stereo3dLines
	Stereo3dColumnns
)

type Stereo3DFlags int

const (
	Stereo3dFlagInvert Stereo3DFlags = C.FFMS_S3D_FLAGS_INVERT
)

type MixingCoefficientType int

const (
	MixingCoefficientQ8 MixingCoefficientType = (iota) + C.FFMS_MIXING_COEFFICIENT_Q8
	MixingCoefficientQ15
	MixingCoefficientFlt
)

type MatrixEncoding int

const (
	MatrixEncodingNone MatrixEncoding = (iota) + C.FFMS_MATRIX_ENCODING_NONE
	MatrixEncodingDobly
	MatrixEncodingProLogicII
	MatrixEncodingProLogicIIX
	MatrixEncodingProLogicIIZ
	MatrixEncodingDolbyEx
	MatrixEncodingDolbyHeadphone
)

type ResampleFilterType int

const (
	ResampleFilterCubic ResampleFilterType = (iota) + C.FFMS_RESAMPLE_FILTER_CUBIC
	ResampleFilterSinc                     /* misnamed as multiple windowsed sinc filters exist, actually called BLACKMAN_NUTTALL */
	ResampleFilterKaiser
)

type AudioDitherMethod int

const (
	ResampleDitherNone AudioDitherMethod = (iota) + C.FFMS_RESAMPLE_DITHER_NONE
	ResampleDitherRectangular
	ResampleDitherTriangular
	ResampleDitherTriangularHighpass
	ResampleDitherTriangularNoiseShaping
)

type LogLevel int

const (
	LogQuiet   LogLevel = C.FFMS_LOG_QUIET
	LogPanic   LogLevel = C.FFMS_LOG_PANIC
	LogFatal   LogLevel = C.FFMS_LOG_FATAL
	LogError   LogLevel = C.FFMS_LOG_ERROR
	LogWarning LogLevel = C.FFMS_LOG_WARNING
	LogInfo    LogLevel = C.FFMS_LOG_INFO
	LogVerbose LogLevel = C.FFMS_LOG_VERBOSE
	LogDebug   LogLevel = C.FFMS_LOG_DEBUG
	LogTrace   LogLevel = C.FFMS_LOG_TRACE
)
