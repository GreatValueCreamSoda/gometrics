package metrics

import (
	"encoding/binary"
	"fmt"
	"io"
	"os/exec"
	"strconv"
	"sync"
	"unsafe"

	"github.com/GreatValueCreamSoda/gometrics/comparator"
)

type MetricWithDistortionMap interface {
	SetDistMapCallback(DistortionMapCallback) error
	GetDistMapResolution() (int, int, error)
	comparator.Metric
}

type DistortionMapCallback func([]float32) error

type HeatmapWriter struct {
	cmd  *exec.Cmd
	pipe io.WriteCloser

	maxValue float32

	normalized []float32
	byteBuf    []byte

	closeOnce sync.Once
}

func WriteDistMapToVideo(metric MetricWithDistortionMap, frameRate float32,
	settings []string, path string, maxValue float32) (*HeatmapWriter,
	error) {

	if maxValue <= 0 {
		return nil, fmt.Errorf("maxValue must be > 0")
	}

	width, height, err := metric.GetDistMapResolution()
	if err != nil {
		return nil, err
	}
	if width <= 0 || height <= 0 {
		return nil, fmt.Errorf("invalid resolution: %dx%d", width, height)
	}

	cmd, pipe, err := startFFmpeg(width, height, frameRate, settings, path)
	if err != nil {
		return nil, err
	}

	writer := &HeatmapWriter{
		cmd:      cmd,
		pipe:     pipe,
		maxValue: maxValue,
	}

	if err := cmd.Start(); err != nil {
		pipe.Close()
		return nil, fmt.Errorf("failed to start ffmpeg: %w", err)
	}

	if err := metric.SetDistMapCallback(writer.WriteDistortion); err != nil {
		_ = writer.Close()
		return nil, err
	}

	return writer, nil
}

func startFFmpeg(width int, height int, frameRate float32, settings []string,
	outputPath string) (*exec.Cmd, io.WriteCloser, error) {

	frameRateStr := strconv.FormatFloat(float64(frameRate), 'f', -1, 64)
	resolution := fmt.Sprintf("%dx%d", width, height)

	filter := "format=rgb24,pseudocolor=p=heat"

	if settings == nil {
		settings = []string{"-c:v", "libx264", "-preset", "fast", "-crf", "18"}
	}

	args := append([]string{
		"-y",
		"-f", "rawvideo",
		"-pixel_format", "grayf32le",
		"-s", resolution,
		"-r", frameRateStr,
		"-i", "-",
		"-vf", filter,
		"-pix_fmt", "yuv420p",
	}, append(settings, outputPath)...)

	cmd := exec.Command("ffmpeg", args...)

	if pipe, err := cmd.StdinPipe(); err != nil {
		return nil, nil, fmt.Errorf("failed to get ffmpeg stdin pipe: %w", err)
	} else {
		return cmd, pipe, nil
	}
}

func (h *HeatmapWriter) WriteDistortion(input []float32) error {
	if len(input) == 0 {
		return nil
	}

	h.ensureBuffers(len(input))
	h.normalize(input)
	return h.writeFloats()
}

func (h *HeatmapWriter) ensureBuffers(n int) {
	if cap(h.normalized) < n {
		h.normalized = make([]float32, n)
		h.byteBuf = make([]byte, n*4)
		return
	}

	h.normalized = h.normalized[:n]
	h.byteBuf = h.byteBuf[:n*4]
}

func (h *HeatmapWriter) normalize(input []float32) {
	scale := float32(1.0) / h.maxValue

	for i, v := range input {
		if v > h.maxValue {
			v = h.maxValue
		}
		h.normalized[i] = v * scale
	}
}

func (h *HeatmapWriter) writeFloats() error {
	for i, v := range h.normalized {
		binary.LittleEndian.PutUint32(
			h.byteBuf[i*4:],
			binary.LittleEndian.Uint32((*[4]byte)(unsafe.Pointer(&v))[:]),
		)
	}
	_, err := h.pipe.Write(h.byteBuf)
	return err
}

func (h *HeatmapWriter) Close() error {
	var waitErr error

	h.closeOnce.Do(func() {
		_ = h.pipe.Close()
		waitErr = h.cmd.Wait()
	})

	if waitErr != nil {
		return fmt.Errorf("ffmpeg failed: %w", waitErr)
	}
	return nil
}
