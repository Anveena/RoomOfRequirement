package ezNetworking

import (
	"errors"
	"io"
	"net"
	"reflect"
)

//todo 重写
type TCPListener struct {
	shouldStop bool
	pCallBack  *NetworkingCallBack
	pListener  *net.TCPListener
	chanSize   int
}

func NewTCPListen(chanSize int, port string) (*TCPListener, error) {
	rs := TCPListener{}
	tcpAddr, err := net.ResolveTCPAddr("tcp", ":"+port)
	if err != nil {
		return nil, err
	}
	rs.pListener, err = net.ListenTCP("tcp", tcpAddr)
	if err != nil {
		return nil, err
	}
	rs.chanSize = chanSize
	return &rs, nil
}
func (l *TCPListener) StartListen(callBack NetworkingCallBack) {
	l.pCallBack = &callBack
	l.shouldStop = false
	for !l.shouldStop {
		tcpConn, err := l.pListener.AcceptTCP()
		if err != nil {
			callBack(nil, nil, nil, err)
			continue
		}
		go l.pTCPConnHandler(tcpConn)
	}
}
func (l *TCPListener) StopListen() {
	err := l.pListener.Close()
	if err != nil {
		(*l.pCallBack)(nil, nil, nil, err)
	}
	l.shouldStop = true
}
func (l *TCPListener) WriteTo(data *[]byte, addr net.Addr, conn net.Conn) error {
	if data == nil {
		return errors.New("empty data")
	}
	if addr == nil || conn == nil {
		return errors.New("what the fuck, conn or addr is nil")
	}
	toWriteLen := len(*data)
	var tcpConn *net.TCPConn
	tcpConn, suc := conn.(*net.TCPConn)
	if !suc {
		return errors.New("what the fucking conn type: " + reflect.TypeOf(conn).String())
	}
	writtenLen, err := tcpConn.Write(*data)
	if err != nil {
		return err
	} else if writtenLen < toWriteLen {
		return errors.New("written len less than data len")
	}
	return nil
}
func (l *TCPListener) pTCPConnHandler(tcpConn *net.TCPConn) {
	receiveChan := make(chan pDataPackage, l.chanSize)
	connClose := false
	go l.receiveHandler(&receiveChan, tcpConn, &connClose)
	mAddr := tcpConn.RemoteAddr()
	for !l.shouldStop {
		var buffer [1024 * 32]byte
		readLen, err := tcpConn.Read(buffer[:])
		if err != nil {
			if err == io.EOF {
				connClose = true
				break
			}
			(*l.pCallBack)(nil, nil, nil, err)
			continue
		}
		data := buffer[:readLen]
		receiveChan <- pDataPackage{mAddr, &data}
	}
	close(receiveChan)
	err := tcpConn.Close()
	if err != nil {
		(*l.pCallBack)(nil, nil, nil, err)
	}
}
func (l *TCPListener) receiveHandler(receiveChan *chan pDataPackage, conn net.Conn, connClose *bool) {
	for !l.shouldStop && !*connClose {
		recvPkg := <-*receiveChan
		(*l.pCallBack)(recvPkg.data, recvPkg.addr, conn, nil)
	}
}
