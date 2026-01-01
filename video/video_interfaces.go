package video

import (
	"errors"
	"fmt"

	vship "github.com/GreatValueCreamSoda/gometrics/c/libvship"
)

// Frame represents a single video Frame's data. It holds the pixel data for
// the three color planes (typically Y, U, V in YUV format) and the line sizes
// (stride) for each plane.
type Frame struct {
	data     [3][]byte // Pixel data for each of the three planes.
	lineSize [3]int    // Line size (stride) for each plane, in bytes.
}

// NewFrame creates a new Frame with the given plane buffers and line sizes.
//
// This is the only supported way to construct a Frame. The provided slices
// become owned by the returned Frame. Callers must not retain references to
// the input slices after this call unless frame lifetime is properly tracked
func NewFrame(data [3][]byte, lineSize [3]int) (Frame, error) {
	for i := 0; i < 3; i++ {
		if len(data[i]) != 0 {
			continue
		}
		return Frame{}, errors.New("plane data must not be nil or zero-length")
	}

	return Frame{data: data, lineSize: lineSize}, nil
}

// Data returns a copy of the array containing the three plane buffers. The
// returned array is safe to read but MUST NOT be modified. The underlying
// slices are still protected by the Frame's ownership.
func (f *Frame) Data() [3][]byte {
	return f.data
}

// LineSizes returns a copy of the array containing the three line sizes
// (strides). The returned array is safe to read and cannot be used to modify
// the Frame.
func (f *Frame) LineSizes() [3]int {
	return f.lineSize
}

// PlaneData returns a read-only view of the data for the requested plane.
func (f *Frame) PlaneData(plane int) []byte {
	if plane < 0 || plane > 2 {
		return nil
	}
	return f.data[plane]
}

// PlaneLineSize returns the line size (stride) in bytes for the requested
// plane.
func (f *Frame) PlaneLineSize(plane int) int {
	if plane < 0 || plane > 2 {
		return 0
	}
	return f.lineSize[plane]
}

// SafeCopyFrom copies pixel data and line sizes from the source frame into
// the receiver frame, preserving the receiver's underlying slice allocations.
// It performs safety checks to prevent incorrect buffer sizes.
//
// Returns an error if any destination plane lacks sufficient capacity.
func (dst *Frame) SafeCopyFrom(src *Frame) error {
	if dst == nil {
		return errors.New("destination frame is nil")
	}
	if src == nil {
		return errors.New("source frame is nil")
	}

	var i int

planeLoop:
	if i >= 3 {
		return nil
	}

	srcPlane, dstPlane := src.data[i], dst.data[i]

	if len(dstPlane) < len(srcPlane) {
		goto ret_error
	}

	copy(dstPlane, srcPlane)
	dst.lineSize[i] = src.lineSize[i]

	i++
	goto planeLoop

ret_error:
	return fmt.Errorf("destination plane %d too small: need %d bytes, have %d",
		i, len(srcPlane), len(dstPlane))
}

type Source interface {
	GetFrame(Frame) error
	GetColorProps() *ColorProperties
	GetNumFrames() int
	GetPlaneSizes() ([3]int, [3]int)
	GetFrameRate() float32
}

// Metric is the interface that every metric must implement
type Metric interface {
	Name() string
	Close()
	Compute(a, b Frame) (map[string]float64, error)
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
