package driver

import (
	"github.com/tiny-go/codec"
)

// global codec registry.
var codecs = NewSmartRegistry()

// Global returns a global instance of SmartRegistry.
func Global() codec.Registry {
	// return an interface (it is safe / read only)
	return codecs
}

// Register makes codec available for provided mime type (using global instance
// of SmartRegistry).
func Register(mime string, f codecInitFunc) {
	if err := codecs.Register(mime, f); err != nil {
		panic(err.Error())
	}
}

// Default allows setting Codec as a default by mime type (using global instance
// of SmartRegistry).
func Default(mime string) error {
	return codecs.Default(mime)
}
