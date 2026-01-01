package libavpixfmts

// #include <libavutil/pixfmt.h>
import "C"

const (
	PaletteSiz     int = C.AVPALETTE_SIZE
	PaletteCount   int = C.AVPALETTE_COUNT
	VideoMaxPlanes int = C.AV_VIDEO_MAX_PLANES
)

type PixelFormat int

const (
	PixFmtNone           PixelFormat = C.AV_PIX_FMT_NONE           // planar YUV 4:2:0, 12bpp, (1 Cr & Cb sample per 2x2 Y samples)
	PixFmtYUV420P        PixelFormat = C.AV_PIX_FMT_YUV420P        // packed YUV 4:2:2, 16bpp, Y0 Cb Y1 Cr
	PixFmtYUV422         PixelFormat = C.AV_PIX_FMT_YUYV422        // packed YUV 4:2:2, 16bpp, Y0 Cb Y1 Cr
	PixFmtRGB24          PixelFormat = C.AV_PIX_FMT_RGB24          // packed RGB 8:8:8, 24bpp, RGBRGB...
	PixFmtBGR24          PixelFormat = C.AV_PIX_FMT_BGR24          // packed RGB 8:8:8, 24bpp, BGRBGR...
	PixFmtYUV422P        PixelFormat = C.AV_PIX_FMT_YUV422P        // planar YUV 4:2:2, 16bpp, (1 Cr & Cb sample per 2x1 Y samples)
	PixFmtYUV444P        PixelFormat = C.AV_PIX_FMT_YUV444P        // planar YUV 4:4:4, 24bpp, (1 Cr & Cb sample per 1x1 Y samples)
	PixFmtYUV410P        PixelFormat = C.AV_PIX_FMT_YUV410P        // planar YUV 4:1:0,  9bpp, (1 Cr & Cb sample per 4x4 Y samples)
	PixFmtYUV411P        PixelFormat = C.AV_PIX_FMT_YUV411P        // planar YUV 4:1:1, 12bpp, (1 Cr & Cb sample per 4x1 Y samples)
	PixFmtGray8          PixelFormat = C.AV_PIX_FMT_GRAY8          //        Y        ,  8bpp
	PixFmtMonoWhite      PixelFormat = C.AV_PIX_FMT_MONOWHITE      //        Y        ,  1bpp, 0 is white, 1 is black, in each byte pixels are ordered from the msb to the lsb
	PixFmtMonoBlack      PixelFormat = C.AV_PIX_FMT_MONOBLACK      //        Y        ,  1bpp, 0 is black, 1 is white, in each byte pixels are ordered from the msb to the lsb
	PixFmtPal8           PixelFormat = C.AV_PIX_FMT_PAL8           // 8 bits with AV_PIX_FMT_RGB32 palette
	PixFmtYUVJ420P       PixelFormat = C.AV_PIX_FMT_YUVJ420P       // planar YUV 4:2:0, 12bpp, full scale (JPEG), deprecated in favor of AV_PIX_FMT_YUV420P and setting color_range
	PixFmtYUVJ422        PixelFormat = C.AV_PIX_FMT_YUVJ422P       // planar YUV 4:2:2, 16bpp, full scale (JPEG), deprecated in favor of AV_PIX_FMT_YUV422P and setting color_range
	PixFmtYUVJ444P       PixelFormat = C.AV_PIX_FMT_YUVJ444P       // planar YUV 4:4:4, 24bpp, full scale (JPEG), deprecated in favor of AV_PIX_FMT_YUV444P and setting color_range
	PixFmUYVY422         PixelFormat = C.AV_PIX_FMT_UYVY422        // packed YUV 4:2:2, 16bpp, Cb Y0 Cr Y1
	PixFmtUYYVYY411      PixelFormat = C.AV_PIX_FMT_UYYVYY411      // packed YUV 4:1:1, 12bpp, Cb Y0 Y1 Cr Y2 Y3
	PixFmtBGR8           PixelFormat = C.AV_PIX_FMT_BGR8           // packed RGB 3:3:2,  8bpp, (msb)2B 3G 3R(lsb)
	PixFmtBGR4           PixelFormat = C.AV_PIX_FMT_BGR4           // packed RGB 1:2:1 bitstream,  4bpp, (msb)1B 2G 1R(lsb), a byte contains two pixels, the first pixel in the byte is the one composed by the 4 msb bits
	PixFmtBGR4_BYTE      PixelFormat = C.AV_PIX_FMT_BGR4_BYTE      // packed RGB 1:2:1,  8bpp, (msb)1B 2G 1R(lsb)
	PixFmtRGB8           PixelFormat = C.AV_PIX_FMT_RGB8           // packed RGB 3:3:2,  8bpp, (msb)3R 3G 2B(lsb)
	PixFmtRGB4           PixelFormat = C.AV_PIX_FMT_RGB4           // packed RGB 1:2:1 bitstream,  4bpp, (msb)1R 2G 1B(lsb), a byte contains two pixels, the first pixel in the byte is the one composed by the 4 msb bits
	PixFmtRGB4_BYTE      PixelFormat = C.AV_PIX_FMT_RGB4_BYTE      // packed RGB 1:2:1,  8bpp, (msb)1R 2G 1B(lsb)
	PixFmtNV12           PixelFormat = C.AV_PIX_FMT_NV12           // planar YUV 4:2:0, 12bpp, 1 plane for Y and 1 plane for the UV components, which are interleaved (first byte U and the following byte V)
	PixFmtNV21           PixelFormat = C.AV_PIX_FMT_NV21           // as above, but U and V bytes are swapped
	PixFmtARGB           PixelFormat = C.AV_PIX_FMT_ARGB           // packed ARGB 8:8:8:8, 32bpp, ARGBARGB...
	PixFmtRGBA           PixelFormat = C.AV_PIX_FMT_RGBA           // packed RGBA 8:8:8:8, 32bpp, RGBARGBA...
	PixFmtARGR           PixelFormat = C.AV_PIX_FMT_ABGR           // packed ABGR 8:8:8:8, 32bpp, ABGRABGR...
	PixFmtBGRA           PixelFormat = C.AV_PIX_FMT_BGRA           // packed BGRA 8:8:8:8, 32bpp, BGRABGRA...
	PixFmtGray16BE       PixelFormat = C.AV_PIX_FMT_GRAY16BE       //        Y        , 16bpp, big-endian
	PixFmtGray16LE       PixelFormat = C.AV_PIX_FMT_GRAY16LE       //        Y        , 16bpp, little-endian
	PixFmtYUV440P        PixelFormat = C.AV_PIX_FMT_YUV440P        // planar YUV 4:4:0 (1 Cr & Cb sample per 1x2 Y samples)
	PixFmtYUVJ440P       PixelFormat = C.AV_PIX_FMT_YUVJ440P       // planar YUV 4:4:0 full scale (JPEG), deprecated in favor of AV_PIX_FMT_YUV440P and setting color_range
	PixFmtYUVA420P       PixelFormat = C.AV_PIX_FMT_YUVA420P       // planar YUV 4:2:0, 20bpp, (1 Cr & Cb sample per 2x2 Y & A samples)
	PixFmtRGB48BE        PixelFormat = C.AV_PIX_FMT_RGB48BE        // packed RGB 16:16:16, 48bpp, 16R, 16G, 16B, the 2-byte value for each R/G/B component is stored as big-endian
	PixFmtRGB48LE        PixelFormat = C.AV_PIX_FMT_RGB48LE        // packed RGB 16:16:16, 48bpp, 16R, 16G, 16B, the 2-byte value for each R/G/B component is stored as little-endian
	PixFmtRGB565BE       PixelFormat = C.AV_PIX_FMT_RGB565BE       // packed RGB 5:6:5, 16bpp, (msb)   5R 6G 5B(lsb), big-endian
	PixFmtRGB565LE       PixelFormat = C.AV_PIX_FMT_RGB565LE       // packed RGB 5:6:5, 16bpp, (msb)   5R 6G 5B(lsb), little-endian
	PixFmtRGB555BE       PixelFormat = C.AV_PIX_FMT_RGB555BE       // packed RGB 5:5:5, 16bpp, (msb)1X 5R 5G 5B(lsb), big-endian   , X=unused/undefined
	PixFmtRGB555LE       PixelFormat = C.AV_PIX_FMT_RGB555LE       // packed RGB 5:5:5, 16bpp, (msb)1X 5R 5G 5B(lsb), little-endian, X=unused/undefined
	PixFmtBGR565BE       PixelFormat = C.AV_PIX_FMT_BGR565BE       // packed BGR 5:6:5, 16bpp, (msb)   5B 6G 5R(lsb), big-endian
	PixFmtBGR565LE       PixelFormat = C.AV_PIX_FMT_BGR565LE       // packed BGR 5:6:5, 16bpp, (msb)   5B 6G 5R(lsb), little-endian
	PixFmtBGR555BE       PixelFormat = C.AV_PIX_FMT_BGR555BE       // packed BGR 5:5:5, 16bpp, (msb)1X 5B 5G 5R(lsb), big-endian   , X=unused/undefined
	PixFmtBGR555LE       PixelFormat = C.AV_PIX_FMT_BGR555LE       // packed BGR 5:5:5, 16bpp, (msb)1X 5B 5G 5R(lsb), little-endian, X=unused/undefined
	PixFmtVAAPI          PixelFormat = C.AV_PIX_FMT_VAAPI          // Hardware acceleration through VA-API, data[3] contains a VASurfaceID.
	PixFmtYUV420P16LE    PixelFormat = C.AV_PIX_FMT_YUV420P16LE    // planar YUV 4:2:0, 24bpp, (1 Cr & Cb sample per 2x2 Y samples), little-endian
	PixFmtYUV420P16BE    PixelFormat = C.AV_PIX_FMT_YUV420P16BE    // planar YUV 4:2:0, 24bpp, (1 Cr & Cb sample per 2x2 Y samples), big-endian
	PixFmtYUV422P16LE    PixelFormat = C.AV_PIX_FMT_YUV422P16LE    // planar YUV 4:2:2, 32bpp, (1 Cr & Cb sample per 2x1 Y samples), little-endian
	PixFmtYUV422P16BE    PixelFormat = C.AV_PIX_FMT_YUV422P16BE    // planar YUV 4:2:2, 32bpp, (1 Cr & Cb sample per 2x1 Y samples), big-endian
	PixFmYUV444P16LE     PixelFormat = C.AV_PIX_FMT_YUV444P16LE    // planar YUV 4:4:4, 48bpp, (1 Cr & Cb sample per 1x1 Y samples), little-endian
	PixFmtYUV444P16      PixelFormat = C.AV_PIX_FMT_YUV444P16BE    // planar YUV 4:4:4, 48bpp, (1 Cr & Cb sample per 1x1 Y samples), big-endian
	PixFmtDXVA2VLD       PixelFormat = C.AV_PIX_FMT_DXVA2_VLD      // HW decoding through DXVA2, Picture.data[3] contains a LPDIRECT3DSURFACE9 pointer
	PixFmtRGB444LE       PixelFormat = C.AV_PIX_FMT_RGB444LE       // packed RGB 4:4:4, 16bpp, (msb)4X 4R 4G 4B(lsb), little-endian, X=unused/undefined
	PixFmtRGB444BE       PixelFormat = C.AV_PIX_FMT_RGB444BE       // packed RGB 4:4:4, 16bpp, (msb)4X 4R 4G 4B(lsb), big-endian,    X=unused/undefined
	PixFmtBGR444LE       PixelFormat = C.AV_PIX_FMT_BGR444LE       // packed BGR 4:4:4, 16bpp, (msb)4X 4B 4G 4R(lsb), little-endian, X=unused/undefined
	PixFmtBGR444BE       PixelFormat = C.AV_PIX_FMT_BGR444BE       // packed BGR 4:4:4, 16bpp, (msb)4X 4B 4G 4R(lsb), big-endian,    X=unused/undefined
	PixFmtYA8            PixelFormat = C.AV_PIX_FMT_YA8            // 8 bits gray, 8 bits alpha
	PixFmtY400A          PixelFormat = C.AV_PIX_FMT_Y400A          // alias for AV_PIX_FMT_YA8
	PixFmtGray8A         PixelFormat = C.AV_PIX_FMT_GRAY8A         // alias for AV_PIX_FMT_YA8
	PixFmtBGR48BE        PixelFormat = C.AV_PIX_FMT_BGR48BE        // packed RGB 16:16:16, 48bpp, 16B, 16G, 16R, the 2-byte value for each R/G/B component is stored as big-endian
	PixFmtBGR48LE        PixelFormat = C.AV_PIX_FMT_BGR48LE        // packed RGB 16:16:16, 48bpp, 16B, 16G, 16R, the 2-byte value for each R/G/B component is stored as little-endian
	PixFmtYUV420P9BE     PixelFormat = C.AV_PIX_FMT_YUV420P9BE     // planar YUV 4:2:0, 13.5bpp, (1 Cr & Cb sample per 2x2 Y samples), big-endian
	PixFmtYUV420P9LE     PixelFormat = C.AV_PIX_FMT_YUV420P9LE     // planar YUV 4:2:0, 13.5bpp, (1 Cr & Cb sample per 2x2 Y samples), little-endian
	PixFmtYUV420P10BE    PixelFormat = C.AV_PIX_FMT_YUV420P10BE    // planar YUV 4:2:0, 15bpp, (1 Cr & Cb sample per 2x2 Y samples), big-endian
	PixFmtYUV420P10LE    PixelFormat = C.AV_PIX_FMT_YUV420P10LE    // planar YUV 4:2:0, 15bpp, (1 Cr & Cb sample per 2x2 Y samples), little-endian
	PixFmtYUV422P10BE    PixelFormat = C.AV_PIX_FMT_YUV422P10BE    // planar YUV 4:2:2, 20bpp, (1 Cr & Cb sample per 2x1 Y samples), big-endian
	PixFmtYUV422P10LE    PixelFormat = C.AV_PIX_FMT_YUV422P10LE    // planar YUV 4:2:2, 20bpp, (1 Cr & Cb sample per 2x1 Y samples), little-endian
	PixFmtYUV444P9BE     PixelFormat = C.AV_PIX_FMT_YUV444P9BE     // planar YUV 4:4:4, 27bpp, (1 Cr & Cb sample per 1x1 Y samples), big-endian
	PixFmtYUV444P9LE     PixelFormat = C.AV_PIX_FMT_YUV444P9LE     // planar YUV 4:4:4, 27bpp, (1 Cr & Cb sample per 1x1 Y samples), little-endian
	PixFmtYUV444P10BE    PixelFormat = C.AV_PIX_FMT_YUV444P10BE    // planar YUV 4:4:4, 30bpp, (1 Cr & Cb sample per 1x1 Y samples), big-endian
	PixFmtYUV444P10LE    PixelFormat = C.AV_PIX_FMT_YUV444P10LE    // planar YUV 4:4:4, 30bpp, (1 Cr & Cb sample per 1x1 Y samples), little-endian
	PixFmtYUV422P9BE     PixelFormat = C.AV_PIX_FMT_YUV422P9BE     // planar YUV 4:2:2, 18bpp, (1 Cr & Cb sample per 2x1 Y samples), big-endian
	PixFmtYUV422P9LE     PixelFormat = C.AV_PIX_FMT_YUV422P9LE     // planar YUV 4:2:2, 18bpp, (1 Cr & Cb sample per 2x1 Y samples), little-endian
	PixFmtGBRP           PixelFormat = C.AV_PIX_FMT_GBRP           // planar GBR 4:4:4 24bpp
	PixFmtGBR24P         PixelFormat = C.AV_PIX_FMT_GBR24P         // alias for #AV_PIX_FMT_GBRP
	PixFmtGBRP9BE        PixelFormat = C.AV_PIX_FMT_GBRP9BE        // planar GBR 4:4:4 27bpp, big-endian
	PixFmtGBRP9LE        PixelFormat = C.AV_PIX_FMT_GBRP9LE        // planar GBR 4:4:4 27bpp, little-endian
	PixFmtGBRP10BE       PixelFormat = C.AV_PIX_FMT_GBRP10BE       // planar GBR 4:4:4 30bpp, big-endian
	PixFmtGBRP10LE       PixelFormat = C.AV_PIX_FMT_GBRP10LE       // planar GBR 4:4:4 30bpp, little-endian
	PixFmtGBRP16BE       PixelFormat = C.AV_PIX_FMT_GBRP16BE       // planar GBR 4:4:4 48bpp, big-endian
	PixFmtGBRP16LE       PixelFormat = C.AV_PIX_FMT_GBRP16LE       // planar GBR 4:4:4 48bpp, little-endian
	PixFmtYUVA422P       PixelFormat = C.AV_PIX_FMT_YUVA422P       // planar YUV 4:2:2 24bpp, (1 Cr & Cb sample per 2x1 Y & A samples)
	PixFmtYUVA444P       PixelFormat = C.AV_PIX_FMT_YUVA444P       // planar YUV 4:4:4 32bpp, (1 Cr & Cb sample per 1x1 Y & A samples)
	PixFmtYUVA420P9BE    PixelFormat = C.AV_PIX_FMT_YUVA420P9BE    // planar YUV 4:2:0 22.5bpp, (1 Cr & Cb sample per 2x2 Y & A samples), big-endian
	PixFmtYUVA420P9LE    PixelFormat = C.AV_PIX_FMT_YUVA420P9LE    // planar YUV 4:2:0 22.5bpp, (1 Cr & Cb sample per 2x2 Y & A samples), little-endian
	PixFmtYUVA422P9BE    PixelFormat = C.AV_PIX_FMT_YUVA422P9BE    // planar YUV 4:2:2 27bpp, (1 Cr & Cb sample per 2x1 Y & A samples), big-endian
	PixFmtYUVA422P9LE    PixelFormat = C.AV_PIX_FMT_YUVA422P9LE    // planar YUV 4:2:2 27bpp, (1 Cr & Cb sample per 2x1 Y & A samples), little-endian
	PixFmtYUVA444P9BE    PixelFormat = C.AV_PIX_FMT_YUVA444P9BE    // planar YUV 4:4:4 36bpp, (1 Cr & Cb sample per 1x1 Y & A samples), big-endian
	PixFmtYUVA444P9LE    PixelFormat = C.AV_PIX_FMT_YUVA444P9LE    // planar YUV 4:4:4 36bpp, (1 Cr & Cb sample per 1x1 Y & A samples), little-endian
	PixFmtYUVA420P10BE   PixelFormat = C.AV_PIX_FMT_YUVA420P10BE   // planar YUV 4:2:0 25bpp, (1 Cr & Cb sample per 2x2 Y & A samples, big-endian)
	PixFmtYUVA420P10LE   PixelFormat = C.AV_PIX_FMT_YUVA420P10LE   // planar YUV 4:2:0 25bpp, (1 Cr & Cb sample per 2x2 Y & A samples, little-endian)
	PixFmtYUVA422P10BE   PixelFormat = C.AV_PIX_FMT_YUVA422P10BE   // planar YUV 4:2:2 30bpp, (1 Cr & Cb sample per 2x1 Y & A samples, big-endian)
	PixFmtYUVA422P10LE   PixelFormat = C.AV_PIX_FMT_YUVA422P10LE   // planar YUV 4:2:2 30bpp, (1 Cr & Cb sample per 2x1 Y & A samples, little-endian)
	PixFmtYUVA444P10BE   PixelFormat = C.AV_PIX_FMT_YUVA444P10BE   // planar YUV 4:4:4 40bpp, (1 Cr & Cb sample per 1x1 Y & A samples, big-endian)
	PixFmtYUVA444P10LE   PixelFormat = C.AV_PIX_FMT_YUVA444P10LE   // planar YUV 4:4:4 40bpp, (1 Cr & Cb sample per 1x1 Y & A samples, little-endian)
	PixFmtYUVA420P16BE   PixelFormat = C.AV_PIX_FMT_YUVA420P16BE   // planar YUV 4:2:0 40bpp, (1 Cr & Cb sample per 2x2 Y & A samples, big-endian)
	PixFmtYUVA420P16LE   PixelFormat = C.AV_PIX_FMT_YUVA420P16LE   // planar YUV 4:2:0 40bpp, (1 Cr & Cb sample per 2x2 Y & A samples, little-endian)
	PixFmtYUVA422P16BE   PixelFormat = C.AV_PIX_FMT_YUVA422P16BE   // planar YUV 4:2:2 48bpp, (1 Cr & Cb sample per 2x1 Y & A samples, big-endian)
	PixFmtYUVA422P16LE   PixelFormat = C.AV_PIX_FMT_YUVA422P16LE   // planar YUV 4:2:2 48bpp, (1 Cr & Cb sample per 2x1 Y & A samples, little-endian)
	PixFmtYUVA444P16BE   PixelFormat = C.AV_PIX_FMT_YUVA444P16BE   // planar YUV 4:4:4 64bpp, (1 Cr & Cb sample per 1x1 Y & A samples, big-endian)
	PixFmtYUVA444P16LE   PixelFormat = C.AV_PIX_FMT_YUVA444P16LE   // planar YUV 4:4:4 64bpp, (1 Cr & Cb sample per 1x1 Y & A samples, little-endian)
	PixFmtVDPAU          PixelFormat = C.AV_PIX_FMT_VDPAU          // HW acceleration through VDPAU, Picture.data[3] contains a VdpVideoSurface
	PixFmtXYZ12LE        PixelFormat = C.AV_PIX_FMT_XYZ12LE        // packed XYZ 4:4:4, 36 bpp, (msb) 12X, 12Y, 12Z (lsb), the 2-byte value for each X/Y/Z is stored as little-endian, the 4 lower bits are set to 0
	PixFmtXYZ12BE        PixelFormat = C.AV_PIX_FMT_XYZ12BE        // packed XYZ 4:4:4, 36 bpp, (msb) 12X, 12Y, 12Z (lsb), the 2-byte value for each X/Y/Z is stored as big-endian, the 4 lower bits are set to 0
	PixFmtNV16           PixelFormat = C.AV_PIX_FMT_NV16           // interleaved chroma YUV 4:2:2, 16bpp, (1 Cr & Cb sample per 2x1 Y samples)
	PixFmtNV20LE         PixelFormat = C.AV_PIX_FMT_NV20LE         // interleaved chroma YUV 4:2:2, 20bpp, (1 Cr & Cb sample per 2x1 Y samples), little-endian
	PixFmtNV20BE         PixelFormat = C.AV_PIX_FMT_NV20BE         // interleaved chroma YUV 4:2:2, 20bpp, (1 Cr & Cb sample per 2x1 Y samples), big-endian
	PixFmtRGBA64BE       PixelFormat = C.AV_PIX_FMT_RGBA64BE       // packed RGBA 16:16:16:16, 64bpp, 16R, 16G, 16B, 16A, the 2-byte value for each R/G/B/A component is stored as big-endian
	PixFmtRGBA64LE       PixelFormat = C.AV_PIX_FMT_RGBA64LE       // packed RGBA 16:16:16:16, 64bpp, 16R, 16G, 16B, 16A, the 2-byte value for each R/G/B/A component is stored as little-endian
	PixFmtBGRA64BE       PixelFormat = C.AV_PIX_FMT_BGRA64BE       // packed RGBA 16:16:16:16, 64bpp, 16B, 16G, 16R, 16A, the 2-byte value for each R/G/B/A component is stored as big-endian
	PixFmtBGRA64LE       PixelFormat = C.AV_PIX_FMT_BGRA64LE       // packed RGBA 16:16:16:16, 64bpp, 16B, 16G, 16R, 16A, the 2-byte value for each R/G/B/A component is stored as little-endian
	PixFmtYVYU422        PixelFormat = C.AV_PIX_FMT_YVYU422        // packed YUV 4:2:2, 16bpp, Y0 Cr Y1 Cb
	PixFmtYA16BE         PixelFormat = C.AV_PIX_FMT_YA16BE         // 16 bits gray, 16 bits alpha (big-endian)
	PixFmtYA16LE         PixelFormat = C.AV_PIX_FMT_YA16LE         // 16 bits gray, 16 bits alpha (little-endian)
	PixFmtGBRAP          PixelFormat = C.AV_PIX_FMT_GBRAP          // planar GBRA 4:4:4:4 32bpp
	PixFmtGBRAP16BE      PixelFormat = C.AV_PIX_FMT_GBRAP16BE      // planar GBRA 4:4:4:4 64bpp, big-endian
	PixFmtGBRAP16LE      PixelFormat = C.AV_PIX_FMT_GBRAP16LE      // planar GBRA 4:4:4:4 64bpp, little-endian
	PixFmtQSV            PixelFormat = C.AV_PIX_FMT_QSV            // TODO:
	PixFmtMMAL           PixelFormat = C.AV_PIX_FMT_MMAL           // HW acceleration though MMAL, data[3] contains a pointer to the MMAL_BUFFER_HEADER_T structure.
	PixFmtD3D11VA_VLD    PixelFormat = C.AV_PIX_FMT_D3D11VA_VLD    // HW decoding through Direct3D11 via old API, Picture.data[3] contains a ID3D11VideoDecoderOutputView pointer
	PixFmtCUDA           PixelFormat = C.AV_PIX_FMT_CUDA           // HW acceleration through CUDA. data[i] contain CUdeviceptr pointers exactly as for system memory frames.
	PixFmt0RGB           PixelFormat = C.AV_PIX_FMT_0RGB           // packed RGB 8:8:8, 32bpp, XRGBXRGB...   X=unused/undefined
	PixFmtRGB0           PixelFormat = C.AV_PIX_FMT_RGB0           // packed RGB 8:8:8, 32bpp, RGBXRGBX...   X=unused/undefined
	PixFmt0BGR           PixelFormat = C.AV_PIX_FMT_0BGR           // packed BGR 8:8:8, 32bpp, XBGRXBGR...   X=unused/undefined
	PixFmtBGR0           PixelFormat = C.AV_PIX_FMT_BGR0           // packed BGR 8:8:8, 32bpp, BGRXBGRX...   X=unused/undefined
	PixFmtYUV420P12BE    PixelFormat = C.AV_PIX_FMT_YUV420P12BE    // planar YUV 4:2:0,18bpp, (1 Cr & Cb sample per 2x2 Y samples), big-endian
	PixFmtYUV420P12LE    PixelFormat = C.AV_PIX_FMT_YUV420P12LE    // planar YUV 4:2:0,18bpp, (1 Cr & Cb sample per 2x2 Y samples), little-endian
	PixFmtYUV420P14BE    PixelFormat = C.AV_PIX_FMT_YUV420P14BE    // planar YUV 4:2:0,21bpp, (1 Cr & Cb sample per 2x2 Y samples), big-endian
	PixFmtYUV420P14LE    PixelFormat = C.AV_PIX_FMT_YUV420P14LE    // planar YUV 4:2:0,21bpp, (1 Cr & Cb sample per 2x2 Y samples), little-endian
	PixFmtYUV422P12BE    PixelFormat = C.AV_PIX_FMT_YUV422P12BE    // planar YUV 4:2:2,24bpp, (1 Cr & Cb sample per 2x1 Y samples), big-endian
	PixFmtYUV422P12LE    PixelFormat = C.AV_PIX_FMT_YUV422P12LE    // planar YUV 4:2:2,24bpp, (1 Cr & Cb sample per 2x1 Y samples), little-endian
	PixFmtYUV422P14BE    PixelFormat = C.AV_PIX_FMT_YUV422P14BE    // planar YUV 4:2:2,28bpp, (1 Cr & Cb sample per 2x1 Y samples), big-endian
	PixFmtYUV422P14LE    PixelFormat = C.AV_PIX_FMT_YUV422P14LE    // planar YUV 4:2:2,28bpp, (1 Cr & Cb sample per 2x1 Y samples), little-endian
	PixFmtYUV444P12BE    PixelFormat = C.AV_PIX_FMT_YUV444P12BE    // planar YUV 4:4:4,36bpp, (1 Cr & Cb sample per 1x1 Y samples), big-endian
	PixFmtYUV444P12LE    PixelFormat = C.AV_PIX_FMT_YUV444P12LE    // planar YUV 4:4:4,36bpp, (1 Cr & Cb sample per 1x1 Y samples), little-endian
	PixFmtYUV444P14BE    PixelFormat = C.AV_PIX_FMT_YUV444P14BE    // planar YUV 4:4:4,42bpp, (1 Cr & Cb sample per 1x1 Y samples), big-endian
	PixFmtYUV444P14LE    PixelFormat = C.AV_PIX_FMT_YUV444P14LE    // planar YUV 4:4:4,42bpp, (1 Cr & Cb sample per 1x1 Y samples), little-endian
	PixFmtGBRP12BE       PixelFormat = C.AV_PIX_FMT_GBRP12BE       // planar GBR 4:4:4 36bpp, big-endian
	PixFmtGBRP12LE       PixelFormat = C.AV_PIX_FMT_GBRP12LE       // planar GBR 4:4:4 36bpp, little-endian
	PixFmtGBRP14BE       PixelFormat = C.AV_PIX_FMT_GBRP14BE       // planar GBR 4:4:4 42bpp, big-endian
	PixFmtGBRP14LE       PixelFormat = C.AV_PIX_FMT_GBRP14LE       // planar GBR 4:4:4 42bpp, little-endian
	PixFmtYUVJ411P       PixelFormat = C.AV_PIX_FMT_YUVJ411P       // planar YUV 4:1:1, 12bpp, (1 Cr & Cb sample per 4x1 Y samples) full scale (JPEG), deprecated in favor of AV_PIX_FMT_YUV411P and setting color_range
	PixFmtBAYER_BGGR8    PixelFormat = C.AV_PIX_FMT_BAYER_BGGR8    // bayer, BGBG..(odd line), GRGR..(even line), 8-bit samples
	PixFmtBAYER_RGGB8    PixelFormat = C.AV_PIX_FMT_BAYER_RGGB8    // bayer, RGRG..(odd line), GBGB..(even line), 8-bit samples
	PixFmtBAYER_GBRG8    PixelFormat = C.AV_PIX_FMT_BAYER_GBRG8    // bayer, GBGB..(odd line), RGRG..(even line), 8-bit samples
	PixFmtBAYER_GRBG8    PixelFormat = C.AV_PIX_FMT_BAYER_GRBG8    // bayer, GRGR..(odd line), BGBG..(even line), 8-bit samples
	PixFmtBAYER_BGGR16LE PixelFormat = C.AV_PIX_FMT_BAYER_BGGR16LE // bayer, BGBG..(odd line), GRGR..(even line), 16-bit samples, little-endian
	PixFmtBAYER_BGGR16BE PixelFormat = C.AV_PIX_FMT_BAYER_BGGR16BE // bayer, BGBG..(odd line), GRGR..(even line), 16-bit samples, big-endian
	PixFmtBAYER_RGGB16LE PixelFormat = C.AV_PIX_FMT_BAYER_RGGB16LE // bayer, RGRG..(odd line), GBGB..(even line), 16-bit samples, little-endian
	PixFmtBAYER_RGGB16BE PixelFormat = C.AV_PIX_FMT_BAYER_RGGB16BE // bayer, RGRG..(odd line), GBGB..(even line), 16-bit samples, big-endian
	PixFmtBAYER_GBRG16LE PixelFormat = C.AV_PIX_FMT_BAYER_GBRG16LE // bayer, GBGB..(odd line), RGRG..(even line), 16-bit samples, little-endian
	PixFmtBAYER_GBRG16BE PixelFormat = C.AV_PIX_FMT_BAYER_GBRG16BE // bayer, GBGB..(odd line), RGRG..(even line), 16-bit samples, big-endian
	PixFmtBAYER_GRBG16LE PixelFormat = C.AV_PIX_FMT_BAYER_GRBG16LE // bayer, GRGR..(odd line), BGBG..(even line), 16-bit samples, little-endian
	PixFmtBAYER_GRBG16BE PixelFormat = C.AV_PIX_FMT_BAYER_GRBG16BE // bayer, GRGR..(odd line), BGBG..(even line), 16-bit samples, big-endian
	PixFmtYUV440P10LE    PixelFormat = C.AV_PIX_FMT_YUV440P10LE    // planar YUV 4:4:0,20bpp, (1 Cr & Cb sample per 1x2 Y samples), little-endian
	PixFmtYUV440P10BE    PixelFormat = C.AV_PIX_FMT_YUV440P10BE    // planar YUV 4:4:0,20bpp, (1 Cr & Cb sample per 1x2 Y samples), big-endian
	PixFmtYUV440P12LE    PixelFormat = C.AV_PIX_FMT_YUV440P12LE    // planar YUV 4:4:0,24bpp, (1 Cr & Cb sample per 1x2 Y samples), little-endian
	PixFmtYUV440P12BE    PixelFormat = C.AV_PIX_FMT_YUV440P12BE    // planar YUV 4:4:0,24bpp, (1 Cr & Cb sample per 1x2 Y samples), big-endian
	PixFmtAYUV64LE       PixelFormat = C.AV_PIX_FMT_AYUV64LE       // packed AYUV 4:4:4,64bpp (1 Cr & Cb sample per 1x1 Y & A samples), little-endian
	PixFmtAYUV64BE       PixelFormat = C.AV_PIX_FMT_AYUV64BE       // packed AYUV 4:4:4,64bpp (1 Cr & Cb sample per 1x1 Y & A samples), big-endian
	PixFmtVIDEOTOOLBOX   PixelFormat = C.AV_PIX_FMT_VIDEOTOOLBOX   // hardware decoding through Videotoolbox
	PixFmtP010LE         PixelFormat = C.AV_PIX_FMT_P010LE         // like NV12, with 10bpp per component, data in the high bits, zeros in the low bits, little-endian
	PixFmtP010BE         PixelFormat = C.AV_PIX_FMT_P010BE         // like NV12, with 10bpp per component, data in the high bits, zeros in the low bits, big-endian
	PixFmtGBRAP12BE      PixelFormat = C.AV_PIX_FMT_GBRAP12BE      // planar GBR 4:4:4:4 48bpp, big-endian
	PixFmtGBRAP12LE      PixelFormat = C.AV_PIX_FMT_GBRAP12LE      // planar GBR 4:4:4:4 48bpp, little-endian
	PixFmtGBRAP10BE      PixelFormat = C.AV_PIX_FMT_GBRAP10BE      // planar GBR 4:4:4:4 40bpp, big-endian
	PixFmtGBRAP10LE      PixelFormat = C.AV_PIX_FMT_GBRAP10LE      // planar GBR 4:4:4:4 40bpp, little-endian
	PixFmtMEDIACODEC     PixelFormat = C.AV_PIX_FMT_MEDIACODEC     // hardware decoding through MediaCodec
	PixFmtGRAY12BE       PixelFormat = C.AV_PIX_FMT_GRAY12BE       //        Y        , 12bpp, big-endian
	PixFmtGRAY12LE       PixelFormat = C.AV_PIX_FMT_GRAY12LE       //        Y        , 12bpp, little-endian
	PixFmtGRAY10BE       PixelFormat = C.AV_PIX_FMT_GRAY10BE       //        Y        , 10bpp, big-endian
	PixFmtGRAY10LE       PixelFormat = C.AV_PIX_FMT_GRAY10LE       //        Y        , 10bpp, little-endian
	PixFmtP016LE         PixelFormat = C.AV_PIX_FMT_P016LE         // like NV12, with 16bpp per component, little-endian
	PixFmtP016BE         PixelFormat = C.AV_PIX_FMT_P016BE         // like NV12, with 16bpp per component, big-endian
	PixFmtD3D11          PixelFormat = C.AV_PIX_FMT_D3D11          // TODO:
	PixFmtGRAY9BE        PixelFormat = C.AV_PIX_FMT_GRAY9BE        //        Y        , 9bpp, big-endian
	PixFmtGRAY9LE        PixelFormat = C.AV_PIX_FMT_GRAY9LE        //        Y        , 9bpp, little-endian
	PixFmtGBRPF32BE      PixelFormat = C.AV_PIX_FMT_GBRPF32BE      // IEEE-754 single precision planar GBR 4:4:4,     96bpp, big-endian
	PixFmtGBRPF32LE      PixelFormat = C.AV_PIX_FMT_GBRPF32LE      // IEEE-754 single precision planar GBR 4:4:4,     96bpp, little-endian
	PixFmtGBRAPF32BE     PixelFormat = C.AV_PIX_FMT_GBRAPF32BE     // IEEE-754 single precision planar GBRA 4:4:4:4, 128bpp, big-endian
	PixFmtGBRAPF32LE     PixelFormat = C.AV_PIX_FMT_GBRAPF32LE     // IEEE-754 single precision planar GBRA 4:4:4:4, 128bpp, little-endian
	PixFmtDRM_PRIME      PixelFormat = C.AV_PIX_FMT_DRM_PRIME      // DRM-managed buffers exposed through PRIME buffer sharing.
	PixFmtOPENCL         PixelFormat = C.AV_PIX_FMT_OPENCL         // Hardware surfaces for OpenCL.
	PixFmtGRAY14BE       PixelFormat = C.AV_PIX_FMT_GRAY14BE       //        Y        , 14bpp, big-endian
	PixFmtGRAY14LE       PixelFormat = C.AV_PIX_FMT_GRAY14LE       //        Y        , 14bpp, little-endian
	PixFmtGRAYF32BE      PixelFormat = C.AV_PIX_FMT_GRAYF32BE      // IEEE-754 single precision Y, 32bpp, big-endian
	PixFmtGRAYF32LE      PixelFormat = C.AV_PIX_FMT_GRAYF32LE      // IEEE-754 single precision Y, 32bpp, little-endian
	PixFmtYUVA422P12BE   PixelFormat = C.AV_PIX_FMT_YUVA422P12BE   // planar YUV 4:2:2,24bpp, (1 Cr & Cb sample per 2x1 Y samples), 12b alpha, big-endian
	PixFmtYUVA422P12LE   PixelFormat = C.AV_PIX_FMT_YUVA422P12LE   // planar YUV 4:2:2,24bpp, (1 Cr & Cb sample per 2x1 Y samples), 12b alpha, little-endian
	PixFmtYUVA444P12BE   PixelFormat = C.AV_PIX_FMT_YUVA444P12BE   // planar YUV 4:4:4,36bpp, (1 Cr & Cb sample per 1x1 Y samples), 12b alpha, big-endian
	PixFmtYUVA444P12LE   PixelFormat = C.AV_PIX_FMT_YUVA444P12LE   // planar YUV 4:4:4,36bpp, (1 Cr & Cb sample per 1x1 Y samples), 12b alpha, little-endian
	PixFmtNV24           PixelFormat = C.AV_PIX_FMT_NV24           // planar YUV 4:4:4, 24bpp, 1 plane for Y and 1 plane for the UV components, which are interleaved (first byte U and the following byte V)
	PixFmtNV42           PixelFormat = C.AV_PIX_FMT_NV42           // as above, but U and V bytes are swapped
	PixFmtVULKAN         PixelFormat = C.AV_PIX_FMT_VULKAN         // Vulkan hardware images.
	PixFmtY210BE         PixelFormat = C.AV_PIX_FMT_Y210BE         // packed YUV 4:2:2 like YUYV422, 20bpp, data in the high bits, big-endian
	PixFmtY210LE         PixelFormat = C.AV_PIX_FMT_Y210LE         // packed YUV 4:2:2 like YUYV422, 20bpp, data in the high bits, little-endian
	PixFmtX2RGB10LE      PixelFormat = C.AV_PIX_FMT_X2RGB10LE      // packed RGB 10:10:10, 30bpp, (msb)2X 10R 10G 10B(lsb), little-endian, X=unused/undefined
	PixFmtX2RGB10BE      PixelFormat = C.AV_PIX_FMT_X2RGB10BE      // packed RGB 10:10:10, 30bpp, (msb)2X 10R 10G 10B(lsb), big-endian, X=unused/undefined
	PixFmtX2BGR10LE      PixelFormat = C.AV_PIX_FMT_X2BGR10LE      // packed BGR 10:10:10, 30bpp, (msb)2X 10B 10G 10R(lsb), little-endian, X=unused/undefined
	PixFmtX2BGR10BE      PixelFormat = C.AV_PIX_FMT_X2BGR10BE      // packed BGR 10:10:10, 30bpp, (msb)2X 10B 10G 10R(lsb), big-endian, X=unused/undefined
	PixFmtP210BE         PixelFormat = C.AV_PIX_FMT_P210BE         // interleaved chroma YUV 4:2:2, 20bpp, data in the high bits, big-endian
	PixFmtP210LE         PixelFormat = C.AV_PIX_FMT_P210LE         // interleaved chroma YUV 4:2:2, 20bpp, data in the high bits, little-endian
	PixFmtP410BE         PixelFormat = C.AV_PIX_FMT_P410BE         // interleaved chroma YUV 4:4:4, 30bpp, data in the high bits, big-endian
	PixFmtP410LE         PixelFormat = C.AV_PIX_FMT_P410LE         // interleaved chroma YUV 4:4:4, 30bpp, data in the high bits, little-endian
	PixFmtP216BE         PixelFormat = C.AV_PIX_FMT_P216BE         // interleaved chroma YUV 4:2:2, 32bpp, big-endian
	PixFmtP216LE         PixelFormat = C.AV_PIX_FMT_P216LE         // interleaved chroma YUV 4:2:2, 32bpp, little-endian
	PixFmtP416BE         PixelFormat = C.AV_PIX_FMT_P416BE         // interleaved chroma YUV 4:4:4, 48bpp, big-endian
	PixFmtP416LE         PixelFormat = C.AV_PIX_FMT_P416LE         // interleaved chroma YUV 4:4:4, 48bpp, little-endian
	PixFmtVUYA           PixelFormat = C.AV_PIX_FMT_VUYA           // packed VUYA 4:4:4:4, 32bpp (1 Cr & Cb sample per 1x1 Y & A samples), VUYAVUYA...
	PixFmtRGBAF16BE      PixelFormat = C.AV_PIX_FMT_RGBAF16BE      // IEEE-754 half precision packed RGBA 16:16:16:16, 64bpp, RGBARGBA..., big-endian
	PixFmtRGBAF16LE      PixelFormat = C.AV_PIX_FMT_RGBAF16LE      // IEEE-754 half precision packed RGBA 16:16:16:16, 64bpp, RGBARGBA..., little-endian
	PixFmtVUYX           PixelFormat = C.AV_PIX_FMT_VUYX           // packed VUYX 4:4:4:4, 32bpp, Variant of VUYA where alpha channel is left undefined
	PixFmtP012LE         PixelFormat = C.AV_PIX_FMT_P012LE         // like NV12, with 12bpp per component, data in the high bits, zeros in the low bits, little-endian
	PixFmtP012BE         PixelFormat = C.AV_PIX_FMT_P012BE         // like NV12, with 12bpp per component, data in the high bits, zeros in the low bits, big-endian
	PixFmtY212BE         PixelFormat = C.AV_PIX_FMT_Y212BE         // packed YUV 4:2:2 like YUYV422, 24bpp, data in the high bits, zeros in the low bits, big-endian
	PixFmtY212LE         PixelFormat = C.AV_PIX_FMT_Y212LE         // packed YUV 4:2:2 like YUYV422, 24bpp, data in the high bits, zeros in the low bits, little-endian
	PixFmtXV30BE         PixelFormat = C.AV_PIX_FMT_XV30BE         // packed XVYU 4:4:4, 32bpp, (msb)2X 10V 10Y 10U(lsb), big-endian, variant of Y410 where alpha channel is left undefined
	PixFmtXV30LE         PixelFormat = C.AV_PIX_FMT_XV30LE         // packed XVYU 4:4:4, 32bpp, (msb)2X 10V 10Y 10U(lsb), little-endian, variant of Y410 where alpha channel is left undefined
	PixFmtXV36BE         PixelFormat = C.AV_PIX_FMT_XV36BE         // packed XVYU 4:4:4, 48bpp, data in the high bits, zeros in the low bits, big-endian, variant of Y412 where alpha channel is left undefined
	PixFmtXV36LE         PixelFormat = C.AV_PIX_FMT_XV36LE         // packed XVYU 4:4:4, 48bpp, data in the high bits, zeros in the low bits, little-endian, variant of Y412 where alpha channel is left undefined
	PixFmtRGBF32BE       PixelFormat = C.AV_PIX_FMT_RGBF32BE       // IEEE-754 single precision packed RGB 32:32:32, 96bpp, RGBRGB..., big-endian
	PixFmtRGBF32LE       PixelFormat = C.AV_PIX_FMT_RGBF32LE       // IEEE-754 single precision packed RGB 32:32:32, 96bpp, RGBRGB..., little-endian
	PixFmtRGBAF32BE      PixelFormat = C.AV_PIX_FMT_RGBAF32BE      // IEEE-754 single precision packed RGBA 32:32:32:32, 128bpp, RGBARGBA..., big-endian
	PixFmtRGBAF32LE      PixelFormat = C.AV_PIX_FMT_RGBAF32LE      // IEEE-754 single precision packed RGBA 32:32:32:32, 128bpp, RGBARGBA..., little-endian
	PixFmtP212BE         PixelFormat = C.AV_PIX_FMT_P212BE         // interleaved chroma YUV 4:2:2, 24bpp, data in the high bits, big-endian
	PixFmtP212LE         PixelFormat = C.AV_PIX_FMT_P212LE         // interleaved chroma YUV 4:2:2, 24bpp, data in the high bits, little-endian
	PixFmtP412BE         PixelFormat = C.AV_PIX_FMT_P412BE         // interleaved chroma YUV 4:4:4, 36bpp, data in the high bits, big-endian
	PixFmtP412LE         PixelFormat = C.AV_PIX_FMT_P412LE         // interleaved chroma YUV 4:4:4, 36bpp, data in the high bits, little-endian
	PixFmtGBRAP14BE      PixelFormat = C.AV_PIX_FMT_GBRAP14BE      // planar GBR 4:4:4:4 56bpp, big-endian
	PixFmtGBRAP14LE      PixelFormat = C.AV_PIX_FMT_GBRAP14LE      // planar GBR 4:4:4:4 56bpp, little-endian
	PixFmtD3D12          PixelFormat = C.AV_PIX_FMT_D3D12          // Hardware surfaces for Direct3D 12. data[0] points to an AVD3D12VAFrame.
	PixFmtAYUV           PixelFormat = C.AV_PIX_FMT_AYUV           // packed AYUV 4:4:4:4, 32bpp (1 Cr & Cb sample per 1x1 Y & A samples), AYUVAYUV...
	PixFmtUYVA           PixelFormat = C.AV_PIX_FMT_UYVA           // packed UYVA 4:4:4:4, 32bpp (1 Cr & Cb sample per 1x1 Y & A samples), UYVAUYVA...
	PixFmtVYU444         PixelFormat = C.AV_PIX_FMT_VYU444         // packed VYU 4:4:4, 24bpp (1 Cr & Cb sample per 1x1 Y), VYUVYU...
	PixFmtV30XBE         PixelFormat = C.AV_PIX_FMT_V30XBE         // packed VYUX 4:4:4 like XV30, 32bpp, (msb)10V 10Y 10U 2X(lsb), big-endian
	PixFmtV30XLE         PixelFormat = C.AV_PIX_FMT_V30XLE         // packed VYUX 4:4:4 like XV30, 32bpp, (msb)10V 10Y 10U 2X(lsb), little-endian
	PixFmtRGBF16BE       PixelFormat = C.AV_PIX_FMT_RGBF16BE       // IEEE-754 half precision packed RGB 16:16:16, 48bpp, RGBRGB..., big-endian
	PixFmtRGBF16LE       PixelFormat = C.AV_PIX_FMT_RGBF16LE       // IEEE-754 half precision packed RGB 16:16:16, 48bpp, RGBRGB..., little-endian
	PixFmtRGBA128BE      PixelFormat = C.AV_PIX_FMT_RGBA128BE      // packed RGBA 32:32:32:32, 128bpp, RGBARGBA..., big-endian
	PixFmtRGBA128LE      PixelFormat = C.AV_PIX_FMT_RGBA128LE      // packed RGBA 32:32:32:32, 128bpp, RGBARGBA..., little-endian
	PixFmtRGB96BE        PixelFormat = C.AV_PIX_FMT_RGB96BE        // packed RGBA 32:32:32, 96bpp, RGBRGB..., big-endian
	PixFmtRGB96LE        PixelFormat = C.AV_PIX_FMT_RGB96LE        // packed RGBA 32:32:32, 96bpp, RGBRGB..., little-endian
	PixFmtY216BE         PixelFormat = C.AV_PIX_FMT_Y216BE         // packed YUV 4:2:2 like YUYV422, 32bpp, big-endian
	PixFmtY216LE         PixelFormat = C.AV_PIX_FMT_Y216LE         // packed YUV 4:2:2 like YUYV422, 32bpp, little-endian
	PixFmtXV48BE         PixelFormat = C.AV_PIX_FMT_XV48BE         // packed XVYU 4:4:4, 64bpp, big-endian, variant of Y416 where alpha channel is left undefined
	PixFmtXV48LE         PixelFormat = C.AV_PIX_FMT_XV48LE         // packed XVYU 4:4:4, 64bpp, little-endian, variant of Y416 where alpha channel is left undefined
	PixFmtGBRPF16BE      PixelFormat = C.AV_PIX_FMT_GBRPF16BE      // IEEE-754 half precision planer GBR 4:4:4, 48bpp, big-endian
	PixFmtGBRPF16LE      PixelFormat = C.AV_PIX_FMT_GBRPF16LE      // IEEE-754 half precision planer GBR 4:4:4, 48bpp, little-endian
	PixFmtGBRAPF16BE     PixelFormat = C.AV_PIX_FMT_GBRAPF16BE     // IEEE-754 half precision planar GBRA 4:4:4:4, 64bpp, big-endian
	PixFmtGBRAPF16LE     PixelFormat = C.AV_PIX_FMT_GBRAPF16LE     // IEEE-754 half precision planar GBRA 4:4:4:4, 64bpp, little-endian
	PixFmtGRAYF16BE      PixelFormat = C.AV_PIX_FMT_GRAYF16BE      // IEEE-754 half precision Y, 16bpp, big-endian
	PixFmtGRAYF16LE      PixelFormat = C.AV_PIX_FMT_GRAYF16LE      // IEEE-754 half precision Y, 16bpp, little-endian
	PixFmtAMF_SURFACE    PixelFormat = C.AV_PIX_FMT_AMF_SURFACE    // HW acceleration through AMF. data[0] contain AMFSurface pointer
	PixFmtGRAY32BE       PixelFormat = C.AV_PIX_FMT_GRAY32BE       //         Y        , 32bpp, big-endian
	PixFmtGRAY32LE       PixelFormat = C.AV_PIX_FMT_GRAY32LE       //         Y        , 32bpp, little-endian
	PixFmtYAF32BE        PixelFormat = C.AV_PIX_FMT_YAF32BE        // IEEE-754 single precision packed YA, 32 bits gray, 32 bits alpha, 64bpp, big-endian
	PixFmtYAF32LE        PixelFormat = C.AV_PIX_FMT_YAF32LE        // IEEE-754 single precision packed YA, 32 bits gray, 32 bits alpha, 64bpp, little-endian
	PixFmtYAF16BE        PixelFormat = C.AV_PIX_FMT_YAF16BE        // IEEE-754 half precision packed YA, 16 bits gray, 16 bits alpha, 32bpp, big-endian
	PixFmtYAF16LE        PixelFormat = C.AV_PIX_FMT_YAF16LE        // IEEE-754 half precision packed YA, 16 bits gray, 16 bits alpha, 32bpp, little-endian
	PixFmtGBRAP32BE      PixelFormat = C.AV_PIX_FMT_GBRAP32BE      // planar GBRA 4:4:4:4 128bpp, big-endian
	PixFmtGBRAP32LE      PixelFormat = C.AV_PIX_FMT_GBRAP32LE      // planar GBRA 4:4:4:4 128bpp, little-endian
	PixFmtYUV444P10MSBBE PixelFormat = C.AV_PIX_FMT_YUV444P10MSBBE // planar YUV 4:4:4, 30bpp, (1 Cr & Cb sample per 1x1 Y samples), lowest bits zero, big-endian
	PixFmtYUV444P10MSBLE PixelFormat = C.AV_PIX_FMT_YUV444P10MSBLE // planar YUV 4:4:4, 30bpp, (1 Cr & Cb sample per 1x1 Y samples), lowest bits zero, little-endian
	PixFmtYUV444P12MSBBE PixelFormat = C.AV_PIX_FMT_YUV444P12MSBBE // planar YUV 4:4:4, 30bpp, (1 Cr & Cb sample per 1x1 Y samples), lowest bits zero, big-endian
	PixFmtYUV444P12MSBLE PixelFormat = C.AV_PIX_FMT_YUV444P12MSBLE // planar YUV 4:4:4, 30bpp, (1 Cr & Cb sample per 1x1 Y samples), lowest bits zero, little-endian
	PixFmtGBRP10MSBBE    PixelFormat = C.AV_PIX_FMT_GBRP10MSBBE    // planar GBR 4:4:4 30bpp, lowest bits zero, big-endian
	PixFmtGBRP10MSBLE    PixelFormat = C.AV_PIX_FMT_GBRP10MSBLE    // planar GBR 4:4:4 30bpp, lowest bits zero, little-endian
	PixFmtGBRP12MSBBE    PixelFormat = C.AV_PIX_FMT_GBRP12MSBBE    // planar GBR 4:4:4 36bpp, lowest bits zero, big-endian
	PixFmtGBRP12MSBLE    PixelFormat = C.AV_PIX_FMT_GBRP12MSBLE    // planar GBR 4:4:4 36bpp, lowest bits zero, little-endian
	PixFmtOHCODEC        PixelFormat = C.AV_PIX_FMT_OHCODEC        // hardware decoding through openharmony
	PixFmtNB             PixelFormat = C.AV_PIX_FMT_NB             // number of pixel formats, DO NOT USE THIS if you want to link with shared libav* because the number of formats might differ between versions
)

