package imagery

import (
	"fmt"

	"github.com/otiai10/gosseract/v2"
)

func ExtractText(img []byte, scale int) ([]*Block, error) {
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
				// TODO: divide by scale factor in caller
				// ExtractText shouldn't care about scaling
				block.Box.Min.X / scale, block.Box.Min.Y / scale,
				block.Box.Max.X / scale, block.Box.Max.Y / scale,
			},
			Text: block.Word,
		})
	}
	return blocks, nil
}
