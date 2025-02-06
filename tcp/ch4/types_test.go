package ch4_test

import (
	"bytes"
	"encoding/binary"
	"net"
	"network/ch4"
	"reflect"
	"testing"
)

func TestPayload(t *testing.T) {
	b1 := ch4.Binary("Clear is better than clever.")
	b2 := ch4.Binary("Don't panic.")
	s1 := ch4.String("Errors are values.")
	payloads := []ch4.Payload{&b1, &s1, &b2}

	listener, err := net.Listen("tcp", "127.0.0.1:")
	if err != nil {
		t.Fatal(err)
	}

	go func() {
		conn, err := listener.Accept()
		if err != nil {
			t.Error(err)
			return
		}
		defer conn.Close()

		for _, p := range payloads {
			_, err = p.WriteTo(conn)
			if err != nil {
				t.Error(err)
				break
			}
		}
	}()

	conn, err := net.Dial("tcp", listener.Addr().String())
	if err != nil {
		t.Fatal(err)
	}
	defer conn.Close()

	for i := 0; i < len(payloads); i++ {
		actual, err := ch4.Decode(conn)
		if err != nil {
			t.Fatal(err)
		}
		if expected := payloads[i]; !reflect.DeepEqual(expected, actual) {
			t.Errorf("value mismatch: %v != %v\n", expected, actual)
			continue
		}
		t.Logf("[%T] %[1]q\n", actual)
	}
}

// the previous test should cover the use case of a payload that is less than 
// the maximum size, but need to modify this test to make sure that is the case
func TestMaxPayloadSize(t *testing.T) {
	buf := new(bytes.Buffer)
	err := buf.WriteByte(ch4.BinaryType)
	if err != nil {
		t.Fatal(err)
	}
	err = binary.Write(buf, binary.BigEndian, uint32(1<<30))
	if err != nil {
		t.Fatal(err)
	}
	var b ch4.Binary
	_, err = b.ReadFrom(buf)
	if err != ch4.ErrMaxPayloadSize {
		t.Fatalf("expected ErrMaxPayloadSize; actual: %v\n", err)
	}
}