// Chromaticity coordinates of the source primaries.
//
// These values match the ones defined by ISO/IEC 23091-2_2019 subclause 8.1 and ITU-T H.273.
type ColorPrimaries int

const (
	ColorPrimariesReserved0    ColorPrimaries = C.AVCOL_PRI_RESERVED0    //
	ColorPrimariesBT709        ColorPrimaries = C.AVCOL_PRI_BT709        // also ITU-R BT1361 / IEC 61966-2-4 / SMPTE RP 177 Annex B
	ColorPrimariesUnspecified  ColorPrimaries = C.AVCOL_PRI_UNSPECIFIED  //
	ColorPrimariesReserved     ColorPrimaries = C.AVCOL_PRI_RESERVED     //
	ColorPrimariesBT470M       ColorPrimaries = C.AVCOL_PRI_BT470M       // also FCC Title 47 Code of Federal Regulations 73.682 (a)(20)
	ColorPrimariesBT470BG      ColorPrimaries = C.AVCOL_PRI_BT470BG      // also ITU-R BT601-6 625 / ITU-R BT1358 625 / ITU-R BT1700 625 PAL & SECAM
	ColorPrimariesSMPTE170M    ColorPrimaries = C.AVCOL_PRI_SMPTE170M    // also ITU-R BT601-6 525 / ITU-R BT1358 525 / ITU-R BT1700 NTSC
	ColorPrimariesSMPTE240M    ColorPrimaries = C.AVCOL_PRI_SMPTE240M    // identical to above, also called "SMPTE C" even though it uses D65
	ColorPrimariesFilm         ColorPrimaries = C.AVCOL_PRI_FILM         // colour filters using Illuminant C
	ColorPrimariesBT2020       ColorPrimaries = C.AVCOL_PRI_BT2020       // ITU-R BT2020
	ColorPrimariesSMPTE428     ColorPrimaries = C.AVCOL_PRI_SMPTE428     // SMPTE ST 428-1 (CIE 1931 XYZ)
	ColorPrimariesSMPTEST428_1 ColorPrimaries = C.AVCOL_PRI_SMPTEST428_1 //
	ColorPrimariesSMPTE431     ColorPrimaries = C.AVCOL_PRI_SMPTE431     // SMPTE ST 431-2 (2011) / DCI P3
	ColorPrimariesSMPTE432     ColorPrimaries = C.AVCOL_PRI_SMPTE432     // SMPTE ST 432-1 (2010) / P3 D65 / Display P3
	ColorPrimariesEBU3213      ColorPrimaries = C.AVCOL_PRI_EBU3213      // EBU Tech. 3213-E (nothing there) / one of JEDEC P22 group phosphors
	ColorPrimariesJEDEC_P22    ColorPrimaries = C.AVCOL_PRI_JEDEC_P22    //
	ColorPrimariesNB           ColorPrimaries = C.AVCOL_PRI_NB           // Not part of ABI
)

