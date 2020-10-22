package index

import (
	"errors"
	"fmt"

	"github.com/blevesearch/bleve"
	"github.com/blevesearch/bleve/analysis/analyzer/keyword"
	"github.com/blevesearch/bleve/mapping"
)

var (
	BaseSearchFields = []string{
		"tags",
		"timestamp",
		"blocks.text",
		"blocks.position",
		"dimensions.height",
		"dimensions.width",
	}
	BaseHighlightFields = []string{
		"blocks.text",
	}
	BaseQueryField = "blocks.text"
)

func CreateIndexMapping() *mapping.IndexMappingImpl {
	fieldMapNumeric := bleve.NewNumericFieldMapping()

	fieldMapKeyword := bleve.NewTextFieldMapping()
	fieldMapKeyword.Analyzer = keyword.Name

	fieldMapTime := bleve.NewDateTimeFieldMapping()

	mappingBlocks := bleve.NewDocumentMapping()
	mappingBlocks.AddFieldMappingsAt("text", fieldMapKeyword)
	mappingBlocks.AddFieldMappingsAt("position", fieldMapNumeric)

	mappingScreenshot := bleve.NewDocumentMapping()
	mappingScreenshot.AddFieldMappingsAt("timestamp", fieldMapTime)
	mappingScreenshot.AddFieldMappingsAt("tags", fieldMapKeyword)
	mappingScreenshot.AddSubDocumentMapping("blocks", mappingBlocks)

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
