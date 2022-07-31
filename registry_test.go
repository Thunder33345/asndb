package asndb

import (
	"net/netip"
	"reflect"
	"sort"
	"testing"
)

func TestNewRegistry(t *testing.T) {
	tests := []struct {
		name        string
		asn         []AS
		wantASN     []int
		wantZoneLen int
		wantASNLen  int
	}{
		{
			name: "test sort",
			asn: []AS{
				{
					StartIP:  netip.MustParseAddr("7.0.0.0"),
					EndIP:    netip.MustParseAddr("8.0.0.0"),
					ASNumber: 4,
				},
				{
					StartIP:  netip.MustParseAddr("2.0.0.0"),
					EndIP:    netip.MustParseAddr("3.0.0.0"),
					ASNumber: 1,
				},
				{
					StartIP:  netip.MustParseAddr("1.0.0.0"),
					EndIP:    netip.MustParseAddr("2.0.0.0"),
					ASNumber: 0,
				},
				{
					StartIP:  netip.MustParseAddr("4.0.0.0"),
					EndIP:    netip.MustParseAddr("5.0.0.0"),
					ASNumber: 3,
				},
				{
					StartIP:  netip.MustParseAddr("3.0.0.0"),
					EndIP:    netip.MustParseAddr("4.0.0.0"),
					ASNumber: 2,
				}, {
					StartIP:  netip.MustParseAddr("5.0.0.0"),
					EndIP:    netip.MustParseAddr("6.0.0.0"),
					ASNumber: 2,
				}, {
					StartIP:  netip.MustParseAddr("6.0.0.0"),
					EndIP:    netip.MustParseAddr("7.0.0.0"),
					ASNumber: 2,
				},
			},
			wantASN:     []int{0, 1, 2, 3, 2, 2, 4},
			wantZoneLen: 7,
			wantASNLen:  5,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := NewRegistry(tt.asn)
			var asl []int
			for i, want := range tt.wantASN {
				gotAsn := got.s[i].ASNumber
				asl = append(asl, got.s[i].ASNumber)
				if gotAsn != want {
					t.Errorf("ASN.s[%d].ASNumber = %v, want %v", i, gotAsn, want)
				}
			}
			if t.Failed() {
				t.Logf("wanted asn: %v", tt.wantASN)
				t.Logf("got asn: %v", asl)
			}
			gotZoneLen := len(got.s)
			if gotZoneLen != tt.wantZoneLen {
				t.Errorf("ASN.ZoneLen() = %v, want %v", gotZoneLen, tt.wantZoneLen)
			}
			gotASNLen := len(got.m)
			if gotASNLen != tt.wantASNLen {
				t.Errorf("ASN.ASLen() = %v, want %v", gotASNLen, tt.wantASNLen)
			}
		})
	}
}

