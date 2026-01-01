package libvship_test

import (
	"testing"

	vship "github.com/GreatValueCreamSoda/gometrics/c/libvship"
)

func Test_ExceptionCode_IsNone(t *testing.T) {
	if !vship.ExceptionCodeNoError.IsNone() {
		t.Fatal("ExceptionCodeNoError should report IsNone() == true")
	}

	if vship.ExceptionCodeBadHandler.IsNone() {
		t.Fatal("non-zero ExceptionCode should report IsNone() == false")
	}
}
