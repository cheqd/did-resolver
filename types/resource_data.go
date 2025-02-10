package types

import (
	"encoding/json"
	"github.com/rs/zerolog/log"
)

type ResourceData struct {
	Content map[string]interface{}
}

func NewResourceData(data []byte) (*ResourceData, error) {
	var parsedData map[string]interface{}

	err := json.Unmarshal(data, &parsedData)
	if err != nil {
		return nil, err
	}
	return &ResourceData{Content: parsedData}, nil
}

func (e *ResourceData) AddContext(newProtocol string) {}
func (e *ResourceData) RemoveContext()                {}
func (e *ResourceData) GetBytes() []byte {
	bytes, err := json.Marshal(e.Content)
	if err != nil {
		log.Info().Msg("Failed to marshal resource")
		return []byte{}
	}
	return bytes
}

func (e *ResourceData) MarshalJSON() ([]byte, error) {
	return json.Marshal(e.Content)
}

func (e *ResourceData) GetContentType() string { return "" }
func (e *ResourceData) IsRedirect() bool       { return false }