// Color Transfer Characteristic.
//
// These values match the ones defined by ISO/IEC 23091-2_2019 subclause 8.2.
type ColorTransferCharacteristic int

const (
	ColorTransferCharacteristicReserved0    ColorTransferCharacteristic = C.AVCOL_TRC_RESERVED0
	ColorTransferCharacteristicBT709        ColorTransferCharacteristic = C.AVCOL_TRC_BT709        ///< also ITU-R BT1361
	ColorTransferCharacteristicUnspecified  ColorTransferCharacteristic = C.AVCOL_TRC_UNSPECIFIED  //
	ColorTransferCharacteristicReserved     ColorTransferCharacteristic = C.AVCOL_TRC_RESERVED     //
	ColorTransferCharacteristicGamma22      ColorTransferCharacteristic = C.AVCOL_TRC_GAMMA22      ///< also ITU-R BT470M / ITU-R BT1700 625 PAL & SECAM
	ColorTransferCharacteristicGamma28      ColorTransferCharacteristic = C.AVCOL_TRC_GAMMA28      ///< also ITU-R BT470BG
	ColorTransferCharacteristicSMPTE170M    ColorTransferCharacteristic = C.AVCOL_TRC_SMPTE170M    ///< also ITU-R BT601-6 525 or 625 / ITU-R BT1358 525 or 625 / ITU-R BT1700 NTSC
	ColorTransferCharacteristicSMPTE240M    ColorTransferCharacteristic = C.AVCOL_TRC_SMPTE240M    //
	ColorTransferCharacteristicLinear       ColorTransferCharacteristic = C.AVCOL_TRC_LINEAR       ///< "Linear transfer characteristics"
	ColorTransferCharacteristicLog          ColorTransferCharacteristic = C.AVCOL_TRC_LOG          ///< "Logarithmic transfer characteristic (100:1 range)"
	ColorTransferCharacteristicLogSqrt      ColorTransferCharacteristic = C.AVCOL_TRC_LOG_SQRT     ///< "Logarithmic transfer characteristic (100 * Sqrt(10) : 1 range)"
	ColorTransferCharacteristicIEC61966_2_4 ColorTransferCharacteristic = C.AVCOL_TRC_IEC61966_2_4 ///< IEC 61966-2-4
	ColorTransferCharacteristicBT1361_ECG   ColorTransferCharacteristic = C.AVCOL_TRC_BT1361_ECG   ///< ITU-R BT1361 Extended Colour Gamut
	ColorTransferCharacteristicIEC61966_2_1 ColorTransferCharacteristic = C.AVCOL_TRC_IEC61966_2_1 ///< IEC 61966-2-1 (sRGB or sYCC)
	ColorTransferCharacteristicBT2020_10    ColorTransferCharacteristic = C.AVCOL_TRC_BT2020_10    ///< ITU-R BT2020 for 10-bit system
	ColorTransferCharacteristicBT2020_12    ColorTransferCharacteristic = C.AVCOL_TRC_BT2020_12    ///< ITU-R BT2020 for 12-bit system
	ColorTransferCharacteristicSMPTE2084    ColorTransferCharacteristic = C.AVCOL_TRC_SMPTE2084    ///< SMPTE ST 2084 for 10-, 12-, 14- and 16-bit systems
	ColorTransferCharacteristicSMPTEST2084  ColorTransferCharacteristic = C.AVCOL_TRC_SMPTEST2084  //
	ColorTransferCharacteristicSMPTE428     ColorTransferCharacteristic = C.AVCOL_TRC_SMPTE428     //< SMPTE ST 428-1
	ColorTransferCharacteristicSMPTEST428_1 ColorTransferCharacteristic = C.AVCOL_TRC_SMPTEST428_1 //
	ColorTransferCharacteristicARIB_STD_B67 ColorTransferCharacteristic = C.AVCOL_TRC_ARIB_STD_B67 //< ARIB STD-B67, known as "Hybrid log-gamma"
	ColorTransferCharacteristicNB           ColorTransferCharacteristic = C.AVCOL_TRC_NB           ///< Not part of ABI
)

