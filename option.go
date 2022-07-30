package asndb

type Option func(r *Registry)

func WithAssumeValid() Option {
	return func(r *Registry) {
		r.assumeValid = true
	}
}

func WithSearchRange(searchRange int) Option {
	return func(r *Registry) {
		r.searchRange = searchRange
	}
}
