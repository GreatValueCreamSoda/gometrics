package libavpixfmts_test

import (
	"testing"

	pixfmts "github.com/GreatValueCreamSoda/gometrics/c/libavpixfmts"
)

func Test_PixFmtDescGet(t *testing.T) {
	yuv420p10le, err := pixfmts.PixFmtDescGet(pixfmts.PixFmtYUV420P10LE)
	if err != nil {
		t.FailNow()
	}

	t.Log(yuv420p10le)
}

func Test_PixFmtDescRef_PixFmtDescID(t *testing.T) {
	yuv420p10le, err := pixfmts.PixFmtDescGet(pixfmts.PixFmtYUV420P10LE)
	if err != nil {
		t.FailNow()
	}

	id, err := yuv420p10le.PixFmtDescID()
	if err != nil {
		t.FailNow()
	}

	t.Log(id)
}

func Test_PixFmtDescRef_Name(t *testing.T) {
	yuv420p10le, err := pixfmts.PixFmtDescGet(pixfmts.PixFmtYUV420P10LE)
	if err != nil {
		t.FailNow()
	}

	name := yuv420p10le.Name()
	if name == "" {
		t.FailNow()
	}

	t.Log(name)
}

func Test_PixFmtDescRef_NbComponents(t *testing.T) {
	yuv420p10le, err := pixfmts.PixFmtDescGet(pixfmts.PixFmtYUV420P10LE)
	if err != nil {
		t.FailNow()
	}

	numComonents := yuv420p10le.NbComponents()
	if numComonents == 0 {
		t.FailNow()
	}

	t.Log(numComonents)
}
