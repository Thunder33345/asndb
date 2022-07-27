package asndb

import (
	"net/netip"
	"sort"
)

func NewRegistry(s []ASN) *Registry {
	sort.Sort(asnList(s))
	return &Registry{s}
}

type Registry struct {
	s asnList
}

func (r *Registry) Lookup(ip netip.Addr) (ASN, bool) {
	index := sort.Search(len(r.s),
		//this function should not be moved into a method
		//otherwise heap allocations will be made
		func(i int) bool {
			return ip.Less(r.s[i].StartIP)
		})
	//index will always be offset by +1
	index--
	if index < len(r.s) {
		return r.s[index], true
	}
	return ASN{}, false
}
