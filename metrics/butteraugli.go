package metrics

import (
	"errors"
	"fmt"
	"unsafe"

	"github.com/GreatValueCreamSoda/gometrics/blockingpool"
	"github.com/GreatValueCreamSoda/gometrics/comparator"
	vship "github.com/GreatValueCreamSoda/govship"
)

const ButteraugliName string = "Butteraugli"

// ButterHandler manages one or more Butteraugli workers and coordinates
// score computation across them.
//
// Internally it owns a blocking pool of vship.ButteraugliHandler instances.
// Each worker is stateful and expensive to create, so handlers are reused
// rather than constructed per-frame.
//
// When retrieveDistortionMap is enabled, only a single worker is allowed.
type ButterHandler struct {
	// pool holds reusable Butteraugli handlers for concurrent scoring.
	pool blockingpool.BlockingPool[*vship.ButteraugliHandler]
	// handlerList tracks all created handlers so they can be closed
	// deterministically when the ButterHandler is shut down.
	handlerList []*vship.ButteraugliHandler
	// dstWidth and dstHeight are the dimensions of the returned distortion
	// map.
	dstWidth, dstHeight int
	// distortionBuffer stores the per-pixel distortion map when requested. It
	// is reused across calls to avoid repeated allocations.
	distortionBuffer []float32
	// callback is a callback function called at the end of .Compute() if it
	// and retrieveDistortionMap are set.
	callback DistortionMapCallback

	numWorkers int
}

func (h *ButterHandler) Name() string { return ButteraugliName }

// NewButterHandler constructs a ButterHandler with the requested number of
// worker instances and configuration parameters.
//
// colorA and colorB define the colorspaces of the reference and test images.
// qNorm specified the p-norm that will be stored in the qnrom score result.
//
// If retrieveDistortionMap is true, a per-pixel distortion map will be
// computed and stored internally. Only a single worker is allowed when
// retrieveDistortionMap is enabled.
func NewButterHandler(numWorkers int, colorA, colorB *vship.Colorspace,
	qNorm int, displayIntensity float32, retrieveDistortionMap bool) (
	MetricWithDistortionMap, error) {
	var handler ButterHandler
	var err error

	handler.pool = blockingpool.NewBlockingPool[*vship.ButteraugliHandler](
		numWorkers)
	handler.dstWidth = int(colorA.TargetWidth)
	handler.dstHeight = int(colorA.TargetHeight)
	handler.numWorkers = numWorkers

	for range numWorkers {
		err = handler.createWorker(colorA, colorB, qNorm, displayIntensity)
		if err == nil {
			continue
		}
		defer handler.Close()
		return nil, err
	}

	return &handler, nil
}

// createWorker instantiates a single Butteraugli handler and registers it
// with both the worker pool and the internal handler list.
//
// Any failure during initialization is wrapped with metric context to make
// upstream error reporting clearer.
func (h *ButterHandler) createWorker(colorA, colorB *vship.Colorspace,
	Qnorm int, DisplayIntensity float32) error {
	vsHandler, exception := vship.NewButteraugliHandler(colorA, colorB,
		Qnorm, DisplayIntensity)
	if exception.IsNone() {
		h.pool.Put(vsHandler)
		h.handlerList = append(h.handlerList, vsHandler)
		return nil
	}
	var err error = exception.GetError()
	return fmt.Errorf("%s initialization failed: %w", ButteraugliName, err)
}

// getDistortionBufferAndSize returns a byte slice pointing to the internal
// distortion buffer along with its stride in bytes.
//
// This method performs an unsafe conversion from []float32 to []byte so that
// the buffer can be passed directly into the underlying C-backed Butteraugli
// implementation without copying.
//
// If distortion maps are disabled, it returns nil and zero.
func (h *ButterHandler) getDistortionBufferAndSize() ([]byte, int64) {
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

// Compute calculates Butteraugli perceptual difference scores between two
// frames.
//
// The method borrows a worker from the pool, computes scalar scores (NormQ,
// Norm3, NormInf), and then returns the worker to the pool.
//
// The returned map keys are prefixed with ButteraugliName to avoid collisions
// with other metrics.
func (h *ButterHandler) Compute(a, b *comparator.Frame) (map[string]float64,
	error) {
	handler := h.pool.Get()
	defer h.pool.Put(handler)
	dstptr, dstStride := h.getDistortionBufferAndSize()

	aData, aLinesize := a.Read()
	bData, bLinesize := b.Read()

	var score vship.ButteraugliScore
	exception := handler.ComputeScore(&score, dstptr, dstStride, aData, bData,
		aLinesize, bLinesize)
	if !exception.IsNone() {
		return nil, fmt.Errorf("%s failed to compute score with error: %w",
			ButteraugliName, exception.GetError())
	}

	if h.callback != nil {
		err := h.callback(h.distortionBuffer)
		if err != nil {
			return nil, err
		}
	}

	scores := map[string]float64{ButteraugliName + "NormQ": score.NormQ,
		ButteraugliName + "Norm3": score.Norm3,
		ButteraugliName + "Inf":   score.NormInf,
	}

	return scores, nil
}

func (h *ButterHandler) SetDistMapCallback(callback DistortionMapCallback) error {
	if h.numWorkers > 1 {
		return errors.New("cannot request more than 1 worker when " +
			"returning a distortion map")
	}
	h.callback = callback
	return nil
}

func (h *ButterHandler) GetDistMapResolution() (int, int, error) {
	return h.dstWidth, h.dstHeight, nil
}

// Close releases all underlying Butteraugli handlers.
//
// After calling Close, the ButterHandler should be considered unusable. This
// method is idempotent and safe to call multiple times.
func (h *ButterHandler) Close() {
	for _, handler := range h.handlerList {
		if handler != nil {
			handler.Close()
		}
	}
	h.handlerList = nil
}
