package gateway

import "net"

type connection struct {
	fd   int
	conn *net.TCPConn
}

func (c *connection) close() error {
	err := c.conn.Close()
	if err != nil {
		return err
	}
	return nil
}
