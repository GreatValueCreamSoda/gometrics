package metrics

import (
	"errors"
	"fmt"
	"os"
	"unsafe"

	"github.com/GreatValueCreamSoda/gometrics/blockingpool"
	"github.com/GreatValueCreamSoda/gometrics/comparator"
	vship "github.com/GreatValueCreamSoda/govship"
)

var CVVDPName string = "CVVDP"

// CVVDPHandler manages one or more CVVDP workers and coordinates score
// computation across them.
//
// Internally it owns a blocking pool of vship.CVVDPHandler instances. Each
// worker is stateful and expensive to create, so handlers are reused rather
// than constructed per-frame.
//
// When retrieveDistortionMap is enabled, only a single worker is allowed.
type CVVDPHandler struct {
	pool blockingpool.BlockingPool[*vship.CVVDPHandler]
	// handlerList tracks all created handlers so they can be closed
	// deterministically when the ButterHandler is shut down.
	handlerList []*vship.CVVDPHandler
	// dstWidth and dstHeight are the dimensions of the returned distortion
	// map.
	dstWidth, dstHeight int
	// distortionBuffer stores the per-pixel distortion map when requested. It
	// is reused across calls to avoid repeated allocations.
	distortionBuffer             []float32
	useTemporal, resizeToDisplay bool
	// callback is a callback function called at the end of .Compute() if it
	// and retrieveDistortionMap are set.
	callback DistortionMapCallback

	numWorkers int
}

// Name returns the metric identifier used as the score key.
func (h *CVVDPHandler) Name() string { return CVVDPName }

// NewCVVDPHandler constructs a CVVDPHandler with the requested number of
// worker instances and configuration parameters.
//
// colorA and colorB define the colorspaces of the reference and test images.
//
// useTemporal defines if temporal weighting will be used for score
// calculations.
//
// resizeToDisplay defines if the content will be resized to the displays
// resolution defined in displayModel
//
// displayModel defines the properties of the final display the
// distorted content will be displayed on.
//
// fps is the fps of the source content being compared. This effects vram
// usage heavily.
//
// If retrieveDistortionMap is true, a per-pixel distortion map will be
// computed and stored internally. Only a single worker is allowed when
// retrieveDistortionMap is enabled.
func NewCVVDPHandler(numWorkers int, a, colorB *vship.Colorspace,
	useTemporal, resizeToDisplay bool, distM vship.DisplayModel, fps float32) (
	MetricWithDistortionMap, error) {

	var h CVVDPHandler

	h.pool = blockingpool.NewBlockingPool[*vship.CVVDPHandler](numWorkers)
	h.useTemporal, h.resizeToDisplay = useTemporal, resizeToDisplay

	if !h.resizeToDisplay {
		h.dstWidth, h.dstHeight = int(a.TargetWidth), int(a.TargetHeight)
	} else {
		h.dstWidth, h.dstHeight = distM.DisplayWidth, distM.DisplayHeight
	}

	h.numWorkers = numWorkers

	tmp, e := os.CreateTemp("", "")
	if e != nil {
		return nil, e
	}
	defer tmp.Close()

	distM.Name = "Custom"

	e = vship.DisplayModelsToCVVDPJSONFile([]vship.DisplayModel{distM},
		tmp.Name())
	if e != nil {
		return nil, e
	}

	defer os.Remove(tmp.Name())

	for range numWorkers {
		err := h.createWorker(a, colorB, tmp.Name(), fps)
		if err != nil {
			defer h.Close()
			return nil, err
		}
	}

	return &h, nil
}

// createWorker instantiates a single CVVDP worker and adds it to the pool.
func (h *CVVDPHandler) createWorker(colorA, colorB *vship.Colorspace,
	jsonPath string, fps float32) error {
	vsHandler, exception := vship.NewCVVDPHandlerWithConfig(
		colorA, colorB, fps, h.resizeToDisplay, "Custom", jsonPath)
	if !exception.IsNone() {
		return fmt.Errorf(
			"%s initialization failed: %w", CVVDPName, exception.GetError())
	}

	h.pool.Put(vsHandler)
	h.handlerList = append(h.handlerList, vsHandler)
	return nil
}

// getDistortionBufferAndSize returns a byte slice pointing to the internal
// distortion buffer along with its stride in bytes.
//
// This method performs an unsafe conversion from []float32 to []byte so that
// the buffer can be passed directly into the underlying C-backed Butteraugli
// implementation without copying.
//
// If distortion maps are disabled, it returns nil and zero.
func (h *CVVDPHandler) getDistortionBufferAndSize() ([]byte, int64) {
	var dstptr []byte = nil
	var dstStride int64 = 0

	if h.callback == nil {
		return nil, 0
	}

	dstStride = int64(h.dstWidth) * int64(unsafe.Sizeof(float32(0)))
	totalSize := h.dstWidth * h.dstHeight

	if h.distortionBuffer == nil || len(h.distortionBuffer) != totalSize {
		h.distortionBuffer = make([]float32, totalSize)
	}

	dstptr = unsafe.Slice((*byte)(unsafe.Pointer(&h.distortionBuffer[0])),
		totalSize*4)

	return dstptr, dstStride
}

// Compute calculates the CVVDP perceptual score between two frames.
//
// The method borrows a worker from the pool, computes the scaler score and
// then returns the worker to the pool.
func (h *CVVDPHandler) Compute(a, b *comparator.Frame) (map[string]float64,
	error) {
	handler := h.pool.Get()
	defer h.pool.Put(handler)

	dstptr, dstStride := h.getDistortionBufferAndSize()
	var code vship.ExceptionCode
	var score float64

	aData, aLinesize := a.Read()
	bData, bLinesize := b.Read()

	if h.useTemporal {
		code = handler.ResetScore()
		if !code.IsNone() {
			return nil, fmt.Errorf(
				"%s ResetScore failed: %w", CVVDPName, code.GetError())
		}
		score, code = handler.ComputeScore(dstptr, dstStride, aData, bData,
			aLinesize, bLinesize)
	} else {
		code = handler.Reset()
		if !code.IsNone() {
			return nil, fmt.Errorf(
				"%s Reset failed: %w", CVVDPName, code.GetError())
		}
		code = handler.ResetScore()
		if !code.IsNone() {
			return nil, fmt.Errorf(
				"%s ResetScore failed: %w", CVVDPName, code.GetError())
		}
		score, code = handler.ComputeScore(dstptr, dstStride, aData, bData,
			aLinesize, bLinesize)
	}

	if !code.IsNone() {
		return nil, fmt.Errorf("%s failed to compute score with error: %w",
			CVVDPName, code.GetError())
	}

	if h.callback != nil {
		err := h.callback(h.distortionBuffer)
		if err != nil {
			return nil, err
		}
	}

	return map[string]float64{CVVDPName: score}, nil
}

func (h *CVVDPHandler) SetDistMapCallback(callback DistortionMapCallback) error {
	if h.numWorkers > 1 {
		return errors.New("cannot request more than 1 worker when " +
			"returning a distortion map")
	}
	h.callback = callback
	return nil
}

func (h *CVVDPHandler) GetDistMapResolution() (int, int, error) {
	return h.dstWidth, h.dstHeight, nil
}

// Close releases all underlying CVVDP workers.
func (h *CVVDPHandler) Close() {
	for _, handler := range h.handlerList {
		if handler != nil {
			handler.Close()
		}
	}
	h.handlerList = nil
}
