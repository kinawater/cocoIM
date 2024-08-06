package tcp

import (
	"fmt"
	"net"
)

func WriteData(data []byte, conn *net.TCPConn) error {
	dataLen := len(data)
	var pos int
	if dataLen <= 0 {
		return fmt.Errorf("no length data,please check")
	}
	for {
		writeLen, err := conn.Write(data[pos:])
		if err != nil {
			return err
		}
		pos = pos + writeLen
		if pos >= dataLen {
			break
		}
	}
	return nil
}
