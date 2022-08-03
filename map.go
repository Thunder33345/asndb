package asndb

import "sort"

type ASNMap struct {
	m map[int][]AS
}

func NewASNMap(s []AS) *ASNMap {
	m := make(map[int][]AS)
	for _, asn := range s {
		m[asn.ASNumber] = append(m[asn.ASNumber], asn)
	}

	return &ASNMap{m: m}
}

// ListAS returns a list of AS zones controlled by given asn.
// The returned slice will be cloned and can be freely edited.
func (m *ASNMap) ListAS(asn int) ([]AS, bool) {
	s, ok := m.m[asn]
	return clone(s), ok
}

// ListASN returns a list of ASN.
// Behaviour of AS's details are undefined if details are inconsistent.
// AS.StartIP and AS.EndIP will not be defined.
func (m *ASNMap) ListASN() []AS {
	s := make([]AS, 0, len(m.m))
	for asn, as := range m.m {
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

type asSortASN []AS

func (a asSortASN) Len() int {
	return len(a)
}

func (a asSortASN) Less(i, j int) bool {
	return a[i].ASNumber < a[j].ASNumber
}

func (a asSortASN) Swap(i, j int) {
	a[i], a[j] = a[j], a[i]
}
