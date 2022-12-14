// Package asndb implement asn lookup and data handling.
// This library will handle tsv data sourced from https://iptoasn.com.
//
// The data can be downloaded via DownloadFromURL(DownloadViaIpToAsn), which then parsed using LoadFromTSV(),
// that can be used to initiate a ASList or ASNMap.
package asndb

import (
	"fmt"
	"net/netip"
)

// AS contains information about an AS zone belonging to an ASNumber.
// the StartIP and EndIP denotes a range that belongs to the AS.
type AS struct {
	StartIP       netip.Addr
	EndIP         netip.Addr
	ASNumber      int
	CountryCode   string
	ASDescription string
}

// String returns a string representation of the AS.
func (a AS) String() string {
	ip := "[invalid]"
	if a.StartIP.IsValid() && a.EndIP.IsValid() {
		ip = fmt.Sprintf("[%s->%s]", a.StartIP, a.EndIP)
	}
	return fmt.Sprintf("AS%d(%s)@%s%s", a.ASNumber, a.ASDescription, a.CountryCode, ip)
}

// Contains checks if an ip is part of this AS zone.
func (a AS) Contains(ip netip.Addr) bool {
	return ip.Compare(a.StartIP) >= 0 && ip.Compare(a.EndIP) <= 0
}

type asSortIP []AS

func (a asSortIP) Len() int {
	return len(a)
}

func (a asSortIP) Less(i, j int) bool {
	return a[i].StartIP.Less(a[j].StartIP)
}

func (a asSortIP) Swap(i, j int) {
	a[i], a[j] = a[j], a[i]
}

func clone[S ~[]E, E any](s S) S {
	if s == nil {
		return nil
	}
	return append(make(S, 0, len(s)), s...)
}
