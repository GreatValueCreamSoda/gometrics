package metrics

import (
	"errors"
	"fmt"

	"github.com/GreatValueCreamSoda/gometrics/blockingpool"
	"github.com/GreatValueCreamSoda/gometrics/comparator"
	vship "github.com/GreatValueCreamSoda/govship"
)

var ErrDistortionMapUnsupported = errors.New("distortion maps are unsupported for this metric.")

// SSIMulacra2Name is the canonical metric name used for score reporting.
var SSIMulacra2Name string = "Ssimulacra2"

// Ssimu2Handler manages one or more SSIMULACRA2 workers and coordinates
// score computation across them.
//
// Internally it owns a blocking pool of vship.SSIMU2Handler instances. Each
// worker is stateful and relatively expensive to create, so handlers are
// reused rather than constructed per-frame.
//
// This handler only produces a single scalar score per comparison and does not
// allocate or retain any per-frame buffers.
type Ssimu2Handler struct {
	pool        blockingpool.BlockingPool[*vship.SSIMU2Handler]
	handlerList []*vship.SSIMU2Handler
}

// Name returns the metric identifier used as the score key.
func (h *Ssimu2Handler) Name() string { return "ssimu2" }

// NewSSIMU2Handler constructs a Ssimu2Handler with the requested number of
// worker instances.
//
// colorA and colorB define the colorspaces of the reference and test images.
func NewSSIMU2Handler(numWorkers int, colorA, colorB *vship.Colorspace) (
	comparator.Metric, error) {
	var h Ssimu2Handler
	h.pool = blockingpool.NewBlockingPool[*vship.SSIMU2Handler](numWorkers)

	for range numWorkers {
		err := h.createWorker(colorA, colorB)
		if err == nil {
			continue
		}
		defer h.Close()
		return nil, err
	}

	return &h, nil
}

// createWorker instantiates a single SSIMULACRA2 handler and registers it with
// both the worker pool and the internal handler list.
//
// Any failure during initialization is wrapped with metric context to make
// upstream error reporting clearer.
func (h *Ssimu2Handler) createWorker(colorA, colorB *vship.Colorspace) error {

	vsHandler, exception := vship.NewSSIMU2Handler(colorA, colorB)
	if !exception.IsNone() {
		defer h.Close()
		var err error = exception.GetError()
		return fmt.Errorf("%s initialization failed: %w", SSIMulacra2Name, err)
	}
	h.pool.Put(vsHandler)
	h.handlerList = append(h.handlerList, vsHandler)
	return nil
}

func (h *Ssimu2Handler) DistortionMap() ([]float32, int, int, error) {
	return nil, 0, 0, ErrDistortionMapUnsupported
}

// Close releases all underlying SSIMULACRA2 handlers.
//
// After calling Close, the Ssimu2Handler should be considered unusable. This
// method is idempotent and safe to call multiple times.
func (h *Ssimu2Handler) Close() {
	for _, handler := range h.handlerList {
		if handler != nil {
			handler.Close()
		}
	}
	h.handlerList = nil
}

// Compute calculates the SSIMULACRA2 perceptual similarity score between two
// frames.
//
// The method borrows a worker from the pool, computes the scalar score, and
// then returns the worker to the pool.
//
// The returned map contains a single entry keyed by Name().
func (h *Ssimu2Handler) Compute(a, b *comparator.Frame) (map[string]float64,
	error) {
	handler := h.pool.Get()
	defer h.pool.Put(handler)

	aData, aLinesize := a.Read()
	bData, bLinesize := b.Read()

	score, code := handler.ComputeScore(aData, bData, aLinesize, bLinesize)

	if !code.IsNone() {
		return nil, fmt.Errorf("%s computation failed: %v", SSIMulacra2Name,
			code.GetError())
	}
	return map[string]float64{h.Name(): score}, nil
}
