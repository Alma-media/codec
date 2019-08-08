package codec

import "io"

// Encoder is responsible for data encoding.
type Encoder interface{ Encode(v interface{}) error }

// EncoderFunc is a wrapper that allows using regular function as Encoder interface.
type EncoderFunc func(interface{}) error

// Encode encodes provided interface to the writer.
func (fn EncoderFunc) Encode(src interface{}) error { return fn(src) }

// Decoder is responsible for data decoding.
type Decoder interface{ Decode(v interface{}) error }

// DecoderFunc is a wrapper that allows using regular function as Decoder interface.
type DecoderFunc func(interface{}) error

// Decode decodes data from the reader to provided receiver.
func (fn DecoderFunc) Decode(recv interface{}) error { return fn(recv) }

// Codec can create Encoder(s) and Decoder(s) with provided io.Reader/io.Writer.
type Codec interface {
	// Ecoder instantiates the ecnoder part of the codec with provided writer.
	Encoder(w io.Writer) Encoder
	// Decoder instantiates the decoder part of the codec with provided reader.
	Decoder(r io.Reader) Decoder
	// MimeType returns the (main) mime type of the codec.
	MimeType() string
}

// Registry represents any kind of codec registry, thus it can be global registry
// or a custom list of codecs that is supposed to be used for a particular case.
type Registry interface {
	// Lookup should find appropriate Codec by MimeType or return nil if not found.
	Lookup(mimeType string) Codec
}
