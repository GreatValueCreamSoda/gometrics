package libffms2

//#cgo LDFLAGS: -lffms2
//#cgo CFLAGS: -I/usr/include
//#include <ffms.h>
//#include <stdlib.h>
//#include "indexer/indexer.h"
import "C"
import (
	"errors"
	"unsafe"
)

var (
	ErrInvalidorNilIndexer error = errors.New("indexer was consumed, failed to create, or was destroyed")
)

var callbackMap map[uintptr]IndexerCallbackFunction = make(
	map[uintptr]IndexerCallbackFunction)

// See Indexer.SetProgressCallback for how the callback works.
type IndexerCallbackFunction func(current, total int64) int

// Private hidden function to help manage removing a indexer callback from
// callbackMap once an indexer is closed or destroyed.
func (i *Indexer) removeCallback() {
	var ICPrivate uintptr = (uintptr)(unsafe.Pointer(i))
	delete(callbackMap, (uintptr)(ICPrivate))
}

// Private method called by C to call back into GO! to execute the Indexers
// IndexerCallbackFunction.

//export goIndexCallback
func goIndexCallback(current, total C.int64_t, ICPriv unsafe.Pointer) C.int {
	if fn, ok := callbackMap[(uintptr)(ICPriv)]; ok {
		return C.int(fn(int64(current), int64(total)))
	}
	return 0
}

type Indexer struct {
	indexer *C.FFMS_Indexer
}

// Creates a Indexer object for the given SourceFile and returns a pointer to
// it. Is basically a shorthand for CreateIndexerWithDemuxer(SourceFile,
// FFMS_SOURCE_DEFAULT, ErrorInfo).
func CreateIndexer(sourceFile string) (*Indexer, *ErrorInfo, error) {
	var sourceFileC *C.char = (*C.char)(C.CString(sourceFile))
	defer C.free(unsafe.Pointer(sourceFileC))

	res, errorInfo, err := withErrorInfo(func(c *C.FFMS_ErrorInfo) *C.FFMS_Indexer {
		return C.FFMS_CreateIndexer(sourceFileC, c)
	})
	if err != nil {
		return nil, errorInfo, err
	}

	return &Indexer{res}, errorInfo, nil
}

// Returns the total number of tracks in the media file represented by the
// given Indexer. In other words, does the same thing as GetNumTracks but does
// not require indexing the entire file first.
func (i *Indexer) GetNumTracks() (int, error) {
	if err := i.checkValidity(); err != nil {
		return 0, err
	}

	var numTracks int = int(C.FFMS_GetNumTracksI(i.indexer))

	return numTracks, nil
}

// Returns an integer representing the FFMS_TrackType of the track number Track
// in the media file represented by the given Indexer. In other words, does the
// same thing as GetTrackType, but does not require having indexed the file
// first. If you have indexed the file, use GetTrackType instead since the
// Indexer object is destructed when the index is created. Note that specifying
// an invalid track number may lead to undefined behavior.
func (i *Indexer) GetTrackType(track TrackType) (int, error) {
	if err := i.checkValidity(); err != nil {
		return 0, err
	}

	var trackNum int = int(C.FFMS_GetTrackTypeI(i.indexer, C.int(track)))

	return trackNum, nil
}

// Returns the human-readable name ("long name" in FFmpeg terms) of the codec
// used in the given track number in the media file represented by the given
// Indexer object. Useful if you want to, say, pop up a menu asking the user
// which tracks he or she wishes to index.
//
// Note: Specifying an invalid track number may lead to undefined behavior.
func (i *Indexer) GetCodecName(track int) (string, error) {
	if err := i.checkValidity(); err != nil {
		return "", err
	}

	var codecName *C.char = C.FFMS_GetCodecNameI(i.indexer, C.int(track))
	if codecName == nil {
		return "", errors.New("failed to get codec name of supplied track")
	}

	return C.GoString(codecName), nil
}

// Returns the human-readable name ("long name" in FFmpeg terms) of the
// container format used by the file represented by the given Indexer.
func (i *Indexer) GetFormatName() (string, error) {
	if err := i.checkValidity(); err != nil {
		return "", err
	}

	var formatName *C.char = C.FFMS_GetFormatNameI(i.indexer)
	if formatName == nil {
		return "", errors.New("failed to get format name of the source file")
	}

	return C.GoString(formatName), nil
}

// If you supply a progress callback, FFMS2 will call it regularly during
// indexing to report progress and give you the chance to interrupt indexing.
//
// The callback function's arguments are as follows:
//
//	current	int64 - indexing progress amount done in bytes
//	total 	int64 - indexing progress total in bytes
//
// Return 0 from the callback function to continue indexing, non-0 to cancel
// indexing (returning non-0 will make DoIndexing fail with the reason
// "indexing cancelled by user").
func (i *Indexer) SetProgressCallback(fn IndexerCallbackFunction) error {
	if err := i.checkValidity(); err != nil {
		return err
	}

	// GO!'s C header parser is unable to parse the `TIndexCallback` from the
	// ffms.h header. So this weird type reflection is required.

	var ICPrivate unsafe.Pointer = unsafe.Pointer(i)
	var callback *[0]byte = (*[0]byte)(C.cIndexingCallback)

	// Registers the c wrapper function for GO! callbacks with the FFMS
	// object.
	C.FFMS_SetProgressCallback(i.indexer, callback, ICPrivate)

	// Stores the callback for the c wrapper functions to call into later via
	// goIndexCallback.
	callbackMap[(uintptr)(ICPrivate)] = fn
	return nil
}

// Runs the passed indexer and returns a Index object representing the file in
// question.
//
// Note: Calling this function is equivalent to calling Indexer.Close().
func (i *Indexer) DoIndexing(errorHandling IndexErrorHandling) (*Index,
	*ErrorInfo, error) {
	if err := i.checkValidity(); err != nil {
		return nil, nil, err
	}

	// FFMS_DoIndexing2 always destorys the Indexer no matter the result. Mark
	// it as invalid as soon as possible.
	res, info, err := withErrorInfo(func(c *C.FFMS_ErrorInfo) *C.FFMS_Index {
		return C.FFMS_DoIndexing2(i.indexer, C.int(errorHandling), c)
	})
	i.indexer = nil // invalid

	if err != nil {
		return nil, info, err
	}

	return &Index{res}, info, nil
}

// checkValidity simply checks if the c ptr to the wrapped *C.FFMS_Indexer is
// nil or not. Any other checks that need to be preformed before the type can
// be used should be added here.
func (i *Indexer) checkValidity() error {
	if i.indexer == nil {
		return ErrInvalidorNilIndexer
	}

	return nil
}

// Destroys the Indexer object if it still exists. Invalidates any further
// usage of the Indexer.
//
// Note: This must be called to avoid memory leaks as the Indexer exists within
// C allocated memory. Therefore it will not be automatically cleaned up by GO!
// once the object leaves scope. (Nor does GO! ever guarentee any finalizer
// will ever be called). If Indexer.DoIndexing() was called this function is
// not required.
func (i *Indexer) Close() {

	if i.indexer != nil {
		C.FFMS_CancelIndexing(i.indexer)
		i.indexer = nil
	}

	i.removeCallback()
}
