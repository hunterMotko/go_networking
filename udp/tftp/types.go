package tftp

import (
	"bytes"
	"encoding/binary"
	"errors"
	"io"
	"strings"
)

const (
	DatagramSize = 516              // the max supported datagram size
	BlockSize    = DatagramSize - 4 // minus the 4 byte header
)

// The first 2 bytes of a tftp packets header is the operation code

type OpCode uint16

const (
	OpRRQ OpCode = iota
	_            // no WRQ support - no write request support
	OpData
	OpAck
	OpErr
)

type ErrCode uint16

const (
	ErrUnkown = iota
	ErrNotFound
	ErrAccessViolation
	ErrDiskFull
	ErrIllegalOp
	ErrUnknownID
	ErrFileExists
	ErrNoUser
)

// Read Request Packet Structure
// --------------------------------------------------------
// | 2-Bytes  |  n bytes   |  1 byte  |  n bytes  | 1 byte |
// | Op-Code  |  Filename  |    0     |   Mode    |   0    |
// --------------------------------------------------------

type ReadReq struct {
	Filename string
	Mode     string
}

// Although not used by our server, a client would make use of this method
func (q ReadReq) MarshalBinary() ([]byte, error) {
	mode := "octet"
	if q.Mode != "" {
		mode = q.Mode
	}
	// operation code + filename + 0 byte + mode + 0 byte
	cap := 2 + 2 + len(q.Filename) + 1 + len(q.Mode) + 1
	b := new(bytes.Buffer)
	b.Grow(cap)
	// write operation code
	err := binary.Write(b, binary.BigEndian, OpRRQ)
	if err != nil {
		return nil, err
	}
	// write filename
	_, err = b.WriteString(q.Filename)
	if err != nil {
		return nil, err
	}
	// write 0 byte
	err = b.WriteByte(0)
	if err != nil {
		return nil, err
	}
	// write mode
	_, err = b.WriteString(mode)
	if err != nil {
		return nil, err
	}
	// write 0 byte
	err = b.WriteByte(0)
	if err != nil {
		return nil, err
	}
	return b.Bytes(), nil
}

var (
  InvaildRRQ error = errors.New("invaild RRQ")
  InvaildData error = errors.New("invaild DATA")
  InvaildAck error = errors.New("invaild ACK")
  InvaildErr error = errors.New("invaild ERROR")
)

func (q *ReadReq) UnmarshallBinary(p []byte) error {
	r := bytes.NewBuffer(p)
	var code OpCode
	err := binary.Read(r, binary.BigEndian, &code)
	if err != nil {
		return err
	}
	if code != OpRRQ {
		return InvaildRRQ
	}
	// Read filename
	q.Filename, err = r.ReadString(0)
	if err != nil {
		return err
	}
	// remove the 0-byte
	q.Filename = strings.TrimRight(q.Filename, "\x00")
	if len(q.Filename) == 0 {
		return InvaildRRQ
	}
	// read mode
	q.Mode, err = r.ReadString(0)
	if err != nil {
		return InvaildRRQ
	}
	// remove the 0-byte
	q.Mode = strings.TrimRight(q.Mode, "\x00")
	if len(q.Mode) == 0 {
		return InvaildRRQ
	}
	// enforce octect mode
	actual := strings.ToLower(q.Mode)
	if actual != "octet" {
		return errors.New("only binary transfers supported")
	}
	return nil
}

// Data Packet Structure
// -----------------------------------
// | 2-Bytes  |   2 bytes  | n bytes |
// | Op-Code  |   Block #  | Payload |
// -----------------------------------

type Data struct {
  Block uint16
  Payload io.Reader
}

func (d *Data) MarshallBinary() ([]byte, error) {
  b := new(bytes.Buffer)
  b.Grow(DatagramSize)
  d.Block++
  err := binary.Write(b, binary.BigEndian, OpData)
  if err != nil {
    return nil, err
  }
  err = binary.Write(b, binary.BigEndian, d.Block)
  if err != nil {
    return nil, err
  }
  // write up to blocksize worth of bytes
  _, err = io.CopyN(b, d.Payload, BlockSize)
  if err != nil {
    return nil, err
  }
  return b.Bytes(), nil
}

func (d *Data) UnmarshallBinary(p []byte) error {
  if l := len(p); l < 4 || l > DatagramSize {
    return InvaildData
  }
  var opcode OpCode
  err := binary.Read(bytes.NewReader(p[:2]), binary.BigEndian, &opcode)
  if err != nil {
   return InvaildData
  }
  err = binary.Read(bytes.NewReader(p[2:4]), binary.BigEndian, &d.Block)
  if err != nil {
   return InvaildData
  }
  d.Payload = bytes.NewBuffer(p[4:])
  return nil
}

// Ack Packet Structure
// -------------------------
// | 2-Bytes  |   2 bytes  |
// | Op-Code  |   Block #  |
// -------------------------

type Ack uint16

func (a Ack) MarshalBinary() ([]byte, error) {
  cap := 2 + 2
  b := new(bytes.Buffer)
  b.Grow(cap)
  err := binary.Write(b, binary.BigEndian, OpAck) //write operation code
  if err != nil {
    return nil, err
  }
  err = binary.Write(b, binary.BigEndian, &a) // write block number
  if err != nil {
    return nil, err
  }
  return b.Bytes(), nil
}

func (a Ack) UnmarshallBinary(p []byte) error {
  var code OpCode
  r := bytes.NewReader(p)
  err := binary.Read(r, binary.BigEndian, &code) // readc operation code
  if err != nil {
    return err
  }
  if code != OpAck {
    return InvaildAck
  }
  return binary.Read(r, binary.BigEndian, a) // read block number
}

// Error Packet Structure
// --------------------------------------------
// | 2-Bytes  |   2 bytes  | n bytes | 1 byte |
// | Op-Code  |   ErrCode  | Message |    0   |
// --------------------------------------------

type Err struct {
  Error ErrCode
  Message string
}

func (e Err) MarshallBinary() ([]byte, error) {
  // operation code + error code + message + 0 byte
  cap := 2 + 2 + len(e.Message) + 1
  b := new(bytes.Buffer)
  b.Grow(cap)
  err := binary.Write(b, binary.BigEndian, OpErr) // write operation code
  if err != nil {
    return nil, err
  }
  err = binary.Write(b, binary.BigEndian, e.Error) // write error code
  if err != nil {
    return nil, err
  }
  _, err = b.WriteString(e.Message) // write message
  if err != nil {
    return nil, err
  }
  err = b.WriteByte(0) // write 0 byte
  if err != nil {
    return nil, err
  }
  return b.Bytes(), nil
}

func (e Err) UnmarshallBinary(p []byte) error {
  r := bytes.NewBuffer(p)
  var code OpCode
  err := binary.Read(r, binary.BigEndian, &code)
  if err != nil {
    return err
  }
  if code != OpErr {
    return InvaildErr
  }
  err = binary.Read(r, binary.BigEndian, &e.Error)
  if err != nil {
    return err
  }
  e.Message, err = r.ReadString(0)
  e.Message = strings.TrimRight(e.Message, "\x00")
  return err
}
