package text

import (
	"bytes"
	"encoding"
	"errors"
	"strings"
	"testing"
)

const testString = "this is a test"

var (
	_ encoding.TextMarshaler   = mockMarshaler{}
	_ encoding.TextUnmarshaler = (*mockUnmarshaler)(nil)
)

type mockMarshaler struct {
	bytes []byte
	err   error
}

func (m mockMarshaler) MarshalText() ([]byte, error) { return m.bytes, m.err }

type mockUnmarshaler []byte

func (m *mockUnmarshaler) UnmarshalText(data []byte) error { *m = data; return nil }

func Test_Encoder(t *testing.T) {
	t.Run("Given text/plain codec", func(t *testing.T) {
		codec := new(Text)

		t.Run("test if variable is encoded to text", func(t *testing.T) {
			buf := bytes.NewBuffer([]byte{})

			if err := codec.Encoder(buf).Encode([]byte(testString)); err != nil {
				t.Errorf("unexpected error: %s", err)
			}
			if buf.String() != testString {
				t.Errorf("unexpected output: %s", buf)
			}
		})

		t.Run("test if TextMarshaler is encoded to text", func(t *testing.T) {
			buf := bytes.NewBuffer([]byte{})

			marshaler := mockMarshaler{bytes: []byte(testString)}

			if err := codec.Encoder(buf).Encode(marshaler); err != nil {
				t.Errorf("unexpected error: %s", err)
			}
			if buf.String() != testString {
				t.Errorf("unexpected output: %s", buf)
			}
		})

		t.Run("test if TextMarshaler error is propagated", func(t *testing.T) {
			buf := bytes.NewBuffer([]byte{})

			expectedErr := errors.New("encoding failed")
			marshaler := mockMarshaler{err: expectedErr}

			err := codec.Encoder(buf).Encode(marshaler)
			if err == nil {
				t.Error("error was expected")
			}
			if err != expectedErr {
				t.Errorf("unexpected error: %s", err)
			}
			if buf.String() != "" {
				t.Error("output was expected to be empty")
			}
		})
	})
}

func Test_Decoder(t *testing.T) {
	t.Run("Given text/plain codec", func(t *testing.T) {
		codec := new(Text)

		t.Run("test if text decoded to *string", func(t *testing.T) {
			var recv string
			if err := codec.Decoder(strings.NewReader(testString)).Decode(&recv); err != nil {
				t.Errorf("unexpected error: %s", err)
			}
			if recv != testString {
				t.Errorf("receiver has unexpected value: %s", recv)
			}
		})

		t.Run("test if text decoded to *[]byte", func(t *testing.T) {
			var recv []byte
			if err := codec.Decoder(strings.NewReader(testString)).Decode(&recv); err != nil {
				t.Errorf("unexpected error: %s", err)
			}
			if string(recv) != testString {
				t.Errorf("variable has unexpected value: %s", recv)
			}
		})

		t.Run("test if text decoded to TextUnmarshaler", func(t *testing.T) {
			var recv mockUnmarshaler
			if err := codec.Decoder(strings.NewReader(testString)).Decode(&recv); err != nil {
				t.Errorf("unexpected error: %s", err)
			}
			if string(recv) != testString {
				t.Errorf("variable has unexpected value: %s", recv)
			}
		})

		t.Run("test if error is returned providing invalid receiver", func(t *testing.T) {
			var recv int
			if err := codec.Decoder(strings.NewReader(testString)).Decode(&recv); err == nil {
				t.Errorf("error was expected")
			}
		})
	})
}

func Test_MimeType(t *testing.T) {
	t.Run("Given text/plain codec", func(t *testing.T) {
		codec := new(Text)

		t.Run("test if codec returns expected MimeType", func(t *testing.T) {
			if contentType := codec.MimeType(); contentType != DataTypeText {
				t.Errorf("codec returns unexpected content type %q", contentType)
			}
		})
	})
}
