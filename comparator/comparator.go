package comparator

import (
	"context"
	"errors"
	"fmt"
	"sync"

	"github.com/GreatValueCreamSoda/gometrics/blockingpool"
	vship "github.com/GreatValueCreamSoda/govship"
	"golang.org/x/sync/errgroup"
)

type ProgressCallback func(done int, total int)

type Source interface {
	GetFrame(*Frame) error
	GetColorspace() *vship.Colorspace
	GetNumFrames() int
	GetPlaneSizes() ([3]int, [3]int)
}

// Metric is the interface that every metric must implement
type Metric interface {
	Name() string
	Close()
	Compute(a, b *Frame) (map[string]float64, error)
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

func (f *Frame) Read() ([3][]byte, [3]int64) { return f.data, f.lineSize }

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
	a, b  *Frame
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
	videoA       Source                            // Source for Video A.
	videoB       Source                            // Source for Video B.
	metrics      []Metric                          // Metrics to compute on each frame pair.
	frameThreads int                               // Number of concurrent metric workers.
	framePoolA   blockingpool.BlockingPool[*Frame] // Reusable frame buffers for video A.
	framePoolB   blockingpool.BlockingPool[*Frame] // Reusable frame buffers for video B.
	numFrames    int                               // Number of frame pairs to process.

	// Internal channels for the pipeline stages.
	videoAFrameChan chan *Frame
	videoBFrameChan chan *Frame
	fPairChan       chan framePair
	scoresChan      chan metricResult

	// finalScores accumulates per-metric lists of per-frame scores. It is
	// populated during Run by the aggregation goroutine.
	finalScores map[string][]float64

	ctx       context.Context
	ctxCancel context.CancelCauseFunc

	progress ProgressCallback
}

// NewComparator creates and initializes a Comparator.
//
// It validates the inputs, pre-allocates reusable frame buffers sized to the
// videos' plane dimensions, and sets up the internal pipeline channels.
//
// The threads parameter controls how many goroutines will concurrently compute
// metrics. If a Metric must process frames sequentially this must be set to 1.
// numFrames specifies how many corresponding frame pairs will be processed
// (must be ≤ the number of frames available in both sources).
//
// Returns an error if any input is invalid or if the sources have fewer frames
// than requested.
func NewComparator(videoA, videoB Source, metrics []Metric, frameThreads int,
	numFrames int) (Comparator, error) {
	var Comparator Comparator
	Comparator.videoA, Comparator.videoB = videoA, videoB
	Comparator.metrics = metrics
	Comparator.frameThreads = frameThreads
	Comparator.numFrames = numFrames

	if Comparator.videoA == nil || Comparator.videoB == nil {
		return Comparator, errors.New("videoA and videoB must be non nil")
	}

	if len(Comparator.metrics) < 1 {
		return Comparator, errors.New("at least one metric must be specifed")
	}

	if Comparator.frameThreads < 1 {
		return Comparator, errors.New("at least 1 thread must be used for " +
			"metric computation")
	}

	if videoA.GetNumFrames() < Comparator.numFrames {
		return Comparator, errors.New("videoa has less frames than numframes")
	}

	if videoB.GetNumFrames() < Comparator.numFrames {
		return Comparator, errors.New("videob has less frames than numframes")
	}

	var totalFrameBuffers int = 0

	Comparator.videoBFrameChan = make(chan *Frame, 1)
	Comparator.videoAFrameChan = make(chan *Frame, 1)
	totalFrameBuffers = totalFrameBuffers + 1

	Comparator.fPairChan = make(chan framePair, frameThreads/2)
	// One frame buffers per queue item. One buffers per working metric thread.
	totalFrameBuffers = totalFrameBuffers + (frameThreads/2 + 1) + frameThreads

	Comparator.framePoolA = blockingpool.NewBlockingPool[*Frame](totalFrameBuffers)
	Comparator.framePoolB = blockingpool.NewBlockingPool[*Frame](totalFrameBuffers)

	for range totalFrameBuffers {
		Comparator.allocateFrameBuffer()
	}

	Comparator.scoresChan = make(chan metricResult, frameThreads)
	Comparator.finalScores = make(map[string][]float64)

	return Comparator, nil
}

