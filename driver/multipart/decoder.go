package multipart

import (
	"io"
	"mime/multipart"
)

// buffer size (2 mb by default)
const maxMemory int64 = 2 << 20

// Decoder is responsible for decoding multipart data to provided receiver.
type Decoder struct {
	mr *multipart.Reader
}

// NewDecoder creates a new instance of multipart.Decoder.
func NewDecoder(r io.Reader, b string) *Decoder {
	return &Decoder{multipart.NewReader(r, b)}
}

// MimeType returns Decoder's mime type.
func (d *Decoder) MimeType() string { return DataTypeMultipart }

// Decode data from multipart form to provided receiver.
func (d *Decoder) Decode(v interface{}) error {

	// TODO: data decoding

	return nil
}
