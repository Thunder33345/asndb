package asndb

type Option func(r *Registry)

func WithAssumeValid() Option {
	return func(r *Registry) {
		r.assumeValid = true
	}
}
