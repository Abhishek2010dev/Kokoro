package kokoro

import (
	"encoding/xml"

	"github.com/bytedance/sonic"
	"github.com/fxamacker/cbor/v2"
	"github.com/goccy/go-yaml"
	"github.com/pelletier/go-toml/v2"
)

// EncoderFunc defines the signature for encoding any Go value into JSON, XML, etc.
// It returns a byte slice that can be directly written to the response.
type EncoderFunc func(v any) ([]byte, error)

// DecoderFunc defines the signature for decoding request body bytes into a Go struct.
type DecoderFunc func(data []byte, v any) error

// JSON
func defaultJsonEncoder(v any) ([]byte, error) {
	return sonic.Marshal(v)
}

func defaultJsonDecoder(data []byte, v any) error {
	return sonic.Unmarshal(data, v)
}

// XML
func defaultXMLEncoder(v any) ([]byte, error) {
	return xml.MarshalIndent(v, "", "  ")
}

func defaultXMLDecoder(data []byte, v any) error {
	return xml.Unmarshal(data, v)
}

// YAML
func defaultYamlEncoder(v any) ([]byte, error) {
	return yaml.Marshal(v)
}

func defaultYamlDecoder(data []byte, v any) error {
	return yaml.Unmarshal(data, v)
}

// TOML
func defaultTomlEncoder(v any) ([]byte, error) {
	return toml.Marshal(v)
}

func defaultTomlDecoder(data []byte, v any) error {
	return toml.Unmarshal(data, v)
}

// CBOR
func defaultCborEncoder(v any) ([]byte, error) {
	return cbor.Marshal(v)
}

func defaultCborDecoder(data []byte, v any) error {
	return cbor.Unmarshal(data, v)
}
