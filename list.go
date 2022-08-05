package asndb

import (
	"net/netip"
	"sort"
)

// NewASList creates a new registry from the given list of AS zones.
// The given slice will be cloned and sorted by StartIP.
func NewASList(s []AS) *ASList {
	s = clone(s)
	sort.Sort(asSortIP(s))
	s = s[:len(s):len(s)]

	r := &ASList{s: s}
	return r
}

// ASList holds a list of AS zones.
type ASList struct {
	s []AS
}

// Find finds and returns the AS zone for a given IP address.
// Bool indicates if AS is valid and found
// Notice: if multiple zones claims an IP, the closest AS zone gets returned.
func (r *ASList) Find(ip netip.Addr) (AS, bool) {
	//get an index
	index := r.Index(ip)
	//when the index is negative its bellow our lower bound
	if index < 0 {
		return AS{}, false
	}

	//we check if the AS actually contains the IP
	if !r.s[index].Contains(ip) {
		return AS{}, false
	}
	return r.s[index], true
}

// FindList attempts to find and return neighbouring AS that contain given ip address.
// search dictate how many invalid AS zones to skip before returning.
// This method is only useful when an IP has been claimed by multiple AS zones.
func (r *ASList) FindList(ip netip.Addr, search uint) []AS {
	//get an index
	index := r.Index(ip)
	//create a slice of AS
	var s []AS

	//loop that counts form 0 till search (acting as a search space)
	//we only search downwards because the slice get sorted by AS.StartIP
	//so it's not possible to have any AS above index that can claim the ip
	for i := 0; i <= int(search); i++ {
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

// Index returns an index closest to AS zone for a given IP address.
// Does not guarantee the AS of said index contains the IP.
// If the index is out of bounds, it returns -1.
func (r *ASList) Index(ip netip.Addr) int {
	//we use sort.Search to find the closest index, using the AS zone's StartIP as comparison
	index := sort.Search(len(r.s),
		func(i int) bool {
			return ip.Less(r.s[i].StartIP)
		})
	//index is actually off by one, so we decrement it
	index--
	if index < 0 || index >= len(r.s) {
		return -1
	}

	return index
}

// FromIndex returns an AS zone at a given index.
// Returns false if the index is out of bounds.
func (r *ASList) FromIndex(i int) (AS, bool) {
	if i < 0 || i >= len(r.s) {
		return AS{}, false
	}
	return r.s[i], true
}

// IndexLen returns the length of the AS zone.
func (r *ASList) IndexLen() int {
	return len(r.s)
}
