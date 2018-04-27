package multipart

import (
	"io"
	"mime"
	"strings"

	"github.com/tiny-go/codec"
	"github.com/tiny-go/codec/driver"
)

// DataTypeMultipart contains multipart form codec data type.
const DataTypeMultipart = "multipart/form-data"

// Multipart codec.
type Multipart struct{ boundary string }

// New creates an instance of multipart codec retrieving boundary from the header.
func New(header string) *Multipart {
	mediatype, params, err := mime.ParseMediaType(header)
	if err == nil && strings.Contains(mediatype, DataTypeMultipart) {
		if boundary, ok := params["boundary"]; ok {
			return &Multipart{boundary}
		}
	}
	return &Multipart{}
}

// Encoder returns multipart encoder.
func (m *Multipart) Encoder(w io.Writer) codec.Encoder { return NewEncoder(w) }

// Decoder returns multipart decoder.
func (m *Multipart) Decoder(r io.Reader) codec.Decoder { return NewDecoder(r, m.boundary) }

// MimeType returns multipart header.
func (m *Multipart) MimeType() string { return DataTypeMultipart }

func init() {
	// parse the boundary and create a new instance for every request
	driver.Register(DataTypeMultipart, func(h string) codec.Codec { return New(h) })
}
