package transport

import (
	"encoding/binary"
	"io"
	"net"
)

type Transport struct {
	conn net.Conn
}

// 建立一个Transport
func NewTransport(conn net.Conn) *Transport {
	return &Transport{conn}
}

// 把数据转成TLV格式，然后发送数据
func (t *Transport) Send(data []byte) error {
	buf := make([]byte, len(data)+4)
	binary.BigEndian.PutUint32(buf[:4], uint32(len(data))) // 把数据长度大端存到字节流中
	copy(buf[4:], data)
	_, err := t.conn.Write(buf)
	if err != nil {
		return err
	}
	return nil
}

// 读取数据
func (t *Transport) Read() ([]byte, error) {
	header := make([]byte, 4)
	_, err := io.ReadFull(t.conn, header)
	if err != nil {
		return nil, err
	}
	dataLength := binary.BigEndian.Uint32(header)
	data := make([]byte, dataLength)
	_, err = io.ReadFull(t.conn, data)
	if err != nil {
		return nil, err
	}
	return data, nil
}
