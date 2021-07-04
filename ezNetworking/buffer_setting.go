package ezNetworking

import (
	"errors"
	"net"
	"os"
)

func SetBufferSize(conn net.Conn, size int, isWrite bool) error {
	if udp, suc := conn.(*net.UDPConn); suc && udp != nil {
		if isWrite {
			if err := udp.SetWriteBuffer(size); err != nil {
				return err
			}
		} else {
			if err := udp.SetReadBuffer(size); err != nil {
				return err
			}
		}
		f, err := udp.File()
		if err != nil {
			return err
		}
		err = checkBufferSize(f, size, isWrite)
		_ = f.Close()
		return err
	}
	if tcp, suc := conn.(*net.TCPConn); suc && tcp != nil {
		if isWrite {
			if err := tcp.SetWriteBuffer(size); err != nil {
				return err
			}
		} else {
			if err := tcp.SetReadBuffer(size); err != nil {
				return err
			}
		}
		f, err := tcp.File()
		if err != nil {
			return err
		}
		err = checkBufferSize(f, size, isWrite)
		_ = f.Close()
		return err
	}
	return errors.New("unknown conn type")
}
func checkBufferSize(f *os.File, size int, isWrite bool) error {
	//TODO 我并不能做这件事情.因为golang有个BUG,等这个BUG修复了,我就可以放开这里的注释
	//https://github.com/golang/go/issues/47050
	//opt := syscall.SO_RCVBUF
	//if isWrite{
	//	opt = syscall.SO_SNDBUF
	//}
	//actuallyValue, err := syscall.GetsockoptInt(int(f.Fd()), syscall.SOL_SOCKET, opt)
	//if err != nil {
	//	return err
	//}
	//if actuallyValue < size {
	//	return errors.New(fmt.Sprintf("set buffer failed,wanted result:%v,actually:%v", size, actuallyValue))
	//}
	return nil
}
