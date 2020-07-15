package ezNetworking

import (
	"errors"
	"net"
	"reflect"
)

//todo 重写
type UDPListener struct {
	shouldStop             bool
	pCallBack              *NetworkingCallBack
	pUDPConn               *net.UDPConn
	chanSize               int
	receiveChan            chan pDataPackage
	receiveChanCloseSignal chan byte
}

func NewUDPListen(chanSize int, port string) (*UDPListener, error) {
	rs := UDPListener{}
	udpAddr, err := net.ResolveUDPAddr("udp", ":"+port)
	if err != nil {
		return nil, err
	}
	rs.pUDPConn, err = net.ListenUDP("udp", udpAddr)
	if err != nil {
		return nil, err
	}
	rs.receiveChan = make(chan pDataPackage, chanSize)
	rs.receiveChanCloseSignal = make(chan byte, 16)
	rs.chanSize = chanSize
	return &rs, nil
}
func (l *UDPListener) StopListen() {
	l.receiveChanCloseSignal <- 0
	l.shouldStop = true
	_ = l.pUDPConn.Close()
	if len(l.receiveChan) > 1 {
		<-l.receiveChan
	}
}
func (l *UDPListener) StartListen(callBack NetworkingCallBack) {
	l.pCallBack = &callBack
	l.shouldStop = false
	go l.receiveHandler(l.pUDPConn)
	for !l.shouldStop {
		var buffer [1024 * 32]byte
		readLen, udpAddr, err := l.pUDPConn.ReadFromUDP(buffer[0:])
		if err != nil {
			(*l.pCallBack)(nil, nil, nil, err)
			break
		}
		data := buffer[:readLen]
		l.receiveChan <- pDataPackage{udpAddr, &data}
	}
	_ = l.pUDPConn.Close()
}
func (l *UDPListener) receiveHandler(conn net.Conn) {
	for !l.shouldStop {
		select {
		case <-l.receiveChanCloseSignal:
			break
		case recvPkg := <-l.receiveChan:
			(*l.pCallBack)(recvPkg.data, recvPkg.addr, conn, nil)
		}
	}
}
func (l *UDPListener) WriteTo(data *[]byte, addr net.Addr, conn net.Conn) error {
	if data == nil {
		return errors.New("empty data")
	}
	if addr == nil || conn == nil {
		return errors.New("what the fuck, conn or addr is nil")
	}
	toWriteLen := len(*data)
	var udpConn *net.UDPConn
	udpConn, suc := conn.(*net.UDPConn)
	if !suc {
		return errors.New("what the fucking conn type: " + reflect.TypeOf(conn).String())
	}
	writtenLen, err := udpConn.WriteTo(*data, addr)
	if err != nil {
		return err
	} else if writtenLen < toWriteLen {
		return errors.New("written len less than data len")
	}
	return nil
}
