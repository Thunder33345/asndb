package asndb

import (
	"bufio"
	"compress/gzip"
	"fmt"
	"io"
	"net/http"
	"net/netip"
	"sort"
	"strconv"
	"strings"
)

const ip2dsnLink = "https://iptoasn.com/data/ip2asn-v4.tsv.gz"

func _() {
	dsnMap, err := Load(ip2dsnLink)
	if err != nil {
		panic(err)
	}

	dsn, ok := dsnMap.Lookup(netip.MustParseAddr("xxx.xxx.xxx.xxx"))
	fmt.Printf("dsn: %+v\n", dsn)
	fmt.Printf("ok: %t", ok)
}

func Load(url string) (*Registry, error) {
	rs, err := http.Get(url)
	if err != nil {
		return nil, err
	}

	defer rs.Body.Close()

	gzipReader, err := gzip.NewReader(rs.Body)
	if err != nil {
		return nil, err
	}
	defer gzipReader.Close()

	return LoadReader(gzipReader)
}

func LoadReader(reader io.Reader) (*Registry, error) {
	var s asnList
	buf := bufio.NewScanner(reader)
	for buf.Scan() {
		parts := strings.Split(buf.Text(), "\t")

		asNumber, _ := strconv.Atoi(parts[2])
		s = append(s, ASN{
			StartIP:       netip.MustParseAddr(parts[0]),
			EndIP:         netip.MustParseAddr(parts[1]),
			ASNumber:      asNumber,
			CountryCode:   parts[3],
			ASDescription: parts[4],
		})
	}
	if !sort.IsSorted(s) {
		sort.Sort(s)
	}

	return &Registry{s: s}, nil
}
