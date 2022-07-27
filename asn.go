package asndb

import (
	"fmt"
	"net/netip"
)

type ASN struct {
	StartIP       netip.Addr
	EndIP         netip.Addr
	ASNumber      int
	CountryCode   string
	ASDescription string
}

func (a ASN) String() string {
	return fmt.Sprintf("AS%d(%s)@%s [%s->%s]", a.ASNumber, a.ASDescription, a.CountryCode, a.StartIP.String(), a.EndIP.String())
}

func (a ASN) Contains(ip netip.Addr) bool {
	return ip.Compare(a.StartIP) >= 0 && ip.Compare(a.EndIP) <= 0
}

type asnList []ASN

func (a asnList) Len() int {
	return len(a)
}

func (a asnList) Less(i, j int) bool {
	return a[i].StartIP.Less(a[j].StartIP)
}

func (a asnList) Swap(i, j int) {
	a[i], a[j] = a[j], a[i]
}
