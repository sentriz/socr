//nolint:gochecknoglobals
package imagery

import (
	"bytes"
	"fmt"
	"image"
	"image/color"
	"image/gif"
	"image/jpeg"
	"image/png"
	"io"

	"github.com/bakape/thumbnailer"
	"github.com/buckket/go-blurhash"
	"github.com/cenkalti/dominantcolor"
	"github.com/nfnt/resize"
	gosseract "github.com/otiai10/gosseract/v2"
)

const VideoThumbWidth = 1080

type Type string

const (
	TypeImage Type = "image"
	TypeVideo Type = "video"
)

type MIME string

const (
	MIMEGIF  MIME = "image/gif"
	MIMEPNG  MIME = "image/png"
	MIMEJPEG MIME = "image/jpeg"

	MIMEWebM MIME = "video/webm"
	MIMEMP4  MIME = "video/mp4"
	MIMEMPEG MIME = "video/mpeg"
)

type Filetype struct {
	Type      Type
	MIME      MIME
	Extension string
}

func (f *Filetype) IsImage() bool { return f.Type == TypeImage }
func (f *Filetype) IsVideo() bool { return f.Type == TypeVideo }

type EncodeFunc func(io.Writer, image.Image) error
type DecodeFunc func(io.Reader) (image.Image, error)

type Format struct {
	Decode DecodeFunc
	Encode EncodeFunc
}

func EncodeGIF(in io.Writer, i image.Image) error  { return gif.Encode(in, i, nil) }
func EncodePNG(in io.Writer, i image.Image) error  { return png.Encode(in, i) }
func EncodeJPEG(in io.Writer, i image.Image) error { return jpeg.Encode(in, i, nil) }

func ReadMIME(in string) (*Filetype, *Format) {
	switch MIME(in) {
	case MIMEGIF:
		return &Filetype{TypeImage, MIMEGIF, "gif"}, &Format{gif.Decode, EncodeGIF}
	case MIMEPNG:
		return &Filetype{TypeImage, MIMEPNG, "png"}, &Format{png.Decode, EncodePNG}
	case MIMEJPEG:
		return &Filetype{TypeImage, MIMEJPEG, "jpg"}, &Format{jpeg.Decode, EncodeJPEG}
	case MIMEWebM:
		return &Filetype{TypeVideo, MIMEWebM, "webm"}, nil
	case MIMEMP4:
		return &Filetype{TypeVideo, MIMEMP4, "mp4"}, nil
	case MIMEMPEG:
		return &Filetype{TypeVideo, MIMEMPEG, "mpeg"}, nil
	default:
		return nil, nil
	}
}

func ExtractText(img []byte) ([]gosseract.BoundingBox, error) {
	client := gosseract.NewClient()
	defer client.Close()
	if err := client.SetImageFromBytes(img); err != nil {
		return nil, fmt.Errorf("set image bytes: %w", err)
	}

	if err := client.SetPageSegMode(gosseract.PSM_AUTO_OSD); err != nil {
		return nil, fmt.Errorf("set page setmentation mode: %w", err)
	}

	boxes, err := client.GetBoundingBoxes(gosseract.RIL_TEXTLINE)
	if err != nil {
		return nil, fmt.Errorf("get bounding boxes: %w", err)
	}

	return boxes, nil
}

const (
	ScaleFactor = 3
)

func Resize(img image.Image, factor int) image.Image {
	return resize.Resize(
		uint(img.Bounds().Max.X*factor), 0,
		img, resize.Lanczos3,
	)
}

func ScaleDownRect(rect image.Rectangle) image.Rectangle {
	return image.Rectangle{
		Min: image.Point{X: rect.Min.X / ScaleFactor, Y: rect.Min.Y / ScaleFactor},
		Max: image.Point{X: rect.Max.X / ScaleFactor, Y: rect.Max.Y / ScaleFactor},
	}
}

func GreyScale(img image.Image) *image.Gray {
	bounds := img.Bounds()
	gray := image.NewGray(bounds)
	for x := 0; x < bounds.Max.X; x++ {
		for y := 0; y < bounds.Max.Y; y++ {
			gray.Set(x, y, img.At(x, y))
		}
	}
	return gray
}

const (
	BlurhashXC = 4
	BlurhashYC = 3
)

func CalculateBlurhash(img image.Image) (string, error) {
	return blurhash.Encode(BlurhashXC, BlurhashXC, img)
}

func DominantColour(img image.Image) (color.Color, string) {
	colour := dominantcolor.Find(img)
	hex := dominantcolor.Hex(colour)
	return colour, hex
}

func VideoThumbnail(data []byte) (image.Image, error) {
	_, thumb, err := thumbnailer.ProcessBuffer(data, thumbnailer.Options{
		ThumbDims: thumbnailer.Dims{Width: VideoThumbWidth},
	})
	if err != nil {
		return nil, fmt.Errorf("process buffer: %w", err)
	}
	buff := bytes.NewBuffer(thumb.Data)
	if thumb.IsPNG {
		return png.Decode(buff)
	}
	return jpeg.Decode(buff)
}