// YUV colorspace type.
//
// These values match the ones defined by ISO/IEC 23091-2_2019 subclause 8.3.
type ColorSpace int

const (
	ColorSpaceRGB                ColorSpace = C.AVCOL_SPC_RGB                // order of coefficients is actually GBR, also IEC 61966-2-1 (sRGB), YZX and ST 428-1
	ColorSpaceBT709              ColorSpace = C.AVCOL_SPC_BT709              // also ITU-R BT1361 / IEC 61966-2-4 xvYCC709 / derived in SMPTE RP 177 Annex B
	ColorSpaceUnspecified        ColorSpace = C.AVCOL_SPC_UNSPECIFIED        //
	ColorSpaceReserved           ColorSpace = C.AVCOL_SPC_RESERVED           // reserved for future use by ITU-T and ISO/IEC just like 15-255 are
	ColorSpaceFCC                ColorSpace = C.AVCOL_SPC_FCC                // FCC Title 47 Code of Federal Regulations 73.682 (a)(20)
	ColorSpaceBT470BG            ColorSpace = C.AVCOL_SPC_BT470BG            // also ITU-R BT601-6 625 / ITU-R BT1358 625 / ITU-R BT1700 625 PAL & SECAM / IEC 61966-2-4 xvYCC601
	ColorSpaceSMPTE170M          ColorSpace = C.AVCOL_SPC_SMPTE170M          // also ITU-R BT601-6 525 / ITU-R BT1358 525 / ITU-R BT1700 NTSC / functionally identical to above
	ColorSpaceSMPTE240M          ColorSpace = C.AVCOL_SPC_SMPTE240M          // derived from 170M primaries and D65 white point, 170M is derived from BT470 System M's primaries
	ColorSpaceYCGCO              ColorSpace = C.AVCOL_SPC_YCGCO              // used by Dirac / VC-2 and H.264 FRext, see ITU-T SG16
	ColorSpaceYCOCG              ColorSpace = C.AVCOL_SPC_YCOCG              //
	ColorSpaceBT2020_NCL         ColorSpace = C.AVCOL_SPC_BT2020_NCL         // ITU-R BT2020 non-constant luminance system
	ColorSpaceBT2020_CL          ColorSpace = C.AVCOL_SPC_BT2020_CL          // ITU-R BT2020 constant luminance system
	ColorSpaceSMPTE2085          ColorSpace = C.AVCOL_SPC_SMPTE2085          // SMPTE 2085, Y'D'zD'x
	ColorSpaceCHROMA_DERIVED_NCL ColorSpace = C.AVCOL_SPC_CHROMA_DERIVED_NCL // Chromaticity-derived non-constant luminance system
	ColorSpaceCHROMA_DERIVED_CL  ColorSpace = C.AVCOL_SPC_CHROMA_DERIVED_CL  // Chromaticity-derived constant luminance system
	ColorSpaceICTCP              ColorSpace = C.AVCOL_SPC_ICTCP              // ITU-R BT.2100-0, ICtCp
	ColorSpaceIPT_C2             ColorSpace = C.AVCOL_SPC_IPT_C2             // SMPTE ST 2128, IPT-C2
	ColorSpaceYCGCO_RE           ColorSpace = C.AVCOL_SPC_YCGCO_RE           // YCgCo-R, even addition of bits
	ColorSpaceYCGCO_RO           ColorSpace = C.AVCOL_SPC_YCGCO_RO           // YCgCo-R, odd addition of bits
	ColorSpaceNB                 ColorSpace = C.AVCOL_SPC_NB                 // Not part of ABI
)

