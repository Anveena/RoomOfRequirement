package ezNetworking

import (
	"fmt"
	"net"
	"strings"
)

func InetNtoA(ip []byte) string {
	return fmt.Sprintf("%d.%d.%d.%d", ip[0], ip[1], ip[2], ip[3])
}

func InetAtoN(ip string) []byte {
	ip_ := strings.Split(ip, ":")
	return net.ParseIP(ip_[0]).To4()
}
