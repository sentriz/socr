package controller

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"go.senan.xyz/socr/imagery"

	"github.com/blevesearch/bleve"
	"github.com/gorilla/websocket"
)

type Screenshot struct {
	ID         string    `json:"id"`
	Timestamp  time.Time `json:"timestamp"`
	Filetype   imagery.Filetype
	Tags       []string   `json:"tags"`
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

type Controller struct {
	ScreenshotsPath string
	ImportPath      string
	Index           bleve.Index
	SocketUpgrader  websocket.Upgrader
	SocketClients   map[*websocket.Conn]struct{}
}

func (c *Controller) ReadAndIndexBytes(raw []byte) (*Screenshot, error) {
	mime := http.DetectContentType(raw)
	format, ok := imagery.FormatFromMIME(mime)
	if !ok {
		return nil, fmt.Errorf("unrecognised format: %s", mime)
	}

	rawReader := bytes.NewReader(raw)
	image, err := format.Decode(rawReader)
	if err != nil {
		return nil, fmt.Errorf("decoding: %s", mime)
	}

	imageGrey := imagery.GreyScale(image)
	imageBig := imagery.Resize(imageGrey, imagery.ScaleFactor)
	imagePNG := &bytes.Buffer{}
	if err := imagery.EncodePNG(imagePNG, imageBig); err != nil {
		return nil, fmt.Errorf("encode scaled and greyed image: %w", err)
	}

	scrotBlocksOrig, err := imagery.ExtractText(imagePNG.Bytes(), imagery.ScaleFactor)
	if err != nil {
		return nil, fmt.Errorf("extract image text: %w", err)
	}

	scrotBlocks := []*Block{}
	for _, block := range scrotBlocksOrig {
		scrotBlocks = append(scrotBlocks, &Block{
			Position: imagery.ScaleDownRect(block.Box),
			Text:     block.Word,
		})
	}

	scrotBlurhash, err := imagery.CalculateBlurhash(image)
	if err != nil {
		return nil, fmt.Errorf("calculate blurhash: %w", err)
	}

	scrotID := IDNew()
	scrotFilename := fmt.Sprintf("%s.%s", scrotID, format.Filetype)
	scrotPath := filepath.Join(c.ScreenshotsPath, scrotFilename)
	if err := ioutil.WriteFile(scrotPath, raw, 0644); err != nil {
		return nil, fmt.Errorf("write processed bytes: %w", err)
	}

	screenshot := &Screenshot{
		ID:       scrotID,
		Filetype: format.Filetype,
		Dimensions: Dimensions{
			Width:  image.Bounds().Size().X,
			Height: image.Bounds().Size().Y,
		},
		Blurhash:  scrotBlurhash,
		Blocks:    scrotBlocks,
		Timestamp: time.Now(),
		Tags:      []string{},
	}

	if err := c.Index.Index(screenshot.ID, screenshot); err != nil {
		return nil, fmt.Errorf("indexing screenshot: %w", err)
	}

	return screenshot, nil
}

func procSuffixHas(in string) bool   { return strings.HasSuffix(in, ".processed") }
func procSuffixAdd(in string) string { return fmt.Sprintf("%s.processed", in) }

func (c *Controller) ReadRenameImportFile(file os.FileInfo) ([]byte, error) {
	if procSuffixHas(file.Name()) {
		return nil, nil
	}

	filePath := filepath.Join(c.ImportPath, file.Name())
	bytes, err := ioutil.ReadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("reading from disk: %v", err)
	}

	if err := os.Rename(filePath, procSuffixAdd(filePath)); err != nil {
		return nil, fmt.Errorf("renaming: %v", err)
	}

	return bytes, nil
}

func (c *Controller) EchoSockets() error {
	for {
		for client := range c.SocketClients {
			err := client.WriteMessage(websocket.TextMessage, []byte("hello"))
			if err != nil {
				log.Printf("error writing to socket client: %v", err)
				client.Close()
				delete(c.SocketClients, client)
			}
		}
		time.Sleep(1 * time.Second)
	}
}
