package main

import (
	"context"
	"fmt"
	"log"

	"github.com/GreatValueCreamSoda/gometrics/comparator"
	"github.com/GreatValueCreamSoda/gometrics/metrics"
	"github.com/GreatValueCreamSoda/gometrics/sources"
	vship "github.com/GreatValueCreamSoda/govship"
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

	var metricHandlers []comparator.Metric
	var heatmapWriters []*metrics.HeatmapWriter

	for _, metric := range settings.metrics {
		metricHandler, heatmapWriter, err := createMetricAndWriter(
			metric, reference.GetColorspace(), distortion.GetColorspace())
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

	scores, err := comp.Run(context.Background())
	if err != nil {
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
	comparator.Metric, *metrics.HeatmapWriter, error) {
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

func newCVVDP(ref, dist *vship.Colorspace) (comparator.Metric,
	*metrics.HeatmapWriter, error) {
	handler, err := metrics.NewCVVDPHandler(settings.frameThreads, ref, dist,
		settings.cvvdpUseTemporalScore, settings.cvvdpReizeToDisplay,
		settings.displayModel, 15,
	)
	if err != nil {
		return nil, nil, fmt.Errorf("cvvdp  creation failed: %w", err)
	}

	writer, err := createHeatmapWriterIfRequested(handler,
		settings.cvvdpDistMapPath)
	if err != nil {
		return nil, nil, err
	}

	return comparator.Metric(handler), writer, nil
}

func newSSIMULACRA2(ref, dist *vship.Colorspace) (comparator.Metric,
	*metrics.HeatmapWriter, error) {
	handler, err := metrics.NewSSIMU2Handler(settings.frameThreads, ref, dist)
	if err != nil {
		return nil, nil, fmt.Errorf("ssimulacra2 creation failed: %w", err)
	}

	return comparator.Metric(handler), nil, nil
}

func newButteraugli(ref, dist *vship.Colorspace) (comparator.Metric,
	*metrics.HeatmapWriter, error) {
	handler, err := metrics.NewButterHandler(settings.frameThreads, ref, dist,
		settings.butteraugliQnormValue,
		settings.displayModel.DisplayMaxLuminance,
	)
	if err != nil {
		return nil, nil, fmt.Errorf("butteraugli creation failed: %w", err)
	}

	writer, err := createHeatmapWriterIfRequested(handler,
		settings.butteraugliDistMapPath)
	if err != nil {
		return nil, nil, err
	}

	return comparator.Metric(handler), writer, nil
}

func createHeatmapWriterIfRequested(metric metrics.MetricWithDistortionMap,
	outputPath string) (*metrics.HeatmapWriter, error) {
	if outputPath == "" {
		return nil, nil
	}

	writer, err := metrics.WriteDistMapToVideo(metric, 25, nil, outputPath,
		15)
	if err != nil {
		return nil, fmt.Errorf(
			"failed to create heatmap writer for %s: %w", outputPath, err)
	}

	return writer, nil
}