func TestRegistry_Lookup2(t *testing.T) {
	type subtest struct {
		name      string
		ip        netip.Addr
		wantEmpty bool
		wantASN   int
		wantAS    AS
	}
	tests := []struct {
		name          string
		r             *Registry
		assumeCorrect bool
		lookups       []subtest
	}{
		{
			//this is a simple gap less lookup, checks for upper and lower bounds
			name: "simple lookup",
			r: NewRegistry([]AS{
				{
					StartIP:  netip.MustParseAddr("1.0.0.0"),
					EndIP:    netip.MustParseAddr("1.255.255.255"),
					ASNumber: 1,
				}, {
					StartIP:  netip.MustParseAddr("2.0.0.0"),
					EndIP:    netip.MustParseAddr("2.255.255.255"),
					ASNumber: 2,
				}, {
					StartIP:  netip.MustParseAddr("3.0.0.0"),
					EndIP:    netip.MustParseAddr("3.1.255.255"),
					ASNumber: 3,
				}, {
					StartIP:  netip.MustParseAddr("3.2.0.0"),
					EndIP:    netip.MustParseAddr("3.2.255.255"),
					ASNumber: 4,
				},
			}),
			lookups: []subtest{
				{
					name:      "not found lower bound",
					ip:        netip.MustParseAddr("0.1.1.1"),
					wantEmpty: true,
				},
				{
					name:    "1",
					ip:      netip.MustParseAddr("1.1.1.1"),
					wantASN: 1,
					wantAS: AS{
						StartIP:  netip.MustParseAddr("1.0.0.0"),
						EndIP:    netip.MustParseAddr("1.255.255.255"),
						ASNumber: 1,
					},
				}, {
					name:    "2",
					ip:      netip.MustParseAddr("2.2.2.2"),
					wantASN: 2,
				}, {
					name:    "3-1",
					ip:      netip.MustParseAddr("3.0.3.3"),
					wantASN: 3,
				}, {
					name:    "3-2",
					ip:      netip.MustParseAddr("3.1.3.3"),
					wantASN: 3,
				}, {
					name:    "4",
					ip:      netip.MustParseAddr("3.2.2.2"),
					wantASN: 4,
				}, {
					name:      "not found higher bound",
					ip:        netip.MustParseAddr("55.55.55.55"),
					wantEmpty: true,
				},
			},
		}, {
			//this is a gap-ed lookup, checks for midpoint
			name: "gap test",
			r: NewRegistry([]AS{
				{
					StartIP:  netip.MustParseAddr("1.0.0.0"),
					EndIP:    netip.MustParseAddr("1.255.255.255"),
					ASNumber: 1,
				}, {
					StartIP:  netip.MustParseAddr("3.0.0.0"),
					EndIP:    netip.MustParseAddr("3.255.255.255"),
					ASNumber: 2,
				},
			}),
			lookups: []subtest{
				{
					name:      "not found lower bound",
					ip:        netip.MustParseAddr("0.1.1.1"),
					wantEmpty: true,
				},
				{
					name:    "1",
					ip:      netip.MustParseAddr("1.1.1.1"),
					wantASN: 1,
				}, {
					name:      "middle gap",
					ip:        netip.MustParseAddr("2.2.2.2"),
					wantEmpty: true,
				}, {
					name:    "3-1",
					ip:      netip.MustParseAddr("3.0.3.3"),
					wantASN: 2,
				}, {
					name:    "3-2",
					ip:      netip.MustParseAddr("3.1.3.3"),
					wantASN: 2,
				}, {
					name:      "not found higher bound",
					ip:        netip.MustParseAddr("55.55.55.55"),
					wantEmpty: true,
				},
			},
		}, {
			//this check midpoint but checks if assumptions is correct
			name: "assume valid gap test",
			r: NewRegistry([]AS{
				{
					StartIP:  netip.MustParseAddr("1.0.0.0"),
					EndIP:    netip.MustParseAddr("1.255.255.255"),
					ASNumber: 1,
				}, {
					StartIP:  netip.MustParseAddr("3.0.0.0"),
					EndIP:    netip.MustParseAddr("3.255.255.255"),
					ASNumber: 2,
				}, {
					StartIP:  netip.MustParseAddr("5.0.0.0"),
					EndIP:    netip.MustParseAddr("5.255.255.255"),
					ASNumber: 3,
				},
			}, WithAssumeValid()),
			assumeCorrect: true,
			lookups: []subtest{
				{
					name:    "1",
					ip:      netip.MustParseAddr("1.1.1.1"),
					wantASN: 1,
				}, {
					name:    "middle gap assume 1",
					ip:      netip.MustParseAddr("2.2.2.2"),
					wantASN: 1,
				}, {
					name:    "2",
					ip:      netip.MustParseAddr("3.0.3.3"),
					wantASN: 2,
				}, {
					name:    "middle gap assume 2",
					ip:      netip.MustParseAddr("4.2.2.2"),
					wantASN: 2,
				}, {
					name:    "3",
					ip:      netip.MustParseAddr("5.0.3.3"),
					wantASN: 3,
				},
			},
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			r := tc.r
			for _, tt := range tc.lookups {
				t.Run(tt.name, func(t *testing.T) {
					gotAS, gotFound := r.Lookup(tt.ip)
					if tt.wantEmpty {
						if gotFound {
							t.Errorf("Registry.Lookup() Found = %v, want false", gotFound)
						}
						if gotAS != (AS{}) {
							t.Errorf("Registry.Lookup() AS = %v, want empty", gotAS)
						}
					}

					if (!tt.wantEmpty) != gotFound {
						t.Errorf("Registry.Lookup() Found = %v, want %t", gotFound, !tt.wantEmpty)
					}
					if gotAS.ASNumber != tt.wantASN {
						t.Errorf("Registry.Lookup() ASN = %v, want %v", gotAS.ASNumber, tt.wantASN)
					}
					if tt.wantAS.StartIP.IsValid() && tt.wantAS.EndIP.IsValid() {
						if !reflect.DeepEqual(gotAS, tt.wantAS) {
							t.Errorf("Registry.Lookup() AS = %v, want %v", gotAS, tt.wantAS)
						}
					}
					if tc.assumeCorrect == false && !tt.wantEmpty {
						if gotAS.Contains(tt.ip) == false {
							t.Errorf("Registry.Lookup() AS does not actually contain %v", tt.ip)
						}
					}
					if t.Failed() {
						t.Logf("Registry.Lookup() AS = %v, Found = %v, ASN = %d", gotAS, gotFound, gotAS.ASNumber)
					}
				})
			}
		})
	}
}

