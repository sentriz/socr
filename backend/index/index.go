package index

import (
	"errors"
	"fmt"
	"strings"

	"github.com/blevesearch/bleve"
	"github.com/blevesearch/bleve/analysis/analyzer/keyword"
	"github.com/blevesearch/bleve/analysis/lang/en"
	"github.com/blevesearch/bleve/mapping"
)

type Field string

const (
	FieldTags             Field = "tags"
	FieldTimestamp        Field = "timestamp"
	FieldBlocks           Field = "blocks"
	FieldBlocksText       Field = "text"
	FieldBlocksPosition   Field = "position"
	FieldDimensions       Field = "dimensions"
	FieldDimensionsHeight Field = "height"
	FieldDimensionsWidth  Field = "width"
	FieldDominantColour   Field = "dominant_colour"
)

type Path []Field

var (
	PathTags             Path = []Field{FieldTags}
	PathTimestamp        Path = []Field{FieldTimestamp}
	PathBlocksText       Path = []Field{FieldBlocks, FieldBlocksText}
	PathBlocksPosition   Path = []Field{FieldBlocks, FieldBlocksPosition}
	PathDimensionsHeight Path = []Field{FieldDimensions, FieldDimensionsHeight}
	PathDimensionsWidth  Path = []Field{FieldDimensions, FieldDimensionsWidth}
	PathDominantColour   Path = []Field{FieldDominantColour}
)

func (p Path) String() string {
	var s []string
	for _, path := range p {
		s = append(s, string(path))
	}

	return strings.Join(s, ".")
}

var (
	BaseSearchFields = []string{
		PathTags.String(),
		PathTimestamp.String(),
		PathBlocksText.String(),
		PathBlocksPosition.String(),
		PathDimensionsHeight.String(),
		PathDimensionsWidth.String(),
		PathDominantColour.String(),
	}
	BaseHighlightFields = []string{PathBlocksText.String()}
	BaseQueryField      = PathBlocksText.String()
)

func CreateIndexMapping() *mapping.IndexMappingImpl {
	fieldMapNumeric := bleve.NewNumericFieldMapping()

	fieldMapEnglish := bleve.NewTextFieldMapping()
	fieldMapEnglish.Analyzer = en.AnalyzerName

	fieldMapKeyword := bleve.NewTextFieldMapping()
	fieldMapKeyword.Analyzer = keyword.Name

	fieldMapTime := bleve.NewDateTimeFieldMapping()

	mappingBlocks := bleve.NewDocumentMapping()
	mappingBlocks.AddFieldMappingsAt(string(FieldBlocksText), fieldMapEnglish)
	mappingBlocks.AddFieldMappingsAt(string(FieldBlocksPosition), fieldMapNumeric)

	mappingScreenshot := bleve.NewDocumentMapping()
	mappingScreenshot.AddFieldMappingsAt(string(FieldTimestamp), fieldMapTime)
	mappingScreenshot.AddFieldMappingsAt(string(FieldTags), fieldMapKeyword)
	mappingScreenshot.AddFieldMappingsAt(string(FieldDominantColour), fieldMapKeyword)
	mappingScreenshot.AddSubDocumentMapping(string(FieldBlocks), mappingBlocks)

	mappingIndex := bleve.NewIndexMapping()
	mappingIndex.DefaultMapping = mappingScreenshot
	mappingIndex.DefaultField = BaseQueryField

	return mappingIndex
}

func GetOrCreateIndex(path string) (bleve.Index, error) {
	index, err := bleve.Open(path)
	switch {
	case
		errors.Is(err, bleve.ErrorIndexMetaMissing),
		errors.Is(err, bleve.ErrorIndexPathDoesNotExist):
		indexMapping := CreateIndexMapping()
		return bleve.New(path, indexMapping)
	case err != nil:
		return nil, fmt.Errorf("open index: %w", err)
	default:
		return index, nil
	}
}
