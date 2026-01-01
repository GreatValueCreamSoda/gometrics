package comparator

import (
	"context"
	"errors"
	"fmt"
	"sync"

	"github.com/GreatValueCreamSoda/gometrics/blockingpool"
	vship "github.com/GreatValueCreamSoda/gometrics/c/libvship"
	"github.com/GreatValueCreamSoda/gometrics/video"
	"golang.org/x/sync/errgroup"
)

type ProgressCallback func(done int, total int)

type Source interface {
	GetFrame(*Frame) error
	GetColorspace() *vship.Colorspace
	GetNumFrames() int
	GetPlaneSizes() ([3]int, [3]int)
	GetFrameRate() float32
}

// Frame represents a single video Frame's data. It holds the pixel data for
// the three color planes (typically Y, U, V in YUV format) and the line sizes
// (stride) for each plane.
type Frame struct {
	data     [3][]byte // Pixel data for each of the three planes.
	lineSize [3]int64  // Line size (stride) for each plane, in bytes.
}

func (f *Frame) Write(data [3][]byte, lineSize [3]int64) error {
	for i := range f.data {
		if len(f.data[i]) != len(data[i]) {
			return errors.New("failed to write frame data. data plane sizes " +
				"do not match")
		}
	}

	for p := 0; p < 3; p++ {
		copy(f.data[p], data[p])
		f.lineSize[p] = int64(lineSize[p])
	}

	return nil
}

// metricResult holds the computed metric scores for a specific frame pair.
// The scores are a map of metric names to their float64 values.
type metricResult struct {
	// The index of the frame pair these scores belong to.
	index  int
	scores map[string]float64 // Map of metric names to computed scores.
}

// framePair represents a paired set of frames from video A and video B, along
// with their indices for tracking.
type framePair struct {
	index int
	a, b  *video.Frame
}

// Comparator orchestrates the concurrent comparison of two video sources using
// a set of metrics.
//
// It reads frames from both sources in parallel, pairs them, computes the
// requested metrics on each pair using a configurable number of worker
// goroutines, and aggregates the results.
//
// The zero value is not valid; use NewComparator to construct an instance.
type Comparator struct {
	// Source video A and B are the two videos that will be compared to each
	// other
	videoA, videoB video.Source
	// List of metrics who scores will be computed on each frame concurrently
	metrics []video.Metric
	// The number of frames that metrics will be ran on concurrently. This is
	// not the number of metric threads as each metric will be called
	// concurrently on each frame.
	frameThreads int // Number of concurrent metric workers.
	// A pool of reusable frames buffers that reader threads will pull from,
	// copy the frame data to, and that metric threads will return.
	framePoolA, framePoolB blockingpool.BlockingPool[*video.Frame]
	// The total number of frames that will be compared between video A and B.
	numFrames int

	// Internal channels for the pipeline stages.

	// videoAFrameChan and videoBFrameChan as the name implies are two channels
	// frame reader thread A and B will write frames squentially to. These are
	// then consumed by the frame pair goroutine.
	videoAFrameChan, videoBFrameChan chan *video.Frame

	// fPairChan is the channel all metric threads will read from. Each
	// framePair will contain one frame from video A and one frame from video B
	// that it will conpute metrics for.
	fPairChan chan framePair

	// scoresChan is the channel metric threads will send their results to that
	// will be consumed by the aggergation goroutine.
	scoresChan chan metricResult

	// finalScores accumulates per-metric lists of per-frame scores. It is
	// populated during Run by the aggregation goroutine.
	finalScores map[string][]float64

	// ctx is the global context that all sub goroutines will run with during
	// .Run(). This is canceled if any error occures within any stage of the
	// pipeline.
	ctx       context.Context
	ctxCancel context.CancelCauseFunc

	// progress is a function that is called every time the score aggergator
	// goroutine receives a metric result from a metric thread. Used to update
	// the user of the total ammount of frames compared relative to the total.
	//
	// Note: Metrics are not always computes in increasing order. The progress
	// callback might be called with a earlier "total" than before, or for a
	// frame before previous frames are done if frame threads is greater than 1
	progress ProgressCallback
}

