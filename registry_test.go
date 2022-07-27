package asndb

import (
	"net/netip"
	"reflect"
	"testing"
)

func TestRegistry_Lookup(t *testing.T) {
	r := NewRegistry([]ASN{
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
	},
	)

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
