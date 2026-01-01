package main

import (
	"context"
	"fmt"
	"log"

	vship "github.com/GreatValueCreamSoda/gometrics/c/libvship"
	"github.com/GreatValueCreamSoda/gometrics/video"
	"github.com/GreatValueCreamSoda/gometrics/video/comparator"
	"github.com/GreatValueCreamSoda/gometrics/video/metrics"
	"github.com/GreatValueCreamSoda/gometrics/video/sources"
	"github.com/schollz/progressbar/v3"
)

func main() {
	reference, err := sources.NewFFms2Reader(settings.referenceVideo)
	if err != nil {
		panic(err)
	}

	distortion, err := sources.NewFFms2Reader(settings.distortionVideo)
	if err != nil {
		panic(err)
	}

	var referenceColorSpace, distortionColorSpace vship.Colorspace
	referenceColorSpace.SetDefaults(0, 0, 0)
	distortionColorSpace.SetDefaults(0, 0, 0)

	if settings.compareHeight > 0 && settings.compareWidth > 0 {
		referenceColorSpace.TargetHeight = settings.compareHeight
		referenceColorSpace.TargetWidth = settings.compareWidth
		distortionColorSpace.TargetHeight = settings.compareHeight
		distortionColorSpace.TargetWidth = settings.compareWidth
	} else {
		referenceColorSpace.TargetHeight = settings.compareHeight
		referenceColorSpace.TargetWidth = settings.compareWidth
		distortionColorSpace.TargetHeight = settings.compareHeight
		distortionColorSpace.TargetWidth = settings.compareWidth

	}

	err = reference.GetColorProps().ToVsHipColorspace(&referenceColorSpace)
	if err != nil {
		panic(err)
	}

	err = distortion.GetColorProps().ToVsHipColorspace(&distortionColorSpace)
	if err != nil {
		panic(err)
	}

	if settings.frameRate < 0 {
		settings.frameRate = reference.GetFrameRate()
	}

	var metricHandlers []video.Metric
	var heatmapWriters []*metrics.HeatmapWriter

	for _, metric := range settings.metrics {
		metricHandler, heatmapWriter, err := createMetricAndWriter(
			metric, &referenceColorSpace, &distortionColorSpace)
		if err != nil {
			panic(err)
		}
		metricHandlers = append(metricHandlers, metricHandler)
		if heatmapWriter != nil {
			heatmapWriters = append(heatmapWriters, heatmapWriter)
		}
	}

	comp, err := comparator.NewComparator(
		reference, distortion, metricHandlers, settings.frameThreads,
		reference.GetNumFrames())
	if err != nil {
		panic(err)
	}

	bar := progressbar.NewOptions(
		reference.GetNumFrames(),
		progressbar.OptionSetDescription("Computing metrics"),
		progressbar.OptionShowCount(),
		progressbar.OptionShowIts(),
	)

	comp.SetProgressCallback(func(done, total int) {
		_ = bar.Add(1)
	})

	var scores map[string][]float64

	if scores, err = comp.Run(context.Background()); err != nil {
		panic(err)
	}

	for _, writer := range heatmapWriters {
		if err := writer.Close(); err != nil {
			log.Fatal("Failed to finalize video:", err)
		}
	}

	printSummary(scores)
}

func createMetricAndWriter(metricName string, ref, dist *vship.Colorspace) (
	video.Metric, *metrics.HeatmapWriter, error) {
	switch metricName {
	case metrics.ButteraugliName:
		return newButteraugli(ref, dist)
	case metrics.SSIMulacra2Name:
		return newSSIMULACRA2(ref, dist)
	case metrics.CVVDPName:
		return newCVVDP(ref, dist)
	default:
		return nil, nil, fmt.Errorf("unsupported metric: %s", metricName)
	}
}

func newCVVDP(ref, dist *vship.Colorspace) (video.Metric,
	*metrics.HeatmapWriter, error) {
	handler, err := metrics.NewCVVDPHandler(settings.frameThreads, ref, dist,
		settings.cvvdpUseTemporalScore, settings.cvvdpReizeToDisplay,
		settings.displayModel, settings.frameRate)
	if err != nil {
		return nil, nil, fmt.Errorf("cvvdp  creation failed: %w", err)
	}

	writer, err := createHeatmapWriterIfRequested(handler,
		settings.cvvdpDistMapPath, settings.cvvdpClipping)
	if err != nil {
		return nil, nil, err
	}

	return video.Metric(handler), writer, nil
}

func newSSIMULACRA2(ref, dist *vship.Colorspace) (video.Metric,
	*metrics.HeatmapWriter, error) {
	handler, err := metrics.NewSSIMU2Handler(settings.frameThreads, ref, dist)
	if err != nil {
		return nil, nil, fmt.Errorf("ssimulacra2 creation failed: %w", err)
	}

	return video.Metric(handler), nil, nil
}

func newButteraugli(ref, dist *vship.Colorspace) (video.Metric,
	*metrics.HeatmapWriter, error) {
	handler, err := metrics.NewButterHandler(settings.frameThreads, ref, dist,
		settings.butteraugliQnormValue,
		settings.displayModel.DisplayMaxLuminance,
	)
	if err != nil {
		return nil, nil, fmt.Errorf("butteraugli creation failed: %w", err)
	}

	writer, err := createHeatmapWriterIfRequested(handler,
		settings.butteraugliDistMapPath, settings.butteraugliClipping)
	if err != nil {
		return nil, nil, err
	}

	return video.Metric(handler), writer, nil
}

func createHeatmapWriterIfRequested(metric metrics.MetricWithDistortionMap,
	outputPath string, clipping float32) (*metrics.HeatmapWriter, error) {
	if outputPath == "" {
		return nil, nil
	}

	writer, err := metrics.WriteDistMapToVideo(metric, settings.frameRate,
		nil, outputPath, clipping)
	if err != nil {
		return nil, fmt.Errorf(
			"failed to create heatmap writer for %s: %w", outputPath, err)
	}

	return writer, nil
}
