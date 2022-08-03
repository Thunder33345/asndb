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
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := NewASList(tt.asn)
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
		})
	}
}

func TestRegistry_Find(t *testing.T) {
	type subtest struct {
		name      string
		ip        netip.Addr
		wantEmpty bool
		wantASN   int
		wantAS    AS
	}
	tests := []struct {
		name          string
		r             *ASList
		assumeCorrect bool
		lookups       []subtest
	}{
		{
			//this is a simple gap less lookup, checks for upper and lower bounds
			name: "simple lookup",
			r: NewASList([]AS{
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
			r: NewASList([]AS{
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
			r: NewASList([]AS{
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
		}, {
			name: "precise list",
			r: NewASList([]AS{
				{
					StartIP:  netip.MustParseAddr("1.0.0.1"),
					EndIP:    netip.MustParseAddr("1.0.0.1"),
					ASNumber: 1,
				}, {
					StartIP:  netip.MustParseAddr("1.0.0.2"),
					EndIP:    netip.MustParseAddr("1.0.0.2"),
					ASNumber: 2,
				}, {
					StartIP:  netip.MustParseAddr("1.0.0.3"),
					EndIP:    netip.MustParseAddr("1.0.0.3"),
					ASNumber: 3,
				}, {
					StartIP:  netip.MustParseAddr("1.0.0.4"),
					EndIP:    netip.MustParseAddr("1.0.0.4"),
					ASNumber: 4,
				}, {
					StartIP:  netip.MustParseAddr("1.0.0.5"),
					EndIP:    netip.MustParseAddr("1.0.0.5"),
					ASNumber: 5,
				},
			}, WithAssumeValid()),
			lookups: []subtest{
				{
					name:    "1",
					ip:      netip.MustParseAddr("1.0.0.1"),
					wantASN: 1,
				}, {
					name:    "2",
					ip:      netip.MustParseAddr("1.0.0.2"),
					wantASN: 2,
				}, {
					name:    "3",
					ip:      netip.MustParseAddr("1.0.0.3"),
					wantASN: 3,
				}, {
					name:    "4",
					ip:      netip.MustParseAddr("1.0.0.4"),
					wantASN: 4,
				}, {
					name:    "5",
					ip:      netip.MustParseAddr("1.0.0.5"),
					wantASN: 5,
				},
			},
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			r := tc.r
			for _, tt := range tc.lookups {
				t.Run(tt.name, func(t *testing.T) {
					gotAS, gotFound := r.Find(tt.ip)
					if tt.wantEmpty {
						if gotFound {
							t.Errorf("ASList.Find() Found = %v, want false", gotFound)
						}
						if gotAS != (AS{}) {
							t.Errorf("ASList.Find() AS = %v, want empty", gotAS)
						}
					}

					if (!tt.wantEmpty) != gotFound {
						t.Errorf("ASList.Find() Found = %v, want %t", gotFound, !tt.wantEmpty)
					}
					if gotAS.ASNumber != tt.wantASN {
						t.Errorf("ASList.Find() ASN = %v, want %v", gotAS.ASNumber, tt.wantASN)
					}
					if tt.wantAS.StartIP.IsValid() && tt.wantAS.EndIP.IsValid() {
						if !reflect.DeepEqual(gotAS, tt.wantAS) {
							t.Errorf("ASList.Find() AS = %v, want %v", gotAS, tt.wantAS)
						}
					}
					if tc.assumeCorrect == false && !tt.wantEmpty {
						if gotAS.Contains(tt.ip) == false {
							t.Errorf("ASList.Find() AS does not actually contain %v", tt.ip)
						}
					}
					if t.Failed() {
						t.Logf("ASList.Find() AS = %v, Found = %v, ASN = %d", gotAS, gotFound, gotAS.ASNumber)
					}
				})
			}
		})
	}
}

func TestRegistry_FindList(t *testing.T) {
	type subtest struct {
		name    string
		ip      netip.Addr
		search  uint
		wantASN []int
		//wantAS  []AS
	}
	tests := []struct {
		name    string
		r       *ASList
		lookups []subtest
	}{
		{
			name: "overlapping 1",
			r: NewASList([]AS{
				{
					ASNumber: 1,
					StartIP:  netip.MustParseAddr("1.0.0.0"),
					EndIP:    netip.MustParseAddr("1.255.255.255"),
				}, {
					ASNumber: 1,
					StartIP:  netip.MustParseAddr("1.5.0.0"),
					EndIP:    netip.MustParseAddr("1.255.255.255"),
				}, {
					ASNumber: 2,
					StartIP:  netip.MustParseAddr("1.30.0.0"),
					EndIP:    netip.MustParseAddr("2.255.255.255"),
				}, {
					ASNumber: 2,
					StartIP:  netip.MustParseAddr("1.50.0.0"),
					EndIP:    netip.MustParseAddr("2.255.255.255"),
				}, {
					ASNumber: 100,
					StartIP:  netip.MustParseAddr("1.50.0.0"),
					EndIP:    netip.MustParseAddr("1.50.255.255"),
				}, {
					ASNumber: 3,
					StartIP:  netip.MustParseAddr("1.60.0.0"),
					EndIP:    netip.MustParseAddr("3.255.255.255"),
				}, {
					ASNumber: 100,
					StartIP:  netip.MustParseAddr("1.61.0.0"),
					EndIP:    netip.MustParseAddr("1.61.255.255"),
				}, {
					ASNumber: 3,
					StartIP:  netip.MustParseAddr("1.62.0.0"),
					EndIP:    netip.MustParseAddr("3.255.255.255"),
				}, {
					ASNumber: 100,
					StartIP:  netip.MustParseAddr("1.63.0.0"),
					EndIP:    netip.MustParseAddr("1.61.255.255"),
				}, {
					ASNumber: 4,
					StartIP:  netip.MustParseAddr("1.70.0.0"),
					EndIP:    netip.MustParseAddr("4.255.255.255"),
				}, {
					ASNumber: 5,
					StartIP:  netip.MustParseAddr("5.0.0.0"),
					EndIP:    netip.MustParseAddr("5.255.255.255"),
				}}),
			lookups: []subtest{
				{
					name:    "first",
					ip:      netip.MustParseAddr("1.80.0.0"),
					search:  3,
					wantASN: []int{4, 3, 3, 2, 2, 1, 1},
				}, {
					name:    "second",
					ip:      netip.MustParseAddr("1.80.0.0"),
					search:  2,
					wantASN: []int{4, 3, 3},
				}, {
					name:    "third",
					ip:      netip.MustParseAddr("1.80.0.0"),
					search:  1,
					wantASN: []int{4, 3},
				}, {
					name:    "fourth",
					ip:      netip.MustParseAddr("1.80.0.0"),
					search:  0,
					wantASN: []int{4},
				}, {
					name:    "last",
					ip:      netip.MustParseAddr("5.5.0.0"),
					search:  1,
					wantASN: []int{5},
				},
			},
		}, {
			name: "overlapping 2",
			r: NewASList([]AS{
				{
					ASNumber: 100,
					StartIP:  netip.MustParseAddr("0.1.0.0"),
					EndIP:    netip.MustParseAddr("0.1.255.255"),
				}, {
					ASNumber: 1,
					StartIP:  netip.MustParseAddr("1.0.0.0"),
					EndIP:    netip.MustParseAddr("1.255.255.255"),
				}, {
					ASNumber: 1,
					StartIP:  netip.MustParseAddr("1.5.0.0"),
					EndIP:    netip.MustParseAddr("1.255.255.255"),
				}, {
					ASNumber: 2,
					StartIP:  netip.MustParseAddr("1.30.0.0"),
					EndIP:    netip.MustParseAddr("2.255.255.255"),
				}, {
					ASNumber: 2,
					StartIP:  netip.MustParseAddr("1.50.0.0"),
					EndIP:    netip.MustParseAddr("2.255.255.255"),
				}, {
					ASNumber: 100,
					StartIP:  netip.MustParseAddr("1.52.0.0"),
					EndIP:    netip.MustParseAddr("1.50.255.255"),
				}, {
					ASNumber: 3,
					StartIP:  netip.MustParseAddr("1.54.0.0"),
					EndIP:    netip.MustParseAddr("3.255.255.255"),
				}, {
					ASNumber: 100,
					StartIP:  netip.MustParseAddr("1.56.0.0"),
					EndIP:    netip.MustParseAddr("1.61.255.255"),
				}, {
					ASNumber: 4,
					StartIP:  netip.MustParseAddr("1.58.0.0"),
					EndIP:    netip.MustParseAddr("4.255.255.255"),
				}, {
					ASNumber: 5,
					StartIP:  netip.MustParseAddr("1.60.0.0"),
					EndIP:    netip.MustParseAddr("5.0.255.255"),
				}, {
					ASNumber: 5,
					StartIP:  netip.MustParseAddr("1.64.0.0"),
					EndIP:    netip.MustParseAddr("5.255.255.255"),
				}, {
					ASNumber: 100,
					StartIP:  netip.MustParseAddr("1.66.0.0"),
					EndIP:    netip.MustParseAddr("1.65.255.255"),
				}, {
					ASNumber: 6,
					StartIP:  netip.MustParseAddr("1.68.0.0"),
					EndIP:    netip.MustParseAddr("6.255.255.255"),
				}, {
					ASNumber: 7,
					StartIP:  netip.MustParseAddr("1.69.0.0"),
					EndIP:    netip.MustParseAddr("7.255.255.255"),
				}, {
					ASNumber: 100,
					StartIP:  netip.MustParseAddr("1.69.5.0"),
					EndIP:    netip.MustParseAddr("1.69.255.255"),
				}, {
					ASNumber: 8,
					StartIP:  netip.MustParseAddr("1.70.0.0"),
					EndIP:    netip.MustParseAddr("8.255.255.255"),
				}, {
					ASNumber: 9,
					StartIP:  netip.MustParseAddr("1.79.0.0"),
					EndIP:    netip.MustParseAddr("9.255.255.255"),
				}, {
					ASNumber: 10,
					StartIP:  netip.MustParseAddr("1.80.0.0"),
					EndIP:    netip.MustParseAddr("10.255.255.255"),
				}, {
					ASNumber: 11,
					StartIP:  netip.MustParseAddr("1.80.0.1"),
					EndIP:    netip.MustParseAddr("11.255.255.255"),
				},
			},
			),
			lookups: []subtest{
				{
					name:    "zero",
					ip:      netip.MustParseAddr("1.80.0.0"),
					search:  0,
					wantASN: []int{10, 9, 8},
				}, {
					name:    "one",
					ip:      netip.MustParseAddr("1.80.0.0"),
					search:  1,
					wantASN: []int{10, 9, 8, 7, 6},
				}, {
					name:    "two",
					ip:      netip.MustParseAddr("1.80.0.0"),
					search:  2,
					wantASN: []int{10, 9, 8, 7, 6, 5, 5, 4},
				}, {
					name:    "three",
					ip:      netip.MustParseAddr("1.80.0.0"),
					search:  3,
					wantASN: []int{10, 9, 8, 7, 6, 5, 5, 4, 3},
				}, {
					name:    "four",
					ip:      netip.MustParseAddr("1.80.0.0"),
					search:  4,
					wantASN: []int{10, 9, 8, 7, 6, 5, 5, 4, 3, 2, 2, 1, 1},
				}, {
					name:    "five",
					ip:      netip.MustParseAddr("1.80.0.0"),
					search:  5,
					wantASN: []int{10, 9, 8, 7, 6, 5, 5, 4, 3, 2, 2, 1, 1},
				}, {
					name:    "ten",
					ip:      netip.MustParseAddr("1.80.0.0"),
					search:  10,
					wantASN: []int{10, 9, 8, 7, 6, 5, 5, 4, 3, 2, 2, 1, 1},
				}, {
					name:    "last",
					ip:      netip.MustParseAddr("10.0.0.0"),
					search:  0,
					wantASN: []int{11, 10},
				}, {
					name:    "5-1",
					ip:      netip.MustParseAddr("1.64.0.0"),
					search:  2,
					wantASN: []int{5, 5, 4, 3, 2, 2, 1, 1},
				},
			},
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			r := tc.r
			for _, tt := range tc.lookups {
				t.Run(tt.name, func(t *testing.T) {
					gotASList := r.FindList(tt.ip, tt.search)
					col := make([]int, 0, len(gotASList))
					if len(gotASList) != len(tt.wantASN) {
						t.Errorf("ASList.FindList() length = %v, want %v", len(gotASList), len(tt.wantASN))
					}

					for i, as := range gotASList {
						if !as.Contains(tt.ip) {
							t.Errorf("ASList.FindList()[%d] does not actually contain %v", i, tt.ip)
						}
					}

					if tt.wantASN != nil {
						for i, as := range gotASList {
							col = append(col, as.ASNumber)
							if i >= len(tt.wantASN) {
								t.Errorf("ASList.FindList()[%d] = %v, but is not wanted", i, as.ASNumber)
								continue
							}
							if as.ASNumber != tt.wantASN[i] {
								t.Errorf("ASList.FindList()[%d] = %v, want %v", i, as.ASNumber, tt.wantASN[i])
							}
						}
						if t.Failed() {
							t.Logf("ASList.FindList() = %v, want %v", col, tt.wantASN)
						}
					}
				})
			}
		})
	}
}

func TestRegistry_GetIndex(t *testing.T) {
	reg1 := NewASList([]AS{
		{ASNumber: 0, StartIP: netip.MustParseAddr("1.1.0.0"), EndIP: netip.MustParseAddr("1.1.255.255")},
		{ASNumber: 1, StartIP: netip.MustParseAddr("1.2.0.0"), EndIP: netip.MustParseAddr("1.2.255.255")},
		{ASNumber: 2, StartIP: netip.MustParseAddr("1.3.0.0"), EndIP: netip.MustParseAddr("1.3.255.255")},
	})
	reg2 := NewASList([]AS{
		{StartIP: netip.MustParseAddr("0.0.0.5"), EndIP: netip.MustParseAddr("1.0.0.0"), ASNumber: 0},
		{StartIP: netip.MustParseAddr("1.0.0.1"), EndIP: netip.MustParseAddr("1.0.0.1"), ASNumber: 1},
		{StartIP: netip.MustParseAddr("1.0.0.2"), EndIP: netip.MustParseAddr("1.0.0.2"), ASNumber: 2},
		{StartIP: netip.MustParseAddr("1.0.0.3"), EndIP: netip.MustParseAddr("1.0.0.3"), ASNumber: 3},
		{StartIP: netip.MustParseAddr("1.0.0.4"), EndIP: netip.MustParseAddr("1.0.0.4"), ASNumber: 4},
		{StartIP: netip.MustParseAddr("1.0.0.5"), EndIP: netip.MustParseAddr("1.0.0.5"), ASNumber: 5},
	}, WithAssumeValid())

	tests := []struct {
		name string
		r    *ASList
		ip   netip.Addr
		want int
	}{
		{name: "-1", r: reg1, ip: netip.MustParseAddr("0.0.0.0"), want: -1},
		{name: "-1", r: reg1, ip: netip.MustParseAddr("0.1.0.0"), want: -1},
		{name: "0", r: reg1, ip: netip.MustParseAddr("1.1.0.0"), want: 0},
		{name: "0-2", r: reg1, ip: netip.MustParseAddr("1.1.1.1"), want: 0},
		{name: "0-2", r: reg1, ip: netip.MustParseAddr("1.1.255.255"), want: 0},
		{name: "1", r: reg1, ip: netip.MustParseAddr("1.2.0.0"), want: 1},
		{name: "1-2", r: reg1, ip: netip.MustParseAddr("1.2.10.0"), want: 1},
		{name: "1-3", r: reg1, ip: netip.MustParseAddr("1.2.255.255"), want: 1},
		{name: "2-1", r: reg1, ip: netip.MustParseAddr("1.3.0.0"), want: 2},
		{name: "2-2", r: reg1, ip: netip.MustParseAddr("1.3.0.5"), want: 2},
		{name: "2-3", r: reg1, ip: netip.MustParseAddr("1.3.255.255"), want: 2},
		{name: "+2", r: reg1, ip: netip.MustParseAddr("1.4.255.255"), want: 2},

		{name: "p-1", r: reg2, ip: netip.MustParseAddr("0.0.0.1"), want: -1},
		{name: "p0", r: reg2, ip: netip.MustParseAddr("0.5.5.5"), want: 0},
		{name: "p1", r: reg2, ip: netip.MustParseAddr("1.0.0.1"), want: 1},
		{name: "p2", r: reg2, ip: netip.MustParseAddr("1.0.0.2"), want: 2},
		{name: "p2", r: reg2, ip: netip.MustParseAddr("1.0.0.3"), want: 3},
		{name: "p2", r: reg2, ip: netip.MustParseAddr("1.0.0.4"), want: 4},
		{name: "p2", r: reg2, ip: netip.MustParseAddr("1.0.0.5"), want: 5},
		{name: "p2", r: reg2, ip: netip.MustParseAddr("1.0.0.6"), want: 5},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.r.getIndex(tt.ip); got != tt.want {
				t.Errorf("ASList.GetIndex(%v) = %v, want %v", tt.ip, got, tt.want)
				const offset = 3
				start := got - offset
				if start < 0 {
					start = 0
				}
				for i := start; i < got+offset+1; i++ {
					s := " "
					switch i {
					case got:
						s = ">"
					case tt.want:
						s = "*"
					}

					t.Logf("%ss[%d] = %v", s, i, tt.r.s[i].StartIP)
				}
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
	r := NewASList(ASNs)
	wantASNOrder := []int{0, 1, 2, 3, 2, 2, 4}
	wantZoneLen := 7

	wantFind := []struct {
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

	t.Run("Find", func(t *testing.T) {
		for _, s := range wantFind {
			a, b := r.Find(s.ip)
			if s.wantASN == -1 {
				d := AS{}
				if a != d {
					t.Errorf("ASN.Find(%v).ASNumber = %#v, wanted empty", s.ip, a)
				}
				if b != false {
					t.Errorf("ASN.Find(%v) b = %v, want false", s.ip, b)
				}
			} else {
				if a.ASNumber != s.wantASN {
					t.Errorf("ASN.Find(%v).ASNumber = %v, want %v", s.ip, a.ASNumber, s.wantASN)
				}
				if b != true {
					t.Errorf("ASN.Find(%v) b = %v, want true", s.ip, b)
				}
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
