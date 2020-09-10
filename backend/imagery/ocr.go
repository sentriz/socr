package imagery

import (
	"fmt"

	"github.com/otiai10/gosseract/v2"
)

func ExtractText(img []byte) ([]*Block, error) {
	client := gosseract.NewClient()
	defer client.Close()
	if err := client.SetImageFromBytes(img); err != nil {
		return nil, fmt.Errorf("set image bytes: %w", err)
	}

	boxes, err := client.GetBoundingBoxes(gosseract.RIL_TEXTLINE)
	if err != nil {
		return nil, fmt.Errorf("get bounding boxes: %w", err)
	}

	var blocks []*Block
	for _, block := range boxes {
		blocks = append(blocks, &Block{
			Position: [...]int{
				block.Box.Min.X, block.Box.Min.Y,
				block.Box.Max.X, block.Box.Max.Y,
			},
			Text: block.Word,
		})
	}
	return blocks, nil
}
