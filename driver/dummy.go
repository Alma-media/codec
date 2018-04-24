package driver

import "github.com/tiny-go/codec"

// DummyRegistry represents simple slice-based codec registry.
type DummyRegistry []codec.Codec

// Lookup returns an appropriate codec from a slice by MimeType (if available) or nil.
func (sr DummyRegistry) Lookup(mime string) codec.Codec {
	for _, codec := range sr {
		if mime == codec.MimeType() {
			return codec
		}
	}
	return nil
}
