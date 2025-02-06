package tftp

import (
	"bytes"
	"errors"
	"log"
	"net"
	"time"
)

type Server struct {
	Payload []byte
	Retires uint8
	Timeout time.Duration
}

func (s Server) ListenAndServe(addr string) error {
	conn, err := net.ListenPacket("udp", addr)
	if err != nil {
		return err
	}
	defer func() { _ = conn.Close() }()
	log.Printf("Listening on %s ...\n", conn.LocalAddr())
	return s.Serve(conn)
}
func (s *Server) Serve(conn net.PacketConn) error {
	if conn == nil {
		return errors.New("nil connection")
	}
	if s.Payload == nil {
		return errors.New("payload is required")
	}
	if s.Retires == 0 {
		s.Retires = 10
	}
	if s.Timeout == 0 {
		s.Timeout = 6 * time.Second
	}
	var rrq ReadReq
	for {
		buf := make([]byte, DatagramSize)
		_, addr, err := conn.ReadFrom(buf)
		if err != nil {
			return err
		}
		err = rrq.UnmarshallBinary(buf)
		if err != nil {
			log.Printf("[%s] bad request: %v\n", addr, err)
			continue
		}
		go s.handle(addr.String(), rrq)
	}
	return nil
}

func (s Server) handle(clientAddr string, rrq ReadReq) {
	log.Printf("[%s] requested file: %s\n", clientAddr, rrq.Filename)
	conn, err := net.Dial("udp", clientAddr)
	if err != nil {
		log.Printf("[%s] dial: %v\n", clientAddr, err)
		return
	}
	defer func() { _ = conn.Close() }()
	var (
		ackPkt  Ack
		errPkt  Err
		dataPkt = Data{Payload: bytes.NewReader(s.Payload)}
		buf     = make([]byte, DatagramSize)
	)

  NEXTPACKET:
	for n := DatagramSize; n == DatagramSize; {
		data, err := dataPkt.MarshallBinary()
		if err != nil {
			log.Printf("[%s] preparing data packet: %v\n", clientAddr, err)
			return
		}
	  RETRY:
		for i := s.Retires; i > 0; i-- {
			n, err = conn.Write(data)
			if err != nil {
				log.Printf("[%s] write: %v\n", clientAddr, err)
				return
			}

			// wait for the clients ACK packet
			_ = conn.SetReadDeadline(time.Now().Add(s.Timeout))
			_, err = conn.Read(buf)
			if err != nil {
				if nErr, ok := err.(net.Error); ok && nErr.Timeout() {
					continue RETRY
				}
				log.Printf("[%s] waiting for ACK: %v\n", clientAddr, err)
				return
			}

			switch {
			case ackPkt.UnmarshallBinary(buf) == nil:
				if uint16(ackPkt) == dataPkt.Block {
					continue NEXTPACKET
				}
			case errPkt.UnmarshallBinary(buf) == nil:
				log.Printf("[%s] received error: %v\n", clientAddr, errPkt.Message)
				return
			default:
				log.Printf("[%s] bad packet\n", clientAddr)
			}
		}

		log.Printf("[%s] exhausted retires\n", clientAddr)
		return
	}

	log.Printf("[%s] sent %d blocks\n", clientAddr, dataPkt.Block)
}
