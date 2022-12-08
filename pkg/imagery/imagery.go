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
	"net/http"
	"strconv"
	"strings"

	"github.com/bakape/thumbnailer"
	"github.com/buckket/go-blurhash"
	"github.com/cenkalti/dominantcolor"
	"github.com/cespare/xxhash"
	"github.com/nfnt/resize"
	gosseract "github.com/otiai10/gosseract/v2"
)

type MediaType string

const (
	TypeImage MediaType = "image"
	TypeVideo MediaType = "video"
)

type Media interface {
	Type() MediaType
	MIME() string
	Hash() string
	Extension() string
	Thumbnail(w, h uint) image.Image
	Image() image.Image
}

const VideoThumbMaxWidth = 1080
const VideoThumbMaxHeight = 1920

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

func ResizeFactor(img image.Image, factor int) image.Image {
	return resize.Resize(
		uint(img.Bounds().Max.X*factor), 0,
		img, resize.Lanczos3,
	)
}

func Resize(img image.Image, width, height uint) image.Image {
	return resize.Resize(width, height, img, resize.Lanczos3)
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
		ThumbDims: thumbnailer.Dims{Width: VideoThumbMaxWidth, Height: VideoThumbMaxHeight},
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

func NewMedia(raw []byte) (Media, error) {
	switch mime := http.DetectContentType(raw); mime {
	case "image/gif", "image/png", "image/jpeg":
		return newMediaImage(raw, mime)
	case "video/webm", "video/mp4", "video/mpeg":
		return newMediaVideo(raw, mime)
	default:
		return nil, fmt.Errorf("unknown image or video mime %q", mime)
	}
}

type mediaImage struct {
	image image.Image
	mime  string
	hash  string
}

func newMediaImage(raw []byte, mime string) (*mediaImage, error) {
	image, err := decodeImage(raw, mime)
	if err != nil {
		return nil, fmt.Errorf("decode: %w", err)
	}
	return &mediaImage{image, mime, hashBytes(raw)}, err
}

func (m *mediaImage) Type() MediaType                 { return TypeImage }
func (m *mediaImage) MIME() string                    { return m.mime }
func (m *mediaImage) Hash() string                    { return m.hash }
func (m *mediaImage) Extension() string               { return mimeExtension(m.mime) }
func (m *mediaImage) Image() image.Image              { return m.image }
func (m *mediaImage) Thumbnail(w, h uint) image.Image { return Resize(m.image, w, h) }

type mediaVideo struct {
	image image.Image
	mime  string
	hash  string
}

func newMediaVideo(raw []byte, mime string) (*mediaVideo, error) {
	image, err := VideoThumbnail(raw)
	if err != nil {
		return nil, fmt.Errorf("get thumbnail: %w", err)
	}
	return &mediaVideo{image, mime, hashBytes(raw)}, err
}

func (m *mediaVideo) Type() MediaType                 { return TypeVideo }
func (m *mediaVideo) MIME() string                    { return m.mime }
func (m *mediaVideo) Hash() string                    { return m.hash }
func (m *mediaVideo) Extension() string               { return mimeExtension(m.mime) }
func (m *mediaVideo) Image() image.Image              { return m.image }
func (m *mediaVideo) Thumbnail(w, h uint) image.Image { return Resize(m.image, w, h) }

func hashBytes(bytes []byte) string {
	sum := xxhash.Sum64(bytes)
	hash := strconv.FormatUint(sum, 16)
	return hash
}

func mimeExtension(mime string) string {
	_, name, _ := strings.Cut(mime, "/")
	return name
}

func decodeImage(raw []byte, mime string) (image.Image, error) {
	switch mime {
	case "image/gif":
		return gif.Decode(bytes.NewReader(raw))
	case "image/png":
		return png.Decode(bytes.NewReader(raw))
	case "image/jpeg":
		return jpeg.Decode(bytes.NewReader(raw))
	default:
		return nil, fmt.Errorf("unknown mime: %q", mime)
	}
}

var _ Media = (*mediaImage)(nil)
var _ Media = (*mediaVideo)(nil)
