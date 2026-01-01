package video

import (
	"fmt"

	pixfmts "github.com/GreatValueCreamSoda/gometrics/c/libavpixfmts"
	vship "github.com/GreatValueCreamSoda/gometrics/c/libvship"
)

type ColorProperties struct {
	Width, Height  int
	PixelFormat    pixfmts.PixelFormat
	ColorRange     pixfmts.ColorRange
	ColorSpace     pixfmts.ColorSpace
	ColorTransfer  pixfmts.ColorTransferCharacteristic
	ColorPrimaries pixfmts.ColorPrimaries
	ChromaLocation pixfmts.ChromaLocation
}

func (cp *ColorProperties) ToVsHipColorspace(cs *vship.Colorspace) error {
	cs.Width, cs.Height = cp.Width, cp.Height

	pixFmtDesc, err := pixfmts.PixFmtDescGet(cp.PixelFormat)
	if err != nil {
		return err
	}

	comp, err := pixFmtDesc.Component(0)
	if err != nil {
		return err
	}

	var pixFmtSamplingFormat vship.SamplingFormat

	switch comp.Depth {
	case 8:
		pixFmtSamplingFormat = vship.SamplingFormatUInt8
	case 9:
		pixFmtSamplingFormat = vship.SamplingFormatUInt9
	case 10:
		pixFmtSamplingFormat = vship.SamplingFormatUInt10
	case 12:
		pixFmtSamplingFormat = vship.SamplingFormatUInt12
	case 14:
		pixFmtSamplingFormat = vship.SamplingFormatUInt14
	case 16:
		pixFmtSamplingFormat = vship.SamplingFormatUInt16
	default:
		return fmt.Errorf("unknown pixel format %s", pixFmtDesc.Name())
	}

	cs.SamplingFormat = pixFmtSamplingFormat

	switch cp.ColorRange {
	case pixfmts.ColorRangeMPEG:
		cs.ColorRange = vship.ColorRangeLimited
	case pixfmts.ColorRangeJPEG:
		cs.ColorRange = vship.ColorRangeFull
	default:
		// return fmt.Errorf("color range is not specified in source properties")
		cs.ColorRange = vship.ColorRangeLimited
	}

	cs.ChromaSubsamplingHeight = pixFmtDesc.Log2ChromaH()
	cs.ChromaSubsamplingWidth = pixFmtDesc.Log2ChromaW()

	switch cp.ChromaLocation {
	case pixfmts.ChromaLocationLeft:
		cs.ChromaLocation = vship.ChromaLocationLeft
	case pixfmts.ChromaLocationCenter:
		cs.ChromaLocation = vship.ChromaLocationCenter
	case pixfmts.ChromaLocationTopLeft:
		cs.ChromaLocation = vship.ChromaLocationTopLeft
	case pixfmts.ChromaLocationTop:
		cs.ChromaLocation = vship.ChromaLocationTop
	default:
		// return fmt.Errorf("chroma location in source props is not supported")
		cs.ChromaLocation = vship.ChromaLocationLeft
	}

	if pixFmtDesc.Flags()&uint64(pixfmts.PixFmtFlagRGB) == 0 {
		cs.ColorFamily = vship.ColorFamilyYUV
	} else {
		cs.ColorFamily = vship.ColorFamilyRGB
	}

	switch cp.ColorSpace {
	case pixfmts.ColorSpaceRGB:
		cs.ColorMatrix = vship.ColorMatrixRGB
	case pixfmts.ColorSpaceBT709:
		cs.ColorMatrix = vship.ColorMatrixBT709
	case pixfmts.ColorSpaceBT470BG:
		cs.ColorMatrix = vship.ColorMatrixBT470BG
	case pixfmts.ColorSpaceSMPTE170M:
		cs.ColorMatrix = vship.ColorMatrixST170M
	case pixfmts.ColorSpaceBT2020_NCL:
		cs.ColorMatrix = vship.ColorMatrixBT2020NCL
	case pixfmts.ColorSpaceBT2020_CL:
		cs.ColorMatrix = vship.ColorMatrixBT2020CL
	case pixfmts.ColorSpaceICTCP:
		cs.ColorMatrix = vship.ColorMatrixBT2100ICTCP
	default:
		// return fmt.Errorf("chroma matrix in source propeties is not supported")
		cs.ColorMatrix = vship.ColorMatrixBT709
	}

	switch cp.ColorTransfer {
	case pixfmts.ColorTransferCharacteristicBT709:
		cs.ColorTransfer = vship.ColorTransferTRCBT709
	case pixfmts.ColorTransferCharacteristicGamma22:
		cs.ColorTransfer = vship.ColorTransferTRCBT470_M
	case pixfmts.ColorTransferCharacteristicGamma28:
		cs.ColorTransfer = vship.ColorTransferTRCBT470_BG
	case pixfmts.ColorTransferCharacteristicSMPTE170M:
		cs.ColorTransfer = vship.ColorTransferTRCBT601
	case pixfmts.ColorTransferCharacteristicLinear:
		cs.ColorTransfer = vship.ColorTransferTRCLinear
	case pixfmts.ColorTransferCharacteristicIEC61966_2_1:
		cs.ColorTransfer = vship.ColorTransferTRCSRGB
	case pixfmts.ColorTransferCharacteristicSMPTE2084:
		cs.ColorTransfer = vship.ColorTransferTRCPQ
	case pixfmts.ColorTransferCharacteristicSMPTE428:
		cs.ColorTransfer = vship.ColorTransferTRCST428
	case pixfmts.ColorTransferCharacteristicARIB_STD_B67:
		cs.ColorTransfer = vship.ColorTransferTRCHLG
	default:
		// return fmt.Errorf("chroma transfer in source props is not supported")
		cs.ColorTransfer = vship.ColorTransferTRCBT709
	}

	switch cp.ColorPrimaries {
	case pixfmts.ColorPrimariesBT709:
		cs.ColorPrimaries = vship.ColorPrimariesBT709
	case pixfmts.ColorPrimariesBT470M:
		cs.ColorPrimaries = vship.ColorPrimariesBT470_M
	case pixfmts.ColorPrimariesBT470BG:
		cs.ColorPrimaries = vship.ColorPrimariesBT470_BG
	case pixfmts.ColorPrimariesBT2020:
		cs.ColorPrimaries = vship.ColorPrimariesBT2020
	default:
		// return fmt.Errorf("chroma primaries in source props is not supported")
		cs.ColorPrimaries = vship.ColorPrimariesBT709
	}

	return nil
}
