package ezNetworking

import (
	"errors"
	"fmt"
	"net"
	"syscall"
)

func SetReadBufferSize(conn net.Conn, size int) error {
	if udp, suc := conn.(*net.UDPConn); suc && udp != nil {
		return setUDPBuffer(udp, size, true, false)
	}
	if tcp, suc := conn.(*net.TCPConn); suc && tcp != nil {
		return setTCPBuffer(tcp, size, true, false)
	}
	return errors.New("unknown conn type")
}
func SetWriteBufferSize(conn net.Conn, size int) error {
	if udp, suc := conn.(*net.UDPConn); suc && udp != nil {
		return setUDPBuffer(udp, size, false, true)
	}
	if tcp, suc := conn.(*net.TCPConn); suc && tcp != nil {
		return setTCPBuffer(tcp, size, false, true)
	}
	return errors.New("unknown conn type")
}
func setTCPBuffer(conn *net.TCPConn, size int, isRead bool, isWrite bool) error {
	if isRead {
		if err := conn.SetReadBuffer(size); err != nil {
			return err
		}
	}
	if isWrite {
		if err := conn.SetWriteBuffer(size); err != nil {
			return err
		}
	}
	f, err := conn.File()
	if err != nil {
		return err
	}
	actuallyValue, err := syscall.GetsockoptInt(int(f.Fd()), syscall.SOL_SOCKET, syscall.SO_RCVBUF)
	if err != nil {
		return err
	}
	if actuallyValue < size {
		return errors.New(fmt.Sprintf("set readbuffer failed,wanted result:%v,actually:%v", size, actuallyValue))
	}
	return nil
}
func setUDPBuffer(conn *net.UDPConn, size int, isRead bool, isWrite bool) error {
	if isRead {
		if err := conn.SetReadBuffer(size); err != nil {
			return err
		}
	}
	if isWrite {
		if err := conn.SetWriteBuffer(size); err != nil {
			return err
		}
	}
	f, err := conn.File()
	if err != nil {
		return err
	}
	actuallyValue, err := syscall.GetsockoptInt(int(f.Fd()), syscall.SOL_SOCKET, syscall.SO_RCVBUF)
	if err != nil {
		return err
	}
	if actuallyValue != size {
		return errors.New(fmt.Sprintf("set readbuffer failed,wanted result:%v,actually:%v", size, actuallyValue))
	}
	return nil
}
