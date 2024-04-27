package ip

import (
	"net"
	"strings"
)

func GetLoacalIPAddresses() string {
	interfaces, err := net.Interfaces()
	if err != nil{
		return ""
	}

	for _, iface := range interfaces {
		addrs, err := iface.Addrs()
		if err != nil{
			continue
		}

		for _, addr := range addrs{
			ipNet, ok := addr.(*net.IPNet)
			if ok && !ipNet.IP.IsLoopback() {
				if ipNet.IP.To4() != nil {
					ip := ipNet.IP.String()
					if strings.Contains(ip, "."){
						return ip
					}					
				}
			}
		}
	}
	return ""
}