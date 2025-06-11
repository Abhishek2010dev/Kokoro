package kokoro

import (
	"encoding/xml"

	"github.com/bytedance/sonic"
)

// EncoderFunc defines the signature for encoding any Go value into JSON, XML, etc.
// It returns a byte slice that can be directly written to the response.
//
// Fast, zero-copy friendly, and works great with fasthttp.
type EncoderFunc func(v any) ([]byte, error)

// DecoderFunc defines the signature for decoding request body bytes into a Go struct.
// You pass the raw body as []byte, and it populates the value.
//
// Zero-copy friendly and avoids io.Reader overhead.
type DecoderFunc func(data []byte, v any) error

// defaultJsonEncoder is Kokoro's default JSON encoder.
// It uses Sonic for blazing fast encoding.
func defaultJsonEncoder(v any) ([]byte, error) {
	return sonic.Marshal(v)
}

// defaultJsonDecoder is Kokoro's default JSON decoder.
// It uses Sonic's default configuration for decoding.
func defaultJsonDecoder(data []byte, v any) error {
	return sonic.Unmarshal(data, v)
}

// defaultXMLEncoder is Kokoro's default XML encoder using encoding/xml.
func defaultXMLEncoder(v any) ([]byte, error) {
	return xml.MarshalIndent(v, "", "  ")
}

// defaultXMLDecoder is Kokoro's default XML decoder using encoding/xml.
func defaultXMLDecoder(data []byte, v any) error {
	return xml.Unmarshal(data, v)
}
