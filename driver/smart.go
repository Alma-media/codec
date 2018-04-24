package driver

import (
	"fmt"
	"strings"
	"sync"

	"github.com/tiny-go/codec"
)

// codecInitFunc is a codec constructor which can either return a singleton instance
// or initialize a new codec for every request (for example multipart, because each
// request has its own unique boundary).
type codecInitFunc func(string) codec.Codec

// SmartRegistry is a smart codec registry that is used to store available codecs
// and retrieve them by mime type.
type SmartRegistry struct {
	sync.Mutex
	codecs map[string]codecInitFunc
}

// NewSmartRegistry is a SmartRegistry constructor func.
func NewSmartRegistry() *SmartRegistry {
	return &SmartRegistry{codecs: make(map[string]codecInitFunc)}
}

// Lookup searches for appropriate codec by MimeType, if there is no exact match
// (that is possible using codecs such as Multipart), it tries to detect a codec
// in an opposite way, comparing default codec MimeType to the provided argument.
func (r *SmartRegistry) Lookup(mime string) codec.Codec {
	r.Lock()
	defer r.Unlock()
	// find strict match
	if codec, ok := r.codecs[mime]; ok {
		return codec(mime)
	}
	// find submatch (for example "multipart/form-data; boundary=-----...")
	for mediatype, codec := range r.codecs {
		// ignore default codec
		if mediatype != "" && mediatype != "*/*" {
			if strings.Contains(mime, mediatype+";") {
				return codec(mime)
			}
		}
	}
	return nil
}

// Register makes codec available for provided mime type.
func (r *SmartRegistry) Register(mime string, f codecInitFunc) error {
	r.Lock()
	defer r.Unlock()

	if _, ok := r.codecs[mime]; ok {
		return fmt.Errorf("codec with mime type %q already registered", mime)
	}
	r.codecs[mime] = f
	return nil
}

// Default allows setting Codec as a default by mime type. If codec was not found
// function returns an error. This function can be called multiple times overriding
// the previous values.
func (r *SmartRegistry) Default(mime string) error {
	r.Lock()
	defer r.Unlock()

	f, ok := r.codecs[mime]
	if !ok {
		return fmt.Errorf("codec %q is not registered", mime)
	}
	r.codecs[""] = f
	r.codecs["*/*"] = f
	return nil
}
