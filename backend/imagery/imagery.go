//nolint:gochecknoglobals
package imagery

import (
	"fmt"
	"image"
	"image/color"
	"image/gif"
	"image/jpeg"
	"image/png"
	"io"

	"github.com/buckket/go-blurhash"
	"github.com/cenkalti/dominantcolor"
	"github.com/nfnt/resize"
	gosseract "github.com/otiai10/gosseract/v2"
)

type Filetype struct {
	MIME      string
	Extension string
}

type EncodeFunc func(io.Writer, image.Image) error
type DecodeFunc func(io.Reader) (image.Image, error)

type Format struct {
	Decode DecodeFunc
	Encode EncodeFunc
}

func EncodeGIF(in io.Writer, i image.Image) error  { return gif.Encode(in, i, nil) }
func EncodePNG(in io.Writer, i image.Image) error  { return png.Encode(in, i) }
func EncodeJPEG(in io.Writer, i image.Image) error { return jpeg.Encode(in, i, nil) }

func ImageFromMIME(in string) (*Filetype, *Format) {
	switch in {
	case "image/gif":
		return &Filetype{"image/gif", "gif"}, &Format{gif.Decode, EncodeGIF}
	case "image/png":
		return &Filetype{"image/png", "png"}, &Format{png.Decode, EncodePNG}
	case "image/jpeg":
		return &Filetype{"image/jpeg", "jpg"}, &Format{jpeg.Decode, EncodeJPEG}
	default:
		return nil, nil
	}
}

func VideoFromMIME(in string) *Filetype {
	switch in {
	case "video/webm":
		return &Filetype{"video/webm", "webm"}
	case "video/mp4":
		return &Filetype{"video/mp4", "mp4"}
	case "video/mpeg":
		return &Filetype{"video/mpeg", "mpeg"}
	default:
		return nil
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
