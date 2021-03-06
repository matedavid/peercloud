package network

import (
	"errors"
	"net"
)

func SendPayload(conn net.Conn, m Model) error {
	mh := MessageHeader{}

	// Receive acknowledge message header
	err := mh.Recv(conn)
	if err != nil {
		return err
	} else if mh.Command != Acknowledge {
		return errors.New("did not receive acknowledge message header")
	}

	// Send content
	content := m.Write()
	_, err = conn.Write(content)
	if err != nil {
		return err
	}

	return nil

}

func ReceivePayload(conn net.Conn, payload uint32, m Model) error {
	// Send acknowledge message header
	ack := MessageHeader{
		NetworkCode: MAIN_NETWORK_CODE,
		Command:     Acknowledge,
		Payload:     0,
	}
	err := ack.Send(conn)
	if err != nil {
		return err
	}

	// Receive data
	buff := make([]byte, payload)
	n, err := conn.Read(buff)
	if err != nil {
		return err
	} else if n != int(payload) {
		return errors.New("length of data received does not match payload")
	}

	m.Read(buff)
	return nil
}
