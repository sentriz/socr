//nolint:gochecknoglobals
package imagery

import (
	"bytes"
	"encoding/csv"
	"errors"
	"fmt"
	"image"
	"image/color"
	"io"
	"net/http"
	"os"
	"os/exec"
	"strconv"
	"strings"

	"github.com/buckket/go-blurhash"
	"github.com/cenkalti/dominantcolor"
	"github.com/cespare/xxhash"
	"github.com/nfnt/resize"
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

type BoundingBox struct {
	Box  image.Rectangle
	Word string
}

func ExtractText(img []byte) ([]BoundingBox, error) {
	cmd := exec.Command("tesseract", "stdin", "stdout", "--psm", "1", "tsv") //nolint:gosec,noctx
	cmd.Stdin = bytes.NewReader(img)

	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	if err := cmd.Run(); err != nil {
		return nil, fmt.Errorf("run tesseract: %w: %s", err, strings.TrimSpace(stderr.String()))
	}

	return parseTesseractTSV(&stdout)
}

func parseTesseractTSV(r io.Reader) ([]BoundingBox, error) {
	cr := csv.NewReader(r)
	cr.Comma = '\t'
	cr.LazyQuotes = true
	cr.FieldsPerRecord = -1

	if _, err := cr.Read(); err != nil { // header
		if errors.Is(err, io.EOF) {
			return nil, nil
		}
		return nil, fmt.Errorf("read tsv header: %w", err)
	}

	type lineKey struct{ block, par, line int }
	type lineData struct {
		box   image.Rectangle
		words []string
	}
	lines := map[lineKey]*lineData{}
	var order []lineKey

	for {
		rec, err := cr.Read()
		if errors.Is(err, io.EOF) {
			break
		}
		if err != nil {
			return nil, fmt.Errorf("read tsv row: %w", err)
		}
		if len(rec) < 12 {
			continue
		}
		level, _ := strconv.Atoi(rec[0])
		block, _ := strconv.Atoi(rec[2])
		par, _ := strconv.Atoi(rec[3])
		line, _ := strconv.Atoi(rec[4])
		left, _ := strconv.Atoi(rec[6])
		top, _ := strconv.Atoi(rec[7])
		width, _ := strconv.Atoi(rec[8])
		height, _ := strconv.Atoi(rec[9])
		text := rec[11]

		key := lineKey{block, par, line}
		ld, ok := lines[key]
		if !ok {
			ld = &lineData{}
			lines[key] = ld
			order = append(order, key)
		}
		switch level {
		case 4: // line
			ld.box = image.Rect(left, top, left+width, top+height)
		case 5: // word
			if strings.TrimSpace(text) != "" {
				ld.words = append(ld.words, text)
			}
		}
	}

	out := make([]BoundingBox, 0, len(order))
	for _, k := range order {
		ld := lines[k]
		word := strings.Join(ld.words, " ")
		if strings.TrimSpace(word) == "" {
			continue
		}
		out = append(out, BoundingBox{Box: ld.box, Word: word})
	}
	return out, nil
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
	tmp, err := os.CreateTemp("", "")
	if err != nil {
		return nil, fmt.Errorf("create temp: %w", err)
	}
	if _, err := tmp.Write(data); err != nil {
		return nil, fmt.Errorf("write video to tmp: %w", err)
	}
	tmp.Close()
	defer os.Remove(tmp.Name())

	cmd := exec.Command("ffmpeg", "-i", tmp.Name(), "-vframes", "1", "-f", "image2pipe", "-") //nolint:gosec,noctx

	var buff bytes.Buffer
	cmd.Stdout = &buff

	if err := cmd.Run(); err != nil {
		return nil, fmt.Errorf("run ffmpeg: %w", err)
	}

	img, _, err := image.Decode(&buff)
	if err != nil {
		return nil, fmt.Errorf("decode thumbnail: %w", err)
	}

	return img, nil
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
	image, _, err := image.Decode(bytes.NewReader(raw))
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

var _ Media = (*mediaImage)(nil)
var _ Media = (*mediaVideo)(nil)
