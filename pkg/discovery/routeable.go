package discovery

import (
	"net"
	"strings"

	externalip "github.com/glendc/go-external-ip"
	"github.com/jackpal/gateway"
)

// GetRouteableInterface returns the address and interface of the internet
// -facing network interface
func GetRouteableInterface() (lanInterface *net.Interface) {
	var gw net.IP
	var err error
	if gw, err = gateway.DiscoverGateway(); err != nil {
		ERROR("gateway error: ", err)
		return nil
	}
	gwMasked := gw.Mask(gw.DefaultMask())
	var ifAddrs []net.Addr
	ifAddrs, err = net.InterfaceAddrs()
	if err != nil {
		ERROR("gateway mask error: ", err)
		return nil
	}
	var ad net.IP
	for _, x := range ifAddrs {
		address := strings.Split(x.String(), "/")[0]
		a := net.ParseIP(address)
		masked := a.Mask(gw.DefaultMask())
		if masked.String() == gwMasked.String() {
			ad = a
		}
	}
	if ad == nil {
		ERROR("somehow didn't find a LAN interface even though we" +
			" have a gateway")
		return nil
	}
	consensus := externalip.DefaultConsensus(nil, nil)
	ip, err := consensus.ExternalIP()
	nat := false
	if err != nil {
		ERROR("could not get external IP, " +
			"probably no network connection")
		return nil
	} else {
		if ip.String() != ad.String() {
			nat = true
		}
	}
	if !nat {
		WARN("we are directly on the internet")
	}
	var interfaces []net.Interface
	interfaces, err = net.Interfaces()
	if err != nil {
		ERROR("error:", err)
	}
	for i := range interfaces {
		if ifs, err := interfaces[i].Addrs(); err == nil {
			for j := range ifs {
				ss := strings.Split(ifs[j].String(), "/")
				if ss[0] == ad.String() {
					lanInterface = &interfaces[i]
				}
			}
		}
	}
	return
}
