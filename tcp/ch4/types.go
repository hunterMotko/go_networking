package ch4

import (
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
	"io"
)

// type-length-value (TLV) encoding
const (
	BinaryType uint8 = iota + 1
	StringType

	MaxPayloadSize uint32 = 10 << 20 // 10 MB
)

var ErrMaxPayloadSize = errors.New("maximum payload size exceeded")

type Payload interface {
	fmt.Stringer
	io.ReaderFrom
	io.WriterTo
	Bytes() []byte
}

type Binary []byte

func (m Binary) Bytes() []byte  { return m }
func (m Binary) String() string { return string(m) }
func (m Binary) WriteTo(w io.Writer) (int64, error) {
	err := binary.Write(w, binary.BigEndian, BinaryType)
	if err != nil {
		return 0, err
	}

	var n int64 = 1
	err = binary.Write(w, binary.BigEndian, uint32(len(m)))
	if err != nil {
		return n, err
	}

	n += 4
	o, err := w.Write(m)
	return n + int64(o), err
}
func (m *Binary) ReadFrom(r io.Reader) (int64, error) {
	var typ uint8
	err := binary.Read(r, binary.BigEndian, &typ)
	if err != nil {
		return 0, err
	}

	var n int64 = 1
	if typ != BinaryType {
		return n, errors.New("invaild binary")
	}

	var size uint32
	err = binary.Read(r, binary.BigEndian, &size)
	if err != nil {
		return n, err
	}

	n += 4
	if size > MaxPayloadSize {
		return n, ErrMaxPayloadSize
	}

	*m = make([]byte, size)
	o, err := r.Read(*m)
	return n + int64(o), err
}

type String string

func (m String) Bytes() []byte  { return []byte(m) }
func (m String) String() string { return string(m) }
func (m String) WriteTo(w io.Writer) (int64, error) {
	err := binary.Write(w, binary.BigEndian, StringType) // 1-byte type
	if err != nil {
		return 0, err
	}

	var n int64 = 1
	err = binary.Write(w, binary.BigEndian, uint32(len(m))) // 4-byte type
	if err != nil {
		return n, err
	}

	n += 4
	o, err := w.Write([]byte(m))
	return n + int64(o), err
}
func (m *String) ReadFrom(r io.Reader) (int64, error) {
	var typ uint8
	err := binary.Read(r, binary.BigEndian, &typ) // 1-byte type
	if err != nil {
		return 0, err
	}

	var n int64 = 1
	if typ != StringType {
		return n, errors.New("invaild String")
	}

	var size uint32
	err = binary.Read(r, binary.BigEndian, &size) // 4-byte type
	if err != nil {
		return n, err
	}

	n += 4
	buf := make([]byte, size)
	o, err := r.Read(buf) // payload
	if err != nil {
		return n, err
	}

	*m = String(buf)
	return n + int64(o), err
}

func Decode(r io.Reader) (Payload, error) {
	var typ uint8
	err := binary.Read(r, binary.BigEndian, &typ)
	if err != nil {
		return nil, err
	}

	var payload Payload
	switch typ {
	case BinaryType:
		payload = new(Binary)
	case StringType:
		payload = new(String)
	default:
		return nil, errors.New("unknown type")
	}

  // aprarently this is not optimal anyway and that there needs to be a refactor 
  // where you do not have to read the first byte to check what type it is 
  // and only be reading the 4 bytes
  // this code needs to be refactored to remove the multireader
	_, err = payload.ReadFrom(io.MultiReader(bytes.NewReader([]byte{typ}), r))
  if err != nil {
    return nil, err
  }

  return payload, err
}
