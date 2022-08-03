package asndb

import (
	"bufio"
	"io"
	"net/netip"
	"strconv"
	"strings"
)

func LoadFromTSV(reader io.Reader) ([]AS, error) {
	var s []AS
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
		s = append(s, AS{
			StartIP:       start,
			EndIP:         end,
			ASNumber:      asNumber,
			CountryCode:   parts[3],
			ASDescription: parts[4],
		})
	}
	return s, nil
}
