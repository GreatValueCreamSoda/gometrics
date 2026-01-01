package libffms2

//#cgo LDFLAGS: -lffms2
//#cgo CFLAGS: -I/usr/include
//#include <ffms.h>
//#include <stdlib.h>
import "C"
import (
	"errors"
	"unsafe"
)

type Index struct {
	index *C.FFMS_Index
}

var (
	ErrInvalidOrNilIndex error = errors.New("index was consumed, failed to create, or was destroyed")
)

// CreateIndex creates an Index from a C.FFMS_Index pointer
func newIndexFromIndexPtr(indexPtr *C.FFMS_Index) *Index {
	return &Index{index: indexPtr}
}

// Returns the total number of tracks in the media file represented by the
// given Index.
func (idx *Index) GetNumTracks() (int, error) {
	if idx.index == nil {
		return 0, ErrInvalidOrNilIndex
	}
	numTracks := int(C.FFMS_GetNumTracks(idx.index))
	if numTracks < 0 {
		return 0, errors.New("failed to get number of tracks")
	}
	return numTracks, nil
}

// Finds the first track of the given TrackType in the given Index and returns
// its track number, suitable for use as an argument to CreateVideoSource or
// CreateAudioSource, as well as to some other functions.
//
// trackType represents the type of track to look for. See the TrackType enum.
//
// Returns the track number (an integer greater than or equal to 0) on success.
// Returns a negative integer and sets ErrorMsg on failure (i.e. if no track of
// the given type was found).
func (idx *Index) GetFirstTrackOfType(trackType TrackType) (int, *ErrorInfo,
	error) {
	if idx.index == nil {
		return 0, nil, ErrInvalidOrNilIndex
	}

	res, errorInfo, err := withErrorInfo(func(c *C.FFMS_ErrorInfo) C.int {
		return C.FFMS_GetFirstTrackOfType(idx.index,
			C.int(trackType), c)
	})

	return int(res), errorInfo, err
}

// Does the exact same thing as GetFirstTrackOfType but ignores tracks that
// have not been indexed.
func (idx *Index) GetFirstIndexedTrackOfType(trackType TrackType) (int,
	*ErrorInfo, error) {
	if idx.index == nil {
		return 0, nil, ErrInvalidOrNilIndex
	}

	res, errorInfo, err := withErrorInfo(func(c *C.FFMS_ErrorInfo) C.int {
		return C.FFMS_GetFirstIndexedTrackOfType(idx.index,
			C.int(trackType), c)
	})

	return int(res), errorInfo, err
}

// Gets which error handling mode was used when creating the given index.
//
// Returns the value of the ErrorHandling parameter which was passed to
// DoIndexing.
func (idx *Index) GetErrorHandling() (IndexErrorHandling, error) {
	if idx.index == nil {
		return 0, ErrInvalidOrNilIndex
	}

	return IndexErrorHandling(C.FFMS_GetErrorHandling(idx.index)), nil
}

// Gets track data for the given track number from the given Index object,
// stores it in a Track object and returns it. Use this function if you don't
// want to (or cannot) open the track with CreateVideoSource or
// CreateAudioSource first. If you already have a VideoSource or AudioSource
// object it's safer to use GetTrackFromVideo or GetTrackFromAudio instead.
// Note that specifying a nonexistent or invalid track number leads to
// undefined behavior (usually an access violation). Also note that the
// returned Track object is only valid until its parent Index object is
// destroyed.
//
// track represents the index of the desired track to be returned.
//
// Returns the Track on success. Note that requesting indexing information for
// a track that has not been indexed will not cause an error, it will just
// return an empty Track (check for >0 frames using GetNumFrames to see if the
// returned object actually contains indexing information).
func (idx *Index) GetTrack(track int) (Track, error) {
	if idx.index == nil {
		return Track{}, ErrInvalidOrNilIndex
	}

	var ptr *C.FFMS_Track = C.FFMS_GetTrackFromIndex(idx.index, C.int(track))
	if ptr == nil {
		return Track{}, errors.New("failed to get specified track from index")
	}

	return Track{ptr}, nil
}

