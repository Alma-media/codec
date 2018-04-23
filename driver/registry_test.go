package driver

import (
	"errors"
	"io"
	"testing"

	"github.com/tiny-go/codec"
)

// compile time type check
var _ codec.Registry = globalRegistry{}

var (
	codecJSON      = &mockJSON{}
	codecXML       = &mockXML{}
	codecMultipart = &mockMultipart{}
)

type mockJSON struct{}

func (j *mockJSON) Encoder(w io.Writer) codec.Encoder { return nil }

func (j *mockJSON) Decoder(r io.Reader) codec.Decoder { return nil }

func (j *mockJSON) MimeType() string { return "application/json" }

type mockXML struct{}

func (x *mockXML) Encoder(w io.Writer) codec.Encoder { return nil }

func (x *mockXML) Decoder(r io.Reader) codec.Decoder { return nil }

func (x *mockXML) MimeType() string { return "application/xml" }

type mockMultipart struct{}

func (m *mockMultipart) Encoder(w io.Writer) codec.Encoder { return nil }

func (m *mockMultipart) Decoder(r io.Reader) codec.Decoder { return nil }

func (m *mockMultipart) MimeType() string { return "multipart/form-data" }

func TestDefault(t *testing.T) {
	t.Run("Given global codec registry", func(t *testing.T) {
		// clean up codec registry
		defer func() { codecs = globalRegistry{} }()

		t.Run("test if error is returned trying to set non-existing codec as a default", func(t *testing.T) {
			if err := Default("application/json"); err == errors.New("codec \"application/json\" is not registered") {
				t.Error("error was expected: codec is not registered")
			}
		})

		t.Run("tests if existing codec can be set as a default", func(t *testing.T) {
			Register("application/json", func(string) codec.Codec { return codecJSON })
			if Default("application/json") != nil {
				t.Error("unexpected error")
			}
		})
	})
}

func TestGlobal(t *testing.T) {
	t.Run("Given global codec registry", func(t *testing.T) {
		// clean up codec registry
		defer func() { codecs = globalRegistry{} }()

		Register("application/json", func(string) codec.Codec { return codecJSON })
		Register("multipart/form-data", func(string) codec.Codec { return codecMultipart })

		t.Run("test if appropriate codec is returned by its mime type if registered", func(t *testing.T) {
			if Global().Lookup("application/json") != codecJSON {
				t.Error("was expected to return JSON codec")
			}
		})

		t.Run("test if error is returned when failed to find a codec by mime type", func(t *testing.T) {
			if Global().Lookup("application/xml") != nil {
				t.Error("was expected to return nil")
			}
		})

		t.Run("test retrieving codec from composite headers", func(t *testing.T) {
			if Global().Lookup("multipart/form-data; boundary=-----...") != codecMultipart {
				t.Error("was expected to return multipart codec")
			}
		})

		t.Run("test returning default codec when mime type is not set or does not matter", func(t *testing.T) {
			if Global().Lookup("") != nil {
				t.Error("was expected to return nil: default codec is not regitered")
			}
			if Global().Lookup("*/*") != nil {
				t.Error("was expected to return nil: default codec is not regitered")
			}
			Default("application/json")
			if Global().Lookup("") != codecJSON {
				t.Error("was expected to return JSON codec")
			}
			if Global().Lookup("*/*") != codecJSON {
				t.Error("was expected to return JSON codec")
			}
		})
	})
}

func TestRegister(t *testing.T) {
	t.Run("Given global codec registry", func(t *testing.T) {
		t.Run("test if codecs with unique mime type are registered successfully", func(t *testing.T) {
			// clean up codec registry
			defer func() { codecs = globalRegistry{} }()

			Register("application/json", func(string) codec.Codec { return codecJSON })
			Register("application/xml", func(string) codec.Codec { return codecXML })
			if len(codecs) != 2 {
				t.Error("registry chould contain exactly 2 codecs")
			}
			jsonFunc, ok := codecs["application/json"]
			if !ok || jsonFunc("") != codecJSON {
				t.Error("codec with mime type \"application/json\" was expected to be available")
			}
			xmlFunc, ok := codecs["application/xml"]
			if !ok || xmlFunc("") != codecXML {
				t.Error("codec with mime type \"application/xml\" was expected to be available")
			}
		})

		t.Run("test if panics trying to register codec with non-unique mime type", func(t *testing.T) {
			// clean up codec registry
			defer func() { codecs = globalRegistry{} }()

			defer func() {
				if r := recover(); r == nil {
					t.Error("was expected to panic but did not")
				}
			}()
			Register("application/json", func(string) codec.Codec { return nil })
			Register("application/json", func(string) codec.Codec { return nil })
			t.Error("should never reach this code")
		})
	})
}
