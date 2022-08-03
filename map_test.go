package asndb

import (
	"reflect"
	"testing"
)

func TestASNMap_Integration(t *testing.T) {
	type wantASList struct {
		asn       int
		wantCount int
		wantFound bool
	}

	tests := []struct {
		name        string
		asList      []AS
		wantASNLen  int
		wantListASN []AS
		wantASList  []wantASList
	}{
		{
			name: "test1",
			asList: []AS{
				{
					ASNumber:      4,
					ASDescription: "the 4th asn",
				},
				{ASNumber: 1},
				{
					ASNumber:      0,
					ASDescription: "the 0th asn",
					CountryCode:   "nil",
				},
				{ASNumber: 3}, {ASNumber: 2}, {ASNumber: 2}, {ASNumber: 2},
			},
			wantASNLen: 5,
			wantListASN: []AS{
				{
					ASNumber:      0,
					ASDescription: "the 0th asn",
					CountryCode:   "nil",
				},
				{ASNumber: 1}, {ASNumber: 2}, {ASNumber: 3},
				{
					ASNumber:      4,
					ASDescription: "the 4th asn",
				},
			},
			wantASList: []wantASList{
				{asn: -1, wantFound: false},
				{asn: 0, wantCount: 1, wantFound: true},
				{asn: 1, wantCount: 1, wantFound: true},
				{asn: 2, wantCount: 3, wantFound: true},
				{asn: 3, wantCount: 1, wantFound: true},
				{asn: 4, wantCount: 1, wantFound: true},
				{asn: 5, wantFound: false},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := NewASNMap(tt.asList)
			if gotASNLen := len(m.m); gotASNLen != tt.wantASNLen {
				t.Errorf("ASNMap.m len() = %v, want %v", gotASNLen, tt.wantASNLen)
			}

			if gotListASN := m.ListASN(); !reflect.DeepEqual(gotListASN, tt.wantListASN) {
				t.Errorf("ASNMap.ListASN() = %v, want %v", gotListASN, tt.wantListASN)
			}

			l := m.ListASN()
			for i, as := range tt.wantListASN {
				if l[i] != as {
					t.Errorf("ASN.ListASN()[%d] = %v, want %v", i, l[i].ASNumber, as.ASNumber)
				}
			}

			for _, tx := range tt.wantASList {
				gotASNs, gotFound := m.ListAS(tx.asn)
				if gotFound != tx.wantFound {
					t.Errorf("ASN.ASList() found = %v, want %v", gotFound, tx.wantFound)
				}
				if tx.wantFound == false {
					continue
				}
				if len(gotASNs) != tx.wantCount {
					t.Errorf("ASN.ASList() count = %v, want %v", gotASNs, tx.wantCount)
				}
			}

			if asl, found := m.ListAS(0); found {
				asl[0].ASNumber = 1

				if asl2, found2 := m.ListAS(0); found2 {
					if asl2[0].ASNumber != 0 {
						t.Errorf("ASN.ListAS(0)!=0: data should not be alterable")
					}
				} else {
					t.Errorf("ASN.ListAS(0) not found, went missing after alteration")
				}
			}
		})
	}
}
