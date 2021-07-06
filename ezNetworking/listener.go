package ezNetworking

import (
	"context"
	"fmt"
	"net"
	"syscall"
)

func ListenUDP4(port int, sendBufferSize int, recvBufferSize int) (*net.UDPConn, *net.OpError) {
	lc := net.ListenConfig{Control: func(network, address string, c syscall.RawConn) error {
		if sendBufferSize > 0 {
			err := setIOBufferSize(c, sendBufferSize, syscall.SO_SNDBUF)
			if err != nil {
				return err
			}
		}
		if recvBufferSize > 0 {
			err := setIOBufferSize(c, recvBufferSize, syscall.SO_RCVBUF)
			if err != nil {
				return err
			}
		}
		return nil
	}}
	l, err := lc.ListenPacket(context.Background(), "udp4", fmt.Sprintf(":%v", port))
	if err != nil {
		return nil, err.(*net.OpError)
	}
	return l.(*net.UDPConn), nil
}
func ListenTCP4(port int, sendBufferSize int, recvBufferSize int) (*net.TCPListener, *net.OpError) {
	lc := net.ListenConfig{Control: func(network, address string, c syscall.RawConn) error {
		if sendBufferSize > 0 {
			err := setIOBufferSize(c, sendBufferSize, syscall.SO_SNDBUF)
			if err != nil {
				return err
			}
		}
		if recvBufferSize > 0 {
			err := setIOBufferSize(c, recvBufferSize, syscall.SO_RCVBUF)
			if err != nil {
				return err
			}
		}
		return nil
	}}
	l, err := lc.Listen(context.Background(), "tcp4", fmt.Sprintf(":%v", port))
	if err != nil {
		return nil, err.(*net.OpError)
	}
	return l.(*net.TCPListener), nil
}