// NewComparator creates a new Comparator instance.
//
// Validates inputs, preallocates reusable frame buffers, and initializes
// channels.
//
// frameThreads controls how many frame pairs are processed concurrently. If
// any metric requires strict sequential processing, set frameThreads = 1.
//
// numFrames specifies how many frame pairs to compare (must not exceed the
// available frames in either source).
func NewComparator(videoA, videoB video.Source, metrics []video.Metric, frameThreads,
	numFrames int) (Comparator, error) {
	c := Comparator{
		videoA:       videoA,
		videoB:       videoB,
		metrics:      metrics,
		frameThreads: frameThreads,
		numFrames:    numFrames,
		finalScores:  make(map[string][]float64),
	}

	if err := c.validateArguments(); err != nil {
		return Comparator{}, err
	}

	totalBuffers := c.calculateTotalNumberOfFrameBuffers()

	c.framePoolA = blockingpool.NewBlockingPool[*video.Frame](totalBuffers)
	c.framePoolB = blockingpool.NewBlockingPool[*video.Frame](totalBuffers)

	for range totalBuffers {
		err := c.allocateFrameBuffer()
		if err != nil {
			return Comparator{}, err
		}
	}

	c.scoresChan = make(chan metricResult, frameThreads)

	return c, nil
}

func (c *Comparator) validateArguments() error {
	if c.videoA == nil || c.videoB == nil {
		return errors.New("either video a or video b was passed as a nil ptr")
	}

	if len(c.metrics) < 1 {
		return errors.New("at least one metric must be passed to measure with")
	}

	if c.frameThreads < 1 {
		return errors.New("at least 1 frame thread must be used to compare")
	}

	if c.videoA.GetNumFrames() < c.numFrames {
		return errors.New("videoa has less frames than number of frames to " +
			" be compared")
	}

	if c.videoB.GetNumFrames() < c.numFrames {
		return errors.New("videob has less frames than number of frames to " +
			" be compared")
	}

	return nil
}

// calculateTotalNumberOfFrameBuffers returns conservative estimate of needed
// buffers accounting for pipeline stages and worker concurrency.
func (c *Comparator) calculateTotalNumberOfFrameBuffers() int {
	c.videoBFrameChan = make(chan *video.Frame, 1)
	c.videoAFrameChan = make(chan *video.Frame, 1)
	var totalFrameBuffers int = 1

	c.fPairChan = make(chan framePair, c.frameThreads/2)
	totalFrameBuffers = totalFrameBuffers + (c.frameThreads/2 + 1) +
		c.frameThreads

	return totalFrameBuffers
}

func (c *Comparator) allocateFrameBuffer() error {
	sA, lA := c.videoA.GetPlaneSizes()
	aData1, code := vship.PinnedMalloc(sA[0])
	if !code.IsNone() {
		return code.GetError()
	}
	aData2, code := vship.PinnedMalloc(sA[1])
	if !code.IsNone() {
		return code.GetError()
	}
	aData3, code := vship.PinnedMalloc(sA[2])
	if !code.IsNone() {
		return code.GetError()
	}

	fA := new(video.Frame)
	fA.Data = [3][]byte{aData1, aData2, aData3}
	fA.LineSize = [3]int64{int64(lA[0]), int64(lA[1]), int64(lA[2])}
	c.framePoolA.Put(fA)

	sB, lB := c.videoA.GetPlaneSizes()
	bData1, code := vship.PinnedMalloc(sB[0])
	if !code.IsNone() {
		return code.GetError()
	}
	bData2, code := vship.PinnedMalloc(sB[1])
	if !code.IsNone() {
		return code.GetError()
	}
	bData3, code := vship.PinnedMalloc(sB[2])
	if !code.IsNone() {
		return code.GetError()
	}

	fB := new(video.Frame)
	fB.Data = [3][]byte{bData1, bData2, bData3}
	fB.LineSize = [3]int64{int64(lB[0]), int64(lB[1]), int64(lB[2])}
	c.framePoolB.Put(fB)

	return nil
}

