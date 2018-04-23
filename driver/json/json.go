package json

import (
	"encoding/json"
	"io"

	"github.com/tiny-go/codec"
	"github.com/tiny-go/codec/driver"
)

// DataTypeJSON contains JSON codec data type.
const DataTypeJSON = "application/json"

var global = &JSON{}

// JSON codec.
type JSON struct{}

// Encoder creates JSON encoder.
func (j *JSON) Encoder(w io.Writer) codec.Encoder { return json.NewEncoder(w) }

// Decoder creates JSON decoder.
func (j *JSON) Decoder(r io.Reader) codec.Decoder { return json.NewDecoder(r) }

// MimeType returns the mime type of this codec (which is application/json).
func (j *JSON) MimeType() string { return DataTypeJSON }

func init() {
	// since JSON is a stateless codec (no need to set boundary per request) we can
	// always return the same instance
	driver.Register(DataTypeJSON, func(string) codec.Codec { return global })
}
