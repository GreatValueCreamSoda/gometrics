package libffms2

//#cgo LDFLAGS: -lffms2
//#cgo CFLAGS: -I/usr/include
//#include <ffms.h>
//#include <stdlib.h>
import "C"
import (
	"errors"
	"reflect"
	"unsafe"
)

var (
	ErrMallocNilPtrReturn = errors.New("malloc returned a nil pointer. are we out of ram?")
	ErrNilPtrFree         = errors.New("attempted to free a nil ptr with safeFree")
	ErrFFmsNilPtrReturn   = errors.New("ffms returned a nill pointer. check error info for more info")
)

func init()                      { C.FFMS_Init(0, 0) }
func GetVersion() int            { return int(C.FFMS_GetVersion()) }
func GetLogLevel() LogLevel      { return LogLevel(C.FFMS_GetLogLevel()) }
func SetLogLevel(level LogLevel) { C.FFMS_SetLogLevel(C.int(level)) }

// safeMalloc is a wrapper around the C-level malloc function.
// The size of element to allocate is passed as type T and numElements
// represents the number of elements of size T to allocate.
//
// Checks if malloc returned a nil pointer and returns an error.
func safeMalloc[T any](numElements uint) (*T, error) {
	var size T
	var bytes C.size_t = C.size_t(unsafe.Sizeof(size) * uintptr(numElements))
	var ptr unsafe.Pointer = C.malloc(bytes)
	var err error = nil
	if ptr == nil {
		err = ErrMallocNilPtrReturn
	}
	return (*T)(ptr), err
}

// safeFree is a wrapper around the C-level free function. It checks for nil
// pointers and will return a error if any are passed.
func safeFree[T any](ptr *T) error {
	if ptr != nil {
		C.free(unsafe.Pointer(ptr))
		return nil
	}

	return ErrNilPtrFree
}

// Error Info

// withErrorInfo wraps a call to a function that accepts a C-level
// FFMS_ErrorInfo structure and returns a value of type T. It allocates and
// initializes the underlying C error info buffers, passes them to the provided
// function, and then converts the resulting C error information into a Go
// ErrorInfo.
//
// Memory allocated for the C error structures is always released before the
// function returns.
//
// If T is a pointer, which is detected using runtime reflection, it is checked
// for nilness. A nil pointer is considered an invalid result and will result
// in a return with a ErrFFmsNilPtrReturn along with a possible ErrorInfo.
//
// ErrorInfo is returned as nil when the wrapper must exit before the
// underlying C function is invoked. Only when the C function executes long
// enough to fill the FFMS_ErrorInfo structure will a non-nil ErrorInfo be
// produced.
func withErrorInfo[T any](fn func(*C.FFMS_ErrorInfo) T) (T, *ErrorInfo, error) {
	// zeroed version of the value if we must return before calling fn
	var zeroRet T

	// malloc memory with c instead of passing go ptrs to avoid issues with
	// memmory pinning
	cErrInfo, err := safeMalloc[C.FFMS_ErrorInfo](1)
	if err != nil {
		return zeroRet, nil, err
	}

	const bufferSize uint = 1024

	cErrInfo.Buffer, err = safeMalloc[C.char](bufferSize)
	cErrInfo.BufferSize = C.int(bufferSize)

	defer func() {
		safeFree(cErrInfo.Buffer)
		safeFree(cErrInfo)
	}()

	var resolvedValue T = fn(cErrInfo)

	goErrInfo := &ErrorInfo{int(cErrInfo.ErrorType), int(cErrInfo.SubType),
		C.GoString(cErrInfo.Buffer)}

	// Use runtime reflection to check if the returned type of T is a pointer.
	// If it is check if it is null as in 99% of use cases including ours,
	// that's bad.
	if reflect.TypeOf(resolvedValue).Kind() == reflect.Pointer {
		if reflect.ValueOf(resolvedValue).IsNil() {
			err = ErrFFmsNilPtrReturn
		}
	}

	// We still want to return the resolvedValue to the user as in some cases
	// it can help idenity what went wrong that we cant identify outselves.
	return resolvedValue, goErrInfo, err
}

func sliceFromCPtr[J, K any](ptr *J, length uint) []K {
	header := unsafe.Slice(ptr, int(length))
	return *(*[]K)(unsafe.Pointer(&header))
}
