package asndb

import (
	"net/netip"
	"sort"
)

//NewRegistry creates a new registry from the given list of AS zones.
//The given slice will be cloned and sorted by StartIP.
func NewRegistry(s []AS, opts ...Option) *Registry {
	s = clone(s)
	sort.Sort(asSortIP(s))
	s = s[:len(s):len(s)]

	m := make(map[int][]AS)
	for _, asn := range s {
		m[asn.ASNumber] = append(m[asn.ASNumber], asn)
	}

	r := &Registry{s: s, m: m}
	for _, opt := range opts {
		opt(r)
	}

	return r
}

//Registry holds a list of AS zones.
type Registry struct {
	s           []AS
	m           map[int][]AS
	assumeValid bool
}

//Lookup finds and returns the AS zone for a given IP address.
//Bool indicates if AS is valid and found
//Notice: if multiple zones claims an IP, the closest AS zone gets returned.
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
		index = len(r.s) - 1
		if r.s[index].Contains(ip) {
			return r.s[index], true
		}
		return AS{}, false
	}

	//if we don't care about possible inaccuracies that will occur in a gap of unclaimed ips between AS zones
	if r.assumeValid {
		return r.s[index], true
	}

	//otherwise we check before returning
	if r.s[index].Contains(ip) {
		return r.s[index], true
	}
	return AS{}, false
}

//MultiLookup attempts to find and return neighbouring AS that contain given ip address.
func (r *Registry) MultiLookup(ip netip.Addr, search int) []AS {
	//get an index
	index := sort.Search(len(r.s),
		func(i int) bool {
			return ip.Less(r.s[i].StartIP)
		})
	//create a slice of AS
	var s []AS
	//remove offset
	index--
	if index > len(r.s)-1 {
		index = len(r.s) - 1
	}

	//loop that counts form 0 till search (acting as a search space)
	//we only search downwards because the slice get sorted by AS.StartIP
	//so it's not possible to have any AS above index that can claim the ip
	for i := 0; i <= search; i++ {
		//create an offset index
		ix := index - i
		if ix < 0 {
			break
		}
		//if the AS contains the IP, we add it to the slice
		if r.s[ix].Contains(ip) {
			s = append(s, r.s[ix])
			//expand the search space for every valid result
			search++
		}
	}

	return s
}

//ListZone returns a list of AS zones controlled by given asn.
//The returned slice will be cloned and can be freely edited.
func (r *Registry) ListZone(asn int) ([]AS, bool) {
	s, ok := r.m[asn]
	return clone(s), ok
}

//ListASN returns a list of ASN.
//Behaviour of AS's details are undefined if details are inconsistent.
//AS.StartIP and AS.EndIP will not be defined.
func (r *Registry) ListASN() []AS {
	s := make([]AS, 0, len(r.m))
	for asn, as := range r.m {
		t := AS{
			ASNumber: asn,
		}
		if len(as) >= 1 {
			t.CountryCode = as[0].CountryCode
			t.ASDescription = as[0].ASDescription
		}
		s = append(s, t)
	}
	sort.Sort(asSortASN(s))
	return s
}

func clone[S ~[]E, E any](s S) S {
	if s == nil {
		return nil
	}
	return append(make(S, 0, len(s)), s...)
}
