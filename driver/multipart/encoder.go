package multipart

import (
	"fmt"
	"io"
	"mime/multipart"
	"reflect"
)

// Namer implies the presence of method Name().
type Namer interface {
	Name() string
}

// Typer implies the presense of ContentType() method.
type Typer interface {
	ContentType() string
}

// Encoder is responsible for encoding data to multipart form on the fly.
type Encoder struct {
	mw       *multipart.Writer
	boundary string
}

// NewEncoder creates a new instance of multipart.Encoder.
func NewEncoder(w io.Writer) *Encoder {
	return &Encoder{multipart.NewWriter(w), ""}
}

// MimeType returns Encoder's mime type with boundary.
func (mpe *Encoder) MimeType() string { return mpe.mw.FormDataContentType() }

// Encode data to multipart form.
func (mpe *Encoder) Encode(v interface{}) error {
	_, kind := ActualType(reflect.ValueOf(v))

	if kind == reflect.Ptr || kind == reflect.Interface || kind == reflect.Invalid {
		return fmt.Errorf("unsupported %T", v)
	}

	// TODO: data encoding

	return nil
}
