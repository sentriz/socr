package imagery

import (
	"bytes"
	"image"
	"image/draw"
	"image/png"
	"os/exec"
	"strings"
	"testing"

	"golang.org/x/image/font"
	"golang.org/x/image/font/basicfont"
	"golang.org/x/image/math/fixed"
)

func TestExtractText(t *testing.T) {
	t.Parallel()

	if _, err := exec.LookPath("tesseract"); err != nil {
		t.Skip("tesseract not found in PATH")
	}

	lines := []string{
		"the quick brown fox j",
		"umps over the lazy dog",
		"hello world",
		"hello",
	}

	img := renderLines(t, lines)
	scaled, _ := ScaleForOCR(img)

	var buf bytes.Buffer
	if err := png.Encode(&buf, scaled); err != nil {
		t.Fatalf("encode png: %v", err)
	}

	boxes, err := ExtractText(buf.Bytes())
	if err != nil {
		t.Fatalf("extract text: %v", err)
	}

	if len(boxes) < len(lines) {
		t.Fatalf("expected at least %d lines, got %d: %+v", len(lines), len(boxes), boxes)
	}

	for _, b := range boxes {
		if strings.TrimSpace(b.Word) == "" {
			t.Errorf("got empty word in box %+v", b)
		}
		if b.Box.Empty() {
			t.Errorf("got empty bounding box for %q", b.Word)
		}
	}

	var joined strings.Builder
	for _, b := range boxes {
		joined.WriteString(" ")
		joined.WriteString(b.Word)
	}
	got := strings.ToLower(joined.String())

	for _, line := range lines {
		for _, want := range strings.Fields(line) {
			if !strings.Contains(got, want) {
				t.Errorf("missing %q in extracted text: %q", want, strings.TrimSpace(got))
			}
		}
	}
}

func renderLines(t *testing.T, lines []string) image.Image {
	t.Helper()

	const (
		w          = 500
		lineHeight = 28
		topPad     = 24
		leftPad    = 20
	)
	h := topPad*2 + lineHeight*len(lines)

	img := image.NewRGBA(image.Rect(0, 0, w, h))
	draw.Draw(img, img.Bounds(), image.White, image.Point{}, draw.Src)

	d := &font.Drawer{
		Dst:  img,
		Src:  image.Black,
		Face: basicfont.Face7x13,
	}
	for i, line := range lines {
		d.Dot = fixed.P(leftPad, topPad+lineHeight*(i+1))
		d.DrawString(line)
	}

	return img
}
