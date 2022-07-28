package asndb

import (
	"net/netip"
	"sort"
)

//NewRegistry creates a new registry from the given list of AS zones.
//The slice will automatically be sorted by StartIP.
//Given slice should not be modified afterwards.
func NewRegistry(s []AS) *Registry {
	sort.Sort(asList(s))
	s = s[:len(s):len(s)]
	return &Registry{s: s}
}

//Registry holds a list of AS zones.
type Registry struct {
	s asList
}

//Lookup finds and returns the AS zone for a given IP address.
//Bool indicates if AS is valid and found
func (r *Registry) Lookup(ip netip.Addr) (AS, bool) {
	index := sort.Search(len(r.s),
		//this function should not be moved into a method
		//otherwise heap allocations will be made
		func(i int) bool {
			return ip.Less(r.s[i].StartIP)
		})
	//index will always be offset by +1 due to our sorting method, so we need to subtract 1
	index--
	//when the index is negative its bellow our lower bound
	if index < 0 {
		return AS{}, false
	}

	//when the index is equal to the length of the slice
	//we have to manually check if the ip is part of the last AS zone
	//or is it actually above our higher bound
	if index >= len(r.s)-1 {
		if r.s[index].Contains(ip) {
			return r.s[index], true
		}
		return AS{}, false
	}
	return r.s[index], true
}
