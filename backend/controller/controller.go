package controller

import (
	"bytes"
	"encoding/json"
	"fmt"
	"image"
	"image/gif"
	"image/jpeg"
	"image/png"
	"io"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"net/http/httputil"
	"net/url"
	"path/filepath"
	"time"

	"github.com/buckket/go-blurhash"
	"github.com/gorilla/mux"

	"github.com/blevesearch/bleve"
	"github.com/otiai10/gosseract/v2"
)

const idLength = 32
const idPool = "" +
	"abcdefghijklmnopqrstuvwxyz" +
	"ABCDEFGHIJKLMNOPQRSTUVWXYZ" +
	"012345679"

func newID() string {
	bytes := make([]byte, idLength)
	for i := 0; i < len(bytes); i++ {
		bytes[i] = idPool[rand.Intn(len(idPool))]
	}
	return string(bytes)
}

func init() {
	rand.Seed(time.Now().UnixNano())
}

type format struct {
	extension string
	decode    func(io.Reader) (image.Image, error)
	encode    func(io.Writer, image.Image) error
}

func encodeGIF(in io.Writer, i image.Image) error  { return gif.Encode(in, i, nil) }
func encodePNG(in io.Writer, i image.Image) error  { return png.Encode(in, i) }
func encodeJPEG(in io.Writer, i image.Image) error { return jpeg.Encode(in, i, nil) }

func formatFromMIME(in string) (format, bool) {
	data := map[string]format{
		"image/gif":  {"gif", gif.Decode, encodeGIF},
		"image/png":  {"png", png.Decode, encodePNG},
		"image/jpeg": {"jpg", jpeg.Decode, encodeJPEG},
	}
	f, ok := data[in]
	return f, ok
}

type Screenshot struct {
	ID         string    `json:"id"`
	Extension  string    `json:"extension"`
	Timestamp  time.Time `json:"timestamp"`
	Blurhash   string    `json:"blurhash"`
	Blocks     []*Block  `json:"blocks"`
	Tags       []string  `json:"tags"`
	Dimensions struct {
		Height int `json:"height"`
		Width  int `json:"width"`
	} `json:"size"`
}

type Block struct {
	Position [4]int `json:"position"`
	Text     string `json:"text"`
}

type Controller struct {
	ScreenshotsDir string
	FrontendDir    string
	FrontendURL    string
	Index          bleve.Index
}

func textBlocksFromBytes(bytes []byte) ([]*Block, error) {
	client := gosseract.NewClient()
	defer client.Close()
	if err := client.SetImageFromBytes(bytes); err != nil {
		return nil, fmt.Errorf("ocr set image bytes: %w", err)
	}
	boxes, err := client.GetBoundingBoxes(gosseract.RIL_TEXTLINE)
	if err != nil {
		return nil, fmt.Errorf("ocr set bounding boxes: %w", err)
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

func (c *Controller) ServeUpload(w http.ResponseWriter, r *http.Request) {
	r.ParseMultipartForm(32 << 20)
	infile, _, err := r.FormFile("image")
	if err != nil {
		http.Error(w, fmt.Sprintf("read form: %v", err), http.StatusInternalServerError)
		return
	}
	defer infile.Close()
	imageBytes, err := ioutil.ReadAll(infile)
	if err != nil {
		http.Error(w, fmt.Sprintf("read form bytes: %v", err), http.StatusInternalServerError)
		return
	}
	mime := http.DetectContentType(imageBytes)
	format, ok := formatFromMIME(mime)
	if !ok {
		http.Error(w, fmt.Sprintf("unrecognised file format: %s", mime), http.StatusInternalServerError)
		return
	}
	imageID := newID()
	imageFilename := fmt.Sprintf("%s.%s", imageID, format.extension)
	imagePath := filepath.Join(c.ScreenshotsDir, imageFilename)
	if err := ioutil.WriteFile(imagePath, imageBytes, 0644); err != nil {
		http.Error(w, fmt.Sprintf("write file bytes: %v", err), http.StatusInternalServerError)
		return
	}
	textBlocks, err := textBlocksFromBytes(imageBytes)
	if err != nil {
		http.Error(w, fmt.Sprintf("get text from image: %v", err), http.StatusInternalServerError)
		return
	}
	hashReader := bytes.NewReader(imageBytes)
	imageRaw, err := format.decode(hashReader)
	if err != nil {
		http.Error(w, fmt.Sprintf("decode image data: %v", err), http.StatusInternalServerError)
		return
	}
	blurhash, err := blurhash.Encode(4, 3, imageRaw)
	if err != nil {
		http.Error(w, fmt.Sprintf("compute blur hash: %v", err), http.StatusInternalServerError)
		return
	}
	screenshot := &Screenshot{
		ID:        imageID,
		Timestamp: time.Now(),
	}
	screenshot.Blurhash = blurhash
	screenshot.Blocks = textBlocks
	screenshot.Dimensions.Width = imageRaw.Bounds().Size().X
	screenshot.Dimensions.Height = imageRaw.Bounds().Size().Y
	c.Index.Index(screenshot.ID, screenshot)
	// json.NewEncoder(w).Encode(screenshot)
	fmt.Fprintf(w, "%s\n", screenshot.ID)
}

func (c *Controller) ServeSearch(w http.ResponseWriter, r *http.Request) {
	queryString := r.URL.Query().Get("q")
	if queryString == "" {
		http.Error(w, "please provide a query", http.StatusBadRequest)
		return
	}
	bleve.NewDocumentStaticMapping()
	query := bleve.NewMatchAllQuery()
	search := bleve.NewSearchRequestOptions(query, 100, 0, false)
	search.Fields = []string{"*"}
	searchResults, err := c.Index.Search(search)
	if err != nil {
		http.Error(w, fmt.Sprintf("searching index: %v", err), http.StatusInternalServerError)
		return
	}
	if err := json.NewEncoder(w).Encode(searchResults); err != nil {
		log.Printf("error encoding response: %v", err)
	}
}

func (c *Controller) ServeImage(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	filename := fmt.Sprintf("%s.png", vars["id"])
	http.ServeFile(w, r, filepath.Join(c.ScreenshotsDir, filename))
}

func (c *Controller) ServeFrontend(w http.ResponseWriter, r *http.Request) {
	switch {
	case c.FrontendURL != "":
		frontendURL, _ := url.Parse(c.FrontendURL)
		proxy := httputil.NewSingleHostReverseProxy(frontendURL)
		proxy.ServeHTTP(w, r)
	case c.FrontendDir != "":
		files := http.FileServer(http.Dir(c.FrontendDir))
		files.ServeHTTP(w, r)
	}
}
