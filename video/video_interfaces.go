package video

import (
	vship "github.com/GreatValueCreamSoda/gometrics/c/libvship"
)

// Frame represents a single video Frame's data. It holds the pixel data for
// the three color planes (typically Y, U, V in YUV format) and the line sizes
// (stride) for each plane.
type Frame struct {
	Data     [3][]byte // Pixel data for each of the three planes.
	LineSize [3]int64  // Line size (stride) for each plane, in bytes.
}

type Source interface {
	GetFrame(*Frame) error
	GetColorProps() *ColorProperties
	GetNumFrames() int
	GetPlaneSizes() ([3]int, [3]int)
	GetFrameRate() float32
}

// Metric is the interface that every metric must implement
type Metric interface {
	Name() string
	Close()
	Compute(a, b *Frame) (map[string]float64, error)
}

type EncoderSettings struct {
	Source     Source
	Output     string
	ColorSpace vship.Colorspace
	Quality    int
	Settings   []string
}

type Encoder interface {
	Encode()
}
