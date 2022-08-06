# ASNDB

This is a package that provides an interface to [iptoasn](https://iptoasn.com)'s database.

This is meant to work with `ip2asn-combined.tsv` dataset, it will also work with the IPv4 and IPv6 only datasets.

This package uses `net/netip`.

See [godoc](https://pkg.go.dev/github.com/thunder33345/asndb) for more reference.

## ASList

ASList facilitates looking up AS zone by IP address using `Find(ip)`.

And viewing neighbour AS zones by `Index(ip)` and `FromIndex(index)`.

## ASNMap

ASNMap facilitates looking up AS zones by ASN using `ListAS(asn)`.

And listing all ASN by `ListASN()`
