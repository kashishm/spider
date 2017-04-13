package debug

import (
	"fmt"
	"net"

	"bytes"
	"time"

	m "github.com/getgauge/spider/gauge_messages"
	"github.com/golang/protobuf/proto"
)

type api struct {
	connection net.Conn
}

func newAPI(host string, port string) (*api, error) {
	c, err := net.Dial("tcp", fmt.Sprintf("%s:%s", host, port))
	return &api{connection: c}, err
}

func (a *api) getResponse(message *m.APIMessage) (*m.APIMessage, error) {
	messageId := time.Now().UnixNano()
	message.MessageId = messageId

	data, err := proto.Marshal(message)
	if err != nil {
		return nil, err
	}
	responseBytes, err := a.writeDataAndGetResponse(data)
	if err != nil {
		return nil, err
	}
	responseMessage := &m.APIMessage{}
	if err := proto.Unmarshal(responseBytes, responseMessage); err != nil {
		return nil, err
	}
	return responseMessage, err
}

func (a *api) writeDataAndGetResponse(messageBytes []byte) ([]byte, error) {
	if err := a.write(messageBytes); err != nil {
		return nil, err
	}

	return a.readResponse()
}

func (a *api) readResponse() ([]byte, error) {
	buffer := new(bytes.Buffer)
	data := make([]byte, 8192)
	for {
		n, err := a.connection.Read(data)
		if err != nil {
			a.connection.Close()
			return nil, fmt.Errorf("Connection closed [%s] cause: %s", a.connection.RemoteAddr(), err.Error())
		}

		buffer.Write(data[0:n])

		messageLength, bytesRead := proto.DecodeVarint(buffer.Bytes())
		if messageLength > 0 && messageLength < uint64(buffer.Len()) {
			return buffer.Bytes()[bytesRead : messageLength+uint64(bytesRead)], nil
		}
	}
}

func (a *api) write(messageBytes []byte) error {
	messageLen := proto.EncodeVarint(uint64(len(messageBytes)))
	data := append(messageLen, messageBytes...)
	_, err := a.connection.Write(data)
	return err
}

func (a *api) close() {
	a.connection.Close()
}