func TestRegistry_Integration(t *testing.T) {
	ASNs := []AS{
		{
			StartIP:       netip.MustParseAddr("7.0.0.0"),
			EndIP:         netip.MustParseAddr("8.0.0.0"),
			ASNumber:      4,
			ASDescription: "the 4th asn",
		},
		{
			StartIP:  netip.MustParseAddr("2.0.0.0"),
			EndIP:    netip.MustParseAddr("3.0.0.0"),
			ASNumber: 1,
		},
		{
			StartIP:       netip.MustParseAddr("1.0.0.0"),
			EndIP:         netip.MustParseAddr("2.0.0.0"),
			ASNumber:      0,
			ASDescription: "the 0th asn",
			CountryCode:   "nil",
		},
		{
			StartIP:  netip.MustParseAddr("4.0.0.0"),
			EndIP:    netip.MustParseAddr("5.0.0.0"),
			ASNumber: 3,
		},
		{
			StartIP:  netip.MustParseAddr("3.0.0.0"),
			EndIP:    netip.MustParseAddr("4.0.0.0"),
			ASNumber: 2,
		}, {
			StartIP:  netip.MustParseAddr("5.0.0.0"),
			EndIP:    netip.MustParseAddr("6.0.0.0"),
			ASNumber: 2,
		}, {
			StartIP:  netip.MustParseAddr("6.0.0.0"),
			EndIP:    netip.MustParseAddr("7.0.0.0"),
			ASNumber: 2,
		},
	}
	r := NewRegistry(ASNs)
	wantASNOrder := []int{0, 1, 2, 3, 2, 2, 4}
	wantZoneLen := 7
	wantASNLen := 5
	wantASList := []struct {
		asn       int
		wantCount int
		wantFound bool
	}{
		{
			asn:       10,
			wantFound: false,
		},
		{
			asn:       0,
			wantCount: 1,
			wantFound: true,
		}, {
			asn:       2,
			wantCount: 3,
			wantFound: true,
		}, {
			asn:       4,
			wantCount: 1,
			wantFound: true,
		},
	}
	wantLookup := []struct {
		ip      netip.Addr
		wantASN int
	}{
		{
			ip:      netip.MustParseAddr("0.0.0.0"),
			wantASN: -1,
		}, {
			ip:      netip.MustParseAddr("1.0.0.0"),
			wantASN: 0,
		}, {
			ip:      netip.MustParseAddr("5.0.0.0"),
			wantASN: 2,
		}, {
			ip:      netip.MustParseAddr("7.0.0.0"),
			wantASN: 4,
		}, {
			ip:      netip.MustParseAddr("200.0.0.0"),
			wantASN: -1,
		},
	}
	wantASNList := []AS{
		{
			ASNumber:      0,
			ASDescription: "the 0th asn",
			CountryCode:   "nil",
		},
		{
			ASNumber: 1,
		}, {
			ASNumber: 2,
		}, {
			ASNumber: 3,
		}, {
			ASNumber:      4,
			ASDescription: "the 4th asn",
		},
	}

	t.Run("Alter Input", func(t *testing.T) {
		for i, a := range ASNs {
			a.ASNumber = -999
			ASNs[i] = a
		}
		for _, a := range r.s {
			if a.ASNumber == -999 {
				t.Errorf("ASN.s should not have been altered")
				return
			}
		}
	})

	t.Run("Order", func(t *testing.T) {
		var order []int
		for i := 0; i < len(r.s); i++ {
			gotAsn := r.s[i].ASNumber
			if gotAsn != wantASNOrder[i] {
				t.Errorf("ASN.s[%d].ASNumber = %v, want %v", i, order[i], wantASNOrder[i])
			}
		}
		if t.Failed() {
			t.Logf("ASN.s order: %v", order)
			t.Logf("Wanted order: %v", wantASNOrder)
		}

	})

	t.Run("ZoneLen", func(t *testing.T) {
		gotZoneLen := len(r.s)
		if gotZoneLen != wantZoneLen {
			t.Errorf("ASN.ZoneLen() = %v, want %v", gotZoneLen, wantZoneLen)
		}
	})

	t.Run("ASNLen", func(t *testing.T) {
		gotASNLen := len(r.m)
		if gotASNLen != wantASNLen {
			t.Errorf("ASN.ASLen() = %v, want %v", gotASNLen, wantASNLen)
		}
	})

	t.Run("ASList", func(t *testing.T) {
		for _, tt := range wantASList {
			gotASNs, gotFound := r.ListZone(tt.asn)
			if gotFound != tt.wantFound {
				t.Errorf("ASN.ASList() found = %v, want %v", gotFound, tt.wantFound)
			}
			if tt.wantFound == false {
				continue
			}
			if len(gotASNs) != tt.wantCount {
				t.Errorf("ASN.ASList() count = %v, want %v", gotASNs, tt.wantCount)
			}
		}
	})

	t.Run("ASList Alter", func(t *testing.T) {
		l, f := r.ListZone(0)
		if !f {
			t.Errorf("ASN.ListZone(0) should exist, but none found")
		}
		l[0].ASNumber = 10
		l2, _ := r.ListZone(0)
		if l2[0].ASNumber != 0 {
			t.Errorf("results of ASN.ListZone() should not be altered")
		}
	})

	t.Run("Lookup", func(t *testing.T) {
		for _, s := range wantLookup {
			a, b := r.Lookup(s.ip)
			if s.wantASN == -1 {
				d := AS{}
				if a != d {
					t.Errorf("ASN.Lookup(%v).ASNumber = %#v, wanted empty", s.ip, a)
				}
				if b != false {
					t.Errorf("ASN.Lookup(%v) b = %v, want false", s.ip, b)
				}
			} else {
				if a.ASNumber != s.wantASN {
					t.Errorf("ASN.Lookup(%v).ASNumber = %v, want %v", s.ip, a.ASNumber, s.wantASN)
				}
				if b != true {
					t.Errorf("ASN.Lookup(%v) b = %v, want true", s.ip, b)
				}
			}
		}
	})

	t.Run("ListASN", func(t *testing.T) {
		l := r.ListASN()
		for i, as := range wantASNList {
			if l[i] != as {
				t.Errorf("ASN.ListASN()[%d] = %v, want %v", i, l[i].ASNumber, as.ASNumber)
			}
		}
	})

	_ = r.s[0].String()
}

