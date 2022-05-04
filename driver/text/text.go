package text

import (
	"bytes"
	"encoding"
	"fmt"
	"io"

	"github.com/tiny-go/codec"
	"github.com/tiny-go/codec/driver"
)

// DataTypeText contains plain text codec data type.
const DataTypeText = "text/plain"

var global = &Text{}

// Text plain codec.
type Text struct{}

// Encoder creates text encoder.
func (Text) Encoder(w io.Writer) codec.Encoder {
	return NewEncoder(w)
}

func NewEncoder(w io.Writer) codec.Encoder {
	return codec.EncoderFunc(func(src interface{}) error {
		if marshaler, ok := src.(encoding.TextMarshaler); ok {
			data, err := marshaler.MarshalText()
			if err != nil {
				return err
			}
			_, err = w.Write(data)
			return err
		}
		fmt.Fprintf(w, "%s", src)
		return nil
	})
}

// Decoder creates text decoder.
func (Text) Decoder(r io.Reader) codec.Decoder {
	return NewDecoder(r)
}

func NewDecoder(r io.Reader) codec.Decoder {
	return codec.DecoderFunc(func(recv interface{}) error {
		var buf = new(bytes.Buffer)
		if _, err := buf.ReadFrom(r); err != nil {
			return err
		}
		if unmarshaler, ok := recv.(encoding.TextUnmarshaler); ok {
			return unmarshaler.UnmarshalText(buf.Bytes())
		}
		switch t := recv.(type) {
		case *string:
			*t = buf.String()
			return nil
		case *[]byte:
			*t = buf.Bytes()
			return nil
		default:
			return fmt.Errorf("cannot unmarshal into variable of type %T", t)
		}
	})
}

// MimeType returns the mime type of the codec.
func (Text) MimeType() string { return DataTypeText }

func init() {
	// since text is a stateless codec we can always return the same instance in order
	// to save some memory
	driver.Register(DataTypeText, func(string) codec.Codec { return global })
}
