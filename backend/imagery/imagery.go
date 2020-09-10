package imagery

import (
	"bytes"
	"fmt"
	"image"
	"image/gif"
	"image/jpeg"
	"image/png"
	"io"
	"net/http"
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
	Decode   func(io.Reader) (image.Image, error)
	Encode   func(io.Writer, image.Image) error
}

func EncodeGIF(in io.Writer, i image.Image) error  { return gif.Encode(in, i, nil) }
func EncodePNG(in io.Writer, i image.Image) error  { return png.Encode(in, i) }
func EncodeJPEG(in io.Writer, i image.Image) error { return jpeg.Encode(in, i, nil) }

func FormatFromMIME(in string) (Format, bool) {
	data := map[string]Format{
		"image/gif":  {FiletypeGIF, gif.Decode, EncodeGIF},
		"image/png":  {FiletypePNG, png.Decode, EncodePNG},
		"image/jpeg": {FiletypeJPEG, jpeg.Decode, EncodeJPEG},
	}
	f, ok := data[in]
	return f, ok
}

type Screenshot struct {
	Filetype   Filetype
	Dimensions Dimensions `json:"dimensions"`
	Blocks     []*Block   `json:"blocks"`
	Blurhash   string     `json:"blurhash"`
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

func ProcessBytes(raw []byte) (*Screenshot, error) {
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
	imagePNG := &bytes.Buffer{}
	if err := EncodePNG(imagePNG, imageBig); err != nil {
		return nil, fmt.Errorf("encode scaled and greyed image: %w", err)
	}

	scrotBlocks, err := ExtractText(imagePNG.Bytes(), ScaleFactor)
	if err != nil {
		return nil, fmt.Errorf("extract image text: %w", err)
	}

	scrotBlurhash, err := CalculateBlurhash(image)
	if err != nil {
		return nil, fmt.Errorf("calculate blurhash: %w", err)
	}

	return &Screenshot{
		Filetype: format.Filetype,
		Dimensions: Dimensions{
			Width:  image.Bounds().Size().X,
			Height: image.Bounds().Size().Y,
		},
		Blurhash: scrotBlurhash,
		Blocks:   scrotBlocks,
	}, nil
}