// Visual content value range.
//
// These values are based on definitions that can be found in multiple
// specifications, such as ITU-T BT.709 (3.4 - Quantization of RGB, luminance
// and colour-difference signals), ITU-T BT.2020 (Table 5 - Digital
// Representation) as well as ITU-T BT.2100 (Table 9 - Digital 10- and 12-bit
// integer representation). At the time of writing, the BT.2100 one is
// recommended, as it also defines the full range representation.
//
// Common definitions:
//   - For RGB and luma planes such as Y in YCbCr and I in ICtCp,
//     'E' is the original value in range of 0.0 to 1.0.
//   - For chroma planes such as Cb,Cr and Ct,Cp, 'E' is the original
//     value in range of -0.5 to 0.5.
//   - 'n' is the output bit depth.
//   - For additional definitions such as rounding and clipping to valid n
//     bit unsigned integer range, please refer to BT.2100 (Table 9).
type ColorRange int

const (
	ColorRangeUnspecified ColorRange = C.AVCOL_RANGE_UNSPECIFIED

	// Narrow or limited range content.
	//
	// - For luma planes:
	//
	//       (219 * E + 16) * 2^(n-8)
	//
	//   F.ex. the range of 16-235 for 8 bits
	//
	// - For chroma planes:
	//
	//       (224 * E + 128) * 2^(n-8)
	//
	//   F.ex. the range of 16-240 for 8 bits
	ColorRangeMPEG ColorRange = C.AVCOL_RANGE_MPEG

	//
	// Full range content.
	//
	// - For RGB and luma planes:
	//
	//       (2^n - 1) * E
	//
	//   F.ex. the range of 0-255 for 8 bits
	//
	// - For chroma planes:
	//
	//       (2^n - 1) * E + 2^(n - 1)
	//
	//   F.ex. the range of 1-255 for 8 bits
	//
	ColorRangeJPEG ColorRange = C.AVCOL_RANGE_JPEG
)

