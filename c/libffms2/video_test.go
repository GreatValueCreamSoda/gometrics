package libffms2_test

import (
	"testing"

	ffms "github.com/GreatValueCreamSoda/gometrics/c/libffms2"
)

func Test_VideoSource(t *testing.T) {
	indexer, _, err := ffms.CreateIndexer("./samples/sample.mkv")
	if err != nil {
		t.FailNow()
	}

	index, _, err := indexer.DoIndexing(ffms.IEHAbort)
	if err != nil {
		t.FailNow()
	}

	trackNum, _, err := index.GetFirstTrackOfType(ffms.TypeVideo)
	if err != nil {
		t.FailNow()
	}

	video, _, err := ffms.CreateVideoSource("./samples/sample.mkv", index,
		trackNum, 16, ffms.SeekLinear)
	if err != nil {
		t.FailNow()
	}

	props, err := video.GetVideoProperties()
	if err != nil {
		t.FailNow()
	}

	t.Log(props.FPSNumerator)

	frame, _, err := video.GetFrame(0)
	if err != nil {
		t.FailNow()
	}

	t.Log(len(frame.Data[0]))
	t.Log(len(frame.Data[1]))
	t.Log(len(frame.Data[2]))
}
