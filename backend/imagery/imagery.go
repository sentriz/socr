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
	"net/http"

	"github.com/buckket/go-blurhash"
	"github.com/cenkalti/dominantcolor"
	"github.com/nfnt/resize"
	gosseract "github.com/otiai10/gosseract/v2"
)

type Filetype string

const (
	FiletypeGIF  Filetype = "gif"
	FiletypePNG  Filetype = "png"
	FiletypeJPEG Filetype = "jpg"
)

type EncodeFunc func(io.Writer, image.Image) error
type DecodeFunc func(io.Reader) (image.Image, error)

type Format struct {
	Filetype Filetype
	Decode   DecodeFunc
	Encode   EncodeFunc
}

func EncodeGIF(in io.Writer, i image.Image) error  { return gif.Encode(in, i, nil) }
func EncodePNG(in io.Writer, i image.Image) error  { return png.Encode(in, i) }
func EncodeJPEG(in io.Writer, i image.Image) error { return jpeg.Encode(in, i, nil) }

var (
	FormatGIF  = Format{FiletypeGIF, gif.Decode, EncodeGIF}
	FormatPNG  = Format{FiletypePNG, png.Decode, EncodePNG}
	FormatJPEG = Format{FiletypeJPEG, jpeg.Decode, EncodeJPEG}
)

func FormatFromMIME(in string) (Format, bool) {
	data := map[string]Format{
		"image/gif":  FormatGIF,
		"image/png":  FormatPNG,
		"image/jpeg": FormatJPEG,
	}
	f, ok := data[in]
	return f, ok
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

func ScaleDownRect(rect image.Rectangle) [4]int {
	return [...]int{
		rect.Min.X / ScaleFactor, rect.Min.Y / ScaleFactor,
		rect.Max.X / ScaleFactor, rect.Max.Y / ScaleFactor,
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

type Dimensions struct {
	Height int `json:"height"`
	Width  int `json:"width"`
}

type Block struct {
	// [x1 y1 x2 y2]
	Position [4]int `json:"position"`
	Text     string `json:"text"`
}

type Properties struct {
	Format         Format     `json:"-"`
	Dimensions     Dimensions `json:"dimensions"`
	DominantColour string     `json:"dominant_colour"`
	Blurhash       string     `json:"blurhash"`
	Blocks         []*Block   `json:"blocks"`
}

func Process(raw []byte) (*Properties, error) {
	mime := http.DetectContentType(raw)
	format, ok := FormatFromMIME(mime)
	if !ok {
		return nil, fmt.Errorf("unrecognised format: %s", mime)
	}

	rawReader := bytes.NewReader(raw)
	image, err := format.Decode(rawReader)
	if err != nil {
		return nil, fmt.Errorf("decoding: %s", mime)
	}

	imageGrey := GreyScale(image)
	imageBig := Resize(imageGrey, ScaleFactor)
	imageEncoded := &bytes.Buffer{}
	if err := FormatPNG.Encode(imageEncoded, imageBig); err != nil {
		return nil, fmt.Errorf("encode scaled and greyed image: %w", err)
	}

	blocks, err := ExtractText(imageEncoded.Bytes(), ScaleFactor)
	if err != nil {
		return nil, fmt.Errorf("extract image text: %w", err)
	}

	propBlocks := []*Block{}
	for _, block := range blocks {
		propBlocks = append(propBlocks, &Block{
			Position: ScaleDownRect(block.Box),
			Text:     block.Word,
		})
	}

	_, propDominantColour := DominantColour(image)

	propBlurhash, err := CalculateBlurhash(image)
	if err != nil {
		return nil, fmt.Errorf("calculate blurhash: %w", err)
	}

	return &Properties{
		Dimensions: Dimensions{
			Width:  image.Bounds().Size().X,
			Height: image.Bounds().Size().Y,
		},
		DominantColour: propDominantColour,
		Blurhash:       propBlurhash,
		Blocks:         propBlocks,
	}, nil
}