// stupid ass go format makes it so this gets auto formated when using //.

/*
*

	Location of chroma samples.

	Illustration showing the location of the first (top left) chroma sample of the
	image, the left shows only luma, the right
	shows the location of the chroma sample, the 2 could be imagined to overlay
	each other but are drawn separately due to limitations of ASCII

	                1st 2nd       1st 2nd horizontal luma sample positions
	                 v   v         v   v
	                 ______        ______
	1st luma line > |X   X ...    |3 4 X ...     X are luma samples,
	                |             |1 2           1-6 are possible chroma positions
	2nd luma line > |X   X ...    |5 6 X ...     0 is undefined/unknown position
*/
type ChromaLocation int

const (
	ChromaLocationUnspecified ChromaLocation = C.AVCHROMA_LOC_UNSPECIFIED //
	ChromaLocationLeft        ChromaLocation = C.AVCHROMA_LOC_LEFT        // MPEG-2/4 4:2:0, H.264 default for 4:2:0
	ChromaLocationCenter      ChromaLocation = C.AVCHROMA_LOC_CENTER      // MPEG-1 4:2:0, JPEG 4:2:0, H.263 4:2:0
	ChromaLocationTopLeft     ChromaLocation = C.AVCHROMA_LOC_TOPLEFT     // ITU-R 601, SMPTE 274M 296M S314M(DV 4:1:1), mpeg2 4:2:2
	ChromaLocationTop         ChromaLocation = C.AVCHROMA_LOC_TOP         //
	ChromaLocationBottomleft  ChromaLocation = C.AVCHROMA_LOC_BOTTOMLEFT  //
	ChromaLocationBottom      ChromaLocation = C.AVCHROMA_LOC_BOTTOM      //
	ChromaLocationNB          ChromaLocation = C.AVCHROMA_LOC_NB          // Not part of ABI
)
