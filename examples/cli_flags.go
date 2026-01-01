package main

import (
	"fmt"
	"os"
	"slices"
	"strings"

	vship "github.com/GreatValueCreamSoda/gometrics/c/libvship"
	"github.com/GreatValueCreamSoda/gometrics/video/metrics"
	"github.com/spf13/pflag"
)

type cliSettings struct {
	referenceVideo, distortionVideo string
	metrics                         []string
	frameThreads                    int
	frameRate                       float32
	compareWidth, compareHeight     int

	butteraugliDistMapPath string
	butteraugliClipping    float32
	cvvdpDistMapPath       string
	cvvdpClipping          float32

	butteraugliQnormValue int

	cvvdpUseTemporalScore bool
	cvvdpReizeToDisplay   bool

	displayModel vship.DisplayModel
}

var settings cliSettings = cliSettings{
	displayModel: vship.DisplayModelPresetStandard4K,
}

func init() {
	pflag.CommandLine.SortFlags = false

	// General Flags
	pflag.StringVarP(&settings.referenceVideo, "reference", "r", "", "The reference video path the distorted video will be compared against")
	pflag.StringVarP(&settings.distortionVideo, "distortion", "d", "", "The distorted video path that will be compared to the reference")
	cliMetrics := pflag.String("metrics", metrics.SSIMulacra2Name, fmt.Sprintf("Comma seperated list of metrics that will be used [%s, %s, %s]", metrics.SSIMulacra2Name, metrics.ButteraugliName, metrics.CVVDPName))
	pflag.IntVar(&settings.frameThreads, "frame-threads", 3, "Number of frames to process in parallel. Set to 2 or 1 with 2 or more metrics")
	pflag.Float32VarP(&settings.frameRate, "fps", "f", -1, "Overide the fps that will be used for temporal scaling. Default is the reference fps")
	pflag.IntVar(&settings.compareWidth, "width", -1, "Overide the resolution to compare at width. -1 defaults to the largest source")
	pflag.IntVar(&settings.compareHeight, "height", -1, "Overide the resolution to compare at height. -1 defaults to the largest source")
	printHelp := pflag.BoolP("help", "h", false, "Show this help message")

	// Output Settings
	var outputsSectionString string = "Output Options"
	pflag.StringVar(&settings.butteraugliDistMapPath, "butteraugli-video-path", "", "Output path for Butterauglis heat map. Empty disables output")
	addFlagToHelpGroup("butteraugli-video-path", outputsSectionString)

	pflag.Float32Var(&settings.butteraugliClipping, "butteraugli-clipping-value", 15, "The clipping value for Butterauglis distortion map.")
	addFlagToHelpGroup("butteraugli-clipping-value", outputsSectionString)

	pflag.StringVar(&settings.cvvdpDistMapPath, "cvvdp-video-path", "", "Output path for CVVDPs heat map. Empty disables output")
	addFlagToHelpGroup("cvvdp-video-path", outputsSectionString)

	pflag.Float32Var(&settings.cvvdpClipping, "cvvdp-clipping-value", 0.75, "The clipping value for CVVDPs distortion map.")
	addFlagToHelpGroup("cvvdp-clipping-value", outputsSectionString)

	// butteraugli settings
	var butteraugliSectionName string = "Butteraugli Options"
	pflag.IntVar(&settings.butteraugliQnormValue, "butteraugli-qnorm", 5, "QNorm value to use for frame quality aggergation")
	addFlagToHelpGroup("butteraugli-qnorm", butteraugliSectionName)

	// CVVDP settings
	var cvvdpSectionName string = "CVVDP Options"
	pflag.BoolVar(&settings.cvvdpUseTemporalScore, "no-cvvdp-temporal", false, "Disable temporal motion for calculating frame scores")
	addFlagToHelpGroup("no-cvvdp-temporal", cvvdpSectionName)

	pflag.BoolVar(&settings.cvvdpReizeToDisplay, "no-resize-to-display", false, "Disable resizing videos to display models resolution")
	addFlagToHelpGroup("no-resize-to-display", cvvdpSectionName)

	// Display Model
	var displayModelSectionName string = "Display Model Options"
	pflag.Float32Var(&settings.displayModel.DisplayMaxLuminance, "display-nits", 203, "The target displays brightness in nits (Used by CVVDP and Butteraugli)")
	addFlagToHelpGroup("display-nits", displayModelSectionName)

	pflag.IntVar(&settings.displayModel.DisplayWidth, "display-width", 3840, "The target displays horizontal resolution in pixels (Used by CVVDP)")
	addFlagToHelpGroup("display-width", displayModelSectionName)

	pflag.IntVar(&settings.displayModel.DisplayHeight, "display-height", 2160, "The target displays vertical resolution in pixels (Used by CVVDP)")
	addFlagToHelpGroup("display-height", displayModelSectionName)

	pflag.Float32Var(&settings.displayModel.DisplayDiagonalSizeInches, "display-size", 32, "The target displays diagonal size in inches (Used by CVVDP)")
	addFlagToHelpGroup("display-size", displayModelSectionName)

	pflag.Float32Var(&settings.displayModel.ViewingDistanceMeters, "display-distance", 0.7472, "The target displays distance away from the viewer in meters (Used by CVVDP)")
	addFlagToHelpGroup("display-distance", displayModelSectionName)

	pflag.IntVar(&settings.displayModel.MonitorContrastRatio, "display-ratio", 10000, "The target displays contrast ratio (Used by CVVDP)")
	addFlagToHelpGroup("display-ratio", displayModelSectionName)

	pflag.IntVar(&settings.displayModel.AmbientLightLevel, "room-brightness", 250, "The rooms ambient lux the target display is in (Used by CVVDP)")
	addFlagToHelpGroup("room-brightness", displayModelSectionName)

	pflag.Parse()

	settings.cvvdpUseTemporalScore = !settings.cvvdpUseTemporalScore
	settings.cvvdpReizeToDisplay = !settings.cvvdpReizeToDisplay

	if *printHelp {
		cliUsage()
		os.Exit(0)
	}

	settings.metrics = strings.Split(*cliMetrics, ",")

	if settings.frameThreads > 1 && settings.cvvdpUseTemporalScore {
		var cvvdp bool = slices.Contains(settings.metrics, metrics.CVVDPName)
		if cvvdp {
			panic("cannot use more than 1 frame thread while using cvvdp with temporal weighting.")
		}
	}
}
