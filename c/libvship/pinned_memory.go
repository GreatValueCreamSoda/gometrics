package libvship

/*
#include <VshipAPI.h>
#include <stdlib.h>
#include "flattened.h"
*/
import "C"
import (
	"unsafe"
)

// PinnedMalloc allocates a block of memory that is page-locked (pinned) in
// physical memory, suitable for high-performance DMA transfers or GPU access.
//
// The returned memory must be freed using PinnedFree().
//
// Parameters:
//   - size: The number of bytes to allocate
//
// Returns:
//   - []byte slice backed by the pinned memory (length and capacity = size)
//   - error if allocation failed
//
// Important:
//   - The returned slice should NOT be appended to (it would cause
//     reallocation)
//   - The memory remains pinned until PinnedFree is called
func PinnedMalloc(size int) ([]byte, ExceptionCode) {
	var ptr unsafe.Pointer
	code := ExceptionCode(C.Vship_PinnedMalloc(&ptr, C.uint64_t(size)))
	if !code.IsNone() {
		return nil, code
	}
	return unsafe.Slice((*byte)(ptr), size), code
}

// PinnedFree releases memory previously allocated with PinnedMalloc.
//
// Passing nil is safe and is a no-op.
//
// This function should be called when the pinned memory is no longer needed.
//
// Returns an error if the operation failed (rare).
func PinnedFree(data []byte) ExceptionCode {
	if len(data) == 0 {
		return ExceptionCodeNoError
	}

	ptr := unsafe.Pointer(&data[0])
	code := ExceptionCode(C.Vship_PinnedFree(ptr))
	if !code.IsNone() {
		return code
	}

	data = nil

	return code
}
