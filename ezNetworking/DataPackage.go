package ezNetworking

import "net"

//todo 重写
type pDataPackage struct {
	addr net.Addr
	data *[]byte
}
