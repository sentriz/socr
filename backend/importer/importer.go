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

type Importer struct {
	DB             *db.DB
	Updates        chan string
	DefaultEncoder imagery.EncodeFunc
	DefaultMIME    imagery.MIME
	ThumbnailWidth uint
}

func (i *Importer) ImportMedia(decoded *Decoded, timestamp time.Time, dirAlias, fileName string) error {
	id, isOld, err := i.importMedia(decoded.Filetype, decoded.Hash, decoded.Image, timestamp)
	if err != nil {
		return fmt.Errorf("import media with props: %w", err)
	}
	if err := i.importDirInfo(id, dirAlias, fileName); err != nil {
		return fmt.Errorf("import dir info: %w", err)
	}
	i.Updates <- decoded.Hash

	if isOld {
		return nil
	}

	if err := i.importThumbnail(id, decoded.Image); err != nil {
		return fmt.Errorf("import thumbnail: %w", err)
	}
	if err := i.importBlocks(id, decoded.Image); err != nil {
		return fmt.Errorf("import blocks: %w", err)
	}
	i.Updates <- decoded.Hash

	return nil
}

func (i *Importer) importMedia(filetype *imagery.Filetype, hash string, image image.Image, timestamp time.Time) (db.MediaID, bool, error) {
	old, err := i.DB.GetMediaByHash(hash)
	if err != nil && !errors.Is(err, pgx.ErrNoRows) {
		return 0, false, fmt.Errorf("getting media by hash: %w", err)
	}
	if err == nil {
		return old.ID, true, nil
	}

	_, propDominantColour := imagery.DominantColour(image)

	propBlurhash, err := imagery.CalculateBlurhash(image)
	if err != nil {
		return 0, false, fmt.Errorf("calculate blurhash: %w", err)
	}

	propDimensions := image.Bounds().Size()
	new, err := i.DB.CreateMedia(&db.Media{
		Hash:           hash,
		Type:           db.MediaType(filetype.Type),
		MIME:           string(filetype.MIME),
		Timestamp:      timestamp,
		DimWidth:       propDimensions.X,
		DimHeight:      propDimensions.Y,
		DominantColour: propDominantColour,
		Blurhash:       propBlurhash,
	})
	if err != nil {
		return 0, false, fmt.Errorf("inserting media: %w", err)
	}

	return new.ID, false, nil
}

func (i *Importer) importBlocks(id db.MediaID, image image.Image) error {
	imageGrey := imagery.GreyScale(image)
	imageBig := imagery.ResizeFactor(imageGrey, imagery.ScaleFactor)
	imageEncoded := &bytes.Buffer{}
	if err := i.DefaultEncoder(imageEncoded, imageBig); err != nil {
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
			MediaID: id,
			Index:   idx,
			MinX:    rect.Min.X,
			MinY:    rect.Min.Y,
			MaxX:    rect.Max.X,
			MaxY:    rect.Max.Y,
			Body:    rawBlock.Word,
		})
	}
	if err := i.DB.CreateBlocks(blocks); err != nil {
		return fmt.Errorf("inserting blocks: %w", err)
	}
	return nil
}

func (i *Importer) importThumbnail(id db.MediaID, image image.Image) error {
	resized := imagery.Resize(image, i.ThumbnailWidth, 0)
	dimensions := resized.Bounds().Size()

	var data bytes.Buffer
	if err := i.DefaultEncoder(&data, resized); err != nil {
		return fmt.Errorf("encoding thumbnail: %w", err)
	}

	thumbnail := &db.Thumbnail{
		MediaID:   id,
		MIME:      string(i.DefaultMIME),
		DimWidth:  dimensions.X,
		DimHeight: dimensions.Y,
		Timestamp: time.Now(),
		Data:      data.Bytes(),
	}
	if _, err := i.DB.CreateThumbnail(thumbnail); err != nil {
		return fmt.Errorf("insert thumbnail: %w", err)
	}
	return nil
}

func (i *Importer) importDirInfo(id db.MediaID, dirAlias string, fileName string) error {
	dirInfo := &db.DirInfo{
		Filename:       fileName,
		DirectoryAlias: dirAlias,
		MediaID:        id,
	}
	if _, err := i.DB.CreateDirInfo(dirInfo); err != nil {
		return fmt.Errorf("insert info dir infos: %w", err)
	}
	return nil
}

func Hash(bytes []byte) string {
	sum := xxhash.Sum64(bytes)
	format := strconv.FormatUint(sum, 16)
	return format
}

type Decoded struct {
	Hash     string
	Data     []byte
	Filetype *imagery.Filetype
	Image    image.Image
}

func DecodeMedia(raw []byte) (*Decoded, error) {
	mime := http.DetectContentType(raw)
	filetype, format := imagery.ReadMIME(mime)
	if filetype == nil {
		return nil, fmt.Errorf("unknown image or video mime %q", mime)
	}

	var hash string
	var image image.Image
	switch {
	// in the case of an image given a raw encoded image, we need to decode and encode it
	// again before calcuation the hash. this ensures the same image given from different
	// sources (eg. uploaded or scanned) has the same hash
	case filetype.IsImage():
		decodedImage, err := format.Decode(bytes.NewReader(raw))
		if err != nil {
			return nil, fmt.Errorf("decoding image: %w", err)
		}
		buff := &bytes.Buffer{}
		if err := format.Encode(buff, decodedImage); err != nil {
			return nil, fmt.Errorf("encoding image: %w", err)
		}
		hash = Hash(buff.Bytes())
		image = decodedImage

	// in the case of a video, we should just use the first frame as the image, but use
	// the original data's hash
	case filetype.IsVideo():
		thumbImage, err := imagery.VideoThumbnail(raw)
		if err != nil {
			return nil, fmt.Errorf("encoding image: %w", err)
		}
		hash = Hash(raw)
		image = thumbImage

	default:
		return nil, fmt.Errorf("unknown filetype %q", mime)
	}

	return &Decoded{
		Hash:     hash,
		Data:     raw,
		Filetype: filetype,
		Image:    image,
	}, nil
}