func TestAsSort(t *testing.T) {
	asl := []AS{
		{
			ASNumber: 0,
			StartIP:  netip.MustParseAddr("0.5.0.0"),
		},
		{
			ASNumber: 1,
			StartIP:  netip.MustParseAddr("1.0.0.0"),
		}, {
			ASNumber: 3,
			StartIP:  netip.MustParseAddr("3.0.0.0"),
		}, {
			ASNumber: 2,
			StartIP:  netip.MustParseAddr("2.0.0.0"),
		}, {
			ASNumber: 4,
			StartIP:  netip.MustParseAddr("4.0.0.0"),
		},
	}
	t.Run("SortIP", func(t *testing.T) {
		c := asSortIP(clone(asl))
		sort.Sort(c)
		for i, as := range c {
			if as.ASNumber != i {
				t.Errorf("list[%d].ASNumber = %v, want %v", i, as.ASNumber, i)
			}
		}
		if t.Failed() {
			t.Logf("asl: %v", c)
		}
		if reflect.DeepEqual(asl, c) {
			t.Errorf("asl and c should not be equal(asl altered?)")
		}
	})
	t.Run("SortASN", func(t *testing.T) {
		c := asSortASN(clone(asl))
		sort.Sort(c)
		for i, as := range c {
			if as.ASNumber != i {
				t.Errorf("list[%d].ASNumber = %v, want %v", i, as.ASNumber, i)
			}
		}
		if t.Failed() {
			t.Logf("asl: %v", c)
		}
		if reflect.DeepEqual(asl, c) {
			t.Errorf("asl and c should not be equal(asl altered?)")
		}
	})
}
