package ezNetworking

import "net"

//todo 重写
type NetworkingCallBack func(dataReceived *[]byte, clientHost net.Addr, conn net.Conn, err error)

type NetworkingListener interface {
	StartListen(callBack NetworkingCallBack)
	StopListen()
	WriteTo(data *[]byte, addr net.Addr, conn net.Conn) error
}