func (c *Comparator) allocateFrameBuffer() {
	var frameA *Frame = new(Frame)
	sA, lA := c.videoA.GetPlaneSizes()
	frameA.data = [3][]byte{
		make([]byte, sA[0]), make([]byte, sA[1]), make([]byte, sA[2])}
	frameA.lineSize = [3]int64{int64(lA[0]), int64(lA[1]), int64(lA[2])}
	c.framePoolA.Put(frameA)

	var frameB *Frame = new(Frame)
	SB, lB := c.videoB.GetPlaneSizes()
	frameB.data = [3][]byte{
		make([]byte, SB[0]), make([]byte, SB[1]), make([]byte, SB[2])}
	frameB.lineSize = [3]int64{int64(lB[0]), int64(lB[1]), int64(lB[2])}
	c.framePoolB.Put(frameB)
}

// Run starts the comparison pipeline.
//
// It spawns reader threads for both videos, a frame-pairing goroutine, the
// requested number of metric computation workers, and a final aggregation
// goroutine. Run blocks until all frames have been processed and results
// aggregated.
//
// Run returns the per frame scores.
func (c *Comparator) Run(parentCtx context.Context) (map[string][]float64,
	error) {
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

// SetProgressCallback registers a progress callback on the Comparator. It must
// be called before Run. Passing nil clears the callback.
func (c *Comparator) SetProgressCallback(cb ProgressCallback) {
	c.progress = cb
}

// ----------------------------------------------------------------------------
// Reader Threads
// ----------------------------------------------------------------------------

// spawnReaderThreads starts two goroutines that decode frames from the two
// video sources and send them on their respective channels.
//
// When both readers finish, the frame channels are closed.
func (c *Comparator) spawnReaderThreads() error {
	group, ctx := errgroup.WithContext(c.ctx)

	group.Go(func() error {
		return c.readerThread(ctx, c.videoA, c.videoAFrameChan, c.framePoolA)
	})
	group.Go(func() error {
		return c.readerThread(ctx, c.videoB, c.videoBFrameChan, c.framePoolB)
	})

	err := group.Wait() // if any reader fails, ctx is cancelled automatically
	return err
}

// readerThread decodes numFrames frames from video into frameChan, reusing
// buffers obtained from framePool.
func (c *Comparator) readerThread(ctx context.Context, video Source,
	frameChan chan *Frame, framePool blockingpool.BlockingPool[*Frame]) error {

	for i := 0; i < c.numFrames; i++ {
		var frame *Frame

		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
			frame = framePool.Get()
		}

		if err := video.GetFrame(frame); err != nil {
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
func (c *Comparator) spawnFramePairThreads() error {
	for i := range make([]struct{}, c.numFrames) {
		var a, b *Frame

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
func (c *Comparator) spawnMetricsThreads() error {
	group, ctx := errgroup.WithContext(c.ctx)

	for range c.frameThreads {
		group.Go(func() error { return c.metricThread(ctx) })
	}

	err := group.Wait()
	return err
}

// metricThread consumes frame pairs from fPairChan, computes all requested
// metrics for each pair, and sends a metricResult on scoresChan.
//
// Returns the first error encountered, which triggers context cancellation
// upstream.
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

// computeFrameMetrics computes all metrics for a single frame pair.
//
// It returns the frames to their respective pools via defer statements before
// returning. It checks for duplicate metric names across the provided metrics
// and returns an error if any are found.
func (c *Comparator) computeFrameMetrics(pair framePair, metrics []Metric) (
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
	metric Metric, mu *sync.Mutex) error {
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
