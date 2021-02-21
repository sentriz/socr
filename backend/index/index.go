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

const (
	FieldTags             = "tags"
	FieldTimestamp        = "timestamp"
	FieldBlocks           = "blocks"
	FieldBlocksText       = "text"
	FieldBlocksPosition   = "position"
	FieldDimensions       = "dimensions"
	FieldDimensionsHeight = "height"
	FieldDimensionsWidth  = "width"
	FieldDominantColour   = "dominant_colour"
)

func path(fields ...string) string {
	return strings.Join(fields, ".")
}

var (
	PathTags             = path(FieldTags)
	PathTimestamp        = path(FieldTimestamp)
	PathBlocksText       = path(FieldBlocks, FieldBlocksText)
	PathBlocksPosition   = path(FieldBlocks, FieldBlocksPosition)
	PathDimensionsHeight = path(FieldDimensions, FieldDimensionsHeight)
	PathDimensionsWidth  = path(FieldDimensions, FieldDimensionsWidth)
	PathDominantColour   = path(FieldDominantColour)
)

var (
	BaseSearchFields = []string{
		PathTags,
		PathTimestamp,
		PathBlocksText,
		PathBlocksPosition,
		PathDimensionsHeight,
		PathDimensionsWidth,
		PathDominantColour,
	}
	BaseHighlightFields = []string{
		PathBlocksText,
	}
	BaseQueryField = PathBlocksText
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
