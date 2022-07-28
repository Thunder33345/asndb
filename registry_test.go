package asndb

import (
	"net/netip"
	"reflect"
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
			for _, i := range tt.wantASN {
				gotAsn := got.s[i].ASNumber
				if gotAsn != i {
					t.Errorf("ASN.s[%d].ASNumber = %v, want %v", i, gotAsn, i)
				}
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

func TestRegistry_Lookup(t *testing.T) {
	r := NewRegistry([]AS{
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
	})

	tests := []struct {
		name      string
		ip        netip.Addr
		want      int
		wantFound bool
	}{
		{
			name:      "not found lower bound",
			ip:        netip.MustParseAddr("0.1.1.1"),
			wantFound: false,
		},
		{
			name:      "1",
			ip:        netip.MustParseAddr("1.1.1.1"),
			want:      1,
			wantFound: true,
		}, {
			name:      "2",
			ip:        netip.MustParseAddr("2.2.2.2"),
			want:      2,
			wantFound: true,
		}, {
			name:      "3-1",
			ip:        netip.MustParseAddr("3.0.3.3"),
			want:      3,
			wantFound: true,
		}, {
			name:      "3-2",
			ip:        netip.MustParseAddr("3.1.3.3"),
			want:      3,
			wantFound: true,
		}, {
			name:      "4",
			ip:        netip.MustParseAddr("3.2.2.2"),
			want:      4,
			wantFound: true,
		}, {
			name:      "not found higher bound",
			ip:        netip.MustParseAddr("55.55.55.55"),
			wantFound: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1 := r.Lookup(tt.ip)
			if got1 != tt.wantFound {
				t.Errorf("Lookup() found = %v, want %v", got1, tt.wantFound)
			}
			if tt.wantFound == false {
				return
			}
			if !reflect.DeepEqual(got.ASNumber, tt.want) {
				t.Errorf("Lookup() ASN = %v, want %v", got.ASNumber, tt.want)
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
				t.Errorf("ASN.ListASN()[%d] = %v, want %v", i, l[i], as)
			}
		}
	})

	_ = r.s[0].String()
}
