package xml

import (
	"encoding/xml"
	"io"

	"github.com/tiny-go/codec"
	"github.com/tiny-go/codec/driver"
)

// DataTypeXML contains XML codec data type.
const DataTypeXML = "application/xml"

var global = &XML{}

// XML codec.
type XML struct{}

// Encoder creates JSON encoder.
func (x *XML) Encoder(w io.Writer) codec.Encoder { return xml.NewEncoder(w) }

// Decoder creates JSON decoder.
func (x *XML) Decoder(r io.Reader) codec.Decoder { return xml.NewDecoder(r) }

// MimeType returns the mime type of this codec (which is application/xml).
func (x *XML) MimeType() string { return DataTypeXML }

func init() {
	// since XML is a stateless codec (no need to set boundary per request) we can
	// always return the same instance
	driver.Register(DataTypeXML, func(string) codec.Codec { return global })
}
