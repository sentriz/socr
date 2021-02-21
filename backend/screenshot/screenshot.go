package screenshot

import (
	"fmt"
	"time"

	"go.senan.xyz/socr/imagery"
)

// alias -> path
type Directories map[string]string

type Screenshot struct {
	ID        string           `json:"id"`
	Timestamp time.Time        `json:"timestamp"`
	Filetype  imagery.Filetype `json:"filetype"`
	Directory string           `json:"directory"`
	Tags      []string         `json:"tags"`
	*imagery.Properties
}

func FromBytesWithHash(bytes []byte, id string, timestamp time.Time) (*Screenshot, error) {
	properties, err := imagery.Process(bytes)
	if err != nil {
		return nil, fmt.Errorf("getting image properties: %w", err)
	}

	return &Screenshot{
		ID:         id,
		Filetype:   properties.Format.Filetype,
		Timestamp:  timestamp,
		Properties: properties,
	}, nil
}
