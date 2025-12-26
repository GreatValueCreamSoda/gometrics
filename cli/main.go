package main

import (
	"context"
	"log"

	"github.com/GreatValueCreamSoda/gometrics/comparator"
	"github.com/GreatValueCreamSoda/gometrics/metrics"
	"github.com/GreatValueCreamSoda/gometrics/sources"
)

func main() {
	sourceA, err := sources.NewFFms2Reader("/home/creamsoda/Desktop/sample.mkv")
	if err != nil {
		panic(err)
	}

	sourceB, err := sources.NewFFms2Reader("/home/creamsoda/Desktop/sample_dst.mkv")
	if err != nil {
		panic(err)
	}

	butter, err := metrics.NewButterHandler(1, sourceA.GetColorspace(), sourceB.GetColorspace(), 5, 203, true)
	if err != nil {
		panic(err)
	}

	writer, err := metrics.WriteDistMapToVideo(butter, 25, nil, "/home/creamsoda/Documents/git/gometrics/encode.mkv", 15)
	if err != nil {
		panic(err)
	}

	comp, err := comparator.NewComparator(sourceA, sourceB, []comparator.Metric{butter}, 1, sourceA.GetNumFrames())
	if err != nil {
		panic(err)
	}

	_, err = comp.Run(context.Background())
	if err != nil {
		panic(err)
	}

	if err := writer.Close(); err != nil {
		log.Fatal("Failed to finalize video:", err)
	}
}