// Makes a heuristic (but very reliable) guess about whether the given Index is
// an index of the given SourceFile or not. Useful to determine if the index
// object you just read with ReadIndex is actually relevant to your interests,
// since the only two ways to pair up index files with source files are a)
// trust the user blindly, or b) comparing the filenames; neither is very
// reliable.
//
// sourceFile represents the path to the file that the given Index will be
// compared against.
//
// Returns 0 if the given index is determined to belong to the given file.
// Returns non-0 and sets ErrorMsg otherwise.
func (idx *Index) BelongsToFile(sourceFile string) (int, *ErrorInfo, error) {
	if idx.index == nil {
		return 0, nil, ErrInvalidOrNilIndex
	}

	var sourceFileC *C.char = (*C.char)(C.CString(sourceFile))
	defer safeFree(sourceFileC)

	res, errorInfo, err := withErrorInfo(func(c *C.FFMS_ErrorInfo) C.int {
		return C.FFMS_IndexBelongsToFile(idx.index, sourceFileC, c)
	})

	return int(res), errorInfo, err
}

// Writes the indexing information from the given Index to the given IndexFile
// (which can be an absolute or relative path; it will be truncated and
// overwritten if it already exists).
//
// IndexFile represents the path that the given Index will be written to.
//
// Returns 0 on success; returns non-0 and sets ErrorMsg on failure.
func (idx *Index) WriteIndex(IndexFile string) (int, *ErrorInfo, error) {
	if err := idx.checkValidity(); err != nil {
		return 0, nil, err
	}

	var IndexFileC *C.char = (*C.char)(C.CString(IndexFile))
	defer safeFree(IndexFileC)

	res, errorInfo, err := withErrorInfo(func(c *C.FFMS_ErrorInfo) C.int {
		return C.FFMS_IndexBelongsToFile(idx.index, IndexFileC, c)
	})

	return int(res), errorInfo, err
}

// Writes the indexing information from the given Index to memory.
//
// Returns 0 on success; returns non-0 and sets ErrorMsg on failure.
func (idx *Index) WriteIndexToByteBuffer() ([]byte, int, *ErrorInfo, error) {
	if err := idx.checkValidity(); err != nil {
		return nil, 0, nil, err
	}

	var buffPtr *C.uint8_t
	var size C.size_t

	res, errorInfo, err := withErrorInfo(func(c *C.FFMS_ErrorInfo) C.int {
		return C.FFMS_WriteIndexToBuffer(&buffPtr, &size, idx.index, c)
	})
	defer func() {
		if buffPtr != nil {
			C.FFMS_FreeIndexBuffer(&buffPtr)
		}
	}()

	return C.GoBytes(unsafe.Pointer(buffPtr), C.int((C.size_t)(unsafe.Sizeof(
		C.uint8_t(0)))*size)), int(res), errorInfo, err
}

// checkValidity simply checks if the c ptr to the wrapped *C.FFMS_Index is nil
// or not. Any other checks that need to be preformed before the type can be
// used should be added here.
func (idx *Index) checkValidity() error {
	var err error
	if idx.index == nil {
		err = ErrInvalidOrNilIndex
	}
	return err
}

// Destroys the index object if it still exists. Invalidates any further usage
// of the index.
//
// Note: This must be called to avoid memory leaks as the index exists within C
// allocated memory. Therefore it will not be automatically cleaned up by GO!
// once the object leaves scope. (Nor does GO! ever guarentee any finalizer
// will ever be called).
func (idx *Index) Close() error {
	if err := idx.checkValidity(); err != nil {
		return err
	}

	C.FFMS_DestroyIndex(idx.index)
	idx.index = nil

	return nil
}
