package asndb

import (
	"fmt"
	"net/netip"
)

//ASN contains information about an ASN.
//the StartIP and EndIP denotes a range that belongs to the ASN.
type ASN struct {
	StartIP       netip.Addr
	EndIP         netip.Addr
	ASNumber      int
	CountryCode   string
	ASDescription string
}

//String returns a string representation of the ASN.
func (a ASN) String() string {
	return fmt.Sprintf("AS%d(%s)@%s [%s->%s]", a.ASNumber, a.ASDescription, a.CountryCode, a.StartIP.String(), a.EndIP.String())
}

//Contains checks if an ip is part of this ASN zone.
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
