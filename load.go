package asndb

import (
	"bufio"
	"compress/gzip"
	"fmt"
	"io"
	"net/http"
	"net/netip"
	"strconv"
	"strings"
)

func LoadFromTSV(reader io.Reader) ([]AS, error) {
	var s []AS
	buf := bufio.NewScanner(reader)
	var i int
	for buf.Scan() {
		parts := strings.Split(buf.Text(), "\t")

		if len(parts) < 5 {
			return s, fmt.Errorf(`invalid data #%d: want 5 parts got %d`, i, len(parts))
		}

		asNumber, _ := strconv.Atoi(parts[2])
		start, err := netip.ParseAddr(parts[0])
		if err != nil {
			return s, fmt.Errorf("invalid start address #%d: %w", i, err)
		}
		end, err := netip.ParseAddr(parts[1])
		if err != nil {
			return s, fmt.Errorf("invalid end address #%d: %w", i, err)
		}
		s = append(s, AS{
			StartIP:       start,
			EndIP:         end,
			ASNumber:      asNumber,
			CountryCode:   parts[3],
			ASDescription: parts[4],
		})
		i++
	}
	return s, nil
}

const DownloadViaIpToAsn = "https://iptoasn.com/data/ip2asn-combined.tsv.gz"

func DownloadFromURL(url string) (io.ReadCloser, error) {
	rs, err := http.Get(url)
	if err != nil {
		return nil, err
	}

	gzipReader, err := gzip.NewReader(rs.Body)
	if err != nil {
		return nil, err
	}

	return gzipReader, nil
}
