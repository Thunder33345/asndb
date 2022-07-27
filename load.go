package asndb

import (
	"bufio"
	"compress/gzip"
	"io"
	"net/http"
	"net/netip"
	"strconv"
	"strings"
)

const _ = "https://iptoasn.com/data/ip2asn-combined.tsv.gz"

func DownloadFromURL(url string) (io.ReadCloser, error) {
	rs, err := http.Get(url)
	if err != nil {
		return nil, err
	}

	gzipReader, err := gzip.NewReader(rs.Body)
	if err != nil {
		return nil, err
	}
	defer gzipReader.Close()

	return gzipReader, nil
}

func LoadFromTSV(reader io.Reader) (*Registry, error) {
	var s []ASN
	buf := bufio.NewScanner(reader)
	for buf.Scan() {
		parts := strings.Split(buf.Text(), "\t")

		asNumber, _ := strconv.Atoi(parts[2])
		start, err := netip.ParseAddr(parts[0])
		if err != nil {
			return nil, err
		}
		end, err := netip.ParseAddr(parts[1])
		if err != nil {
			return nil, err
		}
		s = append(s, ASN{
			StartIP:       start,
			EndIP:         end,
			ASNumber:      asNumber,
			CountryCode:   parts[3],
			ASDescription: parts[4],
		})
	}
	return NewRegistry(s), nil
}
