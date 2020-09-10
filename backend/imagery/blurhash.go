package imagery

import (
	"image"

	"github.com/buckket/go-blurhash"
)

const (
	BlurhashXC = 4
	BlurhashYC = 3
)

func CalculateBlurhash(img image.Image) (string, error) {
	return blurhash.Encode(BlurhashXC, BlurhashXC, img)
}
