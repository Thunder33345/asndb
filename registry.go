package asndb

import (
	"net/netip"
	"sort"
)

func NewRegistry(s []ASN) *Registry {
	sort.Sort(asnList(s))
	s = s[:len(s):len(s)]
	return &Registry{s: s}
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
	//valid indexes will always be offset by +1
	//so anything less than 1 is invalid
	if index <= 0 {
		return ASN{}, false
	}
	index--
	if index < len(r.s) {
		return r.s[index], true
	}
	return ASN{}, false
}
