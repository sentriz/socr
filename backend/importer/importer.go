package importer

import (
	"bytes"
	"errors"
	"fmt"
	"image"
	"net/http"
	"strconv"
	"time"

	"github.com/cespare/xxhash"
	"github.com/jackc/pgx/v4"

	"go.senan.xyz/socr/backend/db"
	"go.senan.xyz/socr/backend/imagery"
)

type Decoded struct {
	Hash   string
	Image  image.Image
	Data   []byte
	Format imagery.Format
}

// DecodeImage takes a raw byte slice of an image (png/jpg/etc) decodes it using an appropriate format
// and encodes it again. Encoding and decoding makes sure the hash will be the same for the same
// image given different sources. clipboard/filesystem/etc
func DecodeImage(raw []byte) (*Decoded, error) {
	mime := http.DetectContentType(raw)
	format, ok := imagery.FormatFromMIME(mime)
	if !ok {
		return nil, fmt.Errorf("unknown image mime %s", mime)
	}
	image, err := format.Decode(bytes.NewReader(raw))
	if err != nil {
		return nil, fmt.Errorf("decoding image %w", err)
	}
	dataBuff := &bytes.Buffer{}
	if err := format.Encode(dataBuff, image); err != nil {
		return nil, fmt.Errorf("encoding image: %w", err)
	}
	data := dataBuff.Bytes()
	hash := Hash(data)
	return &Decoded{
		Hash:   hash,
		Image:  image,
		Data:   data,
		Format: format,
	}, nil
}

type Importer struct {
	DB      *db.DB
	Updates chan string
}

func (i *Importer) ImportScreenshot(decoded *Decoded, timestamp time.Time, dirAlias, fileName string) error {
	// insert screenshot and dir info, alert clients with update
	id, isOld, err := i.importScreenshot(decoded.Hash, decoded.Image, timestamp)
	if err != nil {
		return fmt.Errorf("insert screenshot: %w", err)
	}
	if err := i.importScreenshotDirInfo(id, dirAlias, fileName); err != nil {
		return fmt.Errorf("insert dir info: %w", err)
	}
	i.Updates <- decoded.Hash

	if isOld {
		return nil
	}

	// insert blocks, alert clients with update
	if err := i.importScreenshotBlocks(id, decoded.Image); err != nil {
		return fmt.Errorf("insert blocks: %w", err)
	}
	i.Updates <- decoded.Hash

	return nil
}

func (i *Importer) importScreenshot(hash string, image image.Image, timestamp time.Time) (int, bool, error) {
	old, err := i.DB.GetScreenshotByHash(hash)
	if err != nil && !errors.Is(err, pgx.ErrNoRows) {
		return 0, false, fmt.Errorf("getting screenshot by hash: %w", err)
	}
	if err == nil {
		return old.ID, true, nil
	}

	_, propDominantColour := imagery.DominantColour(image)

	propBlurhash, err := imagery.CalculateBlurhash(image)
	if err != nil {
		return 0, false, fmt.Errorf("calculate blurhash: %w", err)
	}

	propSize := image.Bounds().Size()
	new, err := i.DB.CreateScreenshot(&db.Screenshot{
		Hash:           hash,
		Timestamp:      timestamp,
		DimWidth:       propSize.X,
		DimHeight:      propSize.Y,
		DominantColour: propDominantColour,
		Blurhash:       propBlurhash,
	})
	if err != nil {
		return 0, false, fmt.Errorf("inserting screenshot: %w", err)
	}

	return new.ID, false, nil
}

func (i *Importer) importScreenshotBlocks(id int, image image.Image) error {
	imageGrey := imagery.GreyScale(image)
	imageBig := imagery.Resize(imageGrey, imagery.ScaleFactor)
	imageEncoded := &bytes.Buffer{}
	if err := imagery.FormatPNG.Encode(imageEncoded, imageBig); err != nil {
		return fmt.Errorf("encode scaled and greyed image: %w", err)
	}
	rawBlocks, err := imagery.ExtractText(imageEncoded.Bytes())
	if err != nil {
		return fmt.Errorf("extract image text: %w", err)
	}

	blocks := make([]*db.Block, 0, len(rawBlocks))
	for idx, rawBlock := range rawBlocks {
		rect := imagery.ScaleDownRect(rawBlock.Box)
		blocks = append(blocks, &db.Block{
			ScreenshotID: id,
			Index:        idx,
			MinX:         rect.Min.X,
			MinY:         rect.Min.Y,
			MaxX:         rect.Max.X,
			MaxY:         rect.Max.Y,
			Body:         rawBlock.Word,
		})
	}
	if err := i.DB.CreateBlocks(blocks); err != nil {
		return fmt.Errorf("inserting blocks: %w", err)
	}
	return nil
}

func (i *Importer) importScreenshotDirInfo(id int, dirAlias string, fileName string) error {
	_, err := i.DB.CreateDirInfo(&db.DirInfo{
		ScreenshotID:   id,
		Filename:       fileName,
		DirectoryAlias: dirAlias,
	})
	if err != nil {
		return fmt.Errorf("insert info dir infos: %w", err)
	}
	return nil
}

func Hash(bytes []byte) string {
	sum := xxhash.Sum64(bytes)
	format := strconv.FormatUint(sum, 16)
	return format
}