// Run executes the full comparison pipeline and blocks until completion.
// Returns per-metric arrays of per-frame scores.
func (c *Comparator) Run(parentCtx context.Context) (
	map[string][]float64, error) {
	group, ctx := errgroup.WithContext(parentCtx)
	c.ctx = ctx

	group.Go(func() error {
		defer close(c.videoAFrameChan)
		defer close(c.videoBFrameChan)
		return c.spawnReaderThreads()
	})

	group.Go(func() error {
		defer close(c.fPairChan)
		return c.spawnFramePairThreads()
	})

	group.Go(func() error {
		defer close(c.scoresChan)
		return c.spawnMetricsThreads()
	})

	group.Go(c.aggregateResults)

	return c.finalScores, group.Wait()
}

// SetProgressCallback registers an optional progress callback. Must be called
// before Run(). Pass nil to clear.
func (c *Comparator) SetProgressCallback(cb ProgressCallback) {
	c.progress = cb
}

// ----------------------------------------------------------------------------
// Reader Threads
// ----------------------------------------------------------------------------

// spawnReaderThreads starts two goroutines to read video A and B in parallel.
//
// If any error occures exectuion is terminated early and the error is returned
func (c *Comparator) spawnReaderThreads() error {
	group, ctx := errgroup.WithContext(c.ctx)

	group.Go(func() error {
		return c.readerThread(ctx, c.videoA,
			c.videoAFrameChan, c.framePoolA)
	})
	group.Go(func() error {
		return c.readerThread(ctx, c.videoB,
			c.videoBFrameChan, c.framePoolB)
	})

	err := group.Wait()
	return err
}

// readerThread reads from the supplied video source and sends them to the
// frameChan till the total number of frames is read or the context is canceled
func (c *Comparator) readerThread(ctx context.Context, source video.Source,
	frameChan chan *video.Frame, framePool blockingpool.BlockingPool[*video.Frame]) error {

	for i := 0; i < c.numFrames; i++ {
		var frame *video.Frame

		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
			frame = framePool.Get()
		}

		if err := source.GetFrame(frame); err != nil {
			return err
		}

		select {
		case <-ctx.Done():
			return ctx.Err()
		case frameChan <- frame:
		}
	}

	return nil
}

// ----------------------------------------------------------------------------
// Frame Pair Threads
// ----------------------------------------------------------------------------

// spawnFramePairThreads starts a single goroutine that consumes one frame from
// each video channel, pairs them, and sends the pair on fPairChan.
//
// When the reader channels close, fPairChan is closed.
//
// If any error occures exectuion is terminated early and the error is returned
func (c *Comparator) spawnFramePairThreads() error {
	for i := range make([]struct{}, c.numFrames) {
		var a, b *video.Frame

		select {
		case <-c.ctx.Done():
			return c.ctx.Err()
		case a = <-c.videoAFrameChan:
			if a == nil {
				return nil
			}
		}

		select {
		case <-c.ctx.Done():
			return c.ctx.Err()
		case b = <-c.videoBFrameChan:
			if b == nil {
				return nil
			}
		}

		select {
		case <-c.ctx.Done():
			return c.ctx.Err()
		case c.fPairChan <- framePair{i, a, b}:
		}
	}
	return nil
}

// ----------------------------------------------------------------------------
// Metric Threads
// ----------------------------------------------------------------------------

// spawnMetricsThreads starts metricThreads goroutines that each run
// metricThread, consuming frame pairs and producing metricResult values.
//
// When fPairChan closes, scoresChan is closed.
//
// If any error occures exectuion is terminated early and the error is returned
func (c *Comparator) spawnMetricsThreads() error {
	group, ctx := errgroup.WithContext(c.ctx)

	for range c.frameThreads {
		group.Go(func() error { return c.metricThread(ctx) })
	}

	err := group.Wait()
	return err
}

// metricThread consumes frame pairs from fPairChan, computes all requested
// metrics for each pair in parallel, and sends a metricResult on scoresChan.
//
// If any error occures exectuion is terminated early and the error is returned
func (c *Comparator) metricThread(ctx context.Context) error {
	for pair := range withContext(ctx, c.fPairChan) {
		scores, err := c.computeFrameMetrics(pair, c.metrics)
		if err != nil {
			return err
		}

		select {
		case <-ctx.Done():
			return ctx.Err()
		case c.scoresChan <- metricResult{pair.index, scores}:
		}
	}
	return nil
}

