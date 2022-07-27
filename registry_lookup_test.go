package asndb

import (
	"fmt"
	"net/netip"
	"os"
	"testing"
)

var list *Registry

func loadDB() {
	if list != nil {
		return
	}
	fmt.Printf("Loading database...\n")
	f, err := os.Open("./ip2asn-combined.tsv")
	if err != nil {
		panic(err)
	}
	list, err = LoadReader(f)
	if err != nil {
		panic(err)
	}
	fmt.Printf("Done database...\n")
}

func Benchmark_LookupV4(b *testing.B) {
	loadDB()

	type entry struct {
		addr    netip.Addr
		comment string
	}

	addrList := []entry{
		{netip.MustParseAddr("1.1.1.1"), "cf"},
		{netip.MustParseAddr("223.255.254.5"), "end range"},
		{netip.MustParseAddr("202.1.1.255"), "not routed"},
		{netip.MustParseAddr("31.132.24.5"), ""},
		{netip.MustParseAddr("113.9.3.255"), ""},
	}
	for _, e := range addrList {
		name := e.addr.String()
		if e.comment != "" {
			name += " (" + e.comment + ")"
		}
		b.Run(name, func(b *testing.B) {
			for n := 0; n < b.N; n++ {
				_, _ = list.Lookup(e.addr)
			}
		})
	}
}

func Benchmark_LookupV6(b *testing.B) {
	loadDB()

	type entry struct {
		addr    netip.Addr
		comment string
	}

	addrList := []entry{
		{netip.MustParseAddr("2001:e01::"), "root"},
		{netip.MustParseAddr("2606:4700:4700::1111"), "cf"},
		{netip.MustParseAddr("f4d:c7c7::"), "not routed"},
		{netip.MustParseAddr("2001:4860:4860::8888"), "google"},
	}
	for _, e := range addrList {
		name := e.addr.String()
		if e.comment != "" {
			name += " (" + e.comment + ")"
		}
		b.Run(name, func(b *testing.B) {
			for n := 0; n < b.N; n++ {
				_, _ = list.Lookup(e.addr)
			}
		})
	}
}

func Benchmark_LookupSimple(b *testing.B) {
	loadDB()
	addr := netip.MustParseAddr("1.1.1.1")
	b.Run("1.1.1.1 cf", func(b *testing.B) {
		for n := 0; n < b.N; n++ {
			_, _ = list.Lookup(addr)
		}
	})
	addr = netip.MustParseAddr("223.255.254.5")
	b.Run("223.255.254.5 end range", func(b *testing.B) {
		for n := 0; n < b.N; n++ {
			_, _ = list.Lookup(addr)
		}
	})
}
