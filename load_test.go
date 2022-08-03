package asndb

import (
	"net/netip"
	"strings"
	"testing"
)

func TestDownloadIntegration(t *testing.T) {
	r, err := DownloadFromURL(DownloadViaIpToAsn)
	if err != nil {
		t.Errorf("DownloadFromURL() error = %v", err)
		return
	}
	defer r.Close()

	ls, err := LoadFromTSV(r)
	if err != nil {
		t.Errorf("LoadFromTSV() error = %v", err)
		return
	}
	reg := NewRegistry(ls)
	as, found := reg.Find(netip.MustParseAddr("1.1.1.1"))
	if found == false || as.ASNumber != 13335 {
		t.Errorf("Find() = AS %v Found %t, want 13335, true", as.ASNumber, found)
		return
	}
}

func TestLoadFromTSV(t *testing.T) {
	tests := []struct {
		name      string
		data      string
		wantAS    []AS
		wantError string
	}{
		{
			name: "valid",
			data: "1.0.0.0\t1.0.0.255\t13335\tUS\tCLOUDFLARENET\n1.0.1.0\t1.0.3.255\t0\tNone\tNot routed",
			wantAS: []AS{
				{
					StartIP:       netip.MustParseAddr("1.0.0.0"),
					EndIP:         netip.MustParseAddr("1.0.0.255"),
					ASNumber:      13335,
					CountryCode:   "US",
					ASDescription: "CLOUDFLARENET",
				}, {
					StartIP:       netip.MustParseAddr("1.0.1.0"),
					EndIP:         netip.MustParseAddr("1.0.3.255"),
					ASNumber:      0,
					CountryCode:   "None",
					ASDescription: "Not routed",
				},
			},
		}, {
			name:      "invalid parts",
			data:      "foo\tbar\tbaz",
			wantError: "invalid data #0: want 5 parts got 3",
		}, {
			name:      "invalid start address",
			data:      "foo\tbar\tbaz\tqux\tquux",
			wantError: "invalid start address #0:",
		}, {
			name:      "invalid end address",
			data:      "1.0.0.1\tbar\tbaz\tqux\tquux",
			wantError: "invalid end address #0:",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := strings.NewReader(tt.data)
			ls, err := LoadFromTSV(r)
			if tt.wantError != "" {
				if !strings.HasPrefix(err.Error(), tt.wantError) {
					t.Errorf(`LoadFromTSV() error = %v, want prefix "%v"`, err, tt.wantError)
				}
				return
			}
			if err != nil {
				t.Errorf("LoadFromTSV() error = %v", err)
				return
			}
			if len(ls) != len(tt.wantAS) {
				t.Errorf("LoadFromTSV() = %v, want %v", len(ls), len(tt.wantAS))
				return
			}
			for i, as := range ls {
				if as != tt.wantAS[i] {
					t.Errorf("LoadFromTSV() = %v, want %v", as, tt.wantAS[i])
				}
			}
		})
	}
}