// computeFrameMetrics runs all metrics in parallel for one frame pair. Returns
// frames to pools on exit (via defer).
func (c *Comparator) computeFrameMetrics(pair framePair, metrics []video.Metric) (
	map[string]float64, error) {
	defer c.framePoolA.Put(pair.a)
	defer c.framePoolB.Put(pair.b)

	if len(metrics) == 0 {
		return map[string]float64{}, nil
	}

	result := make(map[string]float64, len(metrics)*3)

	// We let each metric within a fram run in parallel instead of one at a
	// time. This on my machine with ssimu2 + butter increased fps from 85-87
	// to a consistent 90 fps with 1 worker. Should give small gains when
	// generating distortion maps.
	var mu sync.Mutex
	group, _ := errgroup.WithContext(c.ctx)

	// Skip the overhead of spawning a new goroutine and just run it within
	// this one.
	//if len(metrics) == 1 {
	//	return result, c.computeFrameMetric(pair, result, metrics[0], &mu)
	//}

	for _, metric := range metrics {
		group.Go(func() error {
			return c.computeFrameMetric(pair, result, metric, &mu)
		})
	}

	return result, group.Wait()
}

// computeFrameMetric invokes a single Metric's Compute method and merges its
// results into the result map, returning an error on failure or duplicate
// keys.
func (Comparator) computeFrameMetric(pair framePair, res map[string]float64,
	metric video.Metric, mu *sync.Mutex) error {
	scores, err := metric.Compute(pair.a, pair.b)
	if err != nil {
		return fmt.Errorf("%s computation failed: %w", metric.Name(), err)
	}
	mu.Lock()
	defer mu.Unlock()
	for k, v := range scores {
		if _, exists := res[k]; exists {
			return fmt.Errorf("duplicate metric %q from %s", k, metric.Name())
		}
		res[k] = v
	}

	return nil
}

// ----------------------------------------------------------------------------
// Aggergation Threads
// ----------------------------------------------------------------------------

// aggergateResults consumes all metricResult values from scoresChan and
// accumulates them into the Comparator's finalScores map.
func (c *Comparator) aggregateResults() error {
	completed := 0
	for res := range withContext(c.ctx, c.scoresChan) {
		for name, val := range res.scores {
			if res.index < 0 || res.index >= c.numFrames {
				return errors.New("aggergated index outside of numframe")
			}
			if c.finalScores[name] == nil {
				c.finalScores[name] = make([]float64, c.numFrames)
			}
			c.finalScores[name][res.index] = val
		}
		completed++
		if c.progress != nil {
			c.progress(completed, c.numFrames)
		}
	}
	return nil
}

// withContext returns a new read-only channel that mirrors values from the
// input channel ch until either ch is closed or the provided context ctx is
// canceled.
//
// The returned channel will be closed when one of the following occurs:
//   - The input channel ch is closed (all values have been forwarded).
//   - The context ctx is canceled (ctx.Done() becomes readable).
//
// Usage example:
//
//	func processWithTimeout(ctx context.Context, input <-chan WorkItem) {
//	    for item := range withContext(ctx, input) {
//	        // Process item; loop exits cleanly on ctx cancellation or input
//			// close
//	        doWork(item)
//	    }
//	}
//
// Parameters:
//   - ctx context.Context: The context that controls cancellation.
//   - ch <-chan T:        The source channel to mirror.
//
// Returns:
//
//	<-chan T: A new channel that yields values from ch until either terminates.
func withContext[T any](ctx context.Context, ch <-chan T) <-chan T {
	out := make(chan T, 1) // buffered to avoid blocking on send

	go func() {
		defer close(out)
		for {
			select {
			case <-ctx.Done():
				return
			case v, ok := <-ch:
				if !ok {
					return
				}
				select {
				case out <- v:
				case <-ctx.Done():
					return
				}
			}
		}
	}()
	return out
}
